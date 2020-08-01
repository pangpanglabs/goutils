module github.com/hillfolk/goutils

go 1.12

require (
	github.com/Shopify/sarama v1.24.1
	github.com/aws/aws-sdk-go v1.33.17
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-xorm/xorm v0.7.9
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/echo/v4 v4.1.16
	github.com/pangpanglabs/goutils v0.0.0-20200320140103-932a39405894
	github.com/robfig/cron/v3 v3.0.1
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/viper v1.6.1
	github.com/stretchr/testify v1.5.1
	xorm.io/core v0.7.2
)

replace github.com/go-xorm/xorm => github.com/pangpanglabs/xorm v0.6.7-0.20191028024856-98149f1c9e95
