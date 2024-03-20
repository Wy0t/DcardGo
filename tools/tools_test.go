package tools_test

import (
	"os"
	"testing"

	"github.com/Wy0t/DcardGo/adstruct"
	"github.com/Wy0t/DcardGo/database"
	"github.com/Wy0t/DcardGo/tools"
)

func TestContainsString(t *testing.T) {
	slice := []string{"apple", "banana", "orange"}
	str := "banana"
	if !tools.ContainsString(slice, str) {
		t.Errorf("Expected %s to be contained in the slice", str)
	}

	str = "grape"
	if tools.ContainsString(slice, str) {
		t.Errorf("Expected %s not to be contained in the slice", str)
	}
}

func TestValidateConditions(t *testing.T) {
	// 初始化数据库连接
	database.Init()
	conditions := adstruct.Conditions{
		Age:      tools.PtrInt(25),
		Gender:   tools.PtrString("M"),
		Country:  &[]string{"TW"}, // 设置为包含一个国家的切片
		Platform: &[]string{"android", "ios"},
	}

	err := tools.ValidateConditions(conditions)
	if err != nil {
		t.Errorf("Validation failed with error: %v", err)
	}

	// Test invalid age
	conditions.Age = tools.PtrInt(101)
	err = tools.ValidateConditions(conditions)
	expectedErrMsg := "invalid age value. age must be between 1 and 100"
	if err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Expected validation error: %s, but got: %v", expectedErrMsg, err)
	}

	// Test invalid gender
	conditions.Age = tools.PtrInt(25)
	conditions.Gender = tools.PtrString("X")
	expectedErrMsg = "invalid gender value. gender must be M or F"
	err = tools.ValidateConditions(conditions)
	if err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Expected validation error: %s, but got: %v", expectedErrMsg, err)
	}

	// Test invalid country
	conditions.Gender = tools.PtrString("M")
	conditions.Country = &[]string{"TW", "XXX"}
	expectedErrMsg = "invalid country value: XXX"
	err = tools.ValidateConditions(conditions)
	if err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Expected validation error: %s, but got: %v", expectedErrMsg, err)
	}

	// Test invalid platform
	conditions.Country = &[]string{"TW", "JP"}
	conditions.Platform = &[]string{"android", "invalid"}
	expectedErrMsg = "invalid platform value. platform must be android, ios, or web"
	err = tools.ValidateConditions(conditions)
	if err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Expected validation error: %s, but got: %v", expectedErrMsg, err)
	}
	// 关闭数据库连接
	database.CloseDatabase()
}

func TestQueryAdsFromDatabase(t *testing.T) {
	// Initialize database connection
	database.Init()

	// Perform query
	err := database.QueryAdsFromDatabase()
	if err != nil {
		t.Errorf("QueryAdsFromDatabase failed with error: %v", err)
	}

	// Close database connection
	database.CloseDatabase()
}

func TestMain(m *testing.M) {
	// Setup
	database.Init()

	// Run tests
	exitCode := m.Run()

	// Teardown
	database.CloseDatabase()

	// Exit with the same exit code as the tests
	os.Exit(exitCode)
}
