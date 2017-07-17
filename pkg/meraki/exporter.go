package meraki

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/cedriclam/meraki-exporter/pkg/meraki/client"
)

// Exporter represent the Meraki exporter
type Exporter struct {
	config *Config
}

// NewExporter return new Exporter instance
func NewExporter(config *Config) *Exporter {
	return &Exporter{
		config: config,
	}
}

// ListenAndServe listens on the TCP network address addr
// and then calls Serve with handler to handle requests
// on incoming connections.
func (e *Exporter) ListenAndServe() error {
	client, err := client.NewClient(e.config.BaseUrl, e.config.Token, e.config.APIVersion)
	if err != nil {
		fmt.Println("New client error: ", err)
		return err
	}
	if err := test(client); err != nil {
		fmt.Println("Test error: ", err)
	}

	// Finally Expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(e.config.Addr, nil)
}

func test(c *client.Client) error {
	orgas, err := c.Organizations()
	if err != nil {
		return err
	}
	for _, orga := range orgas.Items {
		fmt.Printf("%s \n", orga)
		networks, err := c.Organization(orga.Id).Networks()
		if err != nil {
			continue
		}
		for _, network := range networks.Items {
			fmt.Printf("---  %s \n", network)
			devices, err := c.Organization(orga.Id).Network(network.Id).Devices()
			if err != nil {
				return err
			}
			for _, device := range devices.Items {
				fmt.Printf("--- ---  %s \n", device)
				perf, err := c.Organization(orga.Id).Network(network.Id).Device(device.Serial).Performance()
				if err != nil {
					return err
				}
				fmt.Printf("--- --- ---  %s \n", perf)
			}
		}
	}

	return nil
}
