package vojo

type CheckDownloadIdReq struct {
	DownloadId *string `form:"downloadId" binding:"required"`
}
