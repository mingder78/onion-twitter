package main

import (
	"encoding/base64"
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type TwitterResource struct {
	db gorm.DB
}

func (tr *TwitterResource) CreateTwitter(c *gin.Context) {
	var twitter Twitter

	if c.Bind(&twitter) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "problem decoding body"})
		return
	}
	//twitter.Status = TwitterStatus
	twitter.Ginger_Created = int32(time.Now().Unix())

	tr.db.Save(&twitter)

	c.JSON(http.StatusCreated, twitter)
}

func (tr *TwitterResource) CreateTwitterByUserId(c *gin.Context) {
	var twitter Twitter
	id, err := tr.getUserId(c)
	//fmt.Println(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "problem decoding id sent"})
		return
	}
	//bind twitter
	if c.Bind(&twitter) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "problem decoding body"})
		return
	}

	var user User

	if tr.db.First(&user, id).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"message": "user id not found"})
	} else {

		//create a new twitter
		twitter.Ginger_Created = int32(time.Now().Unix())
		twitter.UserId = int(id)
		tr.db.NewRecord(twitter)
		tr.db.Create(&twitter)
		tr.db.Save(&twitter)

		user.Twitters = append(user.Twitters, twitter)

		tr.db.Save(&user)
		tr.db.Model(&user).Update("twitters", twitter)
		c.JSON(http.StatusOK, twitter)
	}
}

func decodeBasicAuth(token string, tr *TwitterResource) (User, bool) {
	var user User
	encoded := strings.Fields(token)
	result, _ := base64.StdEncoding.DecodeString(encoded[1])
	userPass := strings.Split(string(result[:]), ":")
	userName := userPass[0]
	tr.db.Where("name = ?", userName).First(&user)
	return user, true
}

func (tr *TwitterResource) CreateTwitterWithoutUserId(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	user, _ := decodeBasicAuth(token, tr)

	var twitter Twitter
	//bind twitter
	if c.Bind(&twitter) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "problem decoding body"})
		return
	}

	//create a new twitter
	twitter.Ginger_Created = int32(time.Now().Unix())
	twitter.UserId = int(user.Ginger_Id)
	tr.db.NewRecord(twitter)
	tr.db.Create(&twitter)
	tr.db.Save(&twitter)

	user.Twitters = append(user.Twitters, twitter)

	tr.db.Save(&user)
	tr.db.Model(&user).Update("twitters", twitter)
	c.JSON(http.StatusOK, twitter)
}

func (tr *TwitterResource) GetAllTwitters(c *gin.Context) {
	var twitters []Twitter

	tr.db.Order("ginger__created desc").Find(&twitters)

	c.JSON(http.StatusOK, twitters)
}

func (tr *TwitterResource) GetTwittersByUserId(c *gin.Context) {
	id, err := tr.getId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "problem decoding id sent"})
		return
	}

	var twitters []Twitter
	var user User

	if tr.db.First(&user, id).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
	} else {
		tr.db.Model(&user).Related(&twitters)
		c.JSON(http.StatusOK, twitters)
	}
}

func (tr *TwitterResource) GetTwitter(c *gin.Context) {
	id, err := tr.getId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "problem decoding id sent"})
		return
	}

	var twitter Twitter

	if tr.db.First(&twitter, id).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
	} else {
		c.JSON(http.StatusOK, twitter)
	}
}

func (tr *TwitterResource) UpdateTwitter(c *gin.Context) {
	id, err := tr.getId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "problem decoding id sent"})
		return
	}

	var twitter Twitter

	if c.Bind(&twitter) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "problem decoding body"})
		return
	}
	twitter.Ginger_Id = int32(id)

	var existing Twitter

	if tr.db.First(&existing, id).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
	} else {
		tr.db.Save(&twitter)
		c.JSON(http.StatusOK, twitter)
	}

}

func (tr *TwitterResource) PatchTwitter(c *gin.Context) {
	id, err := tr.getId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "problem decoding id sent"})
		return
	}

	// this is a hack because Gin falsely claims my unmarshalled obj is invalid.
	// recovering from the panic and using my object that already has the json body bound to it.
	var json []Patch

	r := c.Bind(&json)
	if r != nil {
		fmt.Println(r)
	} else {
		if json[0].Op != "replace" && json[0].Path != "/status" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "PATCH support is limited and can only replace the /status path"})
			return
		}
		var twitter Twitter

		if tr.db.First(&twitter, id).RecordNotFound() {
			c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		} else {
			//twitter.Status = json[0].Value

			tr.db.Save(&twitter)
			c.JSON(http.StatusOK, twitter)
		}
	}
}

