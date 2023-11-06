package image

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tklab-group/forge/docker"
	"testing"
)

func TestGetLatestDigest(t *testing.T) {
	t.Parallel()
	t.Run("with digest", func(t *testing.T) {
		t.Parallel()

		dockerClient, err := docker.GetDockerClient()
		require.NoError(t, err)

		got, err := GetLatestDigest(dockerClient, "hello-world@sha256:4f53e2564790c8e7856ec08e384732aa38dc43c52f02952483e3f003afbf23db")
		assert.NoError(t, err)
		assert.Equal(t, "sha256:4f53e2564790c8e7856ec08e384732aa38dc43c52f02952483e3f003afbf23db", got)
	})

	t.Run("with tag", func(t *testing.T) {
		t.Parallel()

		dockerClient, err := docker.GetDockerClient()
		require.NoError(t, err)

		got, err := GetLatestDigest(dockerClient, "alpine:20230901")
		assert.NoError(t, err)
		assert.Equal(t, "sha256:f2d1645cd73c7e54584dc225da0b5229d19223412d719669ebda764f41396853", got)
	})
}
