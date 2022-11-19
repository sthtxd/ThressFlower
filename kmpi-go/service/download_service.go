package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"kmpi-go/config"
	"kmpi-go/dao"
	"kmpi-go/vojo"
)

func CheckDownloadId(c *gin.Context) (*int, error) {

	var req vojo.CheckDownloadIdReq
	err := c.Bind(&req)
	if err != nil {
		return nil, err
	}
	var res dao.DownloadDao
	err = dao.CronDb.Where("download_id=?", req.DownloadId).First(&res).Error
	if err != nil {
		return nil, err
	}
	status := 1
	return &status, nil

}
func DownloadExcel(c *gin.Context) error {
	downloadId := c.Query("downloadId")
	if downloadId == "" {
		return fmt.Errorf("imageName is not illgeal")
	}
dst :=config.FileSaveDir + "/" + downloadId + ".xlsx"
//	dst := filepath.Join(config.FileSaveDir, videoId, videoId)
	//dst := filepath.Join(config.FileSaveDir, "b.mp4")

	//c.Header("content-type","video/mp4");
	c.File(dst)
	return nil

}
