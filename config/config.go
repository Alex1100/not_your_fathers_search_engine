package config

import "os"

type Topics struct {
	UpsertLink string
}

type PubSubConfig struct {
	ProjectID string
	Topics    *Topics
}

type Config struct {
	PubSubConfig *PubSubConfig
}

func ReadConfig() *Config {
	return &Config{
		&PubSubConfig{
			ProjectID: os.Getenv("project_id"),
			Topics: &Topics{
				UpsertLink: os.Getenv("upsert_link_topic_id"),
			},
		},
	}
}
