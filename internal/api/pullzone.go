package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Pullzone struct {
	Id   int64  `json:"Id,omitempty"`
	Name string `json:"Name,omitempty"`

	// sub-resources
	Edgerules        []PullzoneEdgerule       `json:"Edgerules"`
	Hostnames        []PullzoneHostname       `json:"Hostnames"`
	OptimizerClasses []PullzoneOptimizerClass `json:"OptimizerClasses"`

	// origin
	OriginType      uint8  `json:"OriginType"`
	OriginUrl       string `json:"OriginUrl,omitempty"`
	StorageZoneId   int64  `json:"StorageZoneId,omitempty"`
	AddHostHeader   bool   `json:"AddHostHeader"`
	VerifyOriginSSL bool   `json:"VerifyOriginSSL"`
	FollowRedirects bool   `json:"FollowRedirects"`

	// routing
	Type              uint8    `json:"Type"`
	EnableGeoZoneAF   bool     `json:"EnableGeoZoneAF"`
	EnableGeoZoneASIA bool     `json:"EnableGeoZoneASIA"`
	EnableGeoZoneEU   bool     `json:"EnableGeoZoneEU"`
	EnableGeoZoneSA   bool     `json:"EnableGeoZoneSA"`
	EnableGeoZoneUS   bool     `json:"EnableGeoZoneUS"`
	RoutingFilters    []string `json:"RoutingFilters"`

	// optimizer
	OptimizerEnabled                      bool    `json:"OptimizerEnabled"`
	OptimizerMinifyCss                    bool    `json:"OptimizerMinifyCSS"`
	OptimizerMinifyJs                     bool    `json:"OptimizerMinifyJavaScript"`
	OptimizerWebp                         bool    `json:"OptimizerEnableWebP"`
	OptimizerForceClasses                 bool    `json:"OptimizerForceClasses"`
	OptimizerImageOptimization            bool    `json:"OptimizerEnableManipulationEngine"`
	OptimizerAutomaticOptimizationEnabled bool    `json:"OptimizerAutomaticOptimizationEnabled"`
	OptimizerDesktopMaxWidth              uint64  `json:"OptimizerDesktopMaxWidth"`
	OptimizerMobileMaxWidth               uint64  `json:"OptimizerMobileMaxWidth"`
	OptimizerImageQuality                 uint8   `json:"OptimizerImageQuality"`
	OptimizerMobileImageQuality           uint8   `json:"OptimizerMobileImageQuality"`
	OptimizerWatermarkEnabled             bool    `json:"OptimizerWatermarkEnabled"`
	OptimizerWatermarkUrl                 string  `json:"OptimizerWatermarkUrl"`
	OptimizerWatermarkPosition            uint8   `json:"OptimizerWatermarkPosition"`
	OptimizerWatermarkOffset              float64 `json:"OptimizerWatermarkOffset"`
	OptimizerWatermarkMinImageSize        uint64  `json:"OptimizerWatermarkMinImageSize"`
}

func (c *Client) GetPullzone(id int64) (Pullzone, error) {
	var data Pullzone
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/pullzone/%d", c.apiUrl, id), nil)
	if err != nil {
		return data, err
	}

	req.Header.Add("AccessKey", c.apiKey)

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
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

func (c *Client) CreatePullzone(data Pullzone) (Pullzone, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return Pullzone{}, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/pullzone", c.apiUrl), bytes.NewReader(body))
	if err != nil {
		return Pullzone{}, err
	}

	req.Header.Add("AccessKey", c.apiKey)
	req.Header.Add("Content-Type", "application/json")

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return Pullzone{}, err
	}

	if resp.StatusCode != http.StatusCreated {
		bodyResp, err := io.ReadAll(resp.Body)
		if err != nil {
			return Pullzone{}, err
		}

		_ = resp.Body.Close()
		var obj struct {
			Message string `json:"Message"`
		}

		err = json.Unmarshal(bodyResp, &obj)
		if err != nil {
			return Pullzone{}, err
		}

		return Pullzone{}, errors.New(obj.Message)
	}

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return Pullzone{}, err
	}
	_ = resp.Body.Close()

	dataApiResult := Pullzone{}
	err = json.Unmarshal(bodyResp, &dataApiResult)
	if err != nil {
		return dataApiResult, err
	}

	return dataApiResult, nil
}

func (c *Client) UpdatePullzone(dataApi Pullzone) (Pullzone, error) {
	id := dataApi.Id

	body, err := json.Marshal(dataApi)
	if err != nil {
		return Pullzone{}, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/pullzone/%d", c.apiUrl, id), bytes.NewReader(body))
	if err != nil {
		return Pullzone{}, err
	}

	req.Header.Add("AccessKey", c.apiKey)
	req.Header.Add("Content-Type", "application/json")

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return Pullzone{}, err
	}

	if resp.StatusCode != http.StatusOK {
		bodyResp, err := io.ReadAll(resp.Body)
		if err != nil {
			return Pullzone{}, err
		}

		_ = resp.Body.Close()
		var obj struct {
			Message string `json:"Message"`
		}

		err = json.Unmarshal(bodyResp, &obj)
		if err != nil {
			return Pullzone{}, err
		}

		return Pullzone{}, errors.New(obj.Message)
	}

	dataApiResult, err := c.GetPullzone(id)
	if err != nil {
		return dataApiResult, err
	}

	return dataApiResult, nil
}

func (c *Client) DeletePullzone(id int64) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/pullzone/%d", c.apiUrl, id), nil)
	if err != nil {
		return err
	}

	req.Header.Add("AccessKey", c.apiKey)
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return errors.New(resp.Status)
	}

	return nil
}
