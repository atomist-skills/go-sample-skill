package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/pubsub"
	eh "go-sample-skill/handler"
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
	Skill         eh.Skill             `json:"skill"`
	Subscription  SubscriptionIncoming `json:"subscription"`
	WorkspaceId   string               `json:"team_id"`
	LogUrl        string               `json:"log_url"`
	Secrets       []Secret             `json:"secrets"`
}

type EventHandler func([][]map[string]json.RawMessage) eh.Status

var HandlerRegistry = map[string]EventHandler{
	"on_push": eh.PrintCommit,
}

func main() {
	log.Print("Starting server...")
	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
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

	if handle, ok := HandlerRegistry[event.Subscription.Name]; ok {
		log.Printf("Invoking event handler '%s'", event.Subscription.Name)
		status := handle(event.Subscription.Result)

		message := eh.StatusHandlerResponse{
			ApiVersion:    "1",
			CorrelationId: event.CorrelationId,
			Team: eh.Team{
				Id: event.WorkspaceId,
			},
			Skill:  event.Skill,
			Status: status,
		}

		encodedMessage, _ := json.Marshal(message)

		ctx := context.Background()
		client, err := pubsub.NewClient(ctx, "atomist-skill-production")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("Successfully sent message with '%s'", serverId)
		w.WriteHeader(201)
	} else {
		log.Printf("Event handler '%s' not found", event.Subscription.Name)
		w.WriteHeader(404)
	}
}
