
//go:generate ginger $GOFILE
package main

//@ginger
type Twitter struct {
	Ginger_Created int32 `json:"ginger_created"`
	Ginger_Id int32 `json:"ginger_id" gorm:"primary_key"`

	User_id float64 `json:"user_id"`
	Message string `json:"message"`
}