package service

import (
	"fmt"
	"kmpi-go/dao"
	"kmpi-go/vojo"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func AdminLogin(c *gin.Context) (*vojo.AdminLoginRes, error) {
	var req vojo.AdminLoginReq
	err := c.Bind(&req)
	if err != nil {
		return nil, err
	}
	var res dao.UserDao
	err = dao.CronDb.Where("user_account=?", req.AdminAccount).Where("user_password=?", req.AdminPassword).First(&res).Error
	if err != nil {
		return nil, err
	}
	if res.UserAuthority != dao.SUPER_ADMIN_AUTHORITY {
		return nil, fmt.Errorf("not admin")
	}
	adminRes := vojo.AdminLoginRes{}
	adminRes.UserId = &res.UserId
	return &adminRes, nil
}
func Login(c *gin.Context) (*vojo.LoginRes, error) {
	var req vojo.LoginReq
	err := c.Bind(&req)
	if err != nil {
		return nil, err
	}
	var res dao.UserDao
	err = dao.CronDb.Where("user_account=?", req.Name).Where("user_password=?", req.Password).First(&res).Error
	if err != nil {
		return nil, err
	}
	loginRes := vojo.LoginRes{}
	loginRes.UserId = &res.UserId
	loginRes.Authority = &res.UserAuthority
	return &loginRes, nil
}
func DeleteUser(c *gin.Context) error {
	var req vojo.DeleteUserReq
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	var res dao.UserDao
	err = dao.CronDb.Where("user_id=?", req.AdminUserID).First(&res).Error
	if err != nil {
		return err
	}
	if res.UserAuthority != dao.SUPER_ADMIN_AUTHORITY {
		return fmt.Errorf("权限不允许")
	}

	err = dao.CronDb.Where("user_id=?", req.UserID).Delete(dao.UserDao{}).Error
	return err

}
func AddUser(c *gin.Context) error {
	var req vojo.AddUserReq
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	var res dao.UserDao
	err = dao.CronDb.Where("user_id=?", req.AdminUserID).First(&res).Error
	if err != nil {
		return err
	}
	if res.UserAuthority != dao.SUPER_ADMIN_AUTHORITY {
		return fmt.Errorf("权限不允许")
	}
	uuid := uuid.NewV4().String()
	userObj := &dao.UserDao{
		UserAccount:   *req.UserAccount,
		UserPassword:  *req.UserPassword,
		UserAuthority: *req.UserAuthority,
		UserId:        uuid,
	}
	return dao.CronDb.Create(userObj).Error

}
func UpdateUser(c *gin.Context) error {
	var req vojo.UpdateUserReq
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	var res dao.UserDao
	err = dao.CronDb.Where("user_id=?", req.AdminUserID).First(&res).Error
	if err != nil {
		return err
	}
	if res.UserAuthority != dao.SUPER_ADMIN_AUTHORITY {
		return fmt.Errorf("权限不允许")
	}

	userObj := &dao.UserDao{
		UserAccount:   *req.UserAccount,
		UserPassword:  *req.UserPassword,
		UserAuthority: *req.UserAuthority,
	}
	err = dao.CronDb.Model(&dao.UserDao{}).Where("user_id=?", req.UserId).Updates(userObj).Error
	return err
}
func GetAllUser(c *gin.Context) ([]dao.UserDao, error) {
	var req vojo.GetAllUserReq
	err := c.Bind(&req)
	if err != nil {
		return nil, err
	}
	var res dao.UserDao
	err = dao.CronDb.Where("user_id=?", req.AdminUserID).First(&res).Error
	if err != nil {
		return nil, err
	}
	if res.UserAuthority != dao.SUPER_ADMIN_AUTHORITY {
		return nil, fmt.Errorf("权限不允许")
	}

	var users []dao.UserDao
	err = dao.CronDb.Where("user_authority!=100").Find(&users).Error
	return users, err
}
func ModifyPassword(c *gin.Context) error {
	var req vojo.ModifyPasswordReq
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	var res dao.UserDao
	err = dao.CronDb.Where("user_account=? and user_password=?", req.UserAccount, req.UserOldPassword).First(&res).Error
	if err != nil {
		return err
	}
	if res.UserAccount == "" {
		return fmt.Errorf("账号或密码错误")
	}
	if res.UserAuthority == dao.SUPER_ADMIN_AUTHORITY {
		return fmt.Errorf("超级管理员不允许改密码")
	}

	data := map[string]interface{}{
		"user_password": *req.UserNewPassword}

	// realData := dao.HoldingFurnaceDao{
	// 	MaxTemperature:    *req.MaxTemperature,
	// 	MinTemperature:    *req.MinTemperature,
	// 	TargetTemperature: *req.TargetTemperature,
	// }
	var userDao dao.UserDao
	err = dao.CronDb.Model(&userDao).Where("user_account=?", req.UserAccount).Updates(data).Error
	return err
}
