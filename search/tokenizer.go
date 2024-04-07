package search

import (
	"strings"
	"unicode"

	snowballeng "github.com/kljensen/snowball/english"
)

// analyze analyzes the text and returns a slice of tokens.
func analyze(text string) []string {
	tokens := tokenize(text)
	tokens = lowercaseFilter(tokens)
	tokens = stopwordFilter(tokens)
	tokens = stemmerFilter(tokens)
	return tokens
}

// tokenize returns a slice of tokens for the given text.
func tokenize(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		// Split on any character that is not a letter or a number.
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

// lowercaseFilter returns a slice of tokens normalized to lower case.
func lowercaseFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = strings.ToLower(token)
	}
	return r
}

// stopwordFilter returns a slice of tokens with stop words removed.
func stopwordFilter(tokens []string) []string {
	var stopwords = map[string]struct{}{
		"a": {}, "and": {}, "be": {}, "have": {}, "i": {},
		"in": {}, "of": {}, "that": {}, "the": {}, "to": {},
		"it": {}, "for": {}, "not": {}, "on": {}, "with": {},
		"as": {}, "you": {}, "do": {}, "at": {}, "this": {},
		"but": {}, "his": {}, "by": {}, "from": {}, "they": {},
		"we": {}, "say": {}, "her": {}, "she": {}, "or": {},
		"an": {}, "will": {}, "my": {}, "one": {}, "all": {},
		"www": {}, "com": {}, "org": {}, "net": {}, "io": {},
		"https": {}, "http": {}, "html": {}, "php": {}, "asp": {}, "co": {},
	}
	r := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if _, ok := stopwords[token]; !ok {
			r = append(r, token)
		}
	}
	return r
}

// stemmerFilter returns a slice of stemmed tokens.
func stemmerFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = snowballeng.Stem(token, false)
	}
	return r
}
