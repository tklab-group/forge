package parser

import (
	"github.com/kr/pretty"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"testing"
)

const otherTestDataDir = "testdata/other"
const otherGoldenFileDir = "testdata/other/golden"

func TestParseOtherInstruction(t *testing.T) {
	tests := []struct {
		name            string
		fileName        string
		enableMultiline bool
		isError         bool
		err             error
	}{
		{
			name:            "simple case with WORKDIR",
			fileName:        "workdir.mold",
			enableMultiline: false,
			isError:         false,
			err:             nil,
		},
		{
			name:            "comment",
			fileName:        "comment.mold",
			enableMultiline: false,
			isError:         false,
			err:             nil,
		},
		{
			name:            "blank line",
			fileName:        "blankline.mold",
			enableMultiline: false,
			isError:         false,
			err:             nil,
		},
		{
			name:            "multilines with HEALTHCHECK",
			fileName:        "healthcheck.mold",
			enableMultiline: true,
			isError:         false,
			err:             nil,
		},
	}

	g := goldie.New(t, goldie.WithFixtureDir(path.Join(otherGoldenFileDir, "parse")))
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f, err := os.Open(path.Join(otherTestDataDir, test.fileName))
			defer f.Close()
			require.NoError(t, err)

			got, err := ParseOtherInstruction(f, test.enableMultiline)

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

func Test_otherInstruction_ToString(t *testing.T) {
	tests := []struct {
		name            string
		fileName        string
		enableMultiline bool
	}{
		{
			name:            "simple case with WORKDIR",
			fileName:        "workdir.mold",
			enableMultiline: false,
		},
		{
			name:            "comment",
			fileName:        "comment.mold",
			enableMultiline: false,
		},
		{
			name:            "blank line",
			fileName:        "blankline.mold",
			enableMultiline: false,
		},
		{
			name:            "multilines with HEALTHCHECK",
			fileName:        "healthcheck.mold",
			enableMultiline: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filePath := path.Join(otherTestDataDir, test.fileName)
			f, err := os.Open(filePath)
			defer f.Close()
			require.NoError(t, err)

			instruction, err := ParseOtherInstruction(f, test.enableMultiline)
			require.NoError(t, err)

			got := instruction.ToString()

			expected, err := os.ReadFile(filePath)
			require.NoError(t, err)

			assert.Equal(t, string(expected), got)
		})
	}
}
