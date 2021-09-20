package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-playground/webhooks/v6/github"
	log "github.com/sirupsen/logrus"
)

func main() {
	var (
		gitBranch     = os.Getenv("GIT_BRANCH")
		githubSecret  = os.Getenv("GITHUB_SECRET")
		gitProvider   = os.Getenv("GIT_PROVIDER")
		repositoryURL = os.Getenv("REPOSITORY_URL")
		sshPrivateKey = []byte(os.Getenv("SSH_PRIVATE_KEY"))
		workingDir    = os.Getenv("WORKING_DIR")
	)

	if err := checkoutGitRepository(sshPrivateKey, workingDir, repositoryURL, gitBranch); err != nil {
		log.Fatalf("failed to checkout branch %q of repository %q into directory %q: %w", gitBranch, repositoryURL, workingDir, err)
	}

	var webhookHandler http.HandlerFunc
	var err error
	switch gitProvider {
	case "github":
		webhookHandler, err = handleGithubWebhook(gitBranch, githubSecret, workingDir, sshPrivateKey)
		if err != nil {
			log.Fatalf("failed to prepare webhook: %w", err)
		}
	default:
		log.Fatalf("unknown git provider %q; current implementation only accepts \"github\"", gitProvider)
	}

	http.HandleFunc("/webhook", webhookHandler)
	http.HandleFunc("/healthz", healthcheck)

	log.Println("Listening for requests on port 80...")
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatalf("failed to start HTTP server: %w", err)
	}
}

func checkoutGitRepository(sshPrivateKey []byte, workingDir, repositoryURL, gitBranch string) error {
	publicKeys, err := ssh.NewPublicKeys("git", sshPrivateKey, "")
	if err != nil {
		return fmt.Errorf("failed to generate public keys: %w", err)
	}

	cloneOptions := git.CloneOptions{
		Auth: publicKeys,
		URL:  repositoryURL,
	}
	repo, err := git.PlainClone(workingDir, false, &cloneOptions)
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("cloned repository is broken: %w", err)
	}

	checkoutOptions := git.CheckoutOptions{
		Branch: plumbing.ReferenceName(gitBranch),
	}
	if err := worktree.Checkout(&checkoutOptions); err != nil {
		return fmt.Errorf("failed to checkout branch %q: %w", gitBranch, err)
	}

	return nil
}

func handleGithubWebhook(gitBranch, githubSecret, workingDir string, sshPrivateKey []byte) (http.HandlerFunc, error) {
	githubWebhook, err := github.New(github.Options.Secret(githubSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to set up GitHub webhook: %w", err)
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		payload, err := githubWebhook.Parse(r, github.PushEvent)
		switch err {
		case nil: // no error, do nothing.
		case github.ErrEventNotFound:
			log.Info("webhook called for non-existing event, skipping")
		default:
			log.Errorf("could not parse payload: %q", err)
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			return
		}

		pushEvent := payload.(github.PushPayload)
		if pushEvent.Ref != gitBranch {
			log.Info("webhook called for ref %q, skipping", pushEvent.Ref)
			return
		}

		if err := updateRepository(workingDir, sshPrivateKey); err != nil {
			log.Errorf("failed to update repository: %w", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		log.Info("repository updated following event from GitHub")
	}

	return handler, nil
}

func updateRepository(workingDir string, sshPrivateKey []byte) error {
	publicKeys, err := ssh.NewPublicKeys("git", sshPrivateKey, "")
	if err != nil {
		return fmt.Errorf("failed to generate public keys: %w", err)
	}

	repo, err := git.PlainOpen(workingDir)
	if err != nil {
		return fmt.Errorf("repository in %q directory is broken: %w", workingDir, err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("repository in %q directory is broken: %w", workingDir, err)
	}

	pullOptions := git.PullOptions{
		Auth:       publicKeys,
		RemoteName: "origin",
	}
	if err := worktree.Pull(&pullOptions); err != nil {
		return fmt.Errorf("failed to pull repository: %q", err)
	}

	return nil
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	// The HTTP server is always healthy.
}
