package handlers

import (
	"strconv"

	"golang-api/config"
	"golang-api/models"

	"github.com/gin-gonic/gin"
)

func GetFactories(c *gin.Context) {

	var data []models.Factory
	var total int64

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	search := c.Query("search")
	status := c.Query("status")

	offset := (page - 1) * limit

	query := config.DB.Model(&models.Factory{})

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	query.
		Limit(limit).
		Offset(offset).
		Find(&data)

	c.JSON(200, gin.H{
		"data":  data,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}