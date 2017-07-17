package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/cedriclam/meraki-exporter/pkg/meraki/api"
)

const (
	ciscoTokenHeaderKey  = "X-Cisco-Meraki-API-Key"
	contentTypeHeaderKey = "Content-Type"
	acceptHeaderKey      = "Accept"

	contentTypeJsonHeaderValue = "application/json"
)

// Client represent the Meraki client structure
type Client struct {
	baseURL *url.URL
	token   string
	version string

	httpClient *http.Client
}

// ClientInterface regroups all Meraki Client methods
type ClientInterface interface {
	// Organizations returns list of organizations
	Organization(id string) OrganizationInterface
	// Organizations returns list of organizations
	Organizations() (*api.OrganizationList, error)
}

// NewClient returns new instance of a Meraka api client
func NewClient(urlBase *url.URL, token, version string) (*Client, error) {
	versionPath, err := url.Parse("./" + version + "/")
	if err != nil {
		return nil, err
	}
	return &Client{
		baseURL: urlBase.ResolveReference(versionPath),
		token:   token,
		version: version,

		httpClient: &http.Client{},
	}, nil
}

// Organizations returns list of organizations
func (c *Client) Organizations() (*api.OrganizationList, error) {
	req, err := c.newRequest("GET", "./organizations", c.baseURL, nil)
	if err != nil {
		return nil, err
	}

	organizations := &api.OrganizationList{Version: api.Version(c.version)}
	_, err = c.do(req, &organizations.Items)
	if err != nil {
		return nil, err
	}
	return organizations, err
}

// OrganizationInterface interface for the Organization Meraki API endpoint
type OrganizationInterface interface {
	// Get returns the Organization api object corresponding to the ID
	Get() (*api.Organization, error)
	// Networks returns the list of network associated to the Organization ID
	Networks() (*api.NetworkList, error)
	// Network returns the Network API Client endpoint associated to the ID
	Network(id string) NetworkInterface
}

// Organizations returns organization client interface
func (c *Client) Organization(id int) OrganizationInterface {
	rel, _ := url.Parse(fmt.Sprintf("./organizations/%d/", id))
	u := c.baseURL.ResolveReference(rel)

	organizationClient := &organizationClient{c, u}
	return organizationClient
}

type organizationClient struct {
	*Client
	baseURL *url.URL
}

// Get returns the Organization API object corresponding to the ID
func (oc *organizationClient) Get() (*api.Organization, error) {
	req, err := oc.newRequest("GET", "", oc.baseURL, nil)
	if err != nil {
		return nil, err
	}

	organization := &api.Organization{Version: api.Version(oc.version)}
	_, err = oc.do(req, organization)
	if err != nil {
		return nil, err
	}

	return organization, err
}

// Networks returns the list of network associated to the Organization ID
func (oc *organizationClient) Networks() (*api.NetworkList, error) {
	req, err := oc.newRequest("GET", "./networks", oc.baseURL, nil)
	if err != nil {
		return nil, err
	}

	networks := &api.NetworkList{Version: api.Version(oc.version)}
	_, err = oc.do(req, &networks.Items)
	if err != nil {
		return nil, err
	}
	return networks, err
}

// Network returns the Network API Client endpoint associated to the ID
func (oc *organizationClient) Network(id string) NetworkInterface {
	rel, _ := url.Parse(fmt.Sprintf("./networks/%s/", id))
	u := oc.baseURL.ResolveReference(rel)

	newClient := &networkClient{oc, u}
	return newClient
}

type NetworkInterface interface {
	// Get returns the Network api object corresponding to the ID
	Get() (*api.Network, error)
	// Devices returns list of device associated the a Network
	Devices() (*api.DeviceList, error)
	// Device return Device client interface
	Device(id string) DeviceInterface
}

type networkClient struct {
	*organizationClient
	baseURL *url.URL
}

// Get returns the Network api object corresponding to the ID
func (nc *networkClient) Get() (*api.Network, error) {
	req, err := nc.newRequest("GET", "", nc.baseURL, nil)
	if err != nil {
		return nil, err
	}

	network := &api.Network{Version: api.Version(nc.version)}
	_, err = nc.do(req, network)
	if err != nil {
		return nil, err
	}

	return network, err
}

// Devices returns the list of Devices associated to the Network
func (nc *networkClient) Devices() (*api.DeviceList, error) {
	req, err := nc.newRequest("GET", "./devices", nc.baseURL, nil)
	if err != nil {
		return nil, err
	}

	devices := &api.DeviceList{Version: api.Version(nc.version)}
	_, err = nc.do(req, &devices.Items)
	if err != nil {
		return nil, err
	}
	return devices, err
}

// Device returns the Network API Client endpoint associated to the ID
func (nc *networkClient) Device(id string) DeviceInterface {
	rel, _ := url.Parse(fmt.Sprintf("./devices/%s/", id))
	u := nc.baseURL.ResolveReference(rel)

	newClient := &deviceClient{nc, u}
	return newClient
}

type deviceClient struct {
	*networkClient
	baseURL *url.URL
}

type DeviceInterface interface {
	// Get returns the Device api object corresponding to the ID
	Get() (*api.Device, error)
	// Performances returns new Performance client interface
	Performance() (*api.Performance, error)
}

// Get returns the Network api object corresponding to the ID
func (dc *deviceClient) Get() (*api.Device, error) {
	req, err := dc.newRequest("GET", "", dc.baseURL, nil)
	if err != nil {
		return nil, err
	}

	device := &api.Device{Version: api.Version(dc.version)}
	_, err = dc.do(req, device)
	if err != nil {
		return nil, err
	}

	return device, err
}

// Performance returns Performance api object associated to a Device
func (dc *deviceClient) Performance() (*api.Performance, error) {
	req, err := dc.newRequest("GET", "./performance", dc.baseURL, nil)
	if err != nil {
		return nil, err
	}

	performance := &api.Performance{Version: api.Version(dc.version)}
	resp, err := dc.do(req, performance)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Reply status code not 200, current code:%d", resp.StatusCode)
	}
	return performance, err
}

// Generic methods
func (c *Client) addCommonHeaders(header *http.Header) {
	header.Set(contentTypeHeaderKey, contentTypeJsonHeaderValue)
	header.Set(acceptHeaderKey, contentTypeJsonHeaderValue)
	header.Set(ciscoTokenHeaderKey, c.token)
}

func (c *Client) newRequest(method, path string, baseUrl *url.URL, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	u := baseUrl.ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set(contentTypeHeaderKey, contentTypeJsonHeaderValue)
	}
	req.Header.Set(acceptHeaderKey, contentTypeJsonHeaderValue)
	req.Header.Set(ciscoTokenHeaderKey, c.token)
	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return resp, nil
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return resp, err
	}
	return resp, err
}
