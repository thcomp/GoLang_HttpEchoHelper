package HttpEchoHelper

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/labstack/echo/v4"
	apihandler "github.com/thcomp/GoLang_APIHandler"
	awsSDKHelper "github.com/thcomp/GoLang_AwsSDKHelper"
	HttpEntityHelper "github.com/thcomp/GoLang_HttpEntityHelper"
	ThcompUtility "github.com/thcomp/GoLang_Utility"
)

type EchoHelperFunc func(helper *EchoHelper) error
type SubHandlerFunc func(ctx echo.Context, entity HttpEntityHelper.HttpEntity) error
type LambdaTrigger int

const (
	None LambdaTrigger = iota
	APIGateway
	APIGatewayV2
	LambdaFunctionURL
)

type SubHandlerInterface interface {
	NeedAuth() bool
	IsAcceptable(ctx echo.Context) bool
	Entity(ctx echo.Context) (HttpEntityHelper.HttpEntity, error)
	Handler(ctx echo.Context, entity HttpEntityHelper.HttpEntity) error
}

type EchoHelper struct {
	echo          *echo.Echo
	apiManager    *apihandler.APIManager
	subHandlerMap map[string] /*http method*/ ([]SubHandlerInterface)
}

func (helper *EchoHelper) Echo() *echo.Echo {
	if helper.echo == nil {
		helper.echo = echo.New()
	}

	return helper.echo
}

// @deprecated
func (helper *EchoHelper) APIManager() *apihandler.APIManager {
	if helper.apiManager == nil {
		helper.apiManager = apihandler.CreateLocalAPIManager()
	}

	return helper.apiManager
}

// @deprecated
func (helper *EchoHelper) ServeByAPIManager(w http.ResponseWriter, r *http.Request) {
	if helper.apiManager != nil {
		helper.apiManager.ExecuteRequest(r, w)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (helper *EchoHelper) Any(path string, handler echo.HandlerFunc) *EchoHelper {
	helper.Echo().Any(path, handler)
	return helper
}

func (helper *EchoHelper) Get(path string, handler echo.HandlerFunc) *EchoHelper {
	helper.Echo().GET(path, handler)
	return helper
}

func (helper *EchoHelper) Post(path string, handler echo.HandlerFunc) *EchoHelper {
	helper.Echo().POST(path, handler)
	return helper
}

func (helper *EchoHelper) Delete(path string, handler echo.HandlerFunc) *EchoHelper {
	helper.Echo().DELETE(path, handler)
	return helper
}

func (helper *EchoHelper) Put(path string, handler echo.HandlerFunc) *EchoHelper {
	helper.Echo().PUT(path, handler)
	return helper
}

func (helper *EchoHelper) Options(path string, handler echo.HandlerFunc) *EchoHelper {
	helper.Echo().OPTIONS(path, handler)
	return helper
}

func (helper *EchoHelper) Head(path string, handler echo.HandlerFunc) *EchoHelper {
	helper.Echo().HEAD(path, handler)
	return helper
}

func (helper *EchoHelper) AnyWithSubHandler(path string, subHandler SubHandlerInterface) *EchoHelper {
	return helper.withHandler("any", path, subHandler)
}

func (helper *EchoHelper) GetWithSubHandler(path string, subHandler SubHandlerInterface) *EchoHelper {
	return helper.withHandler("get", path, subHandler)
}

func (helper *EchoHelper) PostWithSubHandler(path string, subHandler SubHandlerInterface) *EchoHelper {
	return helper.withHandler("post", path, subHandler)
}

func (helper *EchoHelper) DeleteWithSubHandler(path string, subHandler SubHandlerInterface) *EchoHelper {
	return helper.withHandler("delete", path, subHandler)
}

func (helper *EchoHelper) PutWithSubHandler(path string, subHandler SubHandlerInterface) *EchoHelper {
	return helper.withHandler("put", path, subHandler)
}

func (helper *EchoHelper) OptionsWithSubHandler(path string, subHandler SubHandlerInterface) *EchoHelper {
	return helper.withHandler("options", path, subHandler)
}

func (helper *EchoHelper) HeadWithSubHandler(path string, subHandler SubHandlerInterface) *EchoHelper {
	return helper.withHandler("head", path, subHandler)
}

func (helper *EchoHelper) withHandler(httpMethod, path string, subHandler SubHandlerInterface) *EchoHelper {
	if helper.subHandlerMap == nil {
		helper.subHandlerMap = make(map[string][]SubHandlerInterface)
	}
	if _, exist := helper.subHandlerMap[httpMethod]; !exist {
		helper.subHandlerMap[httpMethod] = make([]SubHandlerInterface, 0)
	}
	helper.subHandlerMap[httpMethod] = append(helper.subHandlerMap[httpMethod], subHandler)

	if httpMethod == "any" {
		helper.Echo().Any(path, helper.firstHandlerForSub)
	} else {
		helper.Echo().Add(httpMethod, path, helper.firstHandlerForSub)
	}

	return helper
}

func (helper *EchoHelper) firstHandlerForSub(c echo.Context) (retErr error) {
	// 先にメソッドが一致するものを優先する
	httpMethod := strings.ToLower(c.Request().Method)
	if handlers, exist := helper.subHandlerMap[httpMethod]; exist {
		for _, handler := range handlers {
			if handler.IsAcceptable(c) {
				if entity, err := handler.Entity(c); err != nil {
					retErr = err
				} else {
					retErr = handler.Handler(c, entity)
				}

				return
			}
		}
	}

	// メソッドが一致しない場合は、全てのハンドラを
	httpMethod = "any"
	if handlers, exist := helper.subHandlerMap[httpMethod]; exist {
		for _, handler := range handlers {
			if handler.IsAcceptable(c) {
				if entity, err := handler.Entity(c); err != nil {
					retErr = err
				} else {
					retErr = handler.Handler(c, entity)
				}

				return
			}
		}
	}

	return c.NoContent(http.StatusNotFound)
}

func (helper *EchoHelper) StartServer(port int) {
	helper.Echo().Logger.Fatal(helper.echo.Start(":" + strconv.FormatInt(int64(port), 10)))
}

func (helper *EchoHelper) StartLambda(trigger LambdaTrigger) bool {
	ret := true
	switch trigger {
	case APIGateway:
		awsSDKHelper.StartLambda(helper.lambdaWithApigwHandler)
	case APIGatewayV2:
		awsSDKHelper.StartLambda(helper.lambdaWithApigwV2Handler)
	case LambdaFunctionURL:
		awsSDKHelper.StartLambda(helper.lambdaWithFunctionURLHandler)
	default:
		ret = false
	}
	return ret
}

func (helper *EchoHelper) lambdaWithApigwHandler(context context.Context, request *events.APIGatewayProxyRequest) (ret *events.APIGatewayProxyResponse, retErr error) {
	if httpRequest, convErr := awsSDKHelper.FromAPIGatewayProxyRequest2HttpRequest(request); convErr == nil {
		httpResponse := ThcompUtility.NewHttpResponseHelper(httpRequest)
		useEcho := false

		if helper.apiManager != nil {
			helper.apiManager.ExecuteRequest(httpRequest, httpResponse)
			if httpResponse.ExportHttpResponse().StatusCode == http.StatusNotFound {
				// JSONなどのAPI Manager側で未登録により処理を行わなかった場合に、初期値に戻す
				httpResponse.WriteHeader(http.StatusOK)
				useEcho = true
			}
		} else {
			useEcho = true
		}

		if useEcho {
			helper.Echo().ServeHTTP(httpResponse, httpRequest)
		}

		if ret, convErr = awsSDKHelper.FromHttpResponse2APIGatewayProxyResponse(httpResponse.ExportHttpResponse()); convErr != nil {
			retErr = convErr
			ret = &events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
			}
		}
	} else {
		retErr = convErr
		ret = &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
		}
	}

	return
}

