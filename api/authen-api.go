package api

import (
	"goworkshop/db"
	"goworkshop/interceptor"
	"goworkshop/model"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func SetupAuthenAPI(router *gin.Engine) {
	authenAPI := router.Group("api/v2")
	{
		authenAPI.POST("/login", login)
		authenAPI.POST("/register", register)
	}
}

func login(c *gin.Context) {
	var user model.User
	if c.ShouldBindJSON(&user) == nil {
		var queryUser model.User
		if err := db.GetDB().First(&queryUser, "username=?", user.Username).Error; err != nil {
			c.JSON(200, gin.H{"result": "nok", "error": err})
		} else if !checkPasswordHash(user.Password, queryUser.Password) {
			c.JSON(200, gin.H{"result": "nok", "error": "invalid password"})
		} else {
			c.JSON(200, gin.H{"result": "ok", "token": interceptor.JwtSign(queryUser)})
		}
	} else {
		c.JSON(401, gin.H{"status": "unable to bind data"})
	}
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func register(c *gin.Context) {
	var user model.User
	if c.ShouldBindJSON(&user) == nil {
		user.Password, _ = hashPassword(user.Password)
		user.CreateAt = time.Now()
		if err := db.GetDB().Create(&user).Error; err != nil {
			panic(err)
		} else {
			c.JSON(200, gin.H{"result": "register", "data": user})
		}
		c.JSON(200, gin.H{"result": "register", "data": user})
	} else {
		c.JSON(401, gin.H{"status": "unable to bind data"})
	}
}
