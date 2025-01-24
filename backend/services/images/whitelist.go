package images

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

type DockerImageWhitelist struct {
	allowedImages map[string]struct{}
}

func NewDockerImageWhitelist() *DockerImageWhitelist {
	return &DockerImageWhitelist{
		allowedImages: make(map[string]struct{}),
	}
}

func (w *DockerImageWhitelist) Add(image string) {
	w.allowedImages[image] = struct{}{}
}

func (w *DockerImageWhitelist) Remove(image string) error {
	if _, exists := w.allowedImages[image]; exists {
		delete(w.allowedImages, image)
		return nil
	}
	return errors.New("image not found in whitelist")
}

func (w *DockerImageWhitelist) List() []string {
	images := make([]string, 0, len(w.allowedImages))
	for image := range w.allowedImages {
		images = append(images, image)
	}
	return images
}

func (w *DockerImageWhitelist) IsAllowed(image string) bool {
	_, exists := w.allowedImages[image]
	return exists
}

func (w *DockerImageWhitelist) SaveToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for image := range w.allowedImages {
		_, err := writer.WriteString(image + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

func (w *DockerImageWhitelist) LoadFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			w.allowedImages = make(map[string]struct{}) // Initialisiere leere Liste
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	w.allowedImages = make(map[string]struct{})
	for scanner.Scan() {
		image := strings.TrimSpace(scanner.Text())
		if image != "" {
			w.allowedImages[image] = struct{}{}
		}
	}
	return scanner.Err()
}
