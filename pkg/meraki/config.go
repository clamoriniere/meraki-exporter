package meraki

import (
	"fmt"
	"strings"
	"time"

	"net/url"

	flag "github.com/spf13/pflag"
)

const (
	BaseUrlDefaut           = "https://dashboard.meraki.com/api/"
	APIVersionDefaut        = "v0"
	APIScrapFrequencyDefaut = time.Duration(5 * time.Second)
)

// Config used to store the exporter configuration
type Config struct {
	Addr       string            // The address to listen on for HTTP requests.
	BaseUrl    *url.URL          // the Meraki Dashboard API Base Url
	APIVersion string            // the Meraki Dashboard API version
	Tokens     []string          // list of Meraki Dashboard API token
	Freq       time.Duration     // the scraping frequency
	Labels     map[string]string // Labels added automaticaly to all metrics
}

// NewConfig returns a new Config instance
func NewConfig() *Config {
	return &Config{
		APIVersion: APIVersionDefaut,
		Labels:     make(map[string]string),
	}
}

// Init used to initialize the configuration
func (c *Config) Init() error {
	flag.StringVar(&c.Addr, "listen-address", ":8080", "The address to listen on for HTTP requests.")
	flag.StringSliceVar(&c.Tokens, "api-token", []string{}, "The Meraki dashboard API tokens.")
	baseUrl := ""
	flag.StringVar(&baseUrl, "api-base-url", BaseUrlDefaut, "The Meraki dashboard API base URL")
	flag.StringVar(&c.APIVersion, "api-version", APIVersionDefaut, "The Meraki dashboard API version")
	flag.DurationVar(&c.Freq, "api-freq", APIScrapFrequencyDefaut, "Time between 2 API scraping")
	tmpLabels := []string{}
	flag.StringSliceVar(&tmpLabels, "label", []string{}, "Used to add label (key:value) to all metrics")

	flag.Parse()

	var err error
	if c.BaseUrl, err = url.Parse(baseUrl); err != nil {
		return err
	}

	for _, val := range tmpLabels {
		keyval := strings.Split(val, ":")
		if len(keyval) != 2 {
			return fmt.Errorf("unable to parse the label %s", val)
		}
		c.Labels[keyval[0]] = keyval[1]
	}

	return nil
}
