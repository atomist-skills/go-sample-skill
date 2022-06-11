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
	"encoding/json"
	"log"
)

type Message struct {
	Data      string `json:"data"`
	MessageId string `json:"messageId"`
}

type MessageEnvelope struct {
	Message Message `json:"message"`
}

type Secret struct {
	Uri   string `json:"uri"`
	Value string `json:"value"`
}

type SubscriptionIncoming struct {
	Name   string                         `json:"name"`
	Tx     int64                          `json:"tx"`
	Result [][]map[string]json.RawMessage `json:"result"`
}

type EventIncoming struct {
	CorrelationId string               `json:"correlation_id"`
	Skill         Skill                `json:"skill"`
	Subscription  SubscriptionIncoming `json:"subscription"`
	WorkspaceId   string               `json:"team_id"`
	LogUrl        string               `json:"log_url"`
	Secrets       []Secret             `json:"secrets"`
}

type Skill struct {
	Id        string `json:"id"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Version   string `json:"version"`
}

type Status struct {
	Code       int8   `json:"code"`
	Reason     string `json:"reason"`
	Visibility string `json:"visibility,omitempty"`
}

type Team struct {
	Id string `json:"id"`
}

type StatusHandlerResponse struct {
	ApiVersion    string `json:"api_version"`
	CorrelationId string `json:"correlation_id"`
	Team          Team   `json:"team"`
	Status        Status `json:"status"`
	Skill         Skill  `json:"skill"`
}

type EventContext struct {
	Data [][]map[string]json.RawMessage
	Log  *log.Logger
}

type EventHandler func(ctx EventContext) Status

type Handlers map[string]EventHandler
