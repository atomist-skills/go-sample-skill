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
