package postapi

import (
	"net/http"

	"github.com/Wy0t/DcardGo/adstruct"
	"github.com/Wy0t/DcardGo/database"
	"github.com/Wy0t/DcardGo/tools"
	"github.com/gin-gonic/gin"
)

// Admin API用於產生管理廣告資源
func PostAD(c *gin.Context) {
	// 檢查是否超過每天的最大廣告創建數量
	if database.IsOverMaxAdsPerDay() {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Exceeded maximum ads creation per day"})
		return
	}
	var newAD adstruct.AD
	if err := c.BindJSON(&newAD); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// 檢查廣告的條件是否合法
	if err := tools.ValidateConditions(newAD.Conditions); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 插入條件信息到 Conditions 資料表
	conditionsID, err := database.InsertConditions(newAD.Conditions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert conditions"})
		return
	}

	// 插入廣告信息到 ADs 資料表
	err = database.InsertAD(newAD, conditionsID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert ad"})
		return
	}

	// 增加今天的廣告創建數量
	database.IncrementAdCreationCount()

	adstruct.ADs = append(adstruct.ADs, newAD)

	c.IndentedJSON(http.StatusCreated, newAD)
}
