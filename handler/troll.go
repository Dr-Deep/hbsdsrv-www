/*
 * troll some folks
 */
package handler

import (
	"fmt"
	"net/http"
	"net/url"
)

const trollGitConfigFile = `[user]
name = your mommy
email = your.mom@outlook.com
`

type HandlerTroll struct{}

func (h *HandlerTroll) IsAble(url *url.URL) bool {
	switch url.Path {
	case "/.ssh/id_ed25519", "/.ssh/id_rsa":
		return true

	case "/.env":
		return true

	case "/.git/config":
		return true

	default:
		return false
	}
}

func (h *HandlerTroll) Handle(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		Error(w, http.StatusMethodNotAllowed)
		return nil
	}

	var msg = ""

	switch r.URL.Path {
	case "/.ssh/id_ed25519", "/.ssh/id_rsa":
		msg = "leck meine eier"

	case "/.env":
		msg = "sibbie"

	case "/.git/config":
		msg = trollGitConfigFile
	}

	fmt.Fprintf(w, "%s", msg)

	return nil
}
