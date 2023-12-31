package gateway

import "bytes"

// swagger:route POST /api/file File upload
// Upload file.
// Security:
//    bearerAuth: []
// responses:
//   200:

// swagger:parameters upload
type UploadFile struct {
	// in:formData
	// swagger:file
	File *bytes.Buffer `json:"files"`
	Tags []string      `json:"tags"`
}

// swagger:route GET /api/file/searchgoogle File searchGoogle
// File searchGoogle
// responses:
//   200: successResponse

// swagger:parameters searchGoogle
type searchGoogleParams struct {
	// in:query
	// type:string
	// required: true
	Query string `json:"query"`

	// in:formData
	// type:file
	// required: true
	File *bytes.Buffer `json:"file"`

	// in:formData
	// type:array
	// items: string
	Tags []string `json:"tags"`
}

// swagger:route GET /api/file File download
// Download file.
// Security:
//    bearerAuth: []
// responses:
//   200:

// swagger:parameters download
type DownloadFile struct {
	// in:formData
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

// swagger:route GET /api/file/list File list
// Its only for admin user.
// Its only for admin user
// Security:
//    bearerAuth: []
// responses:
//   200:
