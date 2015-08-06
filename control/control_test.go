package control

import (
	"testing"
	// "time"
)

func TestGetValidBidWindow(t *testing.T) {

	c := &Controller{}
	from, till := c.GetValidBidWindow(3600, 1800)

	if from.After(till) {
		t.Errorf("Till cannot be before from: ", from)
	}
}

func TestCalculateBidPrice(t *testing.T) {

	c := &Controller{}

	// tests for normal prices
	tests := []struct {
		o   string
		s   string
		r   int
		res string
	}{
		{"0.35", "0.1", 35, "0.135"},
		{"1.2", "0.815", 12, "0.913"},
	}

	for n, cond := range tests {
		if result, _ := c.CalculateBidPrice(cond.o, cond.s, cond.r); result != cond.res {
			t.Errorf("Failed to calculate price in for testcase %v. Result is %v", n+1, result)
		}
	}

	// test for too high prices
	tests = []struct {
		o   string
		s   string
		r   int
		res string
	}{
		{"1.5", "1.0", 50, ""},
		{"0.12", "0.08", 60, ""},
	}

	for n, cond := range tests {
		if result, err := c.CalculateBidPrice(cond.o, cond.s, cond.r); err == nil {
			t.Errorf("Failed to get error on too high price in for testcase %v. Result is %v", n+1, result)
		}
	}

}
