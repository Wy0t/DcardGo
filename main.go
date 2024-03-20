package main

import (
	"github.com/Wy0t/DcardGo/database"
	"github.com/Wy0t/DcardGo/getapi"
	"github.com/Wy0t/DcardGo/postapi"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化資料庫連接
	database.Init()
	database.QueryAdsFromDatabase()

	//初始化Gin路由
	router := gin.Default()

	//訪問/AD路由時,回傳"getAD"JSON數據
	router.GET("/ADs", getapi.GetAD)
	router.POST("ADs", postapi.PostAD)
	// 關閉資料庫連接
	defer database.CloseDatabase()
	//Go應用程序入口點，啟動本地Web服務器
	router.Run("localhost:8080")
}
