package parser

import (
	"github.com/kr/pretty"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"path"
	"strings"
	"testing"
)

const moldfileTestDataDir = "testdata/moldfile"
const moldfileGoldenFileDir = "testdata/moldfile/golden"

func TestParseMoldFile(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		isError  bool
		err      error
	}{
		{
			name:     "simple",
			fileName: "simple.mold",
			isError:  false,
			err:      nil,
		},
	}

	g := goldie.New(t, goldie.WithFixtureDir(path.Join(moldfileGoldenFileDir, "parse")))
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			f, err := os.Open(path.Join(moldfileTestDataDir, test.fileName))
			defer f.Close()
			require.NoError(t, err)

			got, err := ParseMoldFile(f)

			if test.isError {
				assert.Error(t, err)
				// TODO: Check the value of err
			} else {
				assert.NoError(t, err)
			}

			g.Assert(t, test.fileName, []byte(pretty.Sprint(got)))
		})

	}
}

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

			r, err := newReader(strings.NewReader(test.target))
			got, err := checkNextInstructionType(r)

			assert.NoError(t, err)
			assert.Equal(t, test.expected, got)

			// Check if resetting offset is ok
			assert.Equal(t, len(test.target), r.Len())

			fromReader, err := io.ReadAll(r)
			assert.NoError(t, err)
			assert.Equal(t, test.target, string(fromReader))
		})
	}
}
