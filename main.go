package main

import "fmt"

func main() {
	var a, b int

	fmt.Print("Введите первое число: ")
	fmt.Scan(&a)

	fmt.Print("Введите второе число: ")
	fmt.Scan(&b)

	fmt.Println(a, b)
}
