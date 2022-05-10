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

package handler

import (
	"cloud.google.com/go/logging"
	"context"
	"log"
)

func InitLogging(workspaceId string, correlationId string, messageId string, skill Skill) (*log.Logger, *logging.Client) {
	ctx := context.Background()

	client, err := logging.NewClient(ctx, "atomist-skill-production")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	logName := "skills_logging"
	logger := client.Logger(logName, logging.CommonLabels(map[string]string{
		"correlation_id":  correlationId,
		"workspace_id":    workspaceId,
		"event_id":        messageId,
		"skill_id":        skill.Id,
		"skill_namespace": skill.Namespace,
		"skill_name":      skill.Name,
		"skill_version":   skill.Version,
	})).StandardLogger(logging.Info)

	return logger, client
}
