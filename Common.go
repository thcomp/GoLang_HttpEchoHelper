package HttpEchoHelper

import (
	"github.com/labstack/echo/v4"
)

var sEchoHelperMap map[string](*EchoHelper) = map[string](*EchoHelper){}

const sGlobalEchoHelperName = "### global ###"

func GetEchoHelper(params ...interface{}) (ret *EchoHelper) {
	if len(params) > 0 {
		if name, assertionOK := params[0].(string); assertionOK {
			if tempRet, exist := sEchoHelperMap[name]; exist {
				ret = tempRet
			} else {
				ret = &EchoHelper{}
				sEchoHelperMap[name] = ret
			}
		} else if echoIns, assertionOK := params[0].(*echo.Echo); assertionOK {
			for _, echoHelper := range sEchoHelperMap {
				if echoHelper.Echo() == echoIns {
					ret = echoHelper
					break
				}
			}
		}
	} else {
		if tempRet, exist := sEchoHelperMap[sGlobalEchoHelperName]; exist {
			ret = tempRet
		} else {
			ret = &EchoHelper{}
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
				if echoHelper.Echo() == echoIns {
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
