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
	"github.com/atomist-skills/go-skill/util"
	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
	"olympos.io/encoding/edn"
	"reflect"
)



// TransactCommitSignature processed incoming Git pushes and transacts the commit signature
// as returned by GitHub
func TransactCommitSignature(ctx context.Context, req skill.RequestContext) skill.Status {
	result := req.Event.Context.Subscription.Result[0]
	commit := util.Decode[GitCommit](result[0])
	gitCommit, err := getCommit(ctx, req, &commit)

	if err != nil {
		return skill.Status{
			State:  skill.Failed,
			Reason: fmt.Sprintf("Failed to obtain commit signature for %s", commit.Sha),
		}
	}

	err = transactCommitSignature(ctx, req, commit, gitCommit)
	if err != nil {
		return skill.Status{
			State:  skill.Failed,
			Reason: fmt.Sprintf("Failed to transact signature for %s", commit.Sha),
		}
	}

	return skill.Status{
		State:  skill.Completed,
		Reason: fmt.Sprintf("Successfully transacted commit signature for %d commit", len(req.Event.Context.Subscription.Result)),
	}
}

// LogCommitSignature handles new commit signature entities as they are transacted into
// the database and logs the signature
func LogCommitSignature(ctx context.Context, req skill.RequestContext) skill.Status {
	result := req.Event.Context.Subscription.Result[0]
	commit := util.Decode[GitCommit](result[0])
	signature := util.Decode[GitCommitSignature](result[1])

	req.Log.Infof("Commit %s is signed and verified by: %s", commit.Sha, signature.Signature)

	return skill.Status{
		State:  skill.Completed,
		Reason: "Detected signed and verified commit",
	}
}

// LogWebhookBody logs incoming webhook bodies
func LogWebhookBody(ctx context.Context, req skill.RequestContext) skill.Status {
	body := req.Event.Context.Webhook.Request.Body

	req.Log.Infof("Webhook body: %s", body)

	return skill.Status{
		State:  skill.Completed,
		Reason: "Handled incoming webhook event",
	}
}

// transactCommitSignature transact the commit signature facts
func transactCommitSignature(ctx context.Context, req skill.RequestContext, commit GitCommit, gitCommit *github.RepositoryCommit) error {
	var verified edn.Keyword
	if *gitCommit.Commit.Verification.Verified {
		verified = Verified
	} else {
		verified = NotVerified
	}
	var signature string
	verification := *gitCommit.Commit.Verification
	if !reflect.ValueOf(verification.Signature).IsNil() {
		signature = *verification.Signature
	}

	err := req.Transact([]any{GitRepoEntity{
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
		Signature:  signature,
		Status:     verified,
		Reason:     *gitCommit.Commit.Verification.Reason,
	}})
	if err != nil {
		return err
	}

	req.Log.Infof("Transacted commit signature for %s", commit.Sha)
	return err
}

// getCommit obtains commit information from GitHub
func getCommit(ctx context.Context, req skill.RequestContext, commit *GitCommit) (*github.RepositoryCommit, error) {
	var client *github.Client

	if commit.Repo.Org.InstallationToken != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: commit.Repo.Org.InstallationToken},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	} else {
		client = github.NewClient(nil)
	}

	gitCommit, _, err := client.Repositories.GetCommit(ctx, commit.Repo.Org.Name, commit.Repo.Name, commit.Sha, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return gitCommit, err
}
