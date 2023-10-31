package image

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"io"
	"os"
	"os/exec"
	"slices"
	"strings"
)

func GetLatestDigest(imageName string) (digest string, err error) {
	if strings.Contains(imageName, "@") {
		return strings.Split(imageName, "@")[1], nil
	}

	var repository, tag string
	if strings.Contains(imageName, ":") {
		split := strings.Split(imageName, ":")
		repository = split[0]
		tag = split[1]
	} else {
		repository = imageName
		tag = "latest"
	}

	// TODO: Consider how to treat client
	client, err := client.NewClientWithOpts()
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	res, err := client.ImagePull(ctx, imageName, types.ImagePullOptions{})
	defer res.Close()

	if err != nil {
		return "", fmt.Errorf("failed to pull `%s`: %v", imageName, err)
	}

	_, err = io.ReadAll(res)
	if err != nil {
		return "", err
	}

	args := filters.NewArgs(filters.Arg("reference", repository))
	imageSummaries, err := client.ImageList(ctx, types.ImageListOptions{
		All:     true,
		Filters: args,
	})

	if err != nil {
		return "", fmt.Errorf("failed to get image list: %v", err)
	}

	var targetSummary *types.ImageSummary
	repoTag := fmt.Sprintf("%s:%s", repository, tag)
	for _, summary := range imageSummaries {
		if slices.Contains(summary.RepoTags, repoTag) {
			targetSummary = &summary
			break
		}
	}

	if targetSummary == nil {
		return "", fmt.Errorf("target image not found")
	}

	if len(targetSummary.RepoDigests) == 0 {
		return "", fmt.Errorf("image `%s` has no repoDigest", repoTag)
	}

	repoDigest := targetSummary.RepoDigests[0]
	split := strings.Split(repoDigest, "@")
	if len(split) != 2 {
		return "", fmt.Errorf("expected RepoDigest format is `{name}@{repoDigest}`: actual `%s`", repoDigest)
	}

	return split[1], nil
}

func BuildImageWithCLI(dockerfilePath string, buildContext string) (iid string, err error) {
	iidFile, err := os.CreateTemp("/tmp", "iid")
	if err != nil {
		return "", err
	}
	defer os.Remove(iidFile.Name())

	flags := []string{"--iidfile", iidFile.Name()}
	if dockerfilePath != "" {
		flags = append(flags, "--file", dockerfilePath)
	}

	args := append(flags, buildContext)

	cmd := exec.Command("docker", args...)

	// TODO: Output building logs
	err = cmd.Run()

	imageId, err := io.ReadAll(iidFile)
	if err != nil {
		return "", err
	}

	return string(imageId), nil
}
