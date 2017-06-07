# goutils/config

## Getting Started

config.yml:
```yaml
database:
  driver: sqlite3
  connection: pos.db
  showSQL: true
debug: true
httpport: 8080
```

```golang
package main

import (
	"flag"

	"github.com/pangpanglabs/goutils/config"
)

main.go:
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