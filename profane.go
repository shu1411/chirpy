package main

import (
	"strings"
)

func containsProfaneWords(msg string) bool {
	for _, profaneWord := range profaneWords {
		if strings.Contains(msg, profaneWord) {
			return true
		}
	}
	return false
}

func cleanBody(msg string) string {
	words := strings.Split(msg, " ")
	for i, word := range words {
		for _, profaneWord := range profaneWords {
			if strings.ToLower(word) == profaneWord {
				words[i] = "****"
			}
		}
	}

	return strings.Join(words, " ")
}