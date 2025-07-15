package HttpEchoHelper

import (
	"testing"

	"github.com/labstack/echo/v4"
	HttpEntityHelper "github.com/thcomp/GoLang_HttpEntityHelper"
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

func Test_2(t *testing.T) {
	// Example test case for NewJSONRPCHandler
	handler, err := NewJSONRPCHandler(true, func(ctx echo.Context, entity HttpEntityHelper.HttpEntity) error {
		return nil // Example handler logic
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if handler == nil || !handler.NeedAuth() {
		t.Error("Expected JSONRPCHandler with needAuth set to true")
	}

	// Example test case for IsAcceptable
	if !handler.IsAcceptable(nil) {
		t.Error("Expected IsAcceptable to return true")
	}
}
