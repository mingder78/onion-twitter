package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"strconv"
	"time"
)

type TwitterResource struct {
	db gorm.DB
}

func (tr *TwitterResource) CreateTwitter(c *gin.Context) {
	var twitter Twitter

	if !c.Bind(&twitter) {
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
	if !c.Bind(&twitter) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "problem decoding body"})
		return
	}

	var user User

	if tr.db.First(&user, id).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"message": "user id not found"})
	} else {

		//create a new twitter
		twitter.Ginger_Created = int32(time.Now().Unix())
		twitter.User_id = id
		tr.db.NewRecord(twitter)
		tr.db.Create(&twitter)
		b := TwitterId{twitter.Ginger_Id}
		tr.db.Save(&twitter)

		user.Twitters = append(user.Twitters, b)
		spew.Dump(twitter)
		spew.Dump(user)
		tr.db.Save(&user)
		c.JSON(http.StatusOK, user)
	}
}

func (tr *TwitterResource) GetAllTwitters(c *gin.Context) {
	var twitters []Twitter

	tr.db.Order("ginger__created desc").Find(&twitters)

	c.JSON(http.StatusOK, twitters)
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

	if !c.Bind(&twitter) {
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
	if !r {
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

	if !c.Bind(&user) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "problem decoding body"})
		return
	}
	//user.Status = UserStatus
	user.Ginger_Created = int32(time.Now().Unix())
	user.Twitters = make([]TwitterId, 0)
	tr.db.Save(&user)

	c.JSON(http.StatusCreated, user)
}

func (tr *TwitterResource) GetAllUsers(c *gin.Context) {
	var users []User

	tr.db.Order("ginger__created desc").Find(&users)

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

	if !c.Bind(&user) {
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
	if !r {
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
