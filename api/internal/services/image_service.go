package services

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/manuel/shopware-testenv-platform/api/internal/docker"
	"github.com/manuel/shopware-testenv-platform/api/internal/models"
	"github.com/manuel/shopware-testenv-platform/api/internal/repositories"
)

const (
	ThumbnailPublicBasePath = "/thumbnails"
)

var ErrUnsupportedThumbnailFormat = errors.New("unsupported thumbnail format")

type ImageService struct {
	repo         *repositories.ImageRepository
	docker       docker.Client
	thumbnailDir string
}

func NewImageService(repo *repositories.ImageRepository, dockerClient docker.Client, thumbnailDir string) *ImageService {
	service := &ImageService{
		repo:         repo,
		docker:       dockerClient,
		thumbnailDir: thumbnailDir,
	}

	if err := os.MkdirAll(service.thumbnailDir, 0o755); err != nil {
		slog.Error("create thumbnail directory failed", "path", service.thumbnailDir, "cause", err.Error())
	}

	return service
}

func (s *ImageService) ListPublic() ([]models.Image, error) {
	images, err := s.repo.ListPublic()
	if err != nil {
		return nil, err
	}
	return s.attachThumbnailURLs(images), nil
}

func (s *ImageService) ListAll() ([]models.Image, error) {
	images, err := s.repo.ListAll()
	if err != nil {
		return nil, err
	}
	return s.attachThumbnailURLs(images), nil
}

func (s *ImageService) Create(image *models.Image) error {
	image.ID = uuid.New()
	return s.repo.Create(image)
}

func (s *ImageService) FindByID(id uuid.UUID) (*models.Image, error) {
	image, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.attachThumbnailURL(image), nil
}

func (s *ImageService) CreateForUser(
	ctx context.Context,
	userID *uuid.UUID,
	name string,
	tag string,
	title *string,
	description *string,
	isPublic bool,
) (*models.Image, error) {
	fullName := name + ":" + tag
	// Registering an image in the database should only succeed if Docker can
	// actually resolve and fetch that image reference.
	if err := s.docker.EnsureImage(ctx, fullName); err != nil {
		return nil, err
	}

	image := &models.Image{
		ID:              uuid.New(),
		Name:            name,
		Tag:             tag,
		Title:           title,
		Description:     description,
		IsPublic:        isPublic,
		CreatedByUserID: userID,
	}

	if err := s.repo.Create(image); err != nil {
		return nil, err
	}

	return s.attachThumbnailURL(image), nil
}

func (s *ImageService) Update(id uuid.UUID, title, description *string, isPublic bool) (*models.Image, error) {
	image, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	image.Title = title
	image.Description = description
	image.IsPublic = isPublic
	if err := s.repo.Update(image); err != nil {
		return nil, err
	}

	return s.attachThumbnailURL(image), nil
}

func (s *ImageService) SaveThumbnail(id uuid.UUID, file multipart.File, originalFilename, contentType string) (*models.Image, error) {
	image, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	ext, err := thumbnailExtension(file, originalFilename, contentType)
	if err != nil {
		return nil, err
	}

	if err := s.deleteThumbnailFiles(id); err != nil {
		return nil, err
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	targetPath := filepath.Join(s.thumbnailDir, id.String()+ext)
	dst, err := os.Create(targetPath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return nil, err
	}

	return s.attachThumbnailURL(image), nil
}

func (s *ImageService) DeleteThumbnail(id uuid.UUID) (*models.Image, error) {
	image, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if err := s.deleteThumbnailFiles(id); err != nil {
		return nil, err
	}

	return s.attachThumbnailURL(image), nil
}

func (s *ImageService) Delete(ctx context.Context, id uuid.UUID) error {
	image, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// Remove the local Docker image first so the database does not claim an
	// image exists after the runtime artifact was already cleaned up.
	if err := s.docker.RemoveImage(ctx, image.FullName()); err != nil {
		return err
	}

	if err := s.deleteThumbnailFiles(id); err != nil {
		return err
	}

	return s.repo.Delete(id)
}

func (s *ImageService) attachThumbnailURLs(images []models.Image) []models.Image {
	for i := range images {
		s.attachThumbnailURL(&images[i])
	}
	return images
}

func (s *ImageService) attachThumbnailURL(image *models.Image) *models.Image {
	path, err := s.thumbnailPath(image.ID)
	if err != nil {
		slog.Error("resolve thumbnail path failed", "image_id", image.ID.String(), "cause", err.Error())
		image.ThumbnailURL = nil
		return image
	}
	if path == "" {
		image.ThumbnailURL = nil
		return image
	}

	url := ThumbnailPublicBasePath + "/" + filepath.Base(path)
	image.ThumbnailURL = &url
	return image
}

func (s *ImageService) thumbnailPath(id uuid.UUID) (string, error) {
	matches, err := filepath.Glob(filepath.Join(s.thumbnailDir, id.String()+".*"))
	if err != nil {
		return "", err
	}
	if len(matches) == 0 {
		return "", nil
	}
	return matches[0], nil
}

func (s *ImageService) deleteThumbnailFiles(id uuid.UUID) error {
	matches, err := filepath.Glob(filepath.Join(s.thumbnailDir, id.String()+".*"))
	if err != nil {
		return err
	}

	for _, match := range matches {
		if err := os.Remove(match); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	return nil
}

func thumbnailExtension(file multipart.File, originalFilename, contentType string) (string, error) {
	buffer := make([]byte, 512)
	readBytes, err := file.Read(buffer)
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}

	detectedType := http.DetectContentType(buffer[:readBytes])
	if ext := extensionForContentType(detectedType); ext != "" {
		return ext, nil
	}
	if ext := extensionForContentType(contentType); ext != "" {
		return ext, nil
	}

	switch strings.ToLower(filepath.Ext(originalFilename)) {
	case ".jpg", ".jpeg":
		return ".jpg", nil
	case ".png":
		return ".png", nil
	case ".gif":
		return ".gif", nil
	case ".webp":
		return ".webp", nil
	default:
		return "", ErrUnsupportedThumbnailFormat
	}
}

func extensionForContentType(contentType string) string {
	switch strings.ToLower(contentType) {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	default:
		return ""
	}
}
