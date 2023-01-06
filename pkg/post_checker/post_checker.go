package post_checker

import (
	"errors"
	"net/http"
)

// Checkpost tries to get the headers of the target url and returns error if it fails
func CheckPost(url string) (err error) {
	resp, err := http.Head(url)
	if err != nil {
		return
	}
	if resp.StatusCode >= 400 {
		err = errors.New(resp.Status)
	}
	return
}
