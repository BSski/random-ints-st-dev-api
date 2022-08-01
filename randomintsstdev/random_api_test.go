package randomintsstdev

import "testing"

func TestPrepareURLParam(t *testing.T) {
	var tests = []struct {
		input    string
		expected int
	}{
		{"3", 3},
		{"", 1},
		{"20", 20},
	}
	for _, test := range tests {
		if output, _ := prepareURLParam(test.input, "testParam", 50); output != test.expected {
			t.Errorf("Test Failed: %v inputted, %v expected, received: %v", test.input, test.expected, output)
		}
	}
}
