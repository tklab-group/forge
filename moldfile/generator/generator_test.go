package generator

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"regexp"
	"strings"
	"testing"
)

func TestGenerateMoldfile(t *testing.T) {
	type args struct {
		dockerfilePath string
		buildContext   string
	}
	tests := []struct {
		name        string
		args        args
		wantRegexps []string
		wantErr     bool
	}{
		{
			name: "simple",
			args: args{
				dockerfilePath: "testdata/simple.Dockerfile",
				buildContext:   "testdata",
			},
			wantRegexps: []string{
				`FROM ubuntu@sha256:\w+ as base`,
				``,
				`RUN apt-get update && \\`,
				`    apt-get install -y wget=[\w\.\-]+`,
				`RUN echo "test" >> test\.txt`,
				``,
				`# 2nd build stage`,
				`FROM ubuntu@sha256:\w+`,
				``,
				`COPY --from=base test\.txt \.`,
				`RUN apt-get update && \\`,
				`    apt-get install -y wget=[\w\.\-]+`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMoldfile, err := GenerateMoldfile(tt.args.dockerfilePath, tt.args.buildContext)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateMoldfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			output := gotMoldfile.ToString()
			split := strings.Split(output, "\n")

			assert.Len(t, split, len(tt.wantRegexps))
			for i := 0; i < len(split); i++ {
				gotLine := split[i]
				wantLineRegexp := fmt.Sprintf("^%s$", tt.wantRegexps[i])

				assert.Regexp(t, regexp.MustCompile(wantLineRegexp), gotLine)
			}
		})
	}
}
