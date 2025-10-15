package main

import "fmt"

func main() {
	votes := []string{"Анна", "Петр", "Анна", "Мария", "Петр", "Анна"}
	results := make(map[string]int)

	// Подсчет голосов
	for _, candidate := range votes {
		results[candidate]++
	}

	// Находим победителя
	winner := ""
	maxVotes := 0

	for candidate, count := range results {
		fmt.Printf("%s: %d голосов\n", candidate, count)

		if count > maxVotes {
			maxVotes = count
			winner = candidate
		}
	}

	fmt.Printf("\nПобедитель: %s с %d голосами!\n", winner, maxVotes)
}
