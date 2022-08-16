package euchef

import (
	"testing"
)

func TestFetch(t *testing.T) {
	data, _ := FetchData()
	ParseData(data)
}
