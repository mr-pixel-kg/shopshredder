package dto

type ImagePayload struct {
	Name        string  `json:"name"`
	Tag         string  `json:"tag"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	IsPublic    bool    `json:"isPublic"`
}

type CreateImageRequest struct {
	ImagePayload
}

type UpdateImageRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	IsPublic    bool    `json:"isPublic"`
}
