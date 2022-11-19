package vojo

type UpdateUserReq struct {
	UserAccount   *string `form:"userAccount" binding:"required"`
	UserPassword  *string `form:"userPassword" binding:"required"`
	AdminUserID   *string `form:"adminUserID" binding:"required"`
	UserId        *string `form:"userId" binding:"required"`
	UserAuthority *int    `form:"userAuthority" binding:"required"`
}
