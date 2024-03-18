package rs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Deimvis/reactionsstorage/src/models"
)

type HTTPClient struct {
	client  *http.Client
	baseUrl *url.URL
}

func NewClientHTTP(host string, port int, ssl bool) *HTTPClient {
	client := &http.Client{Timeout: 10 * time.Second}
	baseUrl := &url.URL{}
	if ssl {
		baseUrl.Scheme = "https"
	} else {
		baseUrl.Scheme = "http"
	}
	baseUrl.Host = fmt.Sprintf("%s:%d", host, port)
	return &HTTPClient{client: client, baseUrl: baseUrl}
}

func (c *HTTPClient) GetReactions(req *models.ReactionsGETRequest) (models.Response, error) {
	return c.handle(req, []models.Response{
		&models.ReactionsGETResponse200{},
	})
}

func (c *HTTPClient) AddReaction(req *models.ReactionsPOSTRequest) (models.Response, error) {
	return c.handle(req, []models.Response{
		&models.ReactionsPOSTResponse200{},
		&models.ReactionsPOSTResponse403{},
	})
}

func (c *HTTPClient) RemoveReaction(req *models.ReactionsDELETERequest) (models.Response, error) {
	return c.handle(req, []models.Response{
		&models.ReactionsDELETEResponse200{},
		&models.ReactionsDELETEResponse403{},
	})
}

func (c *HTTPClient) UpdateConfiguration(req *models.ConfiguratinPOSTRequest) (models.Response, error) {
	return c.handle(req, []models.Response{
		&models.ConfigurationPOSTResponse200{},
		&models.ConfigurationPOSTResponse415{},
		&models.ConfigurationPOSTResponse422{},
	})
}
func (c *HTTPClient) GetNamespace(req *models.NamespaceGETRequest) (models.Response, error) {
	return c.handle(req, []models.Response{
		&models.NamespaceGETResponse200{},
		&models.NamespaceGETResponse404{},
	})
}

func (c *HTTPClient) GetAvailableReactions(req *models.AvailableReactionsGETRequest) (models.Response, error) {
	return c.handle(req, []models.Response{
		&models.AvailableReactionsGETResponse200{},
		&models.AvailableReactionsGETResponse404{},
	})
}

func (c *HTTPClient) handle(req models.Request, respOptions []models.Response) (models.Response, error) {
	resp, err := c.request(req)
	if err != nil {
		return nil, err
	}
	return handleResponse(respOptions, resp)
}

func (c *HTTPClient) request(req models.Request) (*http.Response, error) {
	reqUrl := *c.baseUrl
	reqUrl.Path = req.Path()
	reqUrl.RawQuery = req.QueryString()

	httpReq, err := http.NewRequest(req.Method(), reqUrl.String(), bytes.NewReader(req.BodyRaw()))
	if err != nil {
		return nil, err
	}
	return c.client.Do(httpReq)
}

func handleResponse(options []models.Response, resp *http.Response) (models.Response, error) {
	defer resp.Body.Close()
	for _, opt := range options {
		if resp.StatusCode == opt.Code() {
			return decodeResponse(resp, opt), nil
		}
	}
	return nil, fmt.Errorf("got unexpected status code: %d", resp.StatusCode)
}

func decodeResponse(resp *http.Response, res models.Response) models.Response {
	err := json.NewDecoder(resp.Body).Decode(res)
	if err != nil {
		panic(fmt.Errorf("failed to decode json body: %w", err))
	}
	return res
}

func decodeResponseT[T models.Response](resp *http.Response) T {
	var res T
	err := json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		panic(fmt.Errorf("failed to decode json body: %w", err))
	}
	return res
}
