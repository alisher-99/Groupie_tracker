package api

import (
	"net/http"
	"strconv"
	"strings"
)

func CheckIndexRequest(r *http.Request) int {
	if r.URL.Path != "/" {
		return http.StatusNotFound
	} else if r.Method != http.MethodGet && r.Method != http.MethodOptions && r.Method != http.MethodHead {
		return http.StatusMethodNotAllowed
	}
	return 0
}

func CheckArtistRequest(r *http.Request) int {
	id := strings.TrimPrefix(r.URL.Path, "/artists/")
	if _, err := strconv.Atoi(id); err != nil {
		return http.StatusNotFound
	}

	if r.Method != http.MethodGet && r.Method != http.MethodOptions && r.Method != http.MethodHead {
		return http.StatusMethodNotAllowed
	}
	return 0
}

func CheckSearchRequest(r *http.Request) int {
	searchinput := r.URL.Query().Get("searchinput")
	searchinput = strings.Replace(searchinput, " ", "", -1)
	searchinput = strings.Replace(searchinput, "\n", "", -1)

	if searchinput == "" {
		return http.StatusNotFound
	}
	return 0
}
