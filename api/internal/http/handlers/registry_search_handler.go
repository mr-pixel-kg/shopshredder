package handlers

import (
	"net/http"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/mr-pixel-kg/shopshredder/api/internal/http/dto"
	"github.com/mr-pixel-kg/shopshredder/api/internal/services"
)

type RegistrySearchHandler struct {
	Search *services.RegistrySearchService
}

func (h RegistrySearchHandler) MountAuthedRoutes(s *fuego.Server) {
	reg := fuego.Group(s, "/registry")
	fuego.Get(reg, "/images/search", h.searchImages,
		option.Summary("Search Docker Hub images"),
		option.Description("Returns matching image repositories from Docker Hub for autocomplete"),
		option.Tags("Registry"),
		option.Query("q", "Search query (e.g. dockware/sh)"),
	)
	fuego.Get(reg, "/tags", h.searchTags,
		option.Summary("Search Docker Hub tags"),
		option.Description("Returns matching tags for a Docker Hub image for autocomplete"),
		option.Tags("Registry"),
		option.Query("image", "Full image name (e.g. dockware/shopware)"),
		option.Query("q", "Tag prefix filter (e.g. 6.7)"),
	)
}

func (h RegistrySearchHandler) searchImages(c fuego.ContextNoBody) (dto.RegistryImageSearchResponse, error) {
	q := c.Request().URL.Query().Get("q")
	if len(q) < 2 {
		return dto.RegistryImageSearchResponse{}, fuego.HTTPError{
			Status: http.StatusBadRequest,
			Detail: "Query must be at least 2 characters",
		}
	}

	results, err := h.Search.SearchImages(c.Request().Context(), q)
	if err != nil {
		return dto.RegistryImageSearchResponse{}, fuego.HTTPError{
			Status: http.StatusInternalServerError,
			Detail: "Image search failed",
		}
	}

	return dto.RegistryImageSearchResponse{Results: results}, nil
}

func (h RegistrySearchHandler) searchTags(c fuego.ContextNoBody) (dto.RegistryTagSearchResponse, error) {
	image := c.Request().URL.Query().Get("image")
	if image == "" {
		return dto.RegistryTagSearchResponse{}, fuego.HTTPError{
			Status: http.StatusBadRequest,
			Detail: "image query parameter is required",
		}
	}

	q := c.Request().URL.Query().Get("q")

	results, err := h.Search.SearchTags(c.Request().Context(), image, q)
	if err != nil {
		return dto.RegistryTagSearchResponse{}, fuego.HTTPError{
			Status: http.StatusInternalServerError,
			Detail: "Tag search failed",
		}
	}

	return dto.RegistryTagSearchResponse{Results: results}, nil
}
