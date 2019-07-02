package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Gist structure
type Gist struct {
	GitPullURL string `json:"git_pull_url"`
}

func main() {
	username := os.Args[1]
	syncFolder := os.Args[2]

	if _, err := os.Stat(syncFolder); os.IsNotExist(err) {
		os.Mkdir(syncFolder, os.ModePerm)
	}

	url := fmt.Sprintf("https://api.github.com/users/%s/gists", username)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var gists []Gist
	err = json.Unmarshal([]byte(body), &gists)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Syncing ", len(gists), " gists")

	for _, gist := range gists {
		repoName := gist.GitPullURL[strings.LastIndex(gist.GitPullURL, "/")+1 : strings.LastIndex(gist.GitPullURL, ".git")]
		repoFolder := filepath.Join(syncFolder, repoName)

		if _, err := os.Stat(repoFolder); err == nil {
			log.Println("Pulling: ", repoName)
			cmd := exec.Command("git", "pull", gist.GitPullURL)
			cmd.Dir = repoFolder
			err = cmd.Run()
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			log.Println("Cloning: ", repoName)
			cmd := exec.Command("git", "clone", gist.GitPullURL)
			cmd.Dir = syncFolder
			err := cmd.Run()
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}
