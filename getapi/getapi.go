package getapi

import (
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/Wy0t/DcardGo/adstruct"
	"github.com/Wy0t/DcardGo/tools"
	"github.com/gin-gonic/gin"
)

// 投放API
func GetAD(c *gin.Context) {
	// offset & limit
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset > 100 || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset value"})
		return
	}
	_ = offset

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit value"})
		return
	}

	// 截取符合 offset 和 limit 的 ADs 數據
	var filteredADs []adstruct.AD
	if offset < len(adstruct.ADs) {
		end := offset + limit
		if end > len(adstruct.ADs) {
			end = len(adstruct.ADs)
		}
		filteredADs = adstruct.ADs[offset:end]
	} else {
		filteredADs = []adstruct.AD{}
	}

	// 支持查詢參數Age,Gender,Country,Platform
	var ageStr = c.Query("age")
	var genderStr = c.Query("gender")
	var countryStr = c.Query("country")
	var platformStr = c.Query("platform")

	// 篩選符合條件的廣告
	var resultADs []adstruct.AD
	for _, ad := range filteredADs {
		if (ageStr == "" || (ad.Conditions.Age != nil && strconv.Itoa(*ad.Conditions.Age) == ageStr)) &&
			(genderStr == "" || (ad.Conditions.Gender != nil && *ad.Conditions.Gender == genderStr)) &&
			(countryStr == "" || (ad.Conditions.Country != nil && tools.ContainsString(*ad.Conditions.Country, countryStr))) &&
			(platformStr == "" || (ad.Conditions.Platform != nil && tools.ContainsString(*ad.Conditions.Platform, platformStr))) {
			resultADs = append(resultADs, ad)
		}
	}

	// 按照 EndAt 升序排序
	sort.Slice(resultADs, func(i, j int) bool {
		return resultADs[i].EndAt.Before(resultADs[j].EndAt)
	})

	// 限制活躍廣告數量不超過1000
	var activeAds []adstruct.AD
	for _, ad := range resultADs {
		if ad.StartAt.Before(time.Now()) && ad.EndAt.After(time.Now()) {
			if len(activeAds) >= 1000 {
				break
			}
			activeAds = append(activeAds, ad)
		}
	}
	c.IndentedJSON(http.StatusOK, activeAds)
}
