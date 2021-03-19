package joblog_test

import (
	"errors"
	"os"
	"testing"

	"github.com/pangpanglabs/goutils/joblog"
	"github.com/pangpanglabs/goutils/test"
)

var url = os.Getenv("JOB_LOG")

func TestLog(t *testing.T) {
	//1.test normal logic
	jobLog := joblog.New(url,
		"test", map[string]interface{}{"log": "this is test"})
	err := jobLog.Info("good")
	test.Ok(t, err)

	err = jobLog.Warning(struct{ Name string }{"xiaoxinmiao"})
	test.Ok(t, err)

	err = jobLog.Error(errors.New("this is bug."))
	test.Ok(t, err)

	err = jobLog.Finish()
	test.Ok(t, err)

	//2.test Disable:this content will not be logged
	jobLog = joblog.New(url,
		"test", map[string]interface{}{"log": "this is test 2"}, func(log *joblog.JobLog) {
			log.Disable = true
		})
	err = jobLog.Info("good 2")
	test.Ok(t, err)

	err = jobLog.Warning(struct{ Name string }{"xiaoxinmiao 2"})
	test.Ok(t, err)

	err = jobLog.Error(errors.New("this is bug. 2"))
	test.Ok(t, err)

	//3.test start param
	type Dto struct{ Message string }
	test.Ok(t, joblog.New(url, "test", Dto{Message: "how are you 1"}).Err)
	test.Ok(t, joblog.New(url, "test", &Dto{Message: "how are you 2"}).Err)
	test.Ok(t, joblog.New(url, "test", map[string]interface{}{"message": "how are you 3"}).Err)

}
