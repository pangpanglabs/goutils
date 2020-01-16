module github.com/pangpanglabs/goutils

go 1.12

require (
	github.com/Shopify/sarama v1.24.1
	github.com/cheekybits/is v0.0.0-20150225183255-68e9c0620927 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-xorm/xorm v0.7.9
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/klauspost/cpuid v1.2.2 // indirect
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/gommon v0.3.0
	github.com/matryer/try v0.0.0-20161228173917-9ac251b645a2 // indirect
	github.com/robfig/cron/v3 v3.0.1
	github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/viper v1.6.1
	github.com/stretchr/testify v1.4.0
	gopkg.in/matryer/try.v1 v1.0.0-20150601225556-312d2599e12e
	xorm.io/core v0.7.2
)

replace github.com/go-xorm/xorm => github.com/pangpanglabs/xorm v0.6.7-0.20191028024856-98149f1c9e95
