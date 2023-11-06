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

const vdiffTestDataDir = "testdata/vdiff"
const vdiffGoldenFileDir = "testdata/vdiff/golden"

func TestVDiffMoldfiles(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		moldfile1 string
		moldfile2 string
		isError   bool
		err       error
	}{
		{
			name:      "simple",
			moldfile1: "simple1.mold",
			moldfile2: "simple2.mold",
			isError:   false,
			err:       nil,
		},
	}

	g := goldie.New(t, goldie.WithFixtureDir(vdiffGoldenFileDir))
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			f1, err := os.Open(path.Join(vdiffTestDataDir, test.moldfile1))
			defer f1.Close()
			require.NoError(t, err)

			f2, err := os.Open(path.Join(vdiffTestDataDir, test.moldfile2))
			defer f2.Close()
			require.NoError(t, err)

			m1, err := ParseMoldFile(f1)
			require.NoError(t, err)

			m2, err := ParseMoldFile(f2)
			require.NoError(t, err)

			got, err := VDiffMoldfiles(m1, m2)

			if test.isError {
				assert.Error(t, err)
				// TODO: Check the value of err
			} else {
				assert.NoError(t, err)
			}

			g.Assert(t, test.name, []byte(pretty.Sprint(got)))
			g.AssertJson(t, test.name+"_json", got)
		})
	}
}
