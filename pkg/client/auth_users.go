package client

import (
	"encoding/json"
	"fmt"
)

// ListAuthUsers retrieves all embedded auth users using Controller REST API.
func (clt *Client) ListAuthUsers() ([]AuthUserResponse, error) {
	body, err := clt.doRequest("GET", "/users", nil)
	if err != nil {
		return nil, err
	}

	var users []AuthUserResponse
	if err = json.Unmarshal(body, &users); err != nil {
		return nil, err
	}
	return users, nil
}

// CreateAuthUser creates an embedded auth user using Controller REST API.
func (clt *Client) CreateAuthUser(request AuthUserCreateRequest) (AuthUserResponse, error) {
	var response AuthUserResponse
	body, err := clt.doRequest("POST", "/users", request)
	if err != nil {
		return response, err
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return response, err
	}
	return response, nil
}

// GetAuthUser retrieves an embedded auth user by ID using Controller REST API.
func (clt *Client) GetAuthUser(id string) (AuthUserResponse, error) {
	var response AuthUserResponse
	body, err := clt.doRequest("GET", fmt.Sprintf("/users/%s", id), nil)
	if err != nil {
		return response, err
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return response, err
	}
	return response, nil
}

// UpdateAuthUser updates an embedded auth user using Controller REST API.
func (clt *Client) UpdateAuthUser(id string, request AuthUserUpdateRequest) (AuthUserResponse, error) {
	var response AuthUserResponse
	body, err := clt.doRequest("PATCH", fmt.Sprintf("/users/%s", id), request)
	if err != nil {
		return response, err
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return response, err
	}
	return response, nil
}

// DeleteAuthUser soft-deletes an embedded auth user using Controller REST API.
func (clt *Client) DeleteAuthUser(id string) error {
	_, err := clt.doRequest("DELETE", fmt.Sprintf("/users/%s", id), nil)
	return err
}

// ResetAuthUserPassword resets an embedded auth user's password to a temporary value.
func (clt *Client) ResetAuthUserPassword(id string) (AuthUserResetPasswordResponse, error) {
	var response AuthUserResetPasswordResponse
	body, err := clt.doRequest("POST", fmt.Sprintf("/users/%s/reset-password", id), nil)
	if err != nil {
		return response, err
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return response, err
	}
	return response, nil
}

// ResetAuthUserToken issues a one-time password reset token for an embedded auth user.
func (clt *Client) ResetAuthUserToken(id string) (AuthUserResetTokenResponse, error) {
	var response AuthUserResetTokenResponse
	body, err := clt.doRequest("POST", fmt.Sprintf("/users/%s/reset-token", id), nil)
	if err != nil {
		return response, err
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return response, err
	}
	return response, nil
}