func (tr *TwitterResource) DeleteTwitter(c *gin.Context) {
	id, err := tr.getId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "problem decoding id sent"})
		return
	}

	var twitter Twitter

	if tr.db.First(&twitter, id).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
	} else {
		tr.db.Delete(&twitter)
		c.Data(http.StatusNoContent, "application/json", make([]byte, 0))
	}
}

func (tr *TwitterResource) getId(c *gin.Context) (int32, error) {
	idStr := c.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Print(err)
		return 0, err
	}
	return int32(id), nil
}

/**
* on patching: http://williamdurand.fr/2014/02/14/please-do-not-patch-like-an-idiot/
 *
  * patch specification https://tools.ietf.org/html/rfc5789
   * json definition http://tools.ietf.org/html/rfc6902
*/

type Patch struct {
	Op    string `json:"op" binding:"required"`
	From  string `json:"from"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

func (tr *TwitterResource) CreateUser(c *gin.Context) {
	var user User

	if c.Bind(&user) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "problem decoding body"})
		return
	}
	//user.Status = UserStatus
	user.Ginger_Created = int32(time.Now().Unix())
	user.Twitters = []Twitter{}
	tr.db.Save(&user)

	c.JSON(http.StatusCreated, user)
}

func (tr *TwitterResource) GetUsers() []User {
	var users []User

	tr.db.Order("ginger__created desc").Find(&users)
	return users
}

func (tr *TwitterResource) GetAllUsers(c *gin.Context) {
	var users []User

	tr.db.Order("ginger__created desc").Find(&users)
	for index, user := range users {
		var twitters []Twitter
		tr.db.Model(&user).Related(&twitters)
		users[index].Twitters = twitters

	}
	c.JSON(http.StatusOK, users)
}

func (tr *TwitterResource) GetUser(c *gin.Context) {
	id, err := tr.getUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "problem decoding id sent"})
		return
	}
	var user User
	if tr.db.First(&user, id).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
	} else {
		var twitters []Twitter
		tr.db.Model(&user).Related(&twitters)
		user.Twitters = twitters
		c.JSON(http.StatusOK, user)
	}
}

func (tr *TwitterResource) UpdateUser(c *gin.Context) {
	id, err := tr.getUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "problem decoding id sent"})
		return
	}

	var user User

	if c.Bind(&user) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "problem decoding body"})
		return
	}
	user.Ginger_Id = int32(id)

	var existing User

	if tr.db.First(&existing, id).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
	} else {
		tr.db.Save(&user)
		c.JSON(http.StatusOK, user)
	}

}

func (tr *TwitterResource) PatchUser(c *gin.Context) {
	id, err := tr.getUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "problem decoding id sent"})
		return
	}

	// this is a hack because Gin falsely claims my unmarshalled obj is invalid.
	// recovering from the panic and using my object that already has the json body bound to it.
	var json []Patch

	r := c.Bind(&json)
	if r != nil {
		fmt.Println(r)
	} else {
		if json[0].Op != "replace" && json[0].Path != "/status" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "PATCH support is limited and can only replace the /status path"})
			return
		}
		var user User

		if tr.db.First(&user, id).RecordNotFound() {
			c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		} else {
			//user.Status = json[0].Value

			tr.db.Save(&user)
			c.JSON(http.StatusOK, user)
		}
	}
}

func (tr *TwitterResource) DeleteUser(c *gin.Context) {
	id, err := tr.getUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "problem decoding id sent"})
		return
	}

	var user User

	if tr.db.First(&user, id).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
	} else {
		tr.db.Delete(&user)
		c.Data(http.StatusNoContent, "application/json", make([]byte, 0))
	}
}

func (tr *TwitterResource) getUserId(c *gin.Context) (int32, error) {
	idStr := c.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Print(err)
		return 0, err
	}
	return int32(id), nil
}

func (tr *TwitterResource) SwaggerCity(c *gin.Context) {
	c.HTML(http.StatusOK, "swagger", gin.H{
		"lowerCase": "city",
		"url":       "api.log4security.com:30194",
		"dataType": `       "name": {
                    "type": "string"
                },
                "age": {
                    "type": "integer",
                    "format": "int"
                },
                "address": {
                    "type": "string"
                }`,
	})
}
