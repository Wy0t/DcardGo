package getapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/Wy0t/DcardGo/adstruct" // 导入所需的结构体
	"github.com/Wy0t/DcardGo/getapi"   // 导入要测试的包
)

func TestGetAD(t *testing.T) {
	// 初始化一個新的gin路由
	router := gin.New()

	// 定義一個使用GetAD處理程序的路由
	router.GET("/get-ad", getapi.GetAD)

	// 設定模擬的廣告數據
	adstruct.ADs = []adstruct.AD{
		{Title: "Ad 1", StartAt: time.Now().Add(-time.Hour), EndAt: time.Now().Add(time.Hour), Conditions: adstruct.Conditions{}},
		{Title: "Ad 2", StartAt: time.Now().Add(-time.Hour), EndAt: time.Now().Add(time.Hour), Conditions: adstruct.Conditions{}},
	}

	// 建立一個HTTP請求
	req, err := http.NewRequest("GET", "/get-ad", nil)
	assert.NoError(t, err)

	// 建立一個記錄器以記錄回應
	resp := httptest.NewRecorder()

	// 提供請求給路由器
	router.ServeHTTP(resp, req)

	// 斷言回應狀態碼為200
	assert.Equal(t, http.StatusOK, resp.Code)

	// 解碼JSON回應並檢查活動廣告數
	var ads []adstruct.AD
	err = json.NewDecoder(resp.Body).Decode(&ads)
	assert.NoError(t, err)

	// 斷言活動廣告的數量
	assert.NotEmpty(t, ads)
}
