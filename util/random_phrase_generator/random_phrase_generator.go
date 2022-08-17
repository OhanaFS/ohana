package random_phrase_generator

import (
	"math/rand"
	"strings"
)

type PhraseGenerator struct {
	Adjectives []string `json:"adjectives"`
	Nouns      []string `json:"nouns"`
}

func New() *PhraseGenerator {

	pg := &PhraseGenerator{
		Adjectives: GetAdjectives(),
		Nouns:      GetNouns(),
	}
	return pg
}

func (pg PhraseGenerator) GenerateRandomPhrase() string {

	// We'll assume AdjectiveAdjectiveNoun

	return StringArrayToPascalCase([]string{
		pg.Adjectives[rand.Intn(len(pg.Adjectives))],
		pg.Adjectives[rand.Intn(len(pg.Adjectives))],
		pg.Nouns[rand.Intn(len(pg.Nouns))],
	})

}

func StringArrayToPascalCase(arr []string) string {

	for i, v := range arr {
		arr[i] = strings.ToUpper(v[:1]) + v[1:]
		arr[i] = strings.Replace(arr[i], " ", "", -1)
	}

	return strings.Join(arr, "")
}
