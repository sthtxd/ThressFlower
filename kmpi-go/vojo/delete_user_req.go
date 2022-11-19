package vojo

type DeleteUserReq struct {
	UserID      *string `form:"userID" binding:"required"`
	AdminUserID *string `form:"adminUserID" binding:"required"`
}
