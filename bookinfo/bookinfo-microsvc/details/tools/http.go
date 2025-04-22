package tools

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
)

type Client struct {
	httpClient *http.Client
}

func NewClient(timeout time.Duration) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}
func (c *Client) Get(url string, header http.Header, callSvcName string, parentSpan opentracing.Span) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating GET request %v", err)
		return nil, err
	}
	req.Header = header
	span := opentracing.StartSpan("call_"+callSvcName, opentracing.ChildOf(parentSpan.Context()))
	defer span.Finish()
	injectErr := opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	if injectErr != nil {
		log.Printf("Could not inject span context info to header failed %v", injectErr)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("Error sending GET request %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading GET response %v", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("GET request failed with status: %d", resp.StatusCode)
		return nil, fmt.Errorf("GET request failed with status: %d", resp.StatusCode)
	}
	log.Printf("GET response Boyd: %s", string(body))
	return body, nil
}
