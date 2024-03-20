package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Wy0t/DcardGo/adstruct"
	_ "github.com/go-sql-driver/mysql" // 导入 MySQL 驱动
)

// SetDB 設置全局數據庫連接
func SetDB(database *sql.DB) {
	DB = database
}

var (
	AdCreationCount map[string]int
	maxAdsPerDay    int
	DB              *sql.DB
)

func Init() {
	// 連接到 MySQL 資料庫
	var err error
	DB, err = sql.Open("mysql", "sql6692347:U3ygT5RDgw@tcp(sql6.freesqldatabase.com:3306)/sql6692347")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Connected to database")

	// 初始化 adCreationCount
	AdCreationCount = make(map[string]int)
	// 初始化 maxAdsPerDay
	err = DB.QueryRow("SELECT max_ads_per_day FROM configuration").Scan(&maxAdsPerDay)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Max ads per day:", maxAdsPerDay)

	// 檢查表是否存在，如果不存在，則創建表
	_, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS ad_creation_count (
            date DATE PRIMARY KEY,
            count INT
        );
    `)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Ad creation count table created")

	// 初始化每天的廣告創建數量
	today := time.Now().Format("2006-01-02")
	var count int
	err = DB.QueryRow("SELECT count FROM ad_creation_count WHERE date = ?", today).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			// 如果沒有記錄，將今天的計數初始化為 0
			AdCreationCount[today] = 0
		} else {
			panic(err.Error())
		}
	} else {
		AdCreationCount[today] = count
	}
	fmt.Println("Ad creation count for today:", AdCreationCount[today])

}

func QueryAdsFromDatabase() error {
	rows, err := DB.Query(`
        SELECT
            ADs.Title AS title,
            ADs.StartAt AS startAt,
            ADs.EndAt AS endAt,
            Conditions.Age AS age,
            Conditions.Gender AS gender,
            GROUP_CONCAT(country.Country) AS countries,
            GROUP_CONCAT(Conditions.Platform) AS platforms
        FROM
            ADs
        JOIN
            Conditions ON ADs.Conditions_ID = Conditions.ID
        LEFT JOIN Condition_Country ON Conditions.ID = Condition_Country.Condition_ID
        LEFT JOIN country ON Condition_Country.Country_ID = country.ID
        GROUP BY
            ADs.Title, ADs.StartAt, ADs.EndAt, Conditions.ID, Conditions.Age, Conditions.Gender
    `)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var AD adstruct.AD
		var startAtStr, endAtStr string
		var countriesStr sql.NullString
		var platformsStr sql.NullString
		if err := rows.Scan(&AD.Title, &startAtStr, &endAtStr, &AD.Conditions.Age, &AD.Conditions.Gender, &countriesStr, &platformsStr); err != nil {
			fmt.Println("掃描失敗:", err)
			return err
		}

		//  NULL 值
		var platforms []string
		if platformsStr.Valid {
			platforms = strings.Split(platformsStr.String, ",")
		}
		var countries []string
		if countriesStr.Valid {
			countries = strings.Split(countriesStr.String, ",")
		}

		// 將條件設為 nil 如果 platforms 為空
		var platformPtr *[]string
		if len(platforms) > 0 {
			platformPtr = &platforms
		}

		AD.Conditions.Platform = platformPtr

		// 將條件設為 nil 如果 countries 為空
		var countryPtr *[]string
		if len(countries) > 0 {
			countryPtr = &countries
		}

		AD.Conditions.Country = countryPtr

		// 將字串解析為 time.Time
		startAt, err := time.Parse("2006-01-02 15:04:05", startAtStr)
		if err != nil {
			fmt.Println("解析 startAt 時間失敗:", err)
			return err
		}

		endAt, err := time.Parse("2006-01-02 15:04:05", endAtStr)
		if err != nil {
			fmt.Println("解析 endAt 時間失敗:", err)
			return err
		}

		AD.StartAt = startAt
		AD.EndAt = endAt

		// 將 AD 結構添加到切片中
		adstruct.ADs = append(adstruct.ADs, AD)
	}

	return nil
}

// 檢查是否超過每天的最大廣告建立數量
func IsOverMaxAdsPerDay() bool {
	today := time.Now().Format("2006-01-02")
	count, exists := AdCreationCount[today]
	if !exists {
		// 如果没有当天的记录，返回 false
		return false
	}
	return count >= maxAdsPerDay
}

// 增加今天的廣告建立數量
func IncrementAdCreationCount() {
	today := time.Now().Format("2006-01-02")
	count := AdCreationCount[today] + 1
	AdCreationCount[today] = count

	// 更新資料庫中的計數
	_, err := DB.Exec("INSERT INTO ad_creation_count (date, count) VALUES (?, ?) ON DUPLICATE KEY UPDATE count = ?", today, count, count)
	if err != nil {
		panic(err.Error())
	}
}

// 插入條件信息到 Conditions 資料表
func InsertConditions(conditions adstruct.Conditions) (int64, error) {
	if conditions.Age == nil || conditions.Gender == nil || conditions.Platform == nil {
		return 0, fmt.Errorf("conditions.Age, conditions.Gender, or conditions.Platform is nil")
	}
	result, err := DB.Exec("INSERT INTO Conditions(age, gender, platform) VALUES (?, ?, ?)", *conditions.Age, *conditions.Gender, strings.Join(*conditions.Platform, ","))
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// 插入廣告信息到 ADs 資料表
func InsertAD(ad adstruct.AD, conditionsID int64) error {
	_, err := DB.Exec("INSERT INTO ADs(Title, StartAt, EndAt, Conditions_ID) VALUES (?, ?, ?, ?)", ad.Title, ad.StartAt, ad.EndAt, conditionsID)
	return err
}

// 檢查國家是否存在於資料庫中
func IsValidCountry(country string) (bool, error) {

	var exists bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM country WHERE Country = ?)", country).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// CloseDatabase 關閉資料庫連接
func CloseDatabase() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed")
	}
}
