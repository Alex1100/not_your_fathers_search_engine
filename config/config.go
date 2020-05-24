package config

var projectID string = "not-your-fathers-search-engine"
var upsertLinkTopicID string = "upsert_link"

type Topics struct {
	UpsertLink string
}

type PubSubConfig struct {
	ProjectID string
	Topics *Topics
}

type Config struct {
	PubSubConfig *PubSubConfig
}

func ReadConfig() *Config {
	return &Config{
		&PubSubConfig{
			ProjectID: projectID,
			Topics: &Topics {
				UpsertLink: upsertLinkTopicID,
			},
		},
	}
}