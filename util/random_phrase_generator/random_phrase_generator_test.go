package random_phrase_generator

import "testing"

func TestRandomPhraseGenerator(t *testing.T) {

	pg := New()

	phrase := pg.GenerateRandomPhrase()
	t.Log(phrase)
}
