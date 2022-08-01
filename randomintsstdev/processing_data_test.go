package randomintsstdev

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRoundedStDev(t *testing.T) {
	var tests = []struct {
		input    []int
		expected float64
	}{
		{[]int{1, 2, 3, 4, 5}, 1.414},
		{[]int{1, 3, 5, 7, 9}, 2.828},
		{[]int{3, 5, 10, 4, 1}, 3.007},
	}
	for _, test := range tests {
		if output, _ := getRoundedStDev(test.input); output != test.expected {
			t.Errorf("Test Failed: %v inputted, %v expected, received: %v", test.input, test.expected, output)
		}
	}
}

func TestGetIntSeqsSum(t *testing.T) {
	var tests = []struct {
		input    [][]int
		expected []int
	}{
		{[][]int{{1, 2}, {3, 4}}, []int{1, 2, 3, 4}},
		{[][]int{{1, 2}, {3, 4}, {5, 6}}, []int{1, 2, 3, 4, 5, 6}},
		{[][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}},
	}
	for _, test := range tests {
		output := getIntSeqsSum(test.input)
		if diff := cmp.Diff(output, test.expected); diff != "" {
			t.Errorf("Test Failed: %v inputted, %v expected, received: %v", test.input, test.expected, output)
		}
	}
}
