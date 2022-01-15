package lokishipper

import (
	"net/url"
	"time"

	"github.com/grafana/dskit/backoff"
	"github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
)

// NOTE the helm chart for promtail and fluent-bit also have defaults for these values, please update to match if you make changes here.
const (
	BatchWait      = 1 * time.Second
	BatchSize  int = 1024 * 1024
	MinBackoff     = 500 * time.Millisecond
	MaxBackoff     = 5 * time.Minute
	MaxRetries int = 10
	Timeout        = 10 * time.Second
)

type HTTPClientConfig = config.HTTPClientConfig
type BackoffConfig = backoff.Config
type LabelSet = model.LabelSet

// Config describes configuration for a HTTP pusher client.
type Config struct {
	URL       *url.URL
	BatchWait time.Duration
	BatchSize int

	Client HTTPClientConfig `yaml:",inline"`

	BackoffConfig BackoffConfig `yaml:"backoff_config"`
	// The labels to add to any time series or alerts when communicating with loki
	ExternalLabels LabelSet      `yaml:"external_labels,omitempty"`
	Timeout        time.Duration `yaml:"timeout"`

	// The tenant ID to use when pushing logs to Loki (empty string means
	// single tenant mode)
	TenantID string `yaml:"tenant_id"`

	StreamLagLabels []string `yaml:"stream_lag_labels"`
}

// UnmarshalYAML implement Yaml Unmarshaler
func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type raw Config
	var cfg raw
	if c.URL != nil {
		// we used flags to set that value, which already has sane default.
		cfg = raw(*c)
	} else {
		// force sane defaults.
		cfg = raw{
			BackoffConfig: backoff.Config{
				MaxBackoff: MaxBackoff,
				MaxRetries: MaxRetries,
				MinBackoff: MinBackoff,
			},
			BatchSize:       BatchSize,
			BatchWait:       BatchWait,
			Timeout:         Timeout,
			StreamLagLabels: []string{"filename"},
		}
	}

	if err := unmarshal(&cfg); err != nil {
		return err
	}

	*c = Config(cfg)
	return nil
}
