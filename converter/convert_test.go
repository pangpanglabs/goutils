package converter_test

import (
	"testing"

	"github.com/pangpanglabs/goutils/converter"
	"github.com/pangpanglabs/goutils/test"
)

func TestConvert(t *testing.T) {
	t.Run("StringToIntSlice", func(t *testing.T) {
		s := "1,2,3,4, 5,6 "
		list := converter.StringToIntSlice(s)
		test.Equals(t, list, []int64{1, 2, 3, 4, 5, 6})
	})
	t.Run("StringToStringSlice", func(t *testing.T) {
		s := "1,2,3,4, 5,6 "
		list := converter.StringToStringSlice(s)
		test.Equals(t, list, []string{"1", "2", "3", "4", "5", "6"})
	})
	t.Run("StringSliceToInt64Slice", func(t *testing.T) {
		s := []string{"1", "2", "3", "4", " 5", "6"}
		list := converter.StringSliceToInt64Slice(s)
		test.Equals(t, list, []int64{1, 2, 3, 4, 5, 6})
	})
	t.Run("Int64SliceToStringSlice", func(t *testing.T) {
		s := []int64{1, 2, 3, 4, 5, 6}
		list := converter.Int64SliceToStringSlice(s)
		test.Equals(t, list, []string{"1", "2", "3", "4", "5", "6"})
	})
	t.Run("Int64SliceToString", func(t *testing.T) {
		s := []int64{1, 2, 3, 4, 5, 6}
		list := converter.Int64SliceToString(s)
		test.Equals(t, list, "1,2,3,4,5,6")
	})
	t.Run("ContainsInt64", func(t *testing.T) {
		s := []int64{1, 2, 3, 4, 5, 6}
		b := converter.ContainsInt64(s, 5)
		test.Equals(t, true, b)
	})
	t.Run("UniqueInt64", func(t *testing.T) {
		s := []int64{-1, 6, 3, 4, 3, 6}
		list := converter.UniqueInt64(s, true)
		test.Equals(t, list, []int64{6, 3, 4})
	})
	t.Run("UniqueString", func(t *testing.T) {
		s := []string{"-1", "6", "3", "4", "3", "6"}
		list := converter.UniqueString(s)
		test.Equals(t, list, []string{"-1", "6", "3", "4"})
	})
}
