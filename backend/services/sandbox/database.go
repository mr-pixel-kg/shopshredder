package sandbox

import (
	"errors"
	"sync"
	"time"
)

type Sandbox struct {
	ID            string    `json:"id"`
	ImageName     string    `json:"image_name"`
	ImageId       string    `json:"image_id"`
	ContainerName string    `json:"container_name"`
	ContainerId   string    `json:"container_id"`
	Url           string    `json:"url"`
	CreatedAt     time.Time `json:"created_at"`
	LifeTime      int64     `json:"lifetime"`
}

type SandboxDatabase struct {
	mutex          sync.RWMutex
	sandboxStorage map[string]Sandbox
}

func NewSandboxDatabase() *SandboxDatabase {
	return &SandboxDatabase{
		sandboxStorage: make(map[string]Sandbox),
	}
}

func (db *SandboxDatabase) Add(sandbox Sandbox) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if _, exists := db.sandboxStorage[sandbox.ID]; exists {
		return errors.New("Sandbox environment with this ID already exists")
	}

	db.sandboxStorage[sandbox.ID] = sandbox
	return nil
}

func (db *SandboxDatabase) GetByID(id string) (Sandbox, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	image, exists := db.sandboxStorage[id]
	if !exists {
		return Sandbox{}, errors.New("Sandbox environment not found")
	}

	return image, nil
}

func (db *SandboxDatabase) Remove(id string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if _, exists := db.sandboxStorage[id]; !exists {
		return errors.New("Sandbox environment not found")
	}

	delete(db.sandboxStorage, id)
	return nil
}

func (db *SandboxDatabase) List() []Sandbox {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	images := make([]Sandbox, 0, len(db.sandboxStorage))
	for _, img := range db.sandboxStorage {
		images = append(images, img)
	}
	return images
}
