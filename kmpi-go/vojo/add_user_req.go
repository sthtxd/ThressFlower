package vojo

type AddUserReq struct {
	UserAccount   *string `form:"userAccount" binding:"required"`
	UserPassword  *string `form:"userPassword" binding:"required"`
	AdminUserID   *string `form:"adminUserID" binding:"required"`
	UserAuthority *int    `form:"userAuthority" binding:"required"`
}
