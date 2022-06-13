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
	"context"
	"fmt"
	"github.com/atomist-skills/go-skill"
	"github.com/google/go-github/v45/github"
	"olympos.io/encoding/edn"
)

type GitCommitAuthor struct {
	Name  string `json:"git.user/name"`
	Login string `json:"git.user/login"`
}

type GitOrg struct {
	Name              string `json:"git.org/name"`
	InstallationToken string `json:"github.org/installation-token"`
	Url               string `json:":git.provider/url"`
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

type GitRepoEntity struct {
	EntityType edn.Keyword `edn:"schema/entity-type"`
	Entity     string      `edn:"schema/entity"`
	SourceId   string      `edn:"git.repo/source-id"`
	Url        string      `edn:"git.provider/url"`
}

type GitCommitEntity struct {
	EntityType edn.Keyword `edn:"schema/entity-type"`
	Entity     string      `edn:"schema/entity"`
	Sha        string      `edn:"git.commit/sha"`
	Repo       string      `edn:"git.commit/repo"`
	Url        string      `edn:"git.provider/url"`
}

type GitCommitSignatureEntity struct {
	EntityType edn.Keyword `edn:"schema/entity-type"`
	Entity     string      `edn:"schema/entity"`
	Commit     string      `edn:"git.commit.signature/commit"`
	Signature  string      `edn:"git.commit.signature/signature"`
	Reason     string      `edn:"git.commit.signature/reason"`
	Verified   edn.Keyword `edn:"git.commit.signature/verified"`
}

const (
	Verified    edn.Keyword = "git.commit.signature/VERIFIED"
	NotVerified             = "git.commit.signature/NOT_VERIFIED"
)

func TransactCommitSignature(ctx skill.EventContext) skill.Status {

	for _, e := range ctx.Event.Subscription.Result {
		commit := skill.Decode[GitCommit](e[0])

		client := github.NewClient(nil)
		gitCommit, _, err := client.Repositories.GetCommit(context.Background(), commit.Repo.Org.Name, commit.Repo.Name, commit.Sha, nil)
		if err != nil {
			fmt.Println(err)
			return skill.Status{
				Code:   1,
				Reason: fmt.Sprintf("Failed to obtain commit for %s", commit.Sha),
			}
		}

		var verified edn.Keyword
		if *gitCommit.Commit.Verification.Verified {
			verified = Verified
		} else {
			verified = NotVerified
		}

		ctx.Transact([]any{GitRepoEntity{
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
			Commit:    "$commit",
			Signature: *gitCommit.Commit.Verification.Signature,
			Verified:  verified,
			Reason:    *gitCommit.Commit.Verification.Reason,
		}})

		ctx.Log.Printf("Transacted commit signature for %s", commit.Sha)
	}

	return skill.Status{
		Code:   0,
		Reason: fmt.Sprintf("Successfully printed %d commits", len(ctx.Event.Subscription.Result)),
	}
}
