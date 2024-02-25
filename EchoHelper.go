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
	echo        *echo.Echo
	echoHandler echo.HandlerFunc
	trigger     LambdaTrigger
	apihandler  *apihandler.APIHandler
}

func (helper *EchoHelper) Echo() *echo.Echo {
	return helper.echo
}

func (helper *EchoHelper) StartLambda(handler echo.HandlerFunc, trigger LambdaTrigger) {
	helper.echoHandler = handler
	switch trigger {
	case APIGateway:
		awsSDKHelper.StartLambda(helper.lambdaWithApigwHandler)
	case APIGatewayV2:
		awsSDKHelper.StartLambda(helper.lambdaWithApigwV2Handler)
	case LambdaFunctionURL:
		awsSDKHelper.StartLambda(helper.lambdaWithFunctionURLHandler)
	default:
		helper.echoHandler = nil
	}
}

func (helper *EchoHelper) lambdaWithApigwHandler(context context.Context, request *events.APIGatewayProxyRequest) (ret *events.APIGatewayProxyResponse, retErr error) {
	if httpRequest, convErr := awsSDKHelper.FromAPIGatewayProxyRequest2HttpRequest(request); convErr == nil {
		httpResponse := ThcompUtility.NewHttpResponseHelper()
		helper.echo.ServeHTTP(&httpResponse, httpRequest)
		if ret, convErr = awsSDKHelper.FromHttpResponse2APIGatewayProxyResponse(httpResponse); convErr != nil {
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
		helper.echo.ServeHTTP(&httpResponse, httpRequest)
		if ret, convErr = awsSDKHelper.FromHttpResponse2APIGatewayV2HTTPResponse(httpResponse); convErr != nil {
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
		helper.echo.ServeHTTP(&httpResponse, httpRequest)
		if ret, convErr = awsSDKHelper.FromHttpResponse2LambdaFunctionURLResponse(httpResponse); convErr != nil {
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

func (helper *EchoHelper) APIHandler() *apihandler.APIHandler {
	if helper.apihandler == nil {
		helper.apihandler = apihandler.CreateLocalAPIManager()
	}

	return helper.apihandler
}

var sEchoHelperMap map[*echo.Echo](*EchoHelper) = map[*echo.Echo](*EchoHelper){}

func GetEchoHelper(echoInfs ...interface{}) (ret *EchoHelper) {
	if echoInfs != nil && len(echoInfs) > 0 {
		if echoIns, assertionOK := echoInfs[0].(*echo.Echo); assertionOK {
			if tempRet, exist := sEchoHelperMap[echoIns]; exist {
				ret = tempRet
			}
		}
	} else {
		ret = &EchoHelper{
			echo: echo.New(),
		}
		sEchoHelperMap[ret.echo] = ret
	}

	return ret
}
