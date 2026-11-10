package model

type PresignedUploadRequest struct {
	Path          string `json:"path"`
	Folder        string `json:"folder"`
	SubFolder     string `json:"sub_folder"`
	Filename      string `json:"filename"`
	Year          string `json:"year"`
	Month         string `json:"month"`
	ExpiryMinutes int    `json:"expiry_minutes"`
}

type PresignedUploadResponse struct {
	UploadURL        string `json:"upload_url"`
	ObjectKey        string `json:"object_key"`
	FullPath         string `json:"full_path"`
	Bucket           string `json:"bucket"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}
