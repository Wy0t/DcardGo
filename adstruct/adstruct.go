package adstruct

import "time"

// ------------------------------定義廣告結構------------------------------
type (
	AD struct {
		Title      string     `json:"title"`
		StartAt    time.Time  `json:"startAt"`
		EndAt      time.Time  `json:"endAt"`
		Conditions Conditions `json:"conditions"`
	}

	// ------------------------------定義顯示條件結構------------------------------
	Conditions struct {
		Age      *int      `json:"age,omitempty"`
		Gender   *string   `json:"gender,omitempty"`
		Country  *[]string `json:"country,omitempty"`
		Platform *[]string `json:"platform,omitempty"`
	}
)

var ADs []AD

//建立一個指向整數的指針
func PtrInt(v int) *int {
	return &v
}

// 建立一個指向字串的指針
func PtrString(v string) *string {
	return &v
}
