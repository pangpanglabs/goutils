# goutils/config

## Getting Started

main.go:
```golang
package main

import (
	"flag"

	"github.com/pangpanglabs/goutils/config"
)

func main() {
	appEnv := flag.String("app-env", os.Getenv("APP_ENV"), "app env")
	flag.Parse()

	var c struct {
		Database struct{ Driver, Connection string }
		Debug    bool
		Httpport string
	}
	if err := config.Read(*appEnv, &c); err != nil {
		panic(err)
	}


        /* ... */
}
```

config.yml:
```yaml
database:
  driver: sqlite3
  connection: pos.db
  showSQL: true
debug: true
httpport: 8080
```

config.test.yml:
```yaml
database:
  connection: test.db
```

config.staging.yml:
```yaml
database:
  driver: mysql
  connection: username:password@tcp(staging.server.com:3307)/db_name?charset=utf8&parseTime=True&loc=UTC
```

config.production.yml:
```yaml
database:
  connection: username:password@tcp(production.server.com:3307)/db_name?charset=utf8&parseTime=True&loc=UTC
  showSQL: false
debug: false
```