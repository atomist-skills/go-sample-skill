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
	"testing"

	"github.com/atomist-skills/go-skill"
	"github.com/atomist-skills/go-skill/test"
	"github.com/atomist-skills/go-skill/util"
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
