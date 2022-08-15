package random_phrase_generator

import (
	"math/rand"
	"strings"
	"time"
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

	// We'll assume 2 x ( 1 adj, 1 noun )

	phrase := ""

	// It will be in Pascal Case
	for i := 0; i < 2; i++ {
		rand.Seed(time.Now().UnixNano())
		word := pg.Adjectives[rand.Intn(len(pg.Adjectives))]
		word = strings.ToUpper(word[:1]) + word[1:]
		word = strings.Replace(word, " ", "", -1)
		phrase += word
		word = pg.Nouns[rand.Intn(len(pg.Nouns))]
		word = strings.ToUpper(word[:1]) + word[1:]
		word = strings.Replace(word, " ", "", -1)
		phrase += word
	}

	return phrase
}
