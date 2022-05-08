package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v44/github"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

func readLinesMustBeNonemptyAndAtLeastLength1(rd io.Reader) ([]string, error) {
	input := bufio.NewReader(os.Stdin)
	lines := []string{}
	for {
		line, err := input.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			return nil, errors.New("repo line is empty")
		}
		lines = append(lines, line)
	}
	if len(lines) < 1 {
		return nil, errors.New("there must be at least one respository")
	}
	return lines, nil
}

func logFatalError[T any](result T, err error) T {
	if err != nil {
		log.Fatal(err)
	}
	return result
}

var (
	sourceUserFlag            = "source-user"
	destinationUserFlag       = "destination-user"
	sourceUserSecretFlag      = "source-user-secret"
	destinationUserSecretFlag = "destination-user-secret"
)

func logFatalBadClient(ctx context.Context, client *github.Client, user string) {
	if !verifyClient(ctx, client, user) {
		log.Fatalf("client for user %v is not valid", user)
	}
}

func createClient(ctx context.Context, secret string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: secret},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client
}

func verifyClient(ctx context.Context, client *github.Client, user string) bool {
	userRes, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return false
	}
	return *userRes.Login == user
}

func transferRepository(ctx context.Context, client *github.Client, sourceUser string, destinationUser string, repo string) error {
	_, _, err := client.Repositories.Transfer(ctx, sourceUser, repo, github.TransferRequest{NewOwner: destinationUser})
	if err != nil {
		return err
	}
	return nil
}

func main() {
	app := &cli.App{
		Name:  "githubtransfer",
		Usage: "Transfer all of the repositories from one user to another user on GitHub.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     sourceUserFlag,
				Required: true,
				Usage:    "Username of the user containing the repositories.",
			},
			&cli.StringFlag{
				Name:     destinationUserFlag,
				Required: true,
				Usage:    "Username of the user to send the repositories to.",
			},
			&cli.StringFlag{
				Name:     sourceUserSecretFlag,
				Required: true,
				Usage:    "GitHub secret for the user containing the repositories.",
			},
			&cli.StringFlag{
				Name:     destinationUserSecretFlag,
				Required: true,
				Usage:    "GitHub secret for the user to send the repositories to.",
			},
		},
		Action: func(cliCtx *cli.Context) error {
			lines := logFatalError(readLinesMustBeNonemptyAndAtLeastLength1(os.Stdin))
			sourceUser := cliCtx.String(sourceUserFlag)
			destinationUser := cliCtx.String(destinationUserFlag)
			sourceSecret := cliCtx.String(sourceUserSecretFlag)
			destinationSecret := cliCtx.String(destinationUserSecretFlag)
			ctx := context.Background()
			sourceClient := createClient(ctx, sourceSecret)
			destinationClient := createClient(ctx, destinationSecret)
			logFatalBadClient(ctx, sourceClient, sourceUser)
			logFatalBadClient(ctx, destinationClient, destinationUser)
			for _, repo := range lines {
				err := transferRepository(ctx, sourceClient, sourceUser, destinationUser, repo)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error in issuing transfer request %v from user %v to user %v: %v\n", repo, sourceUser, destinationUser, err)
				} else {
					fmt.Fprintf(os.Stdout, "issued transfer request of repository %v from user %v to user %v\n", repo, sourceUser, destinationUser)
				}
			}
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
