// Copyright (c) BunnyWay d.o.o.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type DnsZone struct {
	Id                            int64       `json:"Id,omitempty"`
	Domain                        string      `json:"Domain"`
	CustomNameserversEnabled      bool        `json:"CustomNameserversEnabled"`
	Nameserver1                   string      `json:"Nameserver1"`
	Nameserver2                   string      `json:"Nameserver2"`
	SoaEmail                      string      `json:"SoaEmail"`
	LoggingEnabled                bool        `json:"LoggingEnabled"`
	LoggingIPAnonymizationEnabled bool        `json:"LoggingIPAnonymizationEnabled"`
	LogAnonymizationType          uint8       `json:"LogAnonymizationType"`
	Records                       []DnsRecord `json:"Records"`
}

func (c *Client) GetDnsZone(id int64) (DnsZone, error) {
	var data DnsZone
	resp, err := c.doRequest(http.MethodGet, fmt.Sprintf("%s/dnszone/%d", c.apiUrl, id), nil)
	if err != nil {
		return data, err
	}

	if resp.StatusCode != http.StatusOK {
		return data, errors.New(resp.Status)
	}

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}

	_ = resp.Body.Close()
	err = json.Unmarshal(bodyResp, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (c *Client) CreateDnsZone(data DnsZone) (DnsZone, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return DnsZone{}, err
	}

	resp, err := c.doRequest(http.MethodPost, fmt.Sprintf("%s/dnszone", c.apiUrl), bytes.NewReader(body))
	if err != nil {
		return DnsZone{}, err
	}

	if resp.StatusCode != http.StatusCreated {
		bodyResp, err := io.ReadAll(resp.Body)
		if err != nil {
			return DnsZone{}, err
		}

		_ = resp.Body.Close()
		var obj struct {
			Message string `json:"Message"`
		}

		err = json.Unmarshal(bodyResp, &obj)
		if err != nil {
			return DnsZone{}, err
		}

		return DnsZone{}, errors.New(obj.Message)
	}

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return DnsZone{}, err
	}
	_ = resp.Body.Close()

	dataApiResult := DnsZone{}
	err = json.Unmarshal(bodyResp, &dataApiResult)
	if err != nil {
		return dataApiResult, err
	}

	return dataApiResult, nil
}

func (c *Client) UpdateDnsZone(dataApi DnsZone) (DnsZone, error) {
	id := dataApi.Id

	body, err := json.Marshal(dataApi)
	if err != nil {
		return DnsZone{}, err
	}

	resp, err := c.doRequest(http.MethodPost, fmt.Sprintf("%s/dnszone/%d", c.apiUrl, id), bytes.NewReader(body))
	if err != nil {
		return DnsZone{}, err
	}

	if resp.StatusCode != http.StatusOK {
		bodyResp, err := io.ReadAll(resp.Body)
		if err != nil {
			return DnsZone{}, err
		}

		_ = resp.Body.Close()
		var obj struct {
			Message string `json:"Message"`
		}

		err = json.Unmarshal(bodyResp, &obj)
		if err != nil {
			return DnsZone{}, err
		}

		return DnsZone{}, errors.New(obj.Message)
	}

	dataApiResult, err := c.GetDnsZone(id)
	if err != nil {
		return dataApiResult, err
	}

	return dataApiResult, nil
}

func (c *Client) DeleteDnsZone(id int64) error {
	resp, err := c.doRequest(http.MethodDelete, fmt.Sprintf("%s/dnszone/%d", c.apiUrl, id), nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return errors.New(resp.Status)
	}

	return nil
}
