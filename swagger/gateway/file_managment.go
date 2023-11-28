package gateway

// swagger:route GET /api/file File download
// Download file.
// responses:
//   200: Getfile

// swagger:response Getfile
type GetFileResponse struct {
	// in:body
	DiskSpace int
}