func (helper *EchoHelper) lambdaWithApigwV2Handler(context context.Context, request *events.APIGatewayV2HTTPRequest) (ret *events.APIGatewayV2HTTPResponse, retErr error) {
	if httpRequest, convErr := awsSDKHelper.FromAPIGatewayV2HTTPRequest2HttpRequest(request); convErr == nil {
		httpResponse := ThcompUtility.NewHttpResponseHelper(httpRequest)
		useEcho := false

		if helper.apiManager != nil {
			helper.apiManager.ExecuteRequest(httpRequest, httpResponse)
			if httpResponse.ExportHttpResponse().StatusCode == http.StatusNotFound {
				// JSONなどのAPI Manager側で未登録により処理を行わなかった場合に、初期値に戻す
				httpResponse.WriteHeader(http.StatusOK)
				useEcho = true
			}
		} else {
			useEcho = true
		}

		if useEcho {
			helper.Echo().ServeHTTP(httpResponse, httpRequest)
		}

		if ret, convErr = awsSDKHelper.FromHttpResponse2APIGatewayV2HTTPResponse(httpResponse.ExportHttpResponse()); convErr != nil {
			retErr = convErr
			ret = &events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusInternalServerError,
			}
		}
	} else {
		retErr = convErr
		ret = &events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
		}
	}

	return
}

func (helper *EchoHelper) lambdaWithFunctionURLHandler(context context.Context, request *events.LambdaFunctionURLRequest) (ret *events.LambdaFunctionURLResponse, retErr error) {
	if httpRequest, convErr := awsSDKHelper.FromLambdaFunctionURLRequest2HttpRequest(request); convErr == nil {
		httpResponse := ThcompUtility.NewHttpResponseHelper(httpRequest)
		useEcho := false

		if helper.apiManager != nil {
			helper.apiManager.ExecuteRequest(httpRequest, httpResponse)
			if httpResponse.ExportHttpResponse().StatusCode == http.StatusNotFound {
				// JSONなどのAPI Manager側で未登録により処理を行わなかった場合に、初期値に戻す
				httpResponse.WriteHeader(http.StatusOK)
				useEcho = true
			}
		} else {
			useEcho = true
		}

		if useEcho {
			helper.Echo().ServeHTTP(httpResponse, httpRequest)
		}

		if ret, convErr = awsSDKHelper.FromHttpResponse2LambdaFunctionURLResponse(httpResponse.ExportHttpResponse()); convErr != nil {
			retErr = convErr
			ret = &events.LambdaFunctionURLResponse{
				StatusCode: http.StatusInternalServerError,
			}
		}
	} else {
		retErr = convErr
		ret = &events.LambdaFunctionURLResponse{
			StatusCode: http.StatusBadRequest,
		}
	}

	return
}
