package vojo

type SearchFurnaceLogByType struct {
	FurnaceType *int `form:"furnaceType" binding:"required"`
}
