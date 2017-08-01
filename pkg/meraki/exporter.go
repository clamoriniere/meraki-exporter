package meraki

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/cedriclam/meraki-exporter/pkg/meraki/client"
)

// Exporter represent the Meraki exporter
type Exporter struct {
	config  *Config
	clients []*client.Client
}

// NewExporter return new Exporter instance
func NewExporter(config *Config) *Exporter {
	return &Exporter{
		config:  config,
		clients: []*client.Client{},
	}
}

// ListenAndServe listens on the TCP network address addr
// and then calls Serve with handler to handle requests
// on incoming connections.
func (e *Exporter) ListenAndServe() error {
	if err := e.InitExporter(); err != nil {
		fmt.Printf("Unable to init the Exporter, err:%s \n", err)
		return err
	}

	stop := make(chan struct{})
	defer close(stop)

	// start the scraping
	go e.run(stop)

	// Finally Expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("ListenAndServe on ", e.config.Addr)
	return http.ListenAndServe(e.config.Addr, nil)
}

// InitExporter used to init the exporter resources
func (e *Exporter) InitExporter() error {
	for _, token := range e.config.Tokens {
		client, err := client.NewClient(e.config.BaseUrl, token, e.config.APIVersion)
		if err != nil {
			fmt.Println("New client error: ", err)
			return err
		}
		e.clients = append(e.clients, client)
	}

	return nil
}

func (e *Exporter) run(stop chan struct{}) error {
	var err error

	ticker := time.NewTicker(e.config.Freq)
	defer ticker.Stop()
	stopped := false
	for !stopped {
		select {
		case <-stop:
			stopped = true
			break
		case <-ticker.C:
			if err = e.scrapAll(); err != nil {
				fmt.Println("Unable to scrap all, err:", err)
				stopped = true
				break
			}
		}
	}

	return err
}

func (e *Exporter) scrapAll() error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(e.clients))

	for _, ec := range e.clients {
		wg.Add(1)
		go func(c *client.Client, errs chan error) {
			defer wg.Done()
			if err := e.scrap(c); err != nil {
				errs <- err
			}

		}(ec, errChan)
	}

	wg.Wait()
	close(errChan)
	hasErrors := false
	for err := range errChan {
		hasErrors = true
		fmt.Println("Scrap error: ", err)
	}

	if hasErrors {
		return fmt.Errorf("Error during the sraping")
	}
	return nil
}

func init() {
	prometheus.MustRegister(perfGauge)
}

const (
	networkIdLabelKey      = "network_id"
	organizationIdLabelKey = "organization_id"
	deviceIdLabelKey       = "device_id"
)

var perfGauge = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "meraki_device_perf",
		Help: "Meraki device perf gauge",
	},
	[]string{organizationIdLabelKey, networkIdLabelKey, deviceIdLabelKey},
)

func (e *Exporter) scrap(c *client.Client) error {

	return scrapPerf(c, perfGauge)
}

func scrapPerf(c *client.Client, gauge *prometheus.GaugeVec) error {
	fmt.Printf("Scrap at: %s \n", time.Now().String())
	labels := prometheus.Labels{}
	orgas, err := c.Organizations()
	if err != nil {
		return err
	}
	for _, orga := range orgas.Items {
		labels[organizationIdLabelKey] = orga.Name
		networks, err := c.Organization(orga.Id).Networks()
		if err != nil {
			continue
		}
		for _, network := range networks.Items {
			labels[networkIdLabelKey] = network.Name
			devices, err := c.Organization(orga.Id).Network(network.Id).Devices()
			if err != nil {
				return err
			}
			for _, device := range devices.Items {
				labels[deviceIdLabelKey] = device.Name
				perf, err := c.Organization(orga.Id).Network(network.Id).Device(device.Serial).Performance()
				if err != nil {
				} else {
					if g, err := gauge.GetMetricWith(labels); err == nil {
						g.Set(float64(perf.PerfScore))
					}
				}
			}
		}
	}

	return nil
}
