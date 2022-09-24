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
	"reflect"
	"testing"

	"github.com/atomist-skills/go-skill"
	"github.com/atomist-skills/go-skill/test"
	"github.com/atomist-skills/go-skill/util"
	"olympos.io/encoding/edn"
)

func TestSimulateOnPush(t *testing.T) {
	result := test.Simulate(test.SimulateOptions{
		Skill: skill.Skill{
			Id:        "cada301c-0e75-11ed-861d-0242ac120002",
			Namespace: "atomist",
			Name:      "go-sample-skill",
			Version:   "0.0.1",
		},
		Subscription: "datalog/subscription/on_push.edn",
		Schemata:     "datalog/schema/commit_signature.edn",
		TxData:       "testdata/datalog/transaction/push.edn",
		Configuration: skill.Configuration{
			Name: "default",
		},
	}, t)

	if result.Results[0].Subscription != "on_push" {
		t.Errorf("Expected different subscription match")
	}
	if c := len(result.Results[0].Results); c != 1 {
		t.Errorf("Expecting 1 commit result; instead received %d", c)
	}

	commit := util.Decode[GitCommitEntity](result.Results[0].Results[0][0])
	if commit.Sha != "8976e7077a86e2755486eb136103b26cef5c78d7" {
		t.Errorf("Expecting different sha; instead got %s", commit.Sha)
	}
}

func TestProcessCommit(t *testing.T) {
	commit := OnCommit{
		Sha: "d2c6724307f007755fc770944fd7bc5ff55933b0",
		Repo: struct {
			Name          string `edn:"git.repo/name"`
			DefaultBranch string `edn:"git.repo/default-branch"`
			Org           struct {
				Name              string `edn:"git.org/name"`
				InstallationToken string `edn:"github.org/installation-token"`
				Url               string `edn:"git.provider/url"`
			} `edn:"git.repo/org"`
			SourceId string `edn:"git.repo/source-id"`
		}{
			Name:          "go-sample-skill",
			DefaultBranch: "main",
			Org: struct {
				Name              string `edn:"git.org/name"`
				InstallationToken string `edn:"github.org/installation-token"`
				Url               string `edn:"git.provider/url"`
			}{
				Name: "atomist-skills",
				Url:  "https://github.com/",
			},
			SourceId: "123456",
		},
	}
	req := skill.RequestContext{
		Log: skill.Logger{
			Debug:  func(msg string) { t.Log(msg) },
			Debugf: t.Logf,
			Info:   func(msg string) { t.Log(msg) },
			Infof:  t.Logf,
		},
		Transact: func(entities interface{}) error {
			switch reflect.TypeOf(entities).Kind() {
			case reflect.Slice:
				s := reflect.ValueOf(entities)
				if s.Len() != 1 {
					t.Errorf("Expected 3 entities, got %d", s.Len())
				}

				assertRepoEntity := GitRepoEntity{
					Entity: skill.Entity{
						EntityType: edn.Keyword("git/repo"),
					},
					SourceId: "123456",
					Url:      "https://github.com/",
				}
				commitSignatureEntity := s.Index(0).Interface().(GitCommitSignatureEntity)
				commitEntity := commitSignatureEntity.Commit
				commitEntity.Entity.Entity = ""
				commitEntity.Repo.Entity.Entity = ""
				if !reflect.DeepEqual(commitEntity, GitCommitEntity{
					Entity: skill.Entity{
						EntityType: edn.Keyword("git/commit"),
					},
					Repo: assertRepoEntity,
					Sha:  commit.Sha,
					Url:  "https://github.com/",
				}) {
					t.Errorf("Commit entity as not expected")
				}
				repoEntity := commitSignatureEntity.Commit.Repo
				repoEntity.Entity.Entity = ""
				if !reflect.DeepEqual(repoEntity, assertRepoEntity) {
					t.Errorf("Repo entity not  as expected")
				}
				if !(commitSignatureEntity.Signature != "" && commitSignatureEntity.Reason == "valid") {
					t.Errorf("Signature entity not expected")
				}

				return nil
			}
			t.Error("Expected slice of entities")
			return nil
		},
	}
	ctx := context.Background()
	gitCommit, err := getCommit(ctx, req, &commit)
	if err != nil {
		t.Error("getCommit errored")
	}
	err = transactCommitSignature(context.Background(), req, commit, gitCommit)
	if err != nil {
		t.Error("transactCommitSignature errored")
	}
}
