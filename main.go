package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

/*
Response represents the json structure of the response payload
*/
type Response struct {
	Project   string   `json:"project"`
	Owner     string   `json:"owner"`
	Committer string   `json:"committer"`
	Commits   int      `json:"commits"`
	Language  []string `json:"language"`
}

/*
Contributor represents the json structure of contributors taken from https://api.github.com/repos/golang/go/contributors
*/
type Contributor struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	RecievedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
	Contributions     int    `json:"contributions"`
}

func handler(w http.ResponseWriter, r *http.Request) {

	// Split URL into parts and check for correct number of url segments
	parts := strings.Split(r.URL.Path, "/")
	if !(len(parts) == 6 && parts[5] != "" || len(parts) == 7 && parts[6] == "") { // might be 6 or 7 depending on the last slash
		http.Error(w, "Bad Request: the correct url is \"/projectinfo/v1/github.com/<username>/<repo>\"", http.StatusBadRequest)
		return
	}

	// Check that the provided url is formatted correctly
	if parts[1] != "projectinfo" || parts[2] != "v1" || parts[3] != "github.com" {
		http.Error(w, "Bad Request: the correct url is \"/projectinfo/v1/github.com/<username>/<repo>\"", http.StatusBadRequest)
		return
	}

	repo := parts[4] + "/" + parts[5] // <username>/<repo>

	// Request the languages of the repository from the github API
	var langjson map[string]interface{}
	langdata, err := http.Get("https://api.github.com/repos/" + repo + "/languages")
	if err != nil {
		http.Error(w, "Could not reach language of repo\n", http.StatusNotFound)
		return
	}

	// On brand spanking new repos the langauges are empty, but it still passes parsing
	err = json.NewDecoder(langdata.Body).Decode(&langjson)
	if err != nil {
		http.Error(w, "Languages: Error parsing the expected JSON body", http.StatusConflict)
		fmt.Fprintln(w, err.Error())
		return
	}

	// Request the contributors of the repository from the github API
	var contjson []Contributor
	contdata, err := http.Get("https://api.github.com/repos/" + repo + "/contributors")
	if err != nil {
		http.Error(w, "Could not reach contributors of repo\n", http.StatusNotFound)
		return
	}

	// On brand spanking new repos the contributors file is not even there, in that case i get an EOF error, so set a flag to no contributors
	contributors := true
	err = json.NewDecoder(contdata.Body).Decode(&contjson)
	if err != nil {
		if err == io.EOF {
			contributors = false
		} else {
			http.Error(w, "Contributors: Error parsing the expected JSON body", http.StatusNotFound) // will be a result when repo doesn't exist
			return
		}
	}

	// Taken from stackoverflow: golang-getting-a-slice-of-keys-from-a-map
	keys := make([]string, len(langjson))
	i := 0
	for k := range langjson {
		keys[i] = k
		i++
	}

	// Create and fill the response
	resp := new(Response)
	resp.Project = parts[5]
	resp.Owner = parts[4]
	if contributors {
		// The contributors are always sorted by commits in descending order
		resp.Committer = contjson[0].Login
		resp.Commits = contjson[0].Contributions
	}
	resp.Language = keys

	// encode the response to the response writer
	json.NewEncoder(w).Encode(resp)
}

// Initialize API server and handling
func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+port, nil)

}
