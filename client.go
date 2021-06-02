package arpio

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Version is set through ldflags.
var Version = "0.0"

// Commit is set through ldflags.
var Commit = "unknown"

// ArpioURL is the URL of the production Arpio service.
const ArpioURL = "https://api.arpio.io/api"

const userAgentPrefix = "arpio-client-go"

// Client contains Arpio API client state
type Client struct {
	APIUrl       string
	AccountID    string
	APIKeyID     string
	APIKeySecret string
	HTTPClient   *http.Client
}

// ErrorResponse contains the data the Arpio API returns with most error statuses.
type ErrorResponse struct {
	Message         string `json:"message"`
	AuthenticateURL string `json:"authenticateUrl"`
}

// NewClient creates a Client using the specified connection information.
func NewClient(apiURL, apiKeyID, apiKeySecret, accountID string) (*Client, error) {
	if apiURL == "" {
		return nil, fmt.Errorf("apiURL is required")
	}
	if apiKeyID == "" {
		return nil, fmt.Errorf("apiKeyID is required")
	}
	if apiKeySecret == "" {
		return nil, fmt.Errorf("apiKeySecret is required")
	}
	if accountID == "" {
		return nil, fmt.Errorf("accountID is required")
	}

	_, skipVerify := os.LookupEnv(ArpioTlsInsecureSkipVerify)

	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: skipVerify}

	client := Client{
		APIUrl:       apiURL,
		AccountID:    accountID,
		APIKeyID:     apiKeyID,
		APIKeySecret: apiKeySecret,
		HTTPClient: &http.Client{
			Transport: tr,
			Timeout:   60 * time.Second,
		},
	}
	return &client, nil
}

func (c Client) buildApiURL(relativeURL string) (u *url.URL, err error) {
	apiURL := strings.TrimSuffix(c.APIUrl, "/")
	relativeURL = strings.TrimPrefix(relativeURL, "/")
	return url.Parse(fmt.Sprintf("%s/%s", apiURL, relativeURL))
}

func (c Client) apiPost(relativeURL string, requestBody interface{}, responseBody interface{}) (status int, err error) {
	return c.doApiRequest("POST", relativeURL, requestBody, responseBody)
}

func (c Client) apiPut(relativeURL string, requestBody interface{}, responseBody interface{}) (status int, err error) {
	return c.doApiRequest("PUT", relativeURL, requestBody, responseBody)
}

func (c Client) apiGet(relativeURL string, responseBody interface{}) (status int, err error) {
	return c.doApiRequest("GET", relativeURL, nil, responseBody)
}

func (c Client) apiDelete(relativeURL string, responseBody interface{}) (status int, err error) {
	return c.doApiRequest("DELETE", relativeURL, nil, responseBody)
}

// Perform one Arpio API request using the specified HTTP method and optional
// request body (which will be marshaled to JSON).  If a non-nil response body
// is specified, the result will be unmarshaled from JSON and written to it.
//
// An error is included in the returned values if the response status
// code >= 400 (responseBody does not receive the response body when an error
// is returned).
func (c Client) doApiRequest(method, relativeURL string, requestBody interface{}, responseBody interface{}) (status int, err error) {
	u, err := c.buildApiURL(relativeURL)
	if err != nil {
		return status, err
	}

	userAgent := fmt.Sprintf("%s/%s/%s", userAgentPrefix, Version, Commit)
	apiKeyHeader := buildApiKeyHeader(c.APIKeyID, c.APIKeySecret)

	req := &http.Request{
		Method: method,
		URL:    u,
		Header: map[string][]string{
			"Accept":     {"*/*"},
			"User-Agent": {userAgent},
			"X-Api-Key":  {apiKeyHeader},
		},
	}
	if requestBody != nil {
		requestJson, err := json.Marshal(requestBody)
		if err != nil {
			return status, err
		}
		log.Printf("[TRACE] %s", requestJson)
		req.Body = ioutil.NopCloser(bytes.NewReader(requestJson))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Length", fmt.Sprintf("%d", len(requestJson)))
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return status, err
	}
	status = resp.StatusCode
	defer func(b io.ReadCloser) {
		err := b.Close()
		if err != nil {
			log.Printf("[INFO] Error closing response body: %s", err)
		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return status, err
	}

	// Arpio API errors come back in a standard format, as JSON.  Try to unmarshal
	// the response body as that error type in these cases so we can include those
	// details in the error we return.
	if status >= 400 {
		var errorResponse ErrorResponse
		err = json.Unmarshal(body, &errorResponse)

		var msg string
		if err != nil {
			log.Printf("[WARN] Error unmarshaling response body as ErrorResponse: %s", err)
			msg = fmt.Sprintf("Arpio API error: %s", body)
		} else {
			msg = errorResponse.Message
		}

		return status, fmt.Errorf(msg)
	}

	if responseBody != nil {
		err = json.Unmarshal(body, responseBody)
		if err != nil {
			return status, err
		}
	}
	return status, err
}

// buildApiKeyHeader builds the value to use for the X-Api-Key header.
// The format is the same as for HTTP "basic" authentication:
//
//   "X-Api-Key" := BASE64(apiKeyID + ":" + apiKeySecret)
//
// apiKeyID and apiKeySecret must be valid UTF-8 byte sequences.
// apiKeyID must not contain a colon character.
func buildApiKeyHeader(apiKeyID, apiKeySecret string) string {
	tuple := fmt.Sprintf("%s:%s", apiKeyID, apiKeySecret)
	return base64.StdEncoding.EncodeToString([]byte(tuple))
}
