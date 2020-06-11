package euchef

import (
	"testing"
	"time"
)

func TestFetch(t *testing.T) {

	format := "02-01-2006"
	now, _ := time.Parse(format, "11-06-2020")

	data, _ := FetchData(now)

	ParseData(data)
}
