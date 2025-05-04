package api_client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"url_shortener/internal/custom_errors"
)

type APIClient interface {
	Get(ctx context.Context, url string) ([]byte, error)
}

type FakeAPIClient struct {
	Result []byte
	Err    error
}

func (f FakeAPIClient) Get(ctx context.Context, url string) ([]byte, error) {
	return f.Result, f.Err
}

type apiClient struct {
	client  http.Client
	timeout time.Duration
}

type PeopleList struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
}

func (a *apiClient) GetPeopleList(ctx context.Context, url string) (PeopleList, error) {
	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return PeopleList{}, fmt.Errorf(custom_errors.CreateAPIRequestError, url, err)
	}
	response, err := a.client.Do(request)
	if err != nil {
		return PeopleList{}, fmt.Errorf(custom_errors.MakingAPIRequestError, url, err)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return PeopleList{}, fmt.Errorf(custom_errors.ReadBodyError, url, err)
	}

	// NOTE: parse body as interface
	// var data map[string]interface{}
	// err = json.Unmarshal(body, &data)
	// fmt.Printf("Results: %v\n", data)
	// fmt.Println("RESULT", data["count"])
	// if count, ok := data["count"].(float64); ok {
	// 	fmt.Println("FIND count", int(count))
	// }

	var peopleList PeopleList
	err = json.Unmarshal(body, &peopleList)
	if err != nil {
		return PeopleList{}, fmt.Errorf(custom_errors.MakingAPIRequestError, url, err)
	}
	return peopleList, nil
}

func (a *apiClient) Get(ctx context.Context, url string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf(custom_errors.CreateAPIRequestError, url, err)
	}
	response, err := a.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf(custom_errors.MakingAPIRequestError, url, err)
	}
	defer response.Body.Close() // nolint
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(custom_errors.WrongResponseStatusError, response.StatusCode, url)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf(custom_errors.ReadBodyError, url, err)
	}
	return body, nil
}

func New(timeout time.Duration) APIClient {
	return &apiClient{
		client:  http.Client{Transport: &http.Transport{DisableKeepAlives: false}},
		timeout: timeout,
	}
}
