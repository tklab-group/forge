package parser

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func Test_checkNextInstructionType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		target   string
		expected string
	}{
		{
			name:     "run",
			target:   "RUN echo 'hello'",
			expected: "run",
		},
		{
			name:     "from",
			target:   "FROM ubuntu:latest",
			expected: "from",
		},
		{
			name:     "healthcheck",
			target:   "HEALTHCHECK --interval=5m --timeout=3s \\\n  CMD curl -f http://localhost/ || exit 1",
			expected: "healthcheck",
		},
		{
			name:     "comment",
			target:   "# this is comment line",
			expected: "comment",
		},
		{
			name:     "blank line",
			target:   "\n",
			expected: "blank",
		},
		{
			name:     "unknown",
			target:   "UNKNOWN cmd",
			expected: "unknown",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			br := bytes.NewReader([]byte(test.target))
			got, err := checkNextInstructionType(br)

			assert.NoError(t, err)
			assert.Equal(t, test.expected, got)

			// Check if resetting offset is ok
			assert.Equal(t, len(test.target), br.Len())

			fromReader, err := io.ReadAll(br)
			assert.NoError(t, err)
			assert.Equal(t, test.target, string(fromReader))
		})
	}
}
