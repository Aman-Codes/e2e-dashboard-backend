package router

import (
	"github.com/Aman-Codes/e2e-dashboard-backend/pkg/customErrors"
	"github.com/Aman-Codes/e2e-dashboard-backend/pkg/fetchLog"
	"github.com/gin-gonic/gin"
)

func Router() {
	router := gin.Default()
	router.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": customErrors.Success(),
		})
	})
	router.GET("/runs/logs/:id", fetchLog.FetchLogApi)
	router.Run(":8080")
}
