package models

import "mime/multipart"

// File general object contains file details
type File struct {
	Name    string
	UUID    string
	Size    int
	TypeId  string
	UserId  int
	Content *multipart.FileHeader
	Tags    []string
}
