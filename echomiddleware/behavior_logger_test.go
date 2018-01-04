package echomiddleware

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRefactorControllerAndAction(t *testing.T) {
	datas := []struct {
		handlerName, controller, action string
	}{
		{"best/eland-show-service/controllers.(*AccountApiController).GetOpenid-fm", "AccountApiController", "GetOpenid"},
		{"best/eland-show-service/controllers.(RankingApiController).GetByType-fm", "RankingApiController", "GetByType"},
		{"main.(*ProductHandler).SearchSkus-fm", "ProductHandler", "SearchSkus"},
		{"main.(ProductHandler).SearchSkus-fm", "ProductHandler", "SearchSkus"},
		{"main.(*handler).(main.postSystmeContents)-fm", "handler", "postSystmeContents"},
		{"main.(*EventHandler).CreatePriceUpdatedEvent-fm", "EventHandler", "CreatePriceUpdatedEvent"},
		{"order-settings-service/api.CheckShopGrabForLogin", "api", "CheckShopGrabForLogin"},
		{"p2saas/smart-seller/kit.Ping", "kit", "Ping"},
		// {"main.main.func1.3.1", "main.main.func1.3", "1"},
	}
	for _, data := range datas {
		t.Run(data.handlerName, func(t *testing.T) {
			controller, action := echoRouter{}.convertHandlerNameToControllerAndAction(data.handlerName)
			assert.Equal(t, controller, data.controller)
			assert.Equal(t, action, data.action)
			fmt.Println(controller, action)
		})
	}
}
