package HttpEchoHelper

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/thcomp/GoLang_HttpEntityHelper/entity"
	"github.com/thcomp/GoLang_HttpEntityHelper/jsonrpc"
)

var ErrUnknownJSONRPCMethod = fmt.Errorf("handler not known this jsonrpc method")

type Authorizer func(ctx echo.Context, entityIns entity.HttpEntity) *echo.HTTPError

type sMethodHandlerInfo struct {
	subHandler SubHandlerFunc
	authorizer Authorizer
}

type JSONRPCHandler struct {
	createdByNew bool
	methodMap    map[string](*sMethodHandlerInfo)
}

func NewJSONRPCHandler() *JSONRPCHandler {
	return &JSONRPCHandler{
		createdByNew: true,
		methodMap:    map[string](*sMethodHandlerInfo){},
	}
}

func (handler *JSONRPCHandler) IsAcceptable(ctx echo.Context) bool {
	header := (http.Header)(nil)
	if ctx.Request() != nil {
		header = ctx.Request().Header
	}

	if header != nil {
		contentType := header.Get("Content-Type")
		if strings.Contains(contentType, "application/json") {
			return true
		}
	}

	return false
}

func (handler *JSONRPCHandler) IsNeedEntityForAuthorize() bool {
	return true
}

func (handler *JSONRPCHandler) Authorize(ctx echo.Context) (retErr error) {
	helperCtx, _ := ctx.(*EchoHelperContext)

	if helperCtx.entityIns.EntityType() == entity.JSONRPC_Request {
		if jsonrpcReq, assertionOK := helperCtx.entityIns.(*jsonrpc.JSONRPCRequest); assertionOK {
			if subHandlerInfo, exist := handler.methodMap[jsonrpcReq.Method]; exist {
				retErr = subHandlerInfo.authorizer(ctx, jsonrpcReq)
			}
		}
	} else {
		retErr = fmt.Errorf("bad item: %v", helperCtx.entityIns)
	}

	return
}

func (handler *JSONRPCHandler) Entity(ctx echo.Context) (entity.HttpEntity, error) {
	return jsonrpc.NewJSONRPCParser().Parse(ctx.Request()) // Parse the JSON-RPC request
}

func (handler *JSONRPCHandler) Handler(ctx echo.Context, entityIns entity.HttpEntity) error {
	switch entityVaue := entityIns.(type) {
	case *jsonrpc.JSONRPCRequest:
		if subHandlerInfo, exist := handler.methodMap[entityVaue.Method]; exist {
			return subHandlerInfo.subHandler(ctx, entityIns)
		} else {
			return ErrUnknownJSONRPCMethod
		}
	default:
		return ErrNotAcceptable
	}
}

func (handler *JSONRPCHandler) RegisterMethodHandler(jsonrpcMethod string, subHandler SubHandlerFunc, authorizerIfNeed Authorizer) *JSONRPCHandler {
	handler.methodMap[jsonrpcMethod] = &sMethodHandlerInfo{
		subHandler: subHandler,
		authorizer: authorizerIfNeed,
	}

	return handler
}
