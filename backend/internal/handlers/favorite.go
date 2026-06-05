package handlers

import (
	"ninimenu/internal/database"
	"ninimenu/internal/models"
	"ninimenu/internal/services"
	"ninimenu/internal/utils"

	"github.com/gin-gonic/gin"
)

func GetFavorites(c *gin.Context) {
	var favorites []models.Favorite
	database.DB.Order("created_at DESC").Find(&favorites)

	var dishIDs []uint
	for _, f := range favorites {
		dishIDs = append(dishIDs, f.DishID)
	}

	var dishes []models.Dish
	if len(dishIDs) > 0 {
		database.DB.Where("id IN ?", dishIDs).Find(&dishes)
	}

	utils.Success(c, dishes)
}

func AddFavorite(c *gin.Context) {
	dishID := c.Param("dishId")
	var dish models.Dish
	if err := database.DB.First(&dish, dishID).Error; err != nil {
		utils.NotFound(c, "菜品不存在")
		return
	}

	var existing models.Favorite
	if err := database.DB.Where("dish_id = ?", dish.ID).First(&existing).Error; err == nil {
		utils.SuccessMsg(c, "已收藏")
		return
	}

	fav := models.Favorite{DishID: dish.ID}
	database.DB.Create(&fav)
	database.DB.Model(&dish).Update("favorite", true)
	services.QueueAutoAchievementSync()

	utils.SuccessMsg(c, "收藏成功")
}

func RemoveFavorite(c *gin.Context) {
	dishID := c.Param("dishId")

	database.DB.Where("dish_id = ?", dishID).Delete(&models.Favorite{})
	database.DB.Model(&models.Dish{}).Where("id = ?", dishID).Update("favorite", false)
	services.QueueAutoAchievementSync()

	utils.SuccessMsg(c, "取消收藏成功")
}
