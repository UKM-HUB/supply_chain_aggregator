package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReportResponse struct {
	TotalTransaction int `json:"total_transaction"`
	TotalPaid        int `json:"total_paid"`
	TotalPending     int `json:"total_pending"`
}

func dailyReport(c *gin.Context) {
	response := ReportResponse{
		TotalTransaction: 100,
		TotalPaid:        50000000,
		TotalPending:     5,
	}

	c.JSON(http.StatusOK, response)
}

func monthlyReport(c *gin.Context) {
	response := ReportResponse{
		TotalTransaction: 2500,
		TotalPaid:        750000000,
		TotalPending:     20,
	}

	c.JSON(http.StatusOK, response)
}

func exportReport(c *gin.Context) {
	response := gin.H{
		"message": "Report exported successfully",
		"data": ReportResponse{
			TotalTransaction: 100,
			TotalPaid:        50000000,
			TotalPending:     5,
		},
	}

	c.JSON(http.StatusOK, response)
}

func main() {
	router := gin.Default()

	api := router.Group("/api/v1")
	{
		reports := api.Group("/reports")
		{
			reports.GET("/daily", dailyReport)
			reports.GET("/monthly", monthlyReport)
			reports.GET("/export", exportReport)
		}
	}

	router.Run(":8080")
}