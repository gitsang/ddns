package time_test

import (
	"testing"
	"time"

	timex "github.com/gitsang/ddns/pkg/util/time"
)

func TestParseDuration(t *testing.T) {
	dur, err := timex.ParseDuration("30d")
	if err != nil {
		t.Errorf("ParseDuration(`30d`) failed: %v", err)
	}
	if dur != time.Hour*24*30 {
		t.Errorf("ParseDuration(`30d`) = %v, want %v", dur, time.Hour*24*30)
	}
}
