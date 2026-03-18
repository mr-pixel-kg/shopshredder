package dto

type CreateSandboxRequest struct {
	ImageID    string `json:"imageId"`
	TTLMinutes *int   `json:"ttlMinutes"`
}

type SnapshotRequest struct {
	TargetImage string `json:"targetImage"`
}
