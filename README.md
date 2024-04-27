##為何寫此專案
Dcard Intern提供之作業，也是我第一次寫Go語言，甚至是我的第一個作品集，雖然後續被Dcard刷掉沒上
但後來也學到了一些問題，在資料庫方面應該使用Go本身之框架，直接使用MySQL會有注入風險，在資安的方面較差
時間豐富可以寫Docker封裝程式，好讓面試官拿到範例程式時不管是何種系統都能使用，希望未來自己能在資訊方面更精進自己

## 簡介

這個 Go 程式提供了一個用於管理廣告的 API。它允許使用者創建具有特定定位條件的廣告，並根據各種參數（如年齡、性別、國家和平台）檢索活動廣告。

## 設計選擇

資料庫管理：

使用 MySQL 資料庫存儲廣告和條件數據。
設計了表格來存儲廣告、條件以及它們之間的關係。

Gin 框架：

選擇 Gin 作為 Web 框架，因其輕量級和高性能。
它簡化了路由、中間件和錯誤處理。

資料結構：

使用結構來表示廣告及其條件。
AD 結構封裝了廣告的詳細信息，如標題、開始和結束時間以及條件。
Conditions 結構表示定位條件，如年齡、性別、國家和平台。
錯誤處理：

在代碼中實現了錯誤處理，以處理無效的請求或資料庫錯誤。
返回自定義錯誤消息，並附上適當的 HTTP 狀態碼。

安全性：

通過使用參數化查詢來預防 SQL 注入。
執行輸入驗證以確保資料完整性，並防止惡意輸入。

## 使用方法

創建廣告：

使用者可以通過向 /ADs 端點發送帶有 JSON 數據的 POST 請求來創建廣告。
系統會在將廣告插入資料庫之前檢查每天創建的最大廣告數限制和定位條件的有效性。

檢索活動廣告：

可以通過向 /ADs 端點發送帶有可選查詢參數（如年齡、性別、國家和平台）的 GET 請求來檢索活動廣告。
系統根據提供的參數篩選和排序活動廣告，並以 JSON 格式返回結果。

資料庫管理：

程式初始化資料庫表格，並管理每天的廣告創建次數。
它使用 MySQL 查詢與資料庫進行交互，以存儲和檢索廣告數據。

## 總結

這個廣告管理系統提供了一個堅固的平台，用於創建和管理具有靈活定位選項的廣告。通過利用 Go 的並發功能和 Gin 的性能，確保高效處理廣告操作。可以進一步增強安全功能，改進錯誤處理，並根據業務需求實現額外功能。
