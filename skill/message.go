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

package skill

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"log"
	"os"
)

func SendStatus(event EventIncoming, status Status) error {
	message := StatusHandlerResponse{
		ApiVersion:    "1",
		CorrelationId: event.CorrelationId,
		Team: Team{
			Id: event.WorkspaceId,
		},
		Skill:  event.Skill,
		Status: status,
	}

	encodedMessage, _ := json.Marshal(message)

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, "atomist-skill-production")
	if err != nil {
		return err
	}
	defer client.Close()

	t := client.Topic(os.Getenv("ATOMIST_TOPIC"))
	t.EnableMessageOrdering = true

	publishResult := t.Publish(ctx, &pubsub.Message{
		Data:        encodedMessage,
		OrderingKey: event.CorrelationId,
	})

	serverId, err := publishResult.Get(ctx)
	if err != nil {
		return err
	}
	log.Printf("Successfully sent message with '%s'", serverId)
	return nil
}
