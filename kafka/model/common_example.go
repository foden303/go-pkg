package model

type UserFileES struct {
	ID                   string  `json:"id"`
	FolderID             *string `json:"folder_id"`
	OwnerID              string  `json:"owner_id"`
	Name                 string  `json:"name"`
	Type                 string  `json:"type"`
	FileHash             *string `json:"file_hash"`
	FileSize             *int64  `json:"file_size"`
	FileType             *string `json:"file_type"`
	FileExt              *string `json:"file_ext"`
	FileMimeType         *string `json:"file_mime_type"`
	FileVideoResolution  *string `json:"file_video_resolution"`
	TotalChunks          int32   `json:"total_chunks"`
	WrappedFekFilePwd    []uint8 `json:"wrapped_fek_file_pwd"`
	WrappedFekAccountPwd []uint8 `json:"wrapped_fek_account_pwd"`
	Shared               bool    `json:"shared"`
	Favorite             bool    `json:"favorite"`
	DeletedAt            *int64  `json:"deleted_at"`
	DeletedBy            *string `json:"deleted_by"`
	CreatedAt            int64   `json:"created_at"`
	UpdatedAt            int64   `json:"updated_at"`
	FullPath             string  `json:"-"`
	Deleted              string  `json:"__deleted"`
}
