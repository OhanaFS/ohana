package random_phrase_generator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomPhraseGenerator(t *testing.T) {

	pg := New()

	phrase := pg.GenerateRandomPhrase()
	assert.Equal(t, phrase, "QuarrelsomeColorfulGaur")
	t.Log(phrase)
}
