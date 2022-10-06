package main

import "testing"

func TestParseNumber(t *testing.T) {
	type testData struct {
		expression     string
		expectedLength int
		expectedResult float64
	}
	tests := []testData{
		{"123.123", 7, 123.123},
		{"15.15", 5, 15.15},
		{"0.0", 3, 0.0},
	}
	for _, v := range tests {
		result, length, err := ParseNumber(v.expression)
		if err != nil {
			t.Error(err.Error())
		}
		if result != v.expectedResult || length != v.expectedLength {
			t.Error("Test failed for", v.expression, "calculated", result, "with length", length, "expected", v.expectedResult, v.expectedLength)
		}
	}
}

func TestExecuteOperator(t *testing.T) {
	type testData struct {
		leftNumber     float64
		rightNumber    float64
		operator       rune
		expectedResult float64
	}
	tests := []testData{
		{35.35, -10.0, '+', 25.35},
		{35.35, -10.0, '-', 45.35},
		{35.35, -10.0, '*', -353.5},
		{1.0, 3.0, '/', 1.0 / 3.0},
	}
	for _, v := range tests {
		result, err := ExecuteOperator(v.leftNumber, v.rightNumber, v.operator)
		if err != nil {
			t.Error(err.Error())
		}
		if result != v.expectedResult {
			t.Error("Test failed for", v.leftNumber, v.rightNumber, v.operator, "calculated", result, "expected", v.expectedResult)
		}
	}
}

func TestParse(t *testing.T) {
	type testData struct {
		expression     string
		expectedResult float64
	}
	tests := []testData{
		{"(1+2)-3", 0},
		{"(1+2)*3", 9},
		{"1/3", 1.0 / 3.0},
		{"1 + 2 * (3 + 4 / 2 - (1 + 2)) * 2 + 1", 10},
		{"5+(1*(2 + 3)+ 7) *3", 41},
	}
	for _, v := range tests {
		result, err := parse(v.expression)
		if err != nil {
			t.Error(err.Error())
		}
		if result != v.expectedResult {
			t.Error("Test failed for", v.expression, "calculated", result, "expected", v.expectedResult)
		}
	}
}
