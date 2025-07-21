package HttpEchoHelper

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/thcomp/GoLang_HttpEntityHelper/entity"
	"github.com/thcomp/GoLang_HttpEntityHelper/jsonrpc"
	ThcompUtility "github.com/thcomp/GoLang_Utility"
)

type JSONRPCHandler struct {
	createdByNew bool
	method       string
	needAuth     bool
	handler      SubHandlerFunc
}

func NewJSONRPCHandler(jsonrpcMethod string, needAuth bool, handler SubHandlerFunc) *JSONRPCHandler {
	if handler == nil {
		ThcompUtility.LogfE("handler cannot be nil")
		return nil
	} else {
		return &JSONRPCHandler{
			createdByNew: true,
			method:       jsonrpcMethod,
			needAuth:     needAuth,
			handler:      handler,
		}
	}
}

func (handler *JSONRPCHandler) NeedAuth() bool {
	return handler.needAuth
}

func (handler *JSONRPCHandler) IsAcceptable(ctx echo.Context) bool {
	header := (http.Header)(nil)
	if ctx.Request() != nil {
		header = ctx.Request().Header
	} else if ctx.Response() != nil {
		header = ctx.Response().Header()
	}

	if header != nil {
		contentType := header.Get("Content-Type")
		if strings.Contains(contentType, "application/json") {
			return true
		}
	}

	return false
}

func (handler *JSONRPCHandler) Entity(ctx echo.Context) (entity.HttpEntity, error) {
	return jsonrpc.NewJSONRPCParser().Parse(ctx.Request()) // Parse the JSON-RPC request
}

func (handler *JSONRPCHandler) Handler(ctx echo.Context, entity entity.HttpEntity) error {
	return handler.handler(ctx, entity) // Replace with actual handling logic
}
