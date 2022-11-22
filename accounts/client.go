package accounts

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/LeonhardtDavid/form3-client/models"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

type AccountClient interface {
	Create(ctx context.Context, account models.AccountData) (*models.AccountData, error)
	Fetch(ctx context.Context, accountId uuid.UUID) (*models.AccountData, error)
	Delete(ctx context.Context, accountId uuid.UUID, version int) error
}

type AccountClientOptions struct {
	BaseURL string
}

type accountHttpClient struct {
	options AccountClientOptions
	client  *http.Client
}

type accountBody struct {
	Data models.AccountData `json:"data"`
}

type errorBody struct {
	ErrorMessage string `json:"error_message"`
}

func (c *accountHttpClient) Create(ctx context.Context, account models.AccountData) (*models.AccountData, error) {
	url := c.getAccountsURL()

	body, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	return c.doRequestAndDecode(request, http.StatusCreated)
}

func (c *accountHttpClient) Fetch(ctx context.Context, accountId uuid.UUID) (*models.AccountData, error) {
	url := c.getAccountURL(accountId)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return c.doRequestAndDecode(request, http.StatusOK)
}

func (c *accountHttpClient) Delete(ctx context.Context, accountId uuid.UUID, version int) error {
	url := c.getAccountURL(accountId)

	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	queryParams := request.URL.Query()
	queryParams.Add("version", strconv.Itoa(version))
	request.URL.RawQuery = queryParams.Encode()

	if _, err := c.doRequest(request, http.StatusNoContent); err != nil {
		return err
	}

	return nil
}

func (c *accountHttpClient) getAccountsURL() string {
	return fmt.Sprintf("%s/v1/organisation/accounts", c.options.BaseURL)
}

func (c *accountHttpClient) getAccountURL(accountId uuid.UUID) string {
	return fmt.Sprintf("%s/%s", c.getAccountsURL(), accountId)
}

func (c *accountHttpClient) doRequest(request *http.Request, expectedStatus int) (*http.Response, error) {
	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != expectedStatus {
		return nil, handleErrorResponse(response)
	}

	return response, nil
}

func (c *accountHttpClient) doRequestAndDecode(request *http.Request, expectedStatus int) (*models.AccountData, error) {
	response, err := c.doRequest(request, expectedStatus)
	if err != nil {
		return nil, err
	}

	var accountResponse accountBody
	if err := json.NewDecoder(response.Body).Decode(&accountResponse); err != nil {
		return nil, err
	}

	return &accountResponse.Data, nil
}

func handleErrorResponse(response *http.Response) error {
	var errorResponse errorBody
	if err := json.NewDecoder(response.Body).Decode(&errorResponse); err != nil {
		return err
	}

	return errors.New(errorResponse.ErrorMessage)
}

func NewAccountClient(options AccountClientOptions) AccountClient {
	return &accountHttpClient{
		options: options,
		client:  http.DefaultClient,
	}
}
