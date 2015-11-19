package main

type Twitter struct {
	Ginger_Created int32 `json:"ginger_created"`
	Ginger_Id      int32 `json:"ginger_id" gorm:"primary_key"`

	UserId  int    `json:"user_id"`
	Message string `json:"message"`
}

type User struct {
	Ginger_Created int32 `json:"ginger_created"`
	Ginger_Id      int32 `json:"ginger_id" gorm:"primary_key"`

	Name     string    `json:"name"`
	Password string    `json:"password"`
	Twitters []Twitter `json:"twitters"`
}
