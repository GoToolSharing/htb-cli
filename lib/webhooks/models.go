package webhooks

// EmbedField represents a single field in an embed.
type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type EmbedThumbnail struct {
	URL string `json:"url"`
}

// Embed represents a Discord message embed.
type Embed struct {
	Title       string         `json:"title,omitempty"`
	Description string         `json:"description,omitempty"`
	URL         string         `json:"url,omitempty"`
	Color       int            `json:"color,omitempty"`
	Fields      []EmbedField   `json:"fields,omitempty"`
	Thumbnail   EmbedThumbnail `json:"thumbnail,omitempty"`
}
