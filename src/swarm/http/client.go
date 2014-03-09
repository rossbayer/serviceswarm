package http

import (
	gohttp "net/http"
	"net/http/cookiejar"
	"swarm"
	"fmt"
	"time"
	"io"
	"net/url"
	"strings"
)

type HTTPClient struct {
	client *gohttp.Client
	consumer *swarm.Consumer
}

func NewClient(consumer *swarm.Consumer) *HTTPClient {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		panic(fmt.Sprintf("Could not create jar for cookies due to error: %v",
			err))
	}

	return &HTTPClient{
		client: &gohttp.Client{Jar: cookieJar},
		consumer: consumer,
	}
}

func (c *HTTPClient) Do(req *gohttp.Request, taskName string) (resp *gohttp.Response, err error) {
	startTime := time.Now()
	resp, err = c.client.Do(req)
	duration := time.Since(startTime)

	result := c.resultFromResponse(resp, err, taskName, startTime, duration)
	c.consumer.Scenario.Results.Add(result)

	// TODO - Don't know if the wait should be handled here or if it should
	// handled by the task invoker (probably)

	return
}

func (c *HTTPClient) Get(url string, taskName string) (resp *gohttp.Response, err error) {
	req, err := gohttp.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req, taskName)
}

func (c *HTTPClient) Head(url string, taskName string) (resp *gohttp.Response, err error) {
	req, err := gohttp.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req, taskName)
}

func (c *HTTPClient) Post(url string, bodyType string, body io.Reader, taskName string) (resp *gohttp.Response, err error) {
	req, err := gohttp.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", bodyType)
	return c.Do(req, taskName)
}

func (c *HTTPClient) PostForm(url string, data url.Values, taskName string) (resp *gohttp.Response, err error) {
	return c.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()), taskName)
}

func (c *HTTPClient) resultFromResponse(resp *gohttp.Response,
	err error,
	taskName string,
	startTime time.Time,
	duration time.Duration) *swarm.Result {

	result := &swarm.Result{
		StartTime: startTime,
		Duration: duration,
		TaskName: taskName,
		Consumer: c.consumer,
	}

	switch {
	case err != nil:
		result.Status = swarm.Fatal
		result.Details = err
	case resp.StatusCode >= 400:
		result.Status = swarm.Error
		result.Details = resp
	default:
		result.Status = swarm.Success
		result.Details = resp
	}

	return result
}
