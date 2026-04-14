package registry

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/gobwas/glob"
)

type compiledEntry struct {
	glob  glob.Glob
	entry ImageEntry
}

type Resolver struct {
	entries []compiledEntry
}

func NewResolver(reg *ImageRegistry) (*Resolver, error) {
	entries := make([]compiledEntry, 0, len(reg.Images))
	for _, e := range reg.Images {
		g, err := glob.Compile(e.Match)
		if err != nil {
			return nil, fmt.Errorf("compile glob %q: %w", e.Match, err)
		}
		entries = append(entries, compiledEntry{glob: g, entry: e})
	}
	return &Resolver{entries: entries}, nil
}

func (r *Resolver) Resolve(imageName string, ctx TemplateContext) (*ResolvedImage, error) {
	nameOnly := stripTag(imageName)

	for _, ce := range r.entries {
		if !ce.glob.Match(imageName) && !ce.glob.Match(nameOnly) {
			continue
		}
		return renderEntry(ce.entry, ctx)
	}

	return &ResolvedImage{}, nil
}

func (r *Resolver) ResolveEntry(imageName string) *ImageEntry {
	nameOnly := stripTag(imageName)
	for _, ce := range r.entries {
		if ce.glob.Match(imageName) || ce.glob.Match(nameOnly) {
			return &ce.entry
		}
	}
	return nil
}

func stripTag(image string) string {
	if strings.Contains(image, "@") {
		return image
	}
	for i := len(image) - 1; i >= 0; i-- {
		if image[i] == ':' {
			return image[:i]
		}
		if image[i] == '/' {
			break
		}
	}
	return image
}

func renderEntry(entry ImageEntry, ctx TemplateContext) (*ResolvedImage, error) {
	if ctx.Meta == nil {
		ctx.Meta = make(map[string]string)
	}
	for _, it := range entry.Metadata.Items {
		if it.Type != "field" || it.Field == nil || it.Field.Default == "" {
			continue
		}
		if _, ok := ctx.Meta[it.Key]; !ok {
			ctx.Meta[it.Key] = it.Field.Default
		}
	}

	if entry.SSH != nil {
		ctx.SSHPort = strconv.Itoa(entry.SSH.Port)
		ctx.SSHUsername = entry.SSH.Username
		ctx.SSHPassword = entry.SSH.Password
	}

	resolved := &ResolvedImage{
		HealthCheck: entry.HealthCheck,
	}

	if entry.InternalPort != nil {
		resolved.InternalPort = *entry.InternalPort
	}

	keys := make([]string, 0, len(entry.Env))
	for k := range entry.Env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		rendered, err := renderTemplate(k+"="+entry.Env[k], ctx)
		if err != nil {
			return nil, fmt.Errorf("render env %s: %w", k, err)
		}
		resolved.Env = append(resolved.Env, rendered)
	}

	if len(entry.Labels) > 0 {
		resolved.Labels = make(map[string]string, len(entry.Labels))
		for k, v := range entry.Labels {
			rendered, err := renderTemplate(v, ctx)
			if err != nil {
				return nil, fmt.Errorf("render label %s: %w", k, err)
			}
			resolved.Labels[k] = rendered
		}
	}

	for _, cmd := range entry.PostStart {
		rendered, err := renderExecCommand(cmd, ctx)
		if err != nil {
			return nil, fmt.Errorf("render post_start command: %w", err)
		}
		resolved.PostStart = append(resolved.PostStart, rendered)
	}

	for _, cmd := range entry.PreStop {
		rendered, err := renderExecCommand(cmd, ctx)
		if err != nil {
			return nil, fmt.Errorf("render pre_stop command: %w", err)
		}
		resolved.PreStop = append(resolved.PreStop, rendered)
	}

	resolved.SSH = entry.SSH

	for _, ls := range entry.Logs {
		rendered := LogSource{
			Key:   ls.Key,
			Label: ls.Label,
			Type:  ls.Type,
		}
		if ls.Path != "" {
			p, err := renderTemplate(ls.Path, ctx)
			if err != nil {
				return nil, fmt.Errorf("render log path %s: %w", ls.Key, err)
			}
			rendered.Path = p
		}
		resolved.Logs = append(resolved.Logs, rendered)
	}

	return resolved, nil
}

func renderExecCommand(cmd ExecCommand, ctx TemplateContext) (ExecCommand, error) {
	rendered := ExecCommand{
		Label:      cmd.Label,
		Delay:      cmd.Delay,
		Timeout:    cmd.Timeout,
		Retries:    cmd.Retries,
		RetryDelay: cmd.RetryDelay,
	}
	for _, arg := range cmd.Command {
		r, err := renderTemplate(arg, ctx)
		if err != nil {
			return rendered, err
		}
		rendered.Command = append(rendered.Command, r)
	}
	return rendered, nil
}

func (r *Resolver) RenderMetadata(imageName string, values map[string]string, ctx TemplateContext) (*MetadataSchema, error) {
	entry := r.ResolveEntry(imageName)
	if entry == nil {
		return &MetadataSchema{}, nil
	}
	return RenderMetadata(&entry.Metadata, values, ctx)
}

func RenderMetadata(schema *MetadataSchema, values map[string]string, ctx TemplateContext) (*MetadataSchema, error) {
	if schema == nil {
		return &MetadataSchema{}, nil
	}
	out := &MetadataSchema{
		Groups: append([]MetadataGroup(nil), schema.Groups...),
		Items:  make([]MetadataItem, len(schema.Items)),
	}

	meta := make(map[string]string)
	for _, it := range schema.Items {
		if it.Type == "field" && it.Field != nil && it.Field.Default != "" {
			meta[it.Key] = it.Field.Default
		}
	}
	defined := make(map[string]bool, len(schema.Items))
	for _, it := range schema.Items {
		if it.Type == "field" {
			defined[it.Key] = true
		}
	}
	for k, v := range values {
		if defined[k] {
			meta[k] = v
		}
	}
	ctx.Meta = meta

	for i, it := range schema.Items {
		clone := it
		switch it.Type {
		case "field":
			if it.Field != nil {
				f := *it.Field
				if v, ok := values[it.Key]; ok {
					f.Value = v
				} else {
					f.Value = it.Field.Default
				}
				clone.Field = &f
			}
		case "action":
			if it.Action != nil {
				a := *it.Action
				if it.Action.urlTmpl != nil {
					rendered, err := execTemplate(it.Action.urlTmpl, ctx)
					if err != nil {
						return nil, fmt.Errorf("render action.url for %q: %w", it.Key, err)
					}
					a.URL = rendered
				}
				clone.Action = &a
			}
		case "display":
			if it.Display != nil {
				d := *it.Display
				if it.Display.valueTmpl != nil {
					rendered, err := execTemplate(it.Display.valueTmpl, ctx)
					if err != nil {
						return nil, fmt.Errorf("render display.value for %q: %w", it.Key, err)
					}
					d.Value = rendered
				}
				clone.Display = &d
			}
		}
		out.Items[i] = clone
	}

	return out, nil
}

func execTemplate(t *template.Template, ctx TemplateContext) (string, error) {
	var buf bytes.Buffer
	if err := t.Execute(&buf, ctx); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func renderTemplate(text string, ctx TemplateContext) (string, error) {
	tmpl, err := template.New("").Option("missingkey=error").Parse(text)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, ctx); err != nil {
		return "", err
	}
	return buf.String(), nil
}
