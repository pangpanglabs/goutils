module github.com/pangpanglabs/goutils

go 1.12

require (
	github.com/Shopify/sarama v1.24.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-xorm/xorm v0.7.9
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/klauspost/cpuid v1.2.2 // indirect
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/gommon v0.3.0
	github.com/robfig/cron v1.2.0
	github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/viper v1.6.1
	github.com/stretchr/testify v1.4.0
	xorm.io/core v0.7.2
)

replace github.com/go-xorm/xorm => github.com/pangpanglabs/xorm v0.6.7-0.20191028024856-98149f1c9e95
