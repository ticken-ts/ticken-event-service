package sync

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"ticken-event-service/security/auth"
)

type registerUserPayload struct {
	JWT string `json:"jwt"`
}

type ValidatorServiceHTTPClient struct {
	serviceUrl string
	authIssuer *auth.Issuer
}

func NewValidatorServiceHTTPClient(serviceUrl string, authIssuer *auth.Issuer) *ValidatorServiceHTTPClient {
	validatorServiceClient := &ValidatorServiceHTTPClient{
		serviceUrl: serviceUrl,
		authIssuer: authIssuer,
	}
	return validatorServiceClient
}

func (client *ValidatorServiceHTTPClient) RegisterValidator(organizationID uuid.UUID, validatorJWT string) error {
	validatorRegistrationURL, _ := url.JoinPath(
		client.serviceUrl, "api/v", "organizations", organizationID.String(), "validators")

	payload := registerUserPayload{JWT: validatorJWT}
	jsonPayload, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPost, validatorRegistrationURL, bytes.NewReader(jsonPayload))
	if err != nil {
		return err
	}

	rawJWT, err := client.authIssuer.IssueInService(auth.TickenValidatorService)
	if err != nil {
		return fmt.Errorf("failed to get client authentication: %s", err.Error())
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", rawJWT.Token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do validator service request: %s", err.Error())
	}
	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("validator creator failed: %s", readBody(res))
	}

	return nil
}

func (client *ValidatorServiceHTTPClient) SyncTickets(eventID uuid.UUID) error {
	validatorRegistrationURL, _ := url.JoinPath(
		client.serviceUrl, "api/v", "events", eventID.String(), "sync")

	req, err := http.NewRequest(http.MethodPost, validatorRegistrationURL, nil)
	if err != nil {
		return err
	}

	rawJWT, err := client.authIssuer.IssueInService(auth.TickenValidatorService)
	if err != nil {
		return fmt.Errorf("failed to get client authentication: %s", err.Error())
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", rawJWT.Token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to start ticket sync: %s", err.Error())
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to start ticket sync: %s", readBody(res))
	}

	return nil
}
