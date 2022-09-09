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
	"github.com/atomist-skills/go-skill"
	"olympos.io/encoding/edn"
)

// OnCommit maps the incoming commit of the on_push and on_commit_signature to a Go struct
type OnCommit struct {
	Sha     string `edn:"git.commit/sha"`
	Message string `edn:"git.commit/message"`
	Author  struct {
		Name  string `edn:"git.user/name"`
		Login string `edn:"git.user/login"`
	} `edn:"git.commit/author"`
	Repo struct {
		Name          string `edn:"git.repo/name"`
		DefaultBranch string `edn:"git.repo/default-branch"`
		Org           struct {
			Name              string `edn:"git.org/name"`
			InstallationToken string `edn:"github.org/installation-token"`
			Url               string `edn:"git.provider/url"`
		} `edn:"git.repo/org"`
		SourceId string `edn:"git.repo/source-id"`
	} `edn:"git.commit/repo"`
	Refs []struct {
		Name string `edn:"git.ref/name"`
		Type struct {
			Ident edn.Keyword `edn:"db/ident"`
		} `edn:"git.ref/type"`
	} `edn:"git.ref/refs"`
}

// OnCommitSignature maps the incoming commit signature of the on_commit_signature to a Go struct
type OnCommitSignature struct {
	Signature string `edn:"git.commit.signature/signature"`
	Reason    string `edn:"git.commit.signature/reason"`
	Status    struct {
		Ident edn.Keyword `edn:"db/ident"`
	} `edn:"git.commit.signature/status"`
}

// GitRepoEntity provides mappings for a :git/repo entity
type GitRepoEntity struct {
	skill.Entity
	SourceId string `edn:"git.repo/source-id"`
	Url      string `edn:"git.provider/url"`
}

// GitCommitEntity provides mappings for a :git/commit entity
type GitCommitEntity struct {
	skill.Entity
	Sha  string        `edn:"git.commit/sha"`
	Repo GitRepoEntity `edn:"git.commit/repo"`
	Url  string        `edn:"git.provider/url"`
}

// GitCommitSignatureEntity provides mappings for a :git.commit/signature entity
type GitCommitSignatureEntity struct {
	skill.Entity
	Commit    GitCommitEntity `edn:"git.commit.signature/commit"`
	Signature string          `edn:"git.commit.signature/signature,omitempty"`
	Reason    string          `edn:"git.commit.signature/reason,omitempty"`
	Status    edn.Keyword     `edn:"git.commit.signature/status,omitempty"`
}

const (
	Verified    edn.Keyword = "git.commit.signature/VERIFIED"
	NotVerified             = "git.commit.signature/NOT_VERIFIED"
)
