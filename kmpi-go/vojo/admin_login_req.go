package vojo

type AdminLoginReq struct {
	AdminAccount  *string `form:"adminAccount" binding:"required"`
	AdminPassword *string `form:"adminPassword" binding:"required"  json:"adminPassword"`
}
