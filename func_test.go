package HttpEchoHelper

import (
	"testing"
)

func Test_1(t *testing.T) {
	// Example test case for GetEchoHelper
	helper := GetEchoHelper("testHelper")
	if helper == nil {
		t.Error("Expected non-nil EchoHelper")
	}

	// Example test case for DeleteEchoHelper
	if !DeleteEchoHelper("testHelper") {
		t.Error("Expected to delete EchoHelper successfully")
	}

	// Check if the helper is deleted
	if _, exists := sEchoHelperMap["testHelper"]; exists {
		t.Error("Expected EchoHelper to be deleted")
	}
}
