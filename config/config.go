package config

import "os"

// Topics Holds the details about crawl_from_source_topic_id
// and all other topics within Google Cloud Pub/Sub
type Topics struct {
	CrawlFromSource string
}

// PubSubConfig Contains details regarding configuration
// for Google Cloud Pub/Sub
type PubSubConfig struct {
	ProjectID string
	Topics    *Topics
}

// Config Configuration interface for Application wide Configurations
type Config struct {
	PubSubConfig *PubSubConfig
}

// ReadConfig access app config from anywhere
func ReadConfig() *Config {
	return &Config{
		&PubSubConfig{
			ProjectID: os.Getenv("project_id"),
			Topics: &Topics{
				CrawlFromSource: os.Getenv("crawl_from_source_topic_id"),
			},
		},
	}
}
