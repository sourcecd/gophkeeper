package dataparser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParsers(t *testing.T) {
	testCases := []struct {
		name     string
		dType    string
		val      []byte
		expected []byte
	}{
		{
			name:     "textTest",
			dType:    "TEXT",
			val:      []byte("textTest"),
			expected: []byte("textTest"),
		},
		{
			name:     "cardTest",
			dType:    "CARD",
			val:      []byte("4220855426222389"),
			expected: []byte("4220855426222389"),
		},
		{
			name:     "credentialsTest",
			dType:    "CREDENTIALS",
			val:      []byte("login password"),
			expected: []byte("login password"),
		},
		{
			name:     "binaryTest",
			dType:    "BINARY",
			val:      []byte{102, 105},
			expected: []byte{102, 105},
		},
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			parser := Dataparser(v.dType, v.val)
			b, err := parser.Parse()
			require.NoError(t, err)
			require.Equal(t, v.expected, b)
		})
	}
}
