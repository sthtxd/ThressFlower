package vojo

type GetTemperature struct {
	MaxTemperature *int `json:"maxTemperature"`
	MinTemperature *int `json:"minTemperature"`
}
