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
	}

	g := goldie.New(t, goldie.WithFixtureDir(path.Join(runGoldenFileDir, "parse")))
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f, err := os.Open(path.Join(runTestDataDir, test.fileName))
			defer f.Close()
			require.NoError(t, err)

			got, err := ParseRunInstruction(f)

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
