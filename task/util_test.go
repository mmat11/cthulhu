package task

import "testing"

func TestInt64FromStr(t *testing.T) {
	tests := []struct {
		s   string
		res int64
	}{
		{
			s:   "-12345",
			res: -12345,
		},
		{
			s:   "0",
			res: 0,
		},
		{
			s:   "12345",
			res: 12345,
		},
	}

	for _, c := range tests {
		res := int64FromStr(c.s)
		if res != c.res {
			t.Fatalf("unexpected result, want %v, got %v", c.res, res)
		}
	}
}
