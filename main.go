package main

import "fmt"

func describe(value interface{}) {
	switch v := value.(type) {
	case string:
		fmt.Printf("Строка: %s (длина: %d)\n", v, len(v))
	case int:
		fmt.Printf("Целое число: %d\n", v)
	case bool:
		fmt.Printf("Булево значение: %t\n", v)
	default:
		fmt.Printf("Неизвестный тип: %T\n", v)
	}
}

func main() {
	describe("Привет") // Строка: Привет (длина: 6)
	describe(42)       // Целое число: 42
	describe(true)     // Булево значение: true
	describe(3.14)     // Неизвестный тип: float64
}
