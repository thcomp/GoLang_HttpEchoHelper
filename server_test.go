package HttpEchoHelper

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/thcomp/GoLang_HttpEntityHelper/entity"
	"github.com/thcomp/GoLang_HttpEntityHelper/jsonrpc"
)

func Test_server(t *testing.T) {
	echoHelper := GetEchoHelper()
	echoHelper.PostWithSubHandler(
		"/api",
		NewJSONRPCHandler().RegisterMethodHandler(
			"login",
			func(ctx echo.Context, entity entity.HttpEntity) error {
				return nil
			},
			nil,
		),
	)

	go func() {
		client := http.Client{}
		jsonrpcReq, _ := jsonrpc.NewJSONRPCRequest(1, "login", nil)
		buffer := bytes.NewBuffer(nil)
		json.NewEncoder(buffer).Encode(jsonrpcReq)
		client.Post("http://localhost:8080/api", "application/json", buffer)
	}()
	echoHelper.StartServer(8080)
}
