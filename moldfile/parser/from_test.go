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
			name:     "lowercase",
			fileName: "lowercase.mold",
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

			g.Assert(t, test.name, []byte(pretty.Sprint(got)))
		})

	}
}
