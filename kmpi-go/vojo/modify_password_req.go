package vojo

type ModifyPasswordReq struct {
	UserAccount     *string `form:"userAccount" binding:"required"`
	UserOldPassword *string `form:"userOldPassword" binding:"required"  json:"userOldPassword"`
	UserNewPassword *string `form:"userNewPassword" binding:"required"  json:"userNewPassword"`
}
