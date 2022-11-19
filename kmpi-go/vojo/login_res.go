package vojo

type LoginRes struct {
	UserId    *string `json:"userId"`
	Authority *int    `json:"authority"`
}
