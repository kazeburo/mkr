package format

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestISO8601Extended(t *testing.T) {
	now := time.Now()
	expect := now.Format("2006-01-02T15:04:05-07:00") // ISO 8601 extended format
	got := ISO8601Extended(now)
	if got != expect {
		t.Errorf("should be %q got %q", expect, got)
	}
}

func TestReplaceNullIfSlice(t *testing.T) {
	assert.Equal(t, replaceNullIfSlice("null", reflect.TypeOf("0").Kind()), "null")
	assert.Equal(t, replaceNullIfSlice("[]", reflect.TypeOf([]int{}).Kind()), "[]")
	assert.Equal(t, replaceNullIfSlice("null", reflect.TypeOf([]int{}).Kind()), "[]")
	assert.Equal(t, replaceNullIfSlice("null", reflect.TypeOf([0]int{}).Kind()), "[]")

}
