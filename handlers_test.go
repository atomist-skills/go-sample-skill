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
	"github.com/atomist-skills/go-skill"
	"github.com/atomist-skills/go-skill/test"
	"log"
	"olympos.io/encoding/edn"
	"reflect"
	"testing"
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
		Configuration: skill.Configuration{
			Name:       "default",
			Parameters: []skill.ParameterValue{{Name: "test", Value: "test"}},
		},
		TxData: "test/datalog/push.tx",
	}, t)

	if result.Results[0].Subscription != "on_push" {
		t.Errorf("Expected different subscription match")
	}
	if c := len(result.Results[0].Results); c != 1 {
		t.Errorf("Expecting 1 commit result; instead received %d", c)
	}
}

func TestProcessCommit(t *testing.T) {
	commit := GitCommit{
		Sha: "d2c6724307f007755fc770944fd7bc5ff55933b0",
		Repo: GitRepo{
			Name:          "go-sample-skill",
			DefaultBranch: "main",
			Org: GitOrg{
				Name: "atomist-skills",
				Url:  "https://github.com/",
			},
			SourceId: "123456",
		},
	}
	req := skill.RequestContext{
		Log: skill.Logger{
			Debug: func(msg string) {
				log.Print(msg)
			},
			Debugf: func(format string, a ...any) {
				log.Printf(format, a...)
			},
			Info: func(msg string) {
				log.Print(msg)
			},
			Infof: func(format string, a ...any) {
				log.Printf(format, a...)
			},
		},
		Transact: func(entities interface{}) error {
			switch reflect.TypeOf(entities).Kind() {
			case reflect.Slice:
				s := reflect.ValueOf(entities)
				if s.Len() != 3 {
					t.Errorf("Expected 3 entities, got %d", s.Len())
				}

				repoEntity := s.Index(0).Interface().(GitRepoEntity)
				if !reflect.DeepEqual(repoEntity, GitRepoEntity{
					EntityType: edn.Keyword("git/repo"),
					Entity:     "$repo",
					SourceId:   "123456",
					Url:        "https://github.com/",
				}) {
					t.Errorf("Repo entity not expected")
				}
				commitEntity := s.Index(1).Interface().(GitCommitEntity)
				if !reflect.DeepEqual(commitEntity, GitCommitEntity{
					EntityType: edn.Keyword("git/commit"),
					Entity:     "$commit",
					Repo:       "$repo",
					Sha:        commit.Sha,
					Url:        "https://github.com/",
				}) {
					t.Errorf("Commit entity not expected")
				}
				commitSignatureEntity := s.Index(2).Interface().(GitCommitSignatureEntity)
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
