package main

import (
	"strings"
)

const maxChirpLength = 140

var badWords = map[string]struct{}{
	"kerfuffle": {},
	"sharbert":  {},
	"fornax":    {},
}

func isValidChirp(msg string) bool {
	return len(msg) <= maxChirpLength
}

func replaceBadWords(msg string) string {
	words := strings.Split(msg, " ")

	for i, word := range words {
		if _, ok := badWords[strings.ToLower(word)]; ok {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
