package dto

type RegistryImageSuggestion struct {
	Name        string `json:"name" example:"dockware/dev"`
	Description string `json:"description" example:"Dockware development image for Shopware"`
	StarCount   int    `json:"starCount" example:"42"`
	IsOfficial  bool   `json:"isOfficial" example:"false"`
}

type RegistryTagSuggestion struct {
	Name        string `json:"name" example:"6.7.5"`
	FullSize    int64  `json:"fullSize,omitempty" example:"1073741824"`
	LastUpdated string `json:"lastUpdated,omitempty" example:"2026-03-15T10:00:00Z"`
}

type RegistryImageSearchResponse struct {
	Results []RegistryImageSuggestion `json:"results"`
}

type RegistryTagSearchResponse struct {
	Results []RegistryTagSuggestion `json:"results"`
}
