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
	// 初始化一个新的gin路由
	router := gin.New()

	// 定义一个使用GetAD处理程序的路由
	router.GET("/get-ad", getapi.GetAD)

	// 设置模拟的广告数据
	adstruct.ADs = []adstruct.AD{
		{Title: "Ad 1", StartAt: time.Now().Add(-time.Hour), EndAt: time.Now().Add(time.Hour), Conditions: adstruct.Conditions{}},
		{Title: "Ad 2", StartAt: time.Now().Add(-time.Hour), EndAt: time.Now().Add(time.Hour), Conditions: adstruct.Conditions{}},
		// 添加更多广告数据...
	}

	// 创建一个HTTP请求
	req, err := http.NewRequest("GET", "/get-ad", nil)
	assert.NoError(t, err)

	// 创建一个记录器以记录响应
	resp := httptest.NewRecorder()

	// 提供请求给路由器
	router.ServeHTTP(resp, req)

	// 断言响应状态码为200
	assert.Equal(t, http.StatusOK, resp.Code)

	// 解码JSON响应并检查活动广告数
	var ads []adstruct.AD
	err = json.NewDecoder(resp.Body).Decode(&ads)
	assert.NoError(t, err)

	// 断言活动广告的数量
	assert.NotEmpty(t, ads)
}
