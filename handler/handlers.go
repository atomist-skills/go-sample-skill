package handler

import (
	"encoding/json"
	"log"
)

type GitCommitAuthor struct {
	Name  string `json:"git.user/name"`
	Login string `json:"git.user/login"`
}

type GitOrg struct {
	Name              string `json:"git.org/name"`
	InstallationToken string `json:"github.org/installation-token"`
}

type GitRepo struct {
	Name          string `json:"git.repo/name"`
	DefaultBranch string `json:"git.repo/default-branch"`
	Org           GitOrg `json:"git.repo/org"`
}

type GitCommit struct {
	Sha     string          `json:"git.commit/sha"`
	Message string          `json:"git.commit/message"`
	Author  GitCommitAuthor `json:"git.commit/author"`
	Repo    GitRepo         `json:"git.commit/repo"`
}

func Decode[P interface{}](event map[string]json.RawMessage) P {
	jsonbody, _ := json.Marshal(event)
	var decoded P
	json.Unmarshal(jsonbody, &decoded)
	return decoded
}

func PrintCommit(event [][]map[string]json.RawMessage) Status {

	for _, e := range event {
		commit := Decode[GitCommit](e[0])
		log.Printf("Seen commit %s %s", commit.Sha, commit.Message)
	}

	return Status{
		Code:   0,
		Reason: "Successfully invoked scan_image",
	}
}
