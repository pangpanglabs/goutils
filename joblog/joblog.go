package joblog

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/pangpanglabs/goutils/httpreq"
)

type JobLog struct {
	url         string
	serviceName string
	jobId       int64
	Disable     bool
	Err         error
}

type jobStartDto struct {
	Servcie string      `json:"service"`
	Param   interface{} `json:"param"`
}

type messageLevel struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

func New(url, serviceName string, firstMessage interface{}, options ...func(*JobLog)) (jobLog *JobLog) {
	jobLog = &JobLog{
		url:         url,
		serviceName: serviceName,
	}
	for _, option := range options {
		if option == nil {
			continue
		}
		option(jobLog)
	}
	if jobLog.Disable == true {
		return
	}
	if len(jobLog.serviceName) == 0 {
		jobLog.Err = errors.New("serviceName is missing.")
		return
	}
	var result struct {
		Result int64 `json:"result"`
	}
	if ok, kind := validParam(firstMessage); !ok {
		jobLog.Err = fmt.Errorf("firstMessage type act:%v,exp:%v.",
			kind, "Struct, Map, Ptr")
		return
	}
	_, jobLog.Err = httpreq.New(http.MethodPost, jobLog.url, &jobStartDto{
		Servcie: jobLog.serviceName,
		Param:   firstMessage,
	}).Call(&result)
	if jobLog.Err != nil {
		return
	}
	jobLog.jobId = result.Result
	return
}

func (r *JobLog) Info(message interface{}) error {
	return r.write(message, "info")
}

func (r *JobLog) Warning(message interface{}) error {
	return r.write(message, "warning")
}
func (r *JobLog) Error(message interface{}) error {
	return r.write(message, "error")
}

func (r *JobLog) write(message interface{}, level string) (err error) {
	if r.Disable == true {
		return
	}
	if r.Err != nil {
		err = r.Err
		return
	}
	if r.jobId == int64(0) {
		return
	}
	url := fmt.Sprintf("%v/%v/logs", r.url, r.jobId)
	_, err = httpreq.New(http.MethodPost, url, &messageLevel{
		Message: toString(message),
		Level:   level,
	}).Call(nil)
	return
}

func (r *JobLog) Finish() (err error) {
	if r.Disable == true {
		return
	}
	if r.Err != nil {
		err = r.Err
		return
	}
	if r.jobId == int64(0) {
		return
	}
	url := fmt.Sprintf("%v/%v/finish", r.url, r.jobId)
	httpreq.New(http.MethodPost, url, nil).Call(nil)
	return
}

func toString(raw interface{}) string {
	switch v := raw.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, bool, float32, float64:
		return fmt.Sprint(v)
	case string:
		return string(v)
	case error:
		return v.Error()
	}
	val := reflect.ValueOf(raw)
	switch val.Kind() {
	case reflect.Struct, reflect.Map, reflect.Array, reflect.Slice, reflect.Ptr:
		b, _ := json.Marshal(raw)
		return string(b)
	}
	return ""
}
func validParam(param interface{}) (ok bool, kind reflect.Kind) {
	val := reflect.ValueOf(param)
	kind = val.Kind()
	switch kind {
	case reflect.Struct, reflect.Map, reflect.Ptr:
		ok = true
		return
	}
	return
}
