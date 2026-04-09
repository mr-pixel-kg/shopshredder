package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/mr-pixel-kg/shopshredder/api/internal/http/dto"
)

const (
	dockerHubBaseURL = "https://hub.docker.com/v2"
	imageCacheTTL    = 5 * time.Minute
	tagCacheTTL      = 15 * time.Minute
	imageSearchLimit = 15
	tagSearchLimit   = 50
	dockerHubTimeout = 5 * time.Second
)

type cachedResult[T any] struct {
	data      T
	expiresAt time.Time
}

type RegistrySearchService struct {
	client     *http.Client
	imageCache sync.Map
	tagCache   sync.Map
}

func NewRegistrySearchService() *RegistrySearchService {
	return &RegistrySearchService{
		client: &http.Client{Timeout: dockerHubTimeout},
	}
}

func (s *RegistrySearchService) SearchImages(ctx context.Context, query string) ([]dto.RegistryImageSuggestion, error) {
	cacheKey := "images:" + query
	if cached, ok := s.imageCache.Load(cacheKey); ok {
		entry := cached.(cachedResult[[]dto.RegistryImageSuggestion])
		if time.Now().Before(entry.expiresAt) {
			return entry.data, nil
		}
		s.imageCache.Delete(cacheKey)
	}

	u := fmt.Sprintf("%s/search/repositories?query=%s&page_size=%d",
		dockerHubBaseURL, url.QueryEscape(query), imageSearchLimit)

	results, err := s.fetchImageResults(ctx, u)
	if err != nil {
		slog.Warn("docker hub image search failed", "query", query, "error", err)
		return []dto.RegistryImageSuggestion{}, nil
	}

	s.imageCache.Store(cacheKey, cachedResult[[]dto.RegistryImageSuggestion]{
		data:      results,
		expiresAt: time.Now().Add(imageCacheTTL),
	})

	return results, nil
}

func (s *RegistrySearchService) SearchTags(ctx context.Context, image, tagQuery string) ([]dto.RegistryTagSuggestion, error) {
	cacheKey := fmt.Sprintf("tags:%s:%s", image, tagQuery)
	if cached, ok := s.tagCache.Load(cacheKey); ok {
		entry := cached.(cachedResult[[]dto.RegistryTagSuggestion])
		if time.Now().Before(entry.expiresAt) {
			return entry.data, nil
		}
		s.tagCache.Delete(cacheKey)
	}

	ns, repo := parseImageRef(image)

	u := fmt.Sprintf("%s/namespaces/%s/repositories/%s/tags?page_size=%d&ordering=last_updated",
		dockerHubBaseURL, url.PathEscape(ns), url.PathEscape(repo), tagSearchLimit)
	if tagQuery != "" {
		u += "&name=" + url.QueryEscape(tagQuery)
	}

	results, err := s.fetchTagResults(ctx, u, tagQuery)
	if err != nil {
		slog.Warn("docker hub tag search failed", "image", image, "query", tagQuery, "error", err)
		return []dto.RegistryTagSuggestion{}, nil
	}

	s.tagCache.Store(cacheKey, cachedResult[[]dto.RegistryTagSuggestion]{
		data:      results,
		expiresAt: time.Now().Add(tagCacheTTL),
	})

	return results, nil
}

func parseImageRef(image string) (namespace, repo string) {
	parts := strings.SplitN(image, "/", 2)
	if len(parts) == 1 {
		return "library", parts[0]
	}
	return parts[0], parts[1]
}

type hubSearchResponse struct {
	Results []hubSearchResult `json:"results"`
}

type hubSearchResult struct {
	RepoName    string `json:"repo_name"`
	Description string `json:"short_description"`
	StarCount   int    `json:"star_count"`
	IsOfficial  bool   `json:"is_official"`
}

type hubTagResponse struct {
	Results []hubTagResult `json:"results"`
}

type hubTagResult struct {
	Name        string `json:"name"`
	FullSize    int64  `json:"full_size"`
	LastUpdated string `json:"tag_last_pushed"`
}

func (s *RegistrySearchService) fetchImageResults(ctx context.Context, u string) ([]dto.RegistryImageSuggestion, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("docker hub request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("docker hub returned status %d", resp.StatusCode)
	}

	var hubResp hubSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&hubResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	results := make([]dto.RegistryImageSuggestion, len(hubResp.Results))
	for i, r := range hubResp.Results {
		results[i] = dto.RegistryImageSuggestion{
			Name:        r.RepoName,
			Description: r.Description,
			StarCount:   r.StarCount,
			IsOfficial:  r.IsOfficial,
		}
	}
	return results, nil
}

func (s *RegistrySearchService) fetchTagResults(ctx context.Context, u, prefix string) ([]dto.RegistryTagSuggestion, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("docker hub request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("docker hub returned status %d", resp.StatusCode)
	}

	var hubResp hubTagResponse
	if err := json.NewDecoder(resp.Body).Decode(&hubResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	results := make([]dto.RegistryTagSuggestion, 0, len(hubResp.Results))
	for _, r := range hubResp.Results {
		if prefix != "" && !strings.HasPrefix(r.Name, prefix) {
			continue
		}
		results = append(results, dto.RegistryTagSuggestion{
			Name:        r.Name,
			FullSize:    r.FullSize,
			LastUpdated: r.LastUpdated,
		})
	}
	return results, nil
}
