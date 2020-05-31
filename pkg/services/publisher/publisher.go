package publisher

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

func Publish(w io.Writer, projectID, topicID string, payload []byte) error {
	fmt.Println("PAYLOAD IS: ", payload)
	ctx := context.Background()
	// client, err := pubsub.NewClient(ctx, projectID)
	client, err := pubsub.NewClient(ctx, projectID, option.WithCredentialsFile("/Users/alexanderaleksanyan/code/go/cred_vault/not_your_fathers_search_engine/not-your-fathers-search-engine-2a8fa2e76deb.json"))
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}

	t := client.Topic(topicID)
	result := t.Publish(ctx, &pubsub.Message{
		Data: payload,
	})

	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		return fmt.Errorf("Get: %v", err)
	}
	fmt.Fprintf(w, "Published a message; msg ID: %v\n", id)
	return nil
}
