package image

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"slices"
	"strings"
)

func GetLatestDigest(dockerClient *client.Client, imageName string) (digest string, err error) {
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

	ctx := context.Background()
	res, err := dockerClient.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to pull `%s`: %v", imageName, err)
	}

	defer res.Close()

	_, err = io.ReadAll(res)
	if err != nil {
		return "", err
	}

	args := filters.NewArgs(filters.Arg("reference", repository))
	imageSummaries, err := dockerClient.ImageList(ctx, types.ImageListOptions{
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

	args := []string{"build", "--iidfile", iidFile.Name()}
	if dockerfilePath != "" {
		args = append(args, "--file", dockerfilePath)
	}

	args = append(args, buildContext)

	cmd := exec.Command("docker", args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	slog.Info(cmd.String())

	err = cmd.Run()

	imageId, err := os.ReadFile(iidFile.Name())
	if err != nil {
		return "", err
	}

	return string(imageId), nil
}
