package services

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/manuel/shopware-testenv-platform/api/internal/docker"
	"github.com/manuel/shopware-testenv-platform/api/internal/models"
	"github.com/manuel/shopware-testenv-platform/api/internal/repositories"
)

type ImageService struct {
	repo        *repositories.ImageRepository
	sandboxRepo *repositories.SandboxRepository
	docker      docker.Client
	tracker     *docker.PullTracker

	mu          sync.Mutex
	pullCancels map[string]context.CancelFunc
}

func NewImageService(
	repo *repositories.ImageRepository,
	sandboxRepo *repositories.SandboxRepository,
	dockerClient docker.Client,
	tracker *docker.PullTracker,
) *ImageService {
	return &ImageService{
		repo:        repo,
		sandboxRepo: sandboxRepo,
		docker:      dockerClient,
		tracker:     tracker,
		pullCancels: make(map[string]context.CancelFunc),
	}
}

func (s *ImageService) ListPublic() ([]models.Image, error) {
	return s.repo.ListPublic()
}

func (s *ImageService) ListAll() ([]models.Image, error) {
	images, err := s.repo.ListAll()
	if err != nil {
		return nil, err
	}
	for i := range images {
		if images[i].Status == "pulling" {
			p := s.tracker.Progress(images[i].ID.String())
			images[i].PullProgress = &p.Percent
		}
	}
	return images, nil
}

func (s *ImageService) Create(image *models.Image) error {
	image.ID = uuid.New()
	return s.repo.Create(image)
}

func (s *ImageService) FindByID(id uuid.UUID) (*models.Image, error) {
	return s.repo.FindByID(id)
}

func (s *ImageService) CreateForUser(
	ctx context.Context,
	userID *uuid.UUID,
	name string,
	tag string,
	title *string,
	description *string,
	thumbnailURL *string,
	isPublic bool,
) (*models.Image, error) {
	fullName := name + ":" + tag

	status := "pulling"
	if s.docker.ImageExists(ctx, fullName) {
		status = "ready"
	}

	img := &models.Image{
		ID:              uuid.New(),
		Name:            name,
		Tag:             tag,
		Title:           title,
		Description:     description,
		ThumbnailURL:    thumbnailURL,
		IsPublic:        isPublic,
		Status:          status,
		CreatedByUserID: userID,
	}

	if err := s.repo.Create(img); err != nil {
		return nil, err
	}

	if status == "pulling" {
		pullCtx, cancel := context.WithCancel(context.Background())
		s.mu.Lock()
		s.pullCancels[img.ID.String()] = cancel
		s.mu.Unlock()
		go s.pullImage(pullCtx, img.ID, fullName)
	}

	return img, nil
}

func (s *ImageService) pullImage(ctx context.Context, imageID uuid.UUID, fullName string) {
	idStr := imageID.String()
	s.tracker.Start(idStr)

	defer func() {
		s.mu.Lock()
		delete(s.pullCancels, idStr)
		s.mu.Unlock()
		time.AfterFunc(10*time.Second, func() { s.tracker.Remove(idStr) })
	}()

	reader, err := s.docker.PullImage(ctx, fullName)
	if err != nil {
		if ctx.Err() != nil {
			slog.Info("image pull cancelled", "image_id", idStr, "image", fullName)
			s.tracker.Finish(idStr, err)
			return
		}
		slog.Error("image pull failed", "image_id", idStr, "image", fullName, "error", err.Error())
		errMsg := err.Error()
		_ = s.repo.UpdateStatus(imageID, "failed", &errMsg)
		s.tracker.Finish(idStr, err)
		return
	}
	defer reader.Close()

	if err := s.tracker.ConsumePullStream(idStr, reader); err != nil {
		if ctx.Err() != nil {
			slog.Info("image pull cancelled", "image_id", idStr, "image", fullName)
			s.tracker.Finish(idStr, err)
			return
		}
		slog.Error("image pull stream failed", "image_id", idStr, "image", fullName, "error", err.Error())
		errMsg := err.Error()
		_ = s.repo.UpdateStatus(imageID, "failed", &errMsg)
		s.tracker.Finish(idStr, err)
		return
	}

	slog.Info("image pull complete", "image_id", idStr, "image", fullName)
	_ = s.repo.UpdateStatus(imageID, "ready", nil)
	s.tracker.Finish(idStr, nil)
}

func (s *ImageService) WatchPullProgress(imageID string) (<-chan docker.PullProgress, func()) {
	return s.tracker.Watch(imageID)
}

func (s *ImageService) cancelPull(imageID string) {
	s.mu.Lock()
	cancel, ok := s.pullCancels[imageID]
	s.mu.Unlock()
	if ok {
		cancel()
	}
}

func (s *ImageService) Delete(ctx context.Context, id uuid.UUID) error {
	img, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// Cancel any in-progress pull.
	if img.Status == "pulling" {
		s.cancelPull(id.String())
	}

	// Delete all sandboxes that use this image first.
	sandboxes, err := s.sandboxRepo.ListByImageID(id)
	if err != nil {
		return err
	}
	for _, sb := range sandboxes {
		if sb.Status == models.SandboxStatusStarting || sb.Status == models.SandboxStatusRunning {
			if err := s.docker.DeleteContainer(ctx, sb.ContainerID); err != nil {
				slog.Warn("failed to delete sandbox container during image deletion", "container_id", sb.ContainerID, "error", err.Error())
			}
		}
		if err := s.sandboxRepo.DeleteByID(sb.ID); err != nil {
			slog.Warn("failed to delete sandbox during image deletion", "sandbox_id", sb.ID.String(), "error", err.Error())
		}
	}

	// Remove Docker image if it exists locally (covers both ready and partial pulls).
	if s.docker.ImageExists(ctx, img.FullName()) {
		if err := s.docker.RemoveImage(ctx, img.FullName()); err != nil {
			slog.Warn("docker image removal failed, proceeding with db deletion", "image", img.FullName(), "error", err.Error())
		}
	}

	return s.repo.Delete(id)
}

func (s *ImageService) RecoverStalePulls() {
	errMsg := "Pull interrupted by server restart"
	if err := s.repo.ResetStalePulls("pulling", "failed", &errMsg); err != nil {
		slog.Error("failed to recover stale image pulls", "error", err.Error())
	}
}
