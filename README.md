##項目說明

這個專案是一個廣告管理系統的後端服務，使用了Golang和Gin框架來實作RESTful API。 它提供了以下功能：

查詢活躍廣告：根據傳入的查詢參數（如年齡、性別、國家、平台等），傳回符合條件的活躍廣告列表，並依照結束時間升序排序。
建立廣告：透過POST請求，建立新的廣告資源，並對每日建立廣告數量進行了限制。
資料持久化：使用MySQL資料庫儲存廣告和條件訊息，並在程式啟動時從資料庫載入現有的廣告資料。
##使用說明

設定資料庫：在main函數中的資料庫連線部分，修改資料庫連線資訊為你的MySQL資料庫訊息，包括主機、資料庫名稱、使用者名稱和密碼。

初始化資料庫：在程式運行前，請確保資料庫中已經建立了對應的表，可以透過執行專案中的SQL語句來建立。 具體的建立表語句在init()函數中，如果資料庫中已存在對應的表，則會跳過建立。

啟動服務：執行main.go檔案啟動服務，預設監聽在localhost:8080。 可以透過訪問/ADs路由來獲取活躍廣告訊息，並且可以透過發送POST請求到/ADs路由來建立新的廣告。

API文件：API包含兩個端點：

GET /ADs: 查詢活躍廣告，支援傳入查詢參數進行條件篩選。
POST /ADs: 建立新的廣告資源，請求體需要包含新廣告的JSON資料。
##注意事項：

在建立新廣告時，系統會檢查每日建立廣告的數量是否超過預設的最大限制，如果超過則會傳回429狀態碼。
廣告的條件包括年齡、性別、國家和平台，創建廣告時會對條件進行合法性驗證，確保資料的完整性和正確性。
