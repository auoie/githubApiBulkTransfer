# GitHub API Bulk Transfer

## Todo

- Required destination user API key to verify that we have ownership of the destination account
- Check that all of the input repositories are owned by our source user

## Usage

Use `go build` to generate `./githubtransfer`.

There are four required flags:

- `--source-user`. This is the login for the user containing the repositories.
  This is actually not necessary.
  The source user secret is sufficient.
  I am including it for greater redudancy so that it's clear where the repositories are coming from.
- `--destination-user`. This is the login for the user to send the repositories to.
- `--source-user-secret`. This is a personal access token for the source user.
  I gave it the `repo` [scope](https://docs.github.com/en/developers/apps/building-oauth-apps/scopes-for-oauth-apps).
  I'm not entirely sure which scope is sufficient and necessary for transferring repositories.
- `--destination-user-secret`. This is a personal access token for the destination user.
  I didn't give it any scopes.
  This is not necessary.
  I included it just to verify that I have ownership of the account I'm sending the repositories to.

The program reads the repositories from `stdin`.
It expects one repository name per line and expects at least one repository total.

If you have the GitHub CLI `gh` installed, then you can pipe the output into the binary to transfer 1000 of your repositories.

```bash
gh repo list --limit 1000 --json name --jq ".[] | .name" | ./githubtransfer \
--source-user <source-login> \
--destination-user <destination-login> \
--source-user-secret <source-personal-access-token> \
--destination-user-secret <destination-personal-access-token>
```

Suppose you have your repositories listed in a file `repos.txt`.

```txt
repo1
repo2
repo3
```

Then you can `cat` the file and pipe it into the binary.

```bash
cat repos.txt | ./githubtransfer \
--source-user <source-login> \
--destination-user <destination-login> \
--source-user-secret <source-personal-access-token> \
--destination-user-secret <destination-personal-access-token>
```

The response

```
error in issuing transfer request <repository> from user <source-user> to user <destination-user>: job scheduled on GitHub side; try again later
```

is normal.
That just means that the destination user has to accept the transfer request.

After the program is done, the email for the destination user will contain a bunch of verification emails.
You have to click on each of the links in order for the transfers to go through.

## GitHub Secret

To generate a GitHub secret, go to Settings > Developer Settings > Personal Access Tokens.
Depending on what you want, you can select different scopes.
See the source user token, I set the `repo` OAuth scope on.
For the destination user token, I didn't set any of the scopes.

## Notes

- See here for a discussion of how to make a CLI with go. https://news.ycombinator.com/item?id=23318137
- See here for the library I'm using to make the CLI. https://github.com/urfave/cli/blob/main/docs/v2/manual.md
- See here for the basic code I'm following. https://github.com/jdbean/github-api-bulk-transfer
- Convert JSON to Go. https://mholt.github.io/json-to-go/
- Convert curl request to Go. https://mholt.github.io/curl-to-go/
- See here for the documentation for transferring a GitHub repository. https://docs.github.com/en/rest/repos/repos#transfer-a-repository
- I don't know what `context.Context` does.

## Generate `repos.txt`

To generate a list of 100 of your respositories, use

```bash
gh repo list --limit 100 --json name --jq ".[] | .name" > repos.txt
```

## Commands Ran

```bash
go mod init githubtransfer
touch main.go
go get github.com/urfave/cli/v2
go get github.com/google/go-github/v44
go get golang.org/x/oauth2
```

## Read JSON Response Body from HTTP Request

- https://stackoverflow.com/questions/17156371/how-to-get-json-response-from-http-get
- Basically just use `json.NewDecder(res.Body).Decode(...)`
