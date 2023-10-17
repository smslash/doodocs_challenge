package entity

type FileInfo struct {
	FilePath string  `json:"file_path"`
	Size     float64 `json:"size"`
	Mimetype string  `json:"mimetype"`
}
