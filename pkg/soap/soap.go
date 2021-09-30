package soap

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
)

type Client struct{}

func (soap *Client) Request(url string, body interface{}) (*Response, error) {
	request, err := soap.generateRequest(url, body)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil, err
	}

	return soap.do(request)
}

func (soap *Client) generateRequest(url string, body interface{}) (*http.Request, error) {
	t := template.Must(template.ParseFiles("pkg/soap/request.template.xml"))

	data := &bytes.Buffer{}
	if err := t.Execute(data, body); err != nil {
		fmt.Printf("template.Execute error. %s\n", err.Error())
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, url, data)
	request.Header.Set("Content-Type", "text/xml")
	if err != nil {
		fmt.Printf("Error making a request. %s\n", err.Error())
		return nil, err
	}

	return request, nil
}

func (soap *Client) do(request *http.Request) (*Response, error) {
	httpClient := new(http.Client)
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}

	sr := &Response{}
	if err = xml.Unmarshal(body, sr); err != nil {
		return nil, err
	}

	if sr.SoapBody != nil && sr.SoapBody.FaultDetails != nil {
		return nil, fmt.Errorf("%v:%v", sr.SoapBody.FaultDetails.Faultcode, sr.SoapBody.FaultDetails.Faultstring)
	}

	return sr, nil
}
