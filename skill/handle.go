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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func Handle(handlers Handlers) {
	log.Print("Starting server...")
	http.HandleFunc("/", createHttpHandler(handlers))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func createHttpHandler(handlers Handlers) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().UnixMilli()
		traceId := r.Header.Get("x-cloud-trace-context")

		var env MessageEnvelope
		err := json.NewDecoder(r.Body).Decode(&env)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data, _ := base64.StdEncoding.DecodeString(env.Message.Data)
		var event EventIncoming
		err = json.Unmarshal(data, &event)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		logger, loggingClient := InitLogging(event.WorkspaceId, event.CorrelationId, env.Message.MessageId, traceId, event.Subscription.Name, event.Skill)
		defer loggingClient.Close()

		logger.Println("Cloud Run execution started")

		if handle, ok := handlers[event.Subscription.Name]; ok {
			log.Printf("Invoking event handler '%s'", event.Subscription.Name)
			eventContext := EventContext{
				Data: event.Subscription.Result,
				Log:  logger,
			}

			defer func() {
				if err := recover(); err != nil {
					SendStatus(event, Status{
						Code:   1,
						Reason: fmt.Sprintf("Unsuccessfully invoked handler %s/%s@%s", event.Skill.Namespace, event.Skill.Name, event.Subscription.Name),
					})
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					logger.Printf("Cloud Run execution took %d ms, finished with status: 'ok'", time.Now().UnixMilli()-start)
					return
				}
			}()

			status := handle(eventContext)
			SendStatus(event, status)
			w.WriteHeader(201)
		} else {
			log.Printf("Event handler '%s' not found", event.Subscription.Name)
			w.WriteHeader(404)
		}
		logger.Printf("Cloud Run execution took %d ms, finished with status: 'ok'", time.Now().UnixMilli()-start)
	}
}
