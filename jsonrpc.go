package HttpEchoHelper

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/thcomp/GoLang_HttpEntityHelper/entity"
	"github.com/thcomp/GoLang_HttpEntityHelper/jsonrpc"
)

type JSONRPCHandler struct {
	needAuth bool
	handler  SubHandlerFunc
}

func NewJSONRPCHandler(needAuth bool, handler SubHandlerFunc) (*JSONRPCHandler, error) {
	if handler == nil {
		return nil, fmt.Errorf("handler cannot be nil")
	} else {
		return &JSONRPCHandler{
			needAuth: needAuth,
			handler:  handler,
		}, nil
	}
}

func (handler *JSONRPCHandler) NeedAuth() bool {
	return handler.needAuth
}

func (handler *JSONRPCHandler) IsAcceptable(ctx echo.Context) bool {
	// Implement your logic to check if the request is acceptable
	return true
}

func (handler *JSONRPCHandler) Entity(ctx echo.Context) (entity.HttpEntity, error) {
	return jsonrpc.NewJSONRPCParser().Parse(ctx.Request()) // Parse the JSON-RPC request
}

func (handler *JSONRPCHandler) Handler(ctx echo.Context, entity entity.HttpEntity) error {
	return handler.handler(ctx, entity) // Replace with actual handling logic
}
