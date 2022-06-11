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
	"go-sample-skill/skill"
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

func PrintCommit(ctx skill.EventContext) skill.Status {

	for _, e := range ctx.Data {
		commit := skill.Decode[GitCommit](e[0])
		ctx.Log.Printf("Seen commit %s %s", commit.Sha, commit.Message)
	}

	return skill.Status{
		Code:   0,
		Reason: fmt.Sprintf("Successfully printed %d commits", len(ctx.Data)),
	}
}
