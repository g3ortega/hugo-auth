package content

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/BurntSushi/toml"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	c "github.com/otiai10/copy"
)

type source struct {
	GitRepo     string
	ProjectSlug string
	DocsDir     string
	Branch      string
}

type sources struct {
	Sources map[string]source
}

// Import docs from different sources
func UpdateContent() {
	fmt.Println("Importing content...")

	token := os.Getenv("GITHUB_TOKEN")
	var config sources

	if _, err := toml.DecodeFile("sources.toml", &config); err != nil {
		fmt.Println("Something went wrong: ", err)
	}

	for _, repo := range config.Sources {
		os.RemoveAll("./content/en/docs/" + repo.ProjectSlug)
		os.MkdirAll("./content/en/docs/"+repo.ProjectSlug, 0755)

		r, err := git.PlainClone("./"+repo.ProjectSlug, false, &git.CloneOptions{
			Auth: &http.BasicAuth{
				Username: "abc123",
				Password: token,
			},
			URL:      repo.GitRepo,
			Progress: os.Stdout,
		})

		if err != nil {
			fmt.Println(err)
		}

		w, err := r.Worktree()

		err = w.Checkout(&git.CheckoutOptions{
			Branch: plumbing.NewBranchReferenceName(repo.Branch),
		})

		c.Copy("./"+repo.ProjectSlug+"/"+repo.DocsDir, "./content/en/docs/"+repo.ProjectSlug)
		os.RemoveAll("./" + repo.ProjectSlug)

		cmd := exec.Command("hugo")

		fmt.Printf("Running command and waiting for it to finish...")
		cmdErr := cmd.Run()

		if cmdErr != nil {
			fmt.Printf("Command finished with error: %v", cmdErr)
		}

		fmt.Println(r)
	}
}