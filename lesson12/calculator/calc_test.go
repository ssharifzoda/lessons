// calculator_test.go
package calculator

import "testing"

// Тест функции Add
func TestAdd(t *testing.T) {
	result := Add(2, 3)
	expected := 5

	if result != expected {
		t.Errorf("Add(2, 3) = %d; expected %d", result, expected)
	}
}

// Тест функции Multiply
func TestMultiply(t *testing.T) {
	result := Multiply(4, 5)
	expected := 20

	if result != expected {
		t.Errorf("Multiply(4, 5) = %d; expected %d", result, expected)
	}
}
