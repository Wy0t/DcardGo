package tools

import (
	"fmt"

	"github.com/Wy0t/DcardGo/adstruct"
	"github.com/Wy0t/DcardGo/database"
)

// 檢查字串切片中是否包含某個特定的字串
func ContainsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// 驗證廣告的條件是否合法
func ValidateConditions(conditions adstruct.Conditions) error {
	// 檢查年齡是否在範圍內
	if conditions.Age != nil && (*conditions.Age < 1 || *conditions.Age > 100) {
		return fmt.Errorf("invalid age value. age must be between 1 and 100")
	}

	// 檢查性別是否合法
	if conditions.Gender != nil && (*conditions.Gender != "M" && *conditions.Gender != "F") {
		return fmt.Errorf("invalid gender value. gender must be M or F")
	}

	// 檢查國家是否合法
	if conditions.Country != nil {
		for _, country := range *conditions.Country {
			exists, err := database.IsValidCountry(country)
			if err != nil {
				return err
			}
			if !exists {
				return fmt.Errorf("invalid country value: %s", country)
			}
		}
	}

	// 檢查平台是否合法
	if conditions.Platform != nil {
		validPlatforms := map[string]bool{"android": true, "ios": true, "web": true}
		for _, platform := range *conditions.Platform {
			if !validPlatforms[platform] {
				return fmt.Errorf("invalid platform value. platform must be android, ios, or web")
			}
		}
	}

	return nil
}

// 创建一个指向整数的指针
func PtrInt(v int) *int {
	return &v
}

// 创建一个指向字符串的指针
func PtrString(v string) *string {
	return &v
}
