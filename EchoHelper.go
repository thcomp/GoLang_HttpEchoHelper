package HttpEchoHelper

import (
	"context"

	"github.com/labstack/echo/v4"
	apihandler "github.com/thcomp/GoLang_APIHandler"
	awsSDKHelper "github.com/thcomp/GoLang_AwsSDKHelper"
)

type EchoHelper struct {
	echo       *echo.Echo
	apihandler *apihandler.APIHandler
}

func (helper *EchoHelper) Echo() *echo.Echo {
	return helper.echo
}

func (helper *EchoHelper) StartLambda() {
	awsSDKHelper.StartLambda(helper.lambdaHandler)
}

func (helper *EchoHelper) lambdaHandler(context context.Context, event interface{}) (ret interface{}, retErr error) {
	if 
}

func (helper *EchoHelper) APIHandler() *apihandler.APIHandler {
	if helper.apihandler == nil {
		helper.apihandler = apihandler.CreateLocalAPIManager()
	}

	return helper.apihandler
}

var sEchoHelper *EchoHelper

func GetEchoHelper() *EchoHelper {
	if sEchoHelper == nil {
		sEchoHelper = &EchoHelper{
			echo: echo.New(),
		}
	}

	return sEchoHelper
}
