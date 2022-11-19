package dao

type DownloadDao struct {
	ID         int    `gorm:"primary_key;auto_increment" json:"id"`
	DownloadId string `gorm:"type:varchar(50);not null;index:download_id;comment:'下载excel的Id'" json:"downloadId"`
}

func (p DownloadDao) TableName() string {

	return "download_url"
}
