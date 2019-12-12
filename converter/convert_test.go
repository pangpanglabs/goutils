package converter_test

import (
	"testing"

	"github.com/pangpanglabs/goutils/converter"
	"github.com/pangpanglabs/goutils/test"
)

func TestConvert(t *testing.T) {
	t.Run("StringToIntSlice", func(t *testing.T) {
		s := "1,2,3,4,5,6"
		list := converter.StringToIntSlice(s)
		test.Equals(t, list, []int64{1, 2, 3, 4, 5, 6})
	})
	t.Run("StringToStringSlice", func(t *testing.T) {
		s := "1,2,3,4,5,6"
		list := converter.StringToStringSlice(s)
		test.Equals(t, list, []string{"1", "2", "3", "4", "5", "6"})
	})
}
