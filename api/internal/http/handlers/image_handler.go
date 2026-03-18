package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/manuel/shopware-testenv-platform/api/internal/apperror"
	"github.com/manuel/shopware-testenv-platform/api/internal/http/dto"
	mw "github.com/manuel/shopware-testenv-platform/api/internal/http/middleware"
	"github.com/manuel/shopware-testenv-platform/api/internal/http/responses"
	"github.com/manuel/shopware-testenv-platform/api/internal/logging"
	"github.com/manuel/shopware-testenv-platform/api/internal/services"
	"gorm.io/gorm"
)

type ImageHandler struct {
	images *services.ImageService
	audit  *services.AuditService
}

func NewImageHandler(images *services.ImageService, audit *services.AuditService) *ImageHandler {
	return &ImageHandler{images: images, audit: audit}
}

func (h *ImageHandler) ListPublic(c echo.Context) error {
	images, err := h.images.ListPublic()
	if err != nil {
		return responses.FromAppError(c, apperror.Internal("IMAGE_LIST_FAILED", "Could not load public images").WithCause(err))
	}
	slog.Info("listed public images", logging.RequestFields(c, "count", len(images))...)
	return c.JSON(200, images)
}

func (h *ImageHandler) ListAll(c echo.Context) error {
	images, err := h.images.ListAll()
	if err != nil {
		return responses.FromAppError(c, apperror.Internal("IMAGE_LIST_FAILED", "Could not load images").WithCause(err))
	}
	slog.Info("listed all images", logging.RequestFields(c, "count", len(images))...)
	return c.JSON(200, images)
}

func (h *ImageHandler) Create(c echo.Context) error {
	var input dto.CreateImageRequest
	if err := c.Bind(&input); err != nil {
		return responses.FromAppError(c, apperror.BadRequest("VALIDATION_ERROR", "Invalid request body"))
	}

	auth := mw.MustAuth(c)
	slog.Info("image creation requested", logging.RequestFields(c,
		"user_id", auth.UserID.String(),
		"name", input.Name,
		"tag", input.Tag,
		"is_public", input.IsPublic,
	)...)
	image, err := h.images.CreateForUser(
		c.Request().Context(),
		&auth.UserID,
		input.Name,
		input.Tag,
		input.Title,
		input.Description,
		input.IsPublic,
	)
	if err != nil {
		return responses.FromAppError(c, apperror.BadRequest("IMAGE_CREATE_FAILED", err.Error()).WithCause(err))
	}

	slog.Info("image created successfully", logging.RequestFields(c,
		"user_id", auth.UserID.String(),
		"image_id", image.ID.String(),
		"image", image.FullName(),
		"is_public", image.IsPublic,
	)...)
	_ = h.audit.Log(&auth.UserID, "image.created", c.RealIP(), map[string]any{"imageId": image.ID.String()})
	return c.JSON(201, image)
}

func (h *ImageHandler) Update(c echo.Context) error {
	auth := mw.MustAuth(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return responses.FromAppError(c, apperror.BadRequest("VALIDATION_ERROR", "Invalid image id"))
	}

	var input dto.UpdateImageRequest
	if err := c.Bind(&input); err != nil {
		return responses.FromAppError(c, apperror.BadRequest("VALIDATION_ERROR", "Invalid request body"))
	}

	slog.Info("image update requested", logging.RequestFields(c,
		"user_id", auth.UserID.String(),
		"image_id", id.String(),
		"is_public", input.IsPublic,
	)...)
	image, err := h.images.Update(id, input.Title, input.Description, input.IsPublic)
	if err != nil {
		return mapImageError(c, "IMAGE_UPDATE_FAILED", "Could not update image", err)
	}

	slog.Info("image updated successfully", logging.RequestFields(c,
		"user_id", auth.UserID.String(),
		"image_id", image.ID.String(),
		"is_public", image.IsPublic,
		"has_thumbnail", image.ThumbnailURL != nil,
	)...)
	_ = h.audit.Log(&auth.UserID, "image.updated", c.RealIP(), map[string]any{"imageId": image.ID.String()})
	return c.JSON(http.StatusOK, image)
}

func (h *ImageHandler) UploadThumbnail(c echo.Context) error {
	auth := mw.MustAuth(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return responses.FromAppError(c, apperror.BadRequest("VALIDATION_ERROR", "Invalid image id"))
	}

	fileHeader, err := c.FormFile("thumbnail")
	if err != nil {
		return responses.FromAppError(c, apperror.BadRequest("VALIDATION_ERROR", "Missing thumbnail upload"))
	}

	file, err := fileHeader.Open()
	if err != nil {
		return responses.FromAppError(c, apperror.Internal("THUMBNAIL_UPLOAD_FAILED", "Could not open thumbnail upload").WithCause(err))
	}
	defer file.Close()

	slog.Info("thumbnail upload requested", logging.RequestFields(c,
		"user_id", auth.UserID.String(),
		"image_id", id.String(),
		"filename", fileHeader.Filename,
		"size", fileHeader.Size,
	)...)
	image, err := h.images.SaveThumbnail(id, file, fileHeader.Filename, fileHeader.Header.Get(echo.HeaderContentType))
	if err != nil {
		if errors.Is(err, services.ErrUnsupportedThumbnailFormat) {
			return responses.FromAppError(c, apperror.BadRequest("THUMBNAIL_FORMAT_UNSUPPORTED", "Unsupported thumbnail format").WithCause(err))
		}
		return mapImageError(c, "THUMBNAIL_UPLOAD_FAILED", "Could not store thumbnail", err)
	}

	slog.Info("thumbnail uploaded successfully", logging.RequestFields(c,
		"user_id", auth.UserID.String(),
		"image_id", image.ID.String(),
		"thumbnail_url", image.ThumbnailURL,
	)...)
	_ = h.audit.Log(&auth.UserID, "image.thumbnail_uploaded", c.RealIP(), map[string]any{"imageId": image.ID.String()})
	return c.JSON(http.StatusOK, image)
}

func (h *ImageHandler) DeleteThumbnail(c echo.Context) error {
	auth := mw.MustAuth(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return responses.FromAppError(c, apperror.BadRequest("VALIDATION_ERROR", "Invalid image id"))
	}

	slog.Info("thumbnail deletion requested", logging.RequestFields(c, "user_id", auth.UserID.String(), "image_id", id.String())...)
	image, err := h.images.DeleteThumbnail(id)
	if err != nil {
		return mapImageError(c, "THUMBNAIL_DELETE_FAILED", "Could not delete thumbnail", err)
	}

	slog.Info("thumbnail deleted successfully", logging.RequestFields(c,
		"user_id", auth.UserID.String(),
		"image_id", image.ID.String(),
		"has_thumbnail", image.ThumbnailURL != nil,
	)...)
	_ = h.audit.Log(&auth.UserID, "image.thumbnail_deleted", c.RealIP(), map[string]any{"imageId": image.ID.String()})
	return c.NoContent(http.StatusNoContent)
}

func (h *ImageHandler) Delete(c echo.Context) error {
	auth := mw.MustAuth(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return responses.FromAppError(c, apperror.BadRequest("VALIDATION_ERROR", "Invalid image id"))
	}

	slog.Info("image deletion requested", logging.RequestFields(c, "user_id", auth.UserID.String(), "image_id", id.String())...)
	if err := h.images.Delete(c.Request().Context(), id); err != nil {
		return responses.FromAppError(c, apperror.Internal("IMAGE_DELETE_FAILED", "Could not delete image").WithCause(err))
	}

	slog.Info("image deleted successfully", logging.RequestFields(c, "user_id", auth.UserID.String(), "image_id", id.String())...)
	_ = h.audit.Log(&auth.UserID, "image.deleted", c.RealIP(), map[string]any{"imageId": id.String()})
	return c.NoContent(204)
}

func mapImageError(c echo.Context, code, message string, err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return responses.FromAppError(c, apperror.NotFound("IMAGE_NOT_FOUND", "Image not found").WithCause(err))
	}

	return responses.FromAppError(c, apperror.Internal(code, message).WithCause(err))
}
