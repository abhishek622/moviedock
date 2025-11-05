package model

type Metadata struct {
	MetadataID  int64  `json:"metadata_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Director    string `json:"director"`
}
