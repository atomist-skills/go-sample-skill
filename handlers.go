/*
 * Copyright © 2022 Atomist, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"github.com/atomist-skills/go-skill"
	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
	"olympos.io/encoding/edn"
)

// Mapping for types in the incoming event payload
type GitCommitAuthor struct {
	Name  string `json:"git.user/name"`
	Login string `json:"git.user/login"`
}

type GitOrg struct {
	Name              string `json:"git.org/name"`
	InstallationToken string `json:"github.org/installation-token"`
	Url               string `json:"git.provider/url"`
}

type GitRepo struct {
	Name          string `json:"git.repo/name"`
	DefaultBranch string `json:"git.repo/default-branch"`
	Org           GitOrg `json:"git.repo/org"`
	SourceId      string `json:"git.repo/source-id"`
}

type GitCommit struct {
	Sha     string          `json:"git.commit/sha"`
	Message string          `json:"git.commit/message"`
	Author  GitCommitAuthor `json:"git.commit/author"`
	Repo    GitRepo         `json:"git.commit/repo"`
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

// Handler to transact a commit signature on pushes
func TransactCommitSignature(ctx skill.EventContext) skill.Status {

	for _, e := range ctx.Event.Subscription.Result {
		commit := skill.Decode[GitCommit](e[0])
		err := ProcessCommit(ctx, commit)
		if err != nil {
			return skill.Status{
				Code:   1,
				Reason: fmt.Sprintf("Failed to transact signature for %s", commit.Sha),
			}
		}
	}

	return skill.Status{
		Code:   0,
		Reason: fmt.Sprintf("Successfully transacted commit signature for %d commit", len(ctx.Event.Subscription.Result)),
	}
}

func ProcessCommit(ctx skill.EventContext, commit GitCommit) error {
	gitCommit, err := GetCommit(ctx, &commit)
	if err != nil {
		return err
	}

	var verified edn.Keyword
	if *gitCommit.Commit.Verification.Verified {
		verified = Verified
	} else {
		verified = NotVerified
	}

	err = ctx.Transact([]any{GitRepoEntity{
		EntityType: "git/repo",
		Entity:     "$repo",
		SourceId:   commit.Repo.SourceId,
		Url:        commit.Repo.Org.Url,
	}, GitCommitEntity{
		EntityType: "git/commit",
		Entity:     "$commit",
		Sha:        commit.Sha,
		Repo:       "$repo",
		Url:        commit.Repo.Org.Url,
	}, GitCommitSignatureEntity{
		EntityType: "git.commit/signature",
		Commit:     "$commit",
		Signature:  *gitCommit.Commit.Verification.Signature,
		Status:     verified,
		Reason:     *gitCommit.Commit.Verification.Reason,
	}})
	if err != nil {
		return err
	}

	ctx.Log.Printf("Transacted commit signature for %s", commit.Sha)
	return err
}

// Obtain commit information from GitHub
func GetCommit(ctx skill.EventContext, commit *GitCommit) (*github.RepositoryCommit, error) {
	var client *github.Client

	if commit.Repo.Org.InstallationToken != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: commit.Repo.Org.InstallationToken},
		)
		tc := oauth2.NewClient(ctx.Context, ts)
		client = github.NewClient(tc)
	} else {
		client = github.NewClient(nil)
	}

	gitCommit, _, err := client.Repositories.GetCommit(ctx.Context, commit.Repo.Org.Name, commit.Repo.Name, commit.Sha, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return gitCommit, err
}
