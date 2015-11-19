package main

type Twitter struct {
	Ginger_Created int32 `json:"ginger_created"`
	Ginger_Id      int32 `json:"ginger_id" gorm:"primary_key"`

	User_Id int32  `json:"user_id"`
	Message string `json:"message"`
}

type TwitterId struct {
	Twitter_id int32 `json:"twitter_id"`
}

type User struct {
	Ginger_Created int32 `json:"ginger_created"`
	Ginger_Id      int32 `json:"ginger_id" gorm:"primary_key"`

	Name     string `json:"name"`
	Password string `json:"password"`
}
