package euchef

import (
	"testing"
	"time"
)

func TestFetch(t *testing.T) {

	format := "02-01-2006"
	now, _ := time.Parse(format, "29-08-2019")

	data, _ := FetchData(now)

	ParseData(data)
}
