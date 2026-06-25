package client

import (
	"encoding/json"
	"fmt"
	"time"
)

const defaultPlatformPollInterval = 2 * time.Second

// ReconcileAgent enqueues a manual fog platform reconcile retry (POST /iofog/{uuid}/reconcile).
// Resets Failed phase, clears backoff, and returns the fog UUID on success.
func (clt *Client) ReconcileAgent(uuid string) (CreateAgentResponse, error) {
	var response CreateAgentResponse
	if !clt.isLoggedIn() {
		return response, NewError("Controller client must be logged into perform Reconcile Agent request")
	}

	body, err := clt.doRequest("POST", fmt.Sprintf("/iofog/%s/reconcile", uuid), nil)
	if err != nil {
		return response, err
	}

	var respMap map[string]any
	if err = json.Unmarshal(body, &respMap); err != nil {
		return response, err
	}
	reconciledUUID, exists := respMap["uuid"].(string)
	if !exists {
		return response, NewInternalError("Failed to get Agent UUID from reconcile response")
	}

	response.UUID = reconciledUUID
	return response, nil
}

// ReconcileService enqueues a manual service platform reconcile retry (POST /services/{name}/reconcile).
// When provisioningStatus is failed, resets to pending and clears provisioningError before enqueue.
func (clt *Client) ReconcileService(name string) (*ServiceInfo, error) {
	if !clt.isLoggedIn() {
		return nil, NewError("Controller client must be logged into perform Reconcile Service request")
	}

	body, err := clt.doRequest("POST", fmt.Sprintf("/services/%s/reconcile", name), nil)
	if err != nil {
		return nil, err
	}

	service := new(ServiceInfo)
	if err = json.Unmarshal(body, service); err != nil {
		return nil, err
	}
	return service, nil
}

// WaitForAgentPlatformReady polls GET /iofog/{uuid} until platformStatus.phase is Ready.
// Returns an error on Failed phase, timeout, or when lastError is set on a non-Ready phase.
func (clt *Client) WaitForAgentPlatformReady(uuid string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		agent, err := clt.GetAgentByID(uuid)
		if err != nil {
			return err
		}
		if agent.PlatformStatus == nil {
			if time.Now().After(deadline) {
				return NewError(fmt.Sprintf("timed out waiting for agent %s platform status", uuid))
			}
			time.Sleep(defaultPlatformPollInterval)
			continue
		}

		switch agent.PlatformStatus.Phase {
		case PlatformReady:
			return nil
		case PlatformFailed:
			msg := fmt.Sprintf("agent %s platform reconcile failed", uuid)
			if agent.PlatformStatus.LastError != nil && *agent.PlatformStatus.LastError != "" {
				msg = fmt.Sprintf("%s: %s", msg, *agent.PlatformStatus.LastError)
			}
			return NewError(msg)
		case PlatformDeleting:
			return NewError(fmt.Sprintf("agent %s platform is deleting", uuid))
		}

		if time.Now().After(deadline) {
			return NewError(fmt.Sprintf("timed out waiting for agent %s platform ready (phase=%s)", uuid, agent.PlatformStatus.Phase))
		}
		time.Sleep(defaultPlatformPollInterval)
	}
}

// WaitForServiceProvisioningReady polls GET /services/{name} until provisioningStatus is ready.
func (clt *Client) WaitForServiceProvisioningReady(name string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		service, err := clt.GetService(name)
		if err != nil {
			return err
		}

		switch service.ProvisioningStatus {
		case ProvisioningReady:
			return nil
		case ProvisioningFailed:
			msg := fmt.Sprintf("service %s provisioning failed", name)
			if service.ProvisioningError != nil && *service.ProvisioningError != "" {
				msg = fmt.Sprintf("%s: %s", msg, *service.ProvisioningError)
			}
			return NewError(msg)
		}

		if time.Now().After(deadline) {
			return NewError(fmt.Sprintf("timed out waiting for service %s provisioning ready (status=%s)", name, service.ProvisioningStatus))
		}
		time.Sleep(defaultPlatformPollInterval)
	}
}
