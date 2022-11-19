package vojo

type LoginReq struct {
	Name     *string `form:"name" binding:"required"`
	Password *string `form:"password" binding:"required"`
}
