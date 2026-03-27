package handlers

import (
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"
	"github.com/manuel/shopware-testenv-platform/api/internal/models"
	"github.com/manuel/shopware-testenv-platform/api/internal/registry"
	"gorm.io/datatypes"
)

func mergeRegistryAndDB(reg []registry.MetadataItem, dbJSON datatypes.JSON) []registry.MetadataItem {
	var dbItems []registry.MetadataItem
	if len(dbJSON) > 0 {
		_ = json.Unmarshal(dbJSON, &dbItems)
	}

	dbMap := make(map[string]registry.MetadataItem)
	for _, item := range dbItems {
		dbMap[item.Key] = item
	}

	var merged []registry.MetadataItem
	for _, item := range reg {
		if dbItem, ok := dbMap[item.Key]; ok {
			if dbItem.Value != "" {
				item.Value = dbItem.Value
			}
			delete(dbMap, item.Key)
		}
		merged = append(merged, item)
	}
	for _, item := range dbItems {
		if _, ok := dbMap[item.Key]; ok {
			merged = append(merged, item)
		}
	}

	return merged
}

func (h *SandboxHandler) enrichSandboxMetadata(sandboxes []models.Sandbox) {
	imageCache := make(map[uuid.UUID][]registry.MetadataItem)

	for idx := range sandboxes {
		sb := &sandboxes[idx]
		imgMeta, ok := imageCache[sb.ImageID]
		if !ok {
			img, err := h.images.FindByID(sb.ImageID)
			if err != nil {
				slog.Warn("enrich sandbox metadata: image not found", "image_id", sb.ImageID)
				imageCache[sb.ImageID] = nil
			} else {
				reg := h.resolver.ResolveMetadata(img.RegistryName())
				imgMeta = mergeRegistryAndDB(reg, img.Metadata)
				imageCache[sb.ImageID] = imgMeta
			}
		}

		if imgMeta == nil {
			continue
		}

		var values map[string]string
		if len(sb.Metadata) > 0 {
			_ = json.Unmarshal(sb.Metadata, &values)
		}

		enriched := make([]registry.MetadataItem, len(imgMeta))
		copy(enriched, imgMeta)
		for j := range enriched {
			if v, exists := values[enriched[j].Key]; exists {
				enriched[j].Value = v
			}
		}

		data, _ := json.Marshal(enriched)
		sb.Metadata = datatypes.JSON(data)
	}
}

func (h *SandboxHandler) enrichSandbox(sandbox *models.Sandbox) {
	sl := []models.Sandbox{*sandbox}
	h.enrichSandboxMetadata(sl)
	*sandbox = sl[0]
}
