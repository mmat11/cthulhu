package bot

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMarshalUnmarshalCountersMap(t *testing.T) {
	var counters = CountersMap{
		"test": struct{}{},
	}

	encoded := MarshalCountersMap(counters)
	res, err := UnmarshalCountersMap(encoded)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(counters, res) {
		t.Fatalf("error during marshal/unmarshal: want %v, got %v", counters, res)
	}
}
