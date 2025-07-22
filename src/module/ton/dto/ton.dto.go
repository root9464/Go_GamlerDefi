package ton_dto

// Manifest represents a manifest
// @swagger:model Manifest
type Manifest struct {
	// URL of the manifest
	// required: true
	// example: https://example.com/manifest.json
	URL string `json:"url"`

	// Name of the manifest
	// required: true
	// example: Example
	Name string `json:"name"`

	// Icon URL of the manifest
	// required: true
	// example: https://example.com/icon.png
	IconURL string `json:"iconUrl"`
}
