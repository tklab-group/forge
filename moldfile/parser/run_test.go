package parser

import (
	"os"
	"path"
	"testing"

	"github.com/kr/pretty"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const runTestDataDir = "testdata/run"
const runGoldenFileDir = "testdata/run/golden"

func TestParseRunInstruction(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		isError  bool
		err      error
	}{
		{
			name:     "apt simple",
			fileName: "apt-simple.mold",
			isError:  false,
			err:      nil,
		},
		{
			name:     "apt separated with backslash",
			fileName: "apt-multiline.mold",
			isError:  false,
			err:      nil,
		},
		{
			name:     "comment in apt command",
			fileName: "apt-with-comment.mold",
			isError:  false,
			err:      nil,
		},
		{
			name:     "not package manager command",
			fileName: "curl-simple.mold",
			isError:  false,
			err:      nil,
		},
		{
			name:     "not package manager command with backslash and comment",
			fileName: "curl-multiline.mold",
			isError:  false,
			err:      nil,
		},
		{
			name:     "multiple commands with `&&`",
			fileName: "multiple-commands.mold",
			isError:  false,
			err:      nil,
		},
		{
			name:     "multiple commands with `;`",
			fileName: "multiple-commands-using-semicolon.mold",
			isError:  false,
			err:      nil,
		},
		{
			name:     "apt command ends with newline",
			fileName: "apt-with-newline.mold",
			isError:  false,
			err:      nil,
		},
		{
			name:     "curl command ends with newline",
			fileName: "curl-with-newline.mold",
			isError:  false,
			err:      nil,
		},
	}

	g := goldie.New(t, goldie.WithFixtureDir(path.Join(runGoldenFileDir, "parse")))
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			f, err := os.Open(path.Join(runTestDataDir, test.fileName))
			defer f.Close()
			require.NoError(t, err)

			r, err := newReader(f)
			require.NoError(t, err)

			got, err := ParseRunInstruction(r)

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

func Test_runInstruction_ToString(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
	}{
		{
			name:     "apt simple",
			fileName: "apt-simple.mold",
		},
		{
			name:     "apt separated with backslash",
			fileName: "apt-multiline.mold",
		},
		{
			name:     "comment in apt command",
			fileName: "apt-with-comment.mold",
		},
		{
			name:     "not package manager command",
			fileName: "curl-simple.mold",
		},
		{
			name:     "not package manager command with backslash and comment",
			fileName: "curl-multiline.mold",
		},
		{
			name:     "multiple commands with `&&`",
			fileName: "multiple-commands.mold",
		},
		{
			name:     "multiple commands with `;`",
			fileName: "multiple-commands-using-semicolon.mold",
		},
		{
			name:     "apt command ends with newline",
			fileName: "apt-with-newline.mold",
		},
		{
			name:     "curl command ends with newline",
			fileName: "curl-with-newline.mold",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			filePath := path.Join(runTestDataDir, test.fileName)
			f, err := os.Open(filePath)
			defer f.Close()
			require.NoError(t, err)

			r, err := newReader(f)
			require.NoError(t, err)

			instruction, err := ParseRunInstruction(r)
			require.NoError(t, err)

			got := instruction.ToString()

			expected, err := os.ReadFile(filePath)
			require.NoError(t, err)

			assert.Equal(t, string(expected), got)
		})
	}
}
