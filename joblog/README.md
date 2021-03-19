# joblog

log print for eland job monitor

## Getting Started

```golang
jobLog := joblog.New(url,"test", map[string]interface{}{"log": "this is test"})

err := jobLog.Info("good")
test.Ok(t, err)

err = jobLog.Warning(struct{ Name string }{"xiaoxinmiao"})
test.Ok(t, err)

err = jobLog.Error(errors.New("this is bug."))
test.Ok(t, err)
```
### `joblog.New`

Create job log

- `url` - Save data to job monitor via this url,like:`https://xxx.com.cn/batchjob-api/v1/jobs`
- `serviceName` - The name of the microservice to store the log, usually the project name, for example: `ibill-qa`
- `firstMessage` - Print the contents of the first log,it can only pass `struct`, `pointer struct` or `maps`.
- `options` - Optional parameters can be modified,for example:`disable=false`

### `Info Warning Error `

Print info/warning/error log

- `message` - Print the contents of log,It can pass `any type` of parameter.

## View log

https://wiki.elandsystems.cn/display/DBA/Batch-job+Monitor

