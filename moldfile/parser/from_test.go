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

const testDataDir = "testdata/from"
const goldenFileDir = "testdata/from/golden"

func TestParseFromInstruction(t *testing.T) {
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
		{
			name:     "lowercase FROM",
			fileName: "lowercase.mold",
			isError:  false,
			err:      nil,
		},
		{
			name:     "image with tag",
			fileName: "tag.mold",
			isError:  false,
			err:      nil,
		},
		{
			name:     "image with digest",
			fileName: "digest.mold",
			isError:  false,
			err:      nil,
		},
		{
			name:     "with stage name",
			fileName: "stagename.mold",
			isError:  false,
			err:      nil,
		},
		{
			name:     "with platform flag",
			fileName: "platform.mold",
			isError:  false,
			err:      nil,
		},
	}

	g := goldie.New(t, goldie.WithFixtureDir(path.Join(goldenFileDir, "parse")))
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f, err := os.Open(path.Join(testDataDir, test.fileName))
			defer f.Close()
			require.NoError(t, err)

			got, err := ParseFromInstruction(f)

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

func Test_fromInstruction_ToString(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
	}{
		{
			name:     "simple",
			fileName: "simple.mold",
		},
		{
			name:     "lowercase FROM",
			fileName: "lowercase.mold",
		},
		{
			name:     "image with tag",
			fileName: "tag.mold",
		},
		{
			name:     "image with digest",
			fileName: "digest.mold",
		},
		{
			name:     "with stage name",
			fileName: "stagename.mold",
		},
		{
			name:     "with platform flag",
			fileName: "platform.mold",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filePath := path.Join(testDataDir, test.fileName)
			f, err := os.Open(filePath)
			defer f.Close()
			require.NoError(t, err)

			instruction, err := ParseFromInstruction(f)
			require.NoError(t, err)

			got := instruction.ToString()

			expected, err := os.ReadFile(filePath)
			require.NoError(t, err)

			assert.Equal(t, string(expected), got)
		})
	}
}
