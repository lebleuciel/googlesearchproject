package models

// FileType general object contains filetype details
type FileType struct {
	Name        string
	AllowedSize int
	IsBanned    bool
}
