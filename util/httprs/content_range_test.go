package httprs_test

import (
	"testing"

	"github.com/OhanaFS/ohana/util/httprs"
	"github.com/stretchr/testify/assert"
)

func TestContentRange(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		input    string
		expected *httprs.ContentRange
	}{
		{
			input:    "bytes 0-1/2",
			expected: &httprs.ContentRange{Start: 0, End: 1, Total: 2},
		},
		{
			input:    "bytes 0-1/*",
			expected: &httprs.ContentRange{Start: 0, End: 1, Total: -1},
		},
		{
			input:    "bytes */2",
			expected: &httprs.ContentRange{Start: -1, End: -1, Total: 2},
		},
	}

	for _, c := range cases {
		cr, err := httprs.ParseContentRange(c.input)
		assert.NoError(err)
		assert.Equal(c.expected, cr)
	}
}
