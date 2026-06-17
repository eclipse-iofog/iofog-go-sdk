package client

import (
	"encoding/json"
	"fmt"
)

// GetCatalog retrieves all catalog items using Controller REST API
func (clt *Client) GetCatalog() (*CatalogListResponse, error) {
	body, err := clt.doRequest("GET", "/catalog/microservices", nil)
	if err != nil {
		return nil, err
	}

	response := new(CatalogListResponse)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	return response, nil
}

// GetCatalogItem retrieves one catalog item using Controller REST API
func (clt *Client) GetCatalogItem(id string) (*CatalogItemInfo, error) {
	body, err := clt.doRequest("GET", fmt.Sprintf("/catalog/microservices/%s", id), nil)
	if err != nil {
		return nil, err
	}

	response := new(CatalogItemInfo)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	return response, nil
}

// CreateCatalogItem creates one catalog item using Controller REST API
func (clt *Client) CreateCatalogItem(request *CatalogItemCreateRequest) (*CatalogItemInfo, error) {
	if request.RegistryID == 0 {
		request.RegistryID = 1
	}

	body, err := clt.doRequest("POST", "/catalog/microservices", request)
	if err != nil {
		return nil, err
	}
	response := &CatalogItemCreateResponse{}
	if err := json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	return clt.GetCatalogItem(response.ID)
}

// UpdateCatalogItem updates one catalog item using Controller REST API
func (clt *Client) UpdateCatalogItem(request *CatalogItemUpdateRequest) (*CatalogItemInfo, error) {
	_, err := clt.doRequest("PATCH", fmt.Sprintf("/catalog/microservices/%s", request.ID), request)
	if err != nil {
		return nil, err
	}
	return clt.GetCatalogItem(request.ID)
}

// DeleteCatalogItem deletes one catalog item using Controller REST API
func (clt *Client) DeleteCatalogItem(id string) error {
	_, err := clt.doRequest("DELETE", fmt.Sprintf("/catalog/microservices/%s", id), nil)
	return err
}

// GetCatalogItemByName returns a catalog item by listing all catalog items and returning the first occurrence of the specified name
func (clt *Client) GetCatalogItemByName(name string) (*CatalogItemInfo, error) {
	catalog, err := clt.GetCatalog()
	if err != nil {
		return nil, err
	}

	for _, item := range catalog.CatalogItems {
		if item.Name == name {
			return &item, nil
		}
	}

	return nil, NewNotFoundError(fmt.Sprintf("Could not find catalog item %s\n", name))
}
