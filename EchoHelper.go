package HttpEchoHelper

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/labstack/echo/v4"
	apihandler "github.com/thcomp/GoLang_APIHandler"
	awsSDKHelper "github.com/thcomp/GoLang_AwsSDKHelper"
	ThcompUtility "github.com/thcomp/GoLang_Utility"
)

type EchoHelperFunc func(helper *EchoHelper) error

type LambdaTrigger int

const (
	None LambdaTrigger = iota
	APIGateway
	APIGatewayV2
	LambdaFunctionURL
)

type EchoHelper struct {
	echo       *echo.Echo
	apiManager *apihandler.APIManager
}

func (helper *EchoHelper) Echo() *echo.Echo {
	return helper.echo
}

func (helper *EchoHelper) APIManager() *apihandler.APIManager {
	if helper.apiManager == nil {
		helper.apiManager = apihandler.CreateLocalAPIManager()
	}

	return helper.apiManager
}

func (helper *EchoHelper) ServeByAPIManager(w http.ResponseWriter, r *http.Request) {
	if helper.apiManager != nil {
		helper.apiManager.ExecuteRequest(r, w)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
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
		httpResponse := ThcompUtility.NewHttpResponseHelper()
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
			helper.echo.ServeHTTP(httpResponse, httpRequest)
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
		httpResponse := ThcompUtility.NewHttpResponseHelper()
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
			helper.echo.ServeHTTP(httpResponse, httpRequest)
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
		httpResponse := ThcompUtility.NewHttpResponseHelper()
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
			helper.echo.ServeHTTP(httpResponse, httpRequest)
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

var sEchoHelperMap map[string](*EchoHelper) = map[string](*EchoHelper){}

const sGlobalEchoHelperName = "### global ###"

func GetEchoHelper(params ...interface{}) (ret *EchoHelper) {
	if len(params) > 0 {
		if name, assertionOK := params[0].(string); assertionOK {
			if tempRet, exist := sEchoHelperMap[name]; exist {
				ret = tempRet
			}
		} else if echoIns, assertionOK := params[0].(*echo.Echo); assertionOK {
			for _, echoHelper := range sEchoHelperMap {
				if echoHelper.echo == echoIns {
					ret = echoHelper
					break
				}
			}
		}
	} else {
		if tempRet, exist := sEchoHelperMap[sGlobalEchoHelperName]; exist {
			ret = tempRet
		} else {
			ret = &EchoHelper{
				echo: echo.New(),
			}
			sEchoHelperMap[sGlobalEchoHelperName] = ret
		}
	}

	return ret
}

func DeleteEchoHelper(params ...interface{}) (ret bool) {
	if len(params) > 0 {
		if name, assertionOK := params[0].(string); assertionOK {
			if _, exist := sEchoHelperMap[name]; exist {
				delete(sEchoHelperMap, name)
				ret = true
			}
		} else if echoIns, assertionOK := params[0].(*echo.Echo); assertionOK {
			for key, echoHelper := range sEchoHelperMap {
				if echoHelper.echo == echoIns {
					delete(sEchoHelperMap, key)
					ret = true
					break
				}
			}
		}
	} else {
		if _, exist := sEchoHelperMap[sGlobalEchoHelperName]; exist {
			delete(sEchoHelperMap, sGlobalEchoHelperName)
			ret = true
		}
	}

	return ret
}
