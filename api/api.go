package api

import (
	"goworkshop/db"

	"github.com/gin-gonic/gin"
)

func Setup(router *gin.Engine) {

	db.SetupDB()
	SetupAuthenAPI(router)
	SetupProduceAPI(router)
	SetupTransactionAPI(router)
}
