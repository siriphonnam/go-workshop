package api

import (
	"fmt"
	"goworkshop/db"
	"goworkshop/interceptor"
	"goworkshop/model"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupProduceAPI(router *gin.Engine) {
	productAPI := router.Group("api/v2")
	{
		productAPI.GET("/product", interceptor.JwtVerify, getProduct)
		productAPI.POST("/product", createProduct)
	}
}

func getProduct(c *gin.Context) {
	c.JSON(200, gin.H{"result": "get product", "username": c.GetString("jwt_username")})
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func saveImage(image *multipart.FileHeader, product *model.Product, c *gin.Context) {
	if image != nil {
		runningDir, _ := os.Getwd()
		product.Image = image.Filename
		extension := filepath.Ext(image.Filename)
		fileName := fmt.Sprintf("%d%s", product.ID, extension)
		filePath := fmt.Sprintf(`%s/uploaded/images/%s`, runningDir, fileName)

		if fileExists(filePath) {
			os.Remove(filePath)
		}

		c.SaveUploadedFile(image, filePath)
		db.GetDB().Model(&product).Update("image", fileName)
	}
}

func createProduct(c *gin.Context) {

	//create product
	product := model.Product{}
	product.Name = c.PostForm(`name`)
	product.Stock, _ = strconv.ParseInt(c.PostForm(`stock`), 10, 64)
	product.Price, _ = strconv.ParseFloat(c.PostForm(`price`), 64)
	product.CreateAt = time.Now()
	db.GetDB().Create(&product)

	//save image
	image, _ := c.FormFile(`image`)
	saveImage(image, &product, c)

	c.JSON(200, gin.H{"result": product})
}
