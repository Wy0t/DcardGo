package postapi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Wy0t/DcardGo/adstruct"
	"github.com/Wy0t/DcardGo/database"
	"github.com/Wy0t/DcardGo/postapi"
	"github.com/gin-gonic/gin"
)

func init() {
	// 初始化資料庫連接
	database.Init()
}

func TestPostAD(t *testing.T) {
	// 建立測試路由
	router := gin.Default()

	// 添加測試路由處理程序
	router.POST("/postad", postapi.PostAD)

	// 準備測試廣告資料
	newAd := adstruct.AD{
		Title:   "Test Ad",
		StartAt: time.Now(),
		EndAt:   time.Now().Add(time.Hour * 24),
		Conditions: adstruct.Conditions{
			Age:      adstruct.PtrInt(25),
			Gender:   adstruct.PtrString("M"),
			Country:  &[]string{"TW", "JP"},
			Platform: &[]string{"ios"},
		},
	}

	// 將廣告資料轉換為 JSON 格式
	jsonData, err := json.Marshal(newAd)
	if err != nil {
		t.Fatalf("failed to marshal JSON: %v", err)
	}

	// 建立 POST 請求
	req, err := http.NewRequest("POST", "/postad", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 建立 HTTP 測試用的 ResponseRecorder
	rr := httptest.NewRecorder()

	// 發送請求到測試路由
	router.ServeHTTP(rr, req)

	// 檢查 HTTP 狀態碼是否正確
	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %d; got %d", http.StatusCreated, rr.Code)
	}

	// 檢查 response body 是否包含預期的廣告資料
	var responseBody adstruct.AD
	if err := json.Unmarshal(rr.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	// 檢查返回的廣告資料是否與預期一致
	if responseBody.Title != newAd.Title || !responseBody.StartAt.Equal(newAd.StartAt) || !responseBody.EndAt.Equal(newAd.EndAt) {
		t.Error("response body does not match expected ad")
	}
}
