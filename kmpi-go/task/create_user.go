package task

import (
	uuid "github.com/satori/go.uuid"
	"kmpi-go/dao"
	"kmpi-go/log"
)

func init() {
	err := initSuperUser()
	if err != nil {
		log.Error("initSuperUser error", err)
	}
}
func initSuperUser() error {

	var count int64
	err := dao.CronDb.Model(&dao.UserDao{}).Where("user_authority = ?", dao.SUPER_ADMIN_AUTHORITY).Count(&count).Error
	if err != nil {
		return err
	}
	uuid := uuid.NewV4().String()

	if count == 0 {
		userObj := &dao.UserDao{
			UserAccount:   "admin",
			UserPassword:  "2hHY(07}sqF",
			UserAuthority: dao.SUPER_ADMIN_AUTHORITY,
			UserId:        uuid,
		}
		err = dao.CronDb.Create(userObj).Error
	}
	return err
}
