# goutils/test

It is a very simple test helper.

## Getting Started

```
import "pangpanglabs/goutils/test"

func TestXXX(t *testing.T) {
	at, err := time.Parse("2006-01-02", "2017-12-31")
	test.Ok(t, err)
	test.Equals(t, at.Year(), 2017)
	test.Assert(t, at.Month() == 12, "Month should be equals to 12")
}
```