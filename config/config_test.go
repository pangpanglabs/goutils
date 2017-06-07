package config_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/pangpanglabs/goutils/config"
	"github.com/pangpanglabs/goutils/test"
)

var (
	baseConfig = `
database:
  driver: sqlite3
  connection: pos.db
debug: true
httpport: 8080`

	testConfig = `
database:
  connection: test.db`
)

func TestConfig(t *testing.T) {
	var c struct {
		Database struct{ Driver, Connection string }
		Debug    bool
		Httpport string
	}

	err := ioutil.WriteFile("./config.yml", []byte(baseConfig), 0666)
	test.Ok(t, err)
	defer os.Remove("./config.yml")

	err = ioutil.WriteFile("./config.test.yml", []byte(testConfig), 0666)
	test.Ok(t, err)
	defer os.Remove("./config.test.yml")

	err = config.Read("test", &c)
	test.Ok(t, err)

	test.Equals(t, c.Database.Driver, "sqlite3")
	test.Equals(t, c.Database.Connection, "test.db")
	test.Equals(t, c.Debug, true)
	test.Equals(t, c.Httpport, "8080")
}
