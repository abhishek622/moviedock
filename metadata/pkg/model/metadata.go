package model

type Metadata struct {
	MetadataID  int32  `json:"metadata_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Director    string `json:"director"`
	Runtime     int32  `json:"runtime"`
}
