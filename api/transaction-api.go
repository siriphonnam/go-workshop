package api

import "github.com/gin-gonic/gin"

func SetupTransactionAPI(router *gin.Engine) {
	productAPI := router.Group("api/v2")
	{
		productAPI.GET("/transaction", getTransaction)
		productAPI.POST("/transaction", createTranscation)
	}
}

func getTransaction(c *gin.Context) {
	c.JSON(200, gin.H{"result": "list transaction"})
}

func createTranscation(c *gin.Context) {
	c.JSON(200, gin.H{"result": "create transaction"})
}
