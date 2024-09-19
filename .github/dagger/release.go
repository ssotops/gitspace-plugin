package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"dagger.io/dagger"
	"github.com/Masterminds/semver/v3"
	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

func main() {
	if err := publishRelease(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func publishRelease(ctx context.Context) error {
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}
	defer client.Close()

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %v", err)
	}
	projectRoot := filepath.Dir(filepath.Dir(wd))
	fmt.Printf("Project root: %s\n", projectRoot)

	src := client.Host().Directory(projectRoot)

	// Determine the new version
	newVersion, err := determineNewVersion(ctx)
	if err != nil {
		return fmt.Errorf("failed to determine new version: %v", err)
	}

	// Build gitspace-plugin
	build := client.Container().
		From("golang:1.23.0").
		WithDirectory("/src", src).
		WithWorkdir("/src").
		WithExec([]string{
			"go", "build",
			"-ldflags", fmt.Sprintf("-X 'main.Version=%s'", newVersion),
			"-o", "gitspace-plugin.so",
		})

	// Export the built plugin
	_, err = build.File("gitspace-plugin.so").Export(ctx, filepath.Join(projectRoot, "gitspace-plugin.so"))
	if err != nil {
		return fmt.Errorf("failed to export gitspace-plugin.so: %v", err)
	}
	fmt.Println("Successfully built and exported: gitspace-plugin.so")

	// Run tests
	test := client.Container().
		From("golang:1.23.0").
		WithDirectory("/src", src).
		WithWorkdir("/src").
		WithExec([]string{"go", "test", "./..."})

	if _, err := test.Sync(ctx); err == nil {
		fmt.Println("Tests passed. Creating GitHub release...")
		if err := createGitHubRelease(ctx, projectRoot, newVersion); err != nil {
			return fmt.Errorf("failed to create GitHub release: %v", err)
		}
	} else {
		return fmt.Errorf("tests failed: %v", err)
	}

	return nil
}

func determineNewVersion(ctx context.Context) (string, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return "", fmt.Errorf("GITHUB_TOKEN environment variable is not set")
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Fetch the latest release
	latestRelease, _, err := client.Repositories.GetLatestRelease(ctx, "ssotops", "gitspace-plugin")
	if err != nil && err.(*github.ErrorResponse).Response.StatusCode != 404 {
		return "", fmt.Errorf("failed to fetch latest release: %v", err)
	}

	if latestRelease == nil || latestRelease.TagName == nil {
		// If there's no release yet, start with v1.0.0
		return "v1.0.0", nil
	}

	// Parse the latest version and increment the patch number
	v, err := semver.NewVersion(*latestRelease.TagName)
	if err != nil {
		return "", fmt.Errorf("failed to parse latest version: %v", err)
	}
	return fmt.Sprintf("v%d.%d.%d", v.Major(), v.Minor(), v.Patch()+1), nil
}

func createGitHubRelease(ctx context.Context, projectRoot, newVersion string) error {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return fmt.Errorf("GITHUB_TOKEN environment variable is not set")
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	release, _, err := client.Repositories.CreateRelease(ctx, "ssotops", "gitspace-plugin", &github.RepositoryRelease{
		TagName:    github.String(newVersion),
		Name:       github.String(fmt.Sprintf("Release %s", newVersion)),
		Body:       github.String("Description of the release"),
		Draft:      github.Bool(false),
		Prerelease: github.Bool(false),
	})
	if err != nil {
		return fmt.Errorf("failed to create GitHub release: %v", err)
	}

	// Upload gitspace-plugin.so
	filename := "gitspace-plugin.so"
	filepath := filepath.Join(projectRoot, filename)
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open %s: %v", filename, err)
	}
	defer file.Close()

	_, _, err = client.Repositories.UploadReleaseAsset(ctx, "ssotops", "gitspace-plugin", *release.ID, &github.UploadOptions{
		Name: filename,
	}, file)
	if err != nil {
		return fmt.Errorf("failed to upload asset %s: %v", filename, err)
	}
	fmt.Printf("Successfully uploaded: %s\n", filename)

	fmt.Printf("Release %s created: %s\n", newVersion, *release.HTMLURL)
	return nil
}
