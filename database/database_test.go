package database_test

import (
	"os"
	"testing"
	"time"

	"github.com/Wy0t/DcardGo/adstruct"
	"github.com/Wy0t/DcardGo/database"

	_ "github.com/go-sql-driver/mysql"
)

func TestInitDatabase(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Expected database connection to be initialized, got nil")
		}
	}()

	// 初始化資料庫
	database.Init()

	// 檢查是否成功初始化資料庫
	if database.DB == nil {
		t.Errorf("Expected database connection to be initialized, got nil")
	}

	// 確保測試結束時關閉資料庫連接
	defer database.CloseDatabase()
}

func TestQueryAdsFromDatabase(t *testing.T) {
	// 初始化資料庫
	database.Init()
	defer database.CloseDatabase()

	// 清空廣告切片以確保測試結果正確性
	adstruct.ADs = nil

	// 假設已有廣告數據存儲於資料庫中，進行廣告查詢
	err := database.QueryAdsFromDatabase()

	if err != nil {
		t.Errorf("Error querying ads from database: %v", err)
	}

	// 確認是否成功獲取廣告
	if len(adstruct.ADs) == 0 {
		t.Errorf("Expected to retrieve ads from database, got empty result")
	}
}

func TestIncrementAdCreationCount(t *testing.T) {
	// 初始化資料庫
	database.Init()
	defer database.CloseDatabase()

	// 增加廣告創建計數前，獲取當前計數
	today := time.Now().Format("2006-01-02")
	initialCount := database.AdCreationCount[today]

	// 假設增加一個廣告創建計數
	database.IncrementAdCreationCount()

	// 獲取增加後的計數
	updatedCount := database.AdCreationCount[today]

	// 檢查是否成功增加計數
	if updatedCount != initialCount+1 {
		t.Errorf("Expected ad creation count to be incremented by 1, got %d", updatedCount-initialCount)
	}

}

func TestMain(m *testing.M) {
	// 初始化資料庫
	database.Init()

	// 執行測試
	exitVal := m.Run()

	// 關閉資料庫連接
	database.CloseDatabase()

	// 退出測試
	os.Exit(exitVal)
}
