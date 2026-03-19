package dto

type CreateSandboxRequest struct {
	ImageID    string `json:"imageId"`
	TTLMinutes *int   `json:"ttlMinutes"`
}

type ExtendTTLRequest struct {
	TTLMinutes int `json:"ttlMinutes"`
}

type CreateSnapshotRequest struct {
	ImagePayload
}
