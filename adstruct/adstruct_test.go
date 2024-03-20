package adstruct

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestADJSONSerialization(t *testing.T) {
	startAt := time.Date(2024, time.March, 20, 0, 0, 0, 0, time.UTC)
	endAt := time.Date(2024, time.March, 21, 0, 0, 0, 0, time.UTC)
	ad := AD{
		Title:   "Test Ad",
		StartAt: startAt,
		EndAt:   endAt,
		Conditions: Conditions{
			Age:      IntPointer(18),
			Gender:   StringPointer("Male"),
			Country:  &[]string{"USA", "Canada"},
			Platform: &[]string{"iOS", "Android"},
		},
	}

	// 序列化AD結構到JSON
	jsonData, err := json.Marshal(ad)
	if err != nil {
		t.Errorf("Error marshalling AD struct to JSON: %v", err)
	}

	// 反序列化JSON到AD結構
	var adFromJSON AD
	err = json.Unmarshal(jsonData, &adFromJSON)
	if err != nil {
		t.Errorf("Error unmarshalling JSON to AD struct: %v", err)
	}

	// 確保原始結構和從JSON反序列化的結構相等
	if !reflect.DeepEqual(ad, adFromJSON) {
		t.Errorf("AD struct after unmarshalling JSON does not match original AD struct")
	}
}

// 輔助函數：建立整數
func IntPointer(i int) *int {
	return &i
}

// 輔助函數：建立字串
func StringPointer(s string) *string {
	return &s
}
