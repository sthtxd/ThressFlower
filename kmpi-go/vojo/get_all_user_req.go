package vojo

type GetAllUserReq struct {
	AdminUserID *string `form:"adminUserID" binding:"required"`
}
