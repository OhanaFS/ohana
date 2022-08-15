package random_phrase_generator

import "testing"

func TestRandomPhraseGenerator(t *testing.T) {

	pg, err := New()

	if err != nil {
		t.Error(err)
	}

	phrase := pg.GenerateRandomPhrase()
	t.Log(phrase)
}
