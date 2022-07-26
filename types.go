package main

import "olympos.io/encoding/edn"

// Mapping for types in the incoming event payload
type GitCommitAuthor struct {
	Name  string `edn:"git.user/name"`
	Login string `edn:"git.user/login"`
}

type GitOrg struct {
	Name              string `edn:"git.org/name"`
	InstallationToken string `edn:"github.org/installation-token"`
	Url               string `edn:"git.provider/url"`
}

type GitRepo struct {
	Name          string `edn:"git.repo/name"`
	DefaultBranch string `edn:"git.repo/default-branch"`
	Org           GitOrg `edn:"git.repo/org"`
	SourceId      string `edn:"git.repo/source-id"`
}

type GitCommit struct {
	Sha     string          `edn:"git.commit/sha"`
	Message string          `edn:"git.commit/message"`
	Author  GitCommitAuthor `edn:"git.commit/author"`
	Repo    GitRepo         `edn:"git.commit/repo"`
}

type GitCommitSignature struct {
	Signature string      `edn:"git.commit.signature/signature"`
	Reason    string      `edn:"git.commit.signature/reason"`
	Status    edn.Keyword `edn:"git.commit.signature/status"`
}

// Mapping for entities that we want to transact
type GitRepoEntity struct {
	EntityType edn.Keyword `edn:"schema/entity-type"`
	Entity     string      `edn:"schema/entity,omitempty"`
	SourceId   string      `edn:"git.repo/source-id"`
	Url        string      `edn:"git.provider/url"`
}

type GitCommitEntity struct {
	EntityType edn.Keyword `edn:"schema/entity-type"`
	Entity     string      `edn:"schema/entity,omitempty"`
	Sha        string      `edn:"git.commit/sha"`
	Repo       string      `edn:"git.commit/repo"`
	Url        string      `edn:"git.provider/url"`
}

type GitCommitSignatureEntity struct {
	EntityType edn.Keyword `edn:"schema/entity-type"`
	Entity     string      `edn:"schema/entity,omitempty"`
	Commit     string      `edn:"git.commit.signature/commit"`
	Signature  string      `edn:"git.commit.signature/signature,omitempty"`
	Reason     string      `edn:"git.commit.signature/reason,omitempty"`
	Status     edn.Keyword `edn:"git.commit.signature/status,omitempty"`
}

const (
	Verified    edn.Keyword = "git.commit.signature/VERIFIED"
	NotVerified             = "git.commit.signature/NOT_VERIFIED"
)
