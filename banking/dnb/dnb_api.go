package dnb

import (
	"bytes"
	"fmt"
	"github.com/aeolus3000/lendo-sdk/banking"
	"github.com/aeolus3000/lendo-sdk/utility"
	"google.golang.org/protobuf/encoding/protojson"
	"io"
	"net"
	"net/http"
	"strings"
)

type Dnb struct {
	client *http.Client
	config banking.Configuration
}

func NewDnbBanking(configuration banking.Configuration) banking.BankingApi {
	var netTransport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: configuration.TcpConnectTimeout,
		}).DialContext,
		TLSHandshakeTimeout: configuration.TlsHandshakeTimeout,
	}
	var netClient = &http.Client{
		Timeout:   configuration.RequestTimeout,
		Transport: netTransport,
	}
	return Dnb{netClient, configuration}
}

func (d Dnb) Create(application *banking.Application) (*banking.Application, error) {
	request := translateToDnbCreateRequest(application)
	json, err := utility.MarshalToJson(&request)
	if err != nil {
		return nil, err
	}
	response, err := d.client.Post(d.createUrl(), d.config.ContentType, bytes.NewBuffer(json))
	if err != nil {
		return nil, err
	}
	applicationResponse, err := readApplicationsResponse(response)
	if err != nil {
		return nil, err
	}
	return translateFromDnbCreateResponse(applicationResponse), nil
}

func (d Dnb) CheckStatus(applicationId string) (*banking.Application, error) {
	response, err := d.client.Get(d.checkStatusUrl(applicationId))
	if err != nil {
		return nil, err
	}
	jobsResponse, err := readJobsResponse(response)
	if err != nil {
		return nil, err
	}
	return translateFromDnbCheckStatusResponse(jobsResponse), nil
}

func (d Dnb) createUrl() string {
	return fmt.Sprintf("%s/%s", d.createEndpoint(), d.config.CreateSlug)
}

func (d Dnb) checkStatusUrl(applicationId string) string {
	return fmt.Sprintf("%s/%s?%s=%s", d.createEndpoint(),
		d.config.CheckStatusSlug, d.config.CheckStatusParameterName, applicationId)
}

func (d Dnb) createEndpoint() string {
	return fmt.Sprintf("http://%s:%s", d.config.Host, d.config.Port)
}

func translateToDnbCreateRequest(application *banking.Application) DnbApplicationsRequest {
	return DnbApplicationsRequest{
		Id:        application.Id,
		FirstName: application.FirstName,
		LastName:  application.LastName,
	}
}

func translateFromDnbCreateResponse(application *DnbApplicationsResponse) *banking.Application {
	return &banking.Application{
		Id:        application.Id,
		FirstName: application.FirstName,
		LastName:  application.LastName,
		Status:    application.Status,
	}
}

func translateFromDnbCheckStatusResponse(jobsResponse *DnbJobsResponse) *banking.Application {
	return &banking.Application{
		Id:     jobsResponse.ApplicationId,
		Status: jobsResponse.Status,
		JobId:  jobsResponse.Id,
	}
}

func readApplicationsResponse(response *http.Response) (*DnbApplicationsResponse, error) {
	defer response.Body.Close()
	body, err := ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("cant read response body; %v", err)
	}
	switch response.StatusCode {
	case 201:
		application := DnbApplicationsResponse{}
		if err := ConvertToApplicationsResponse(body, &application); err != nil {
			return nil, err
		}
		return &application, nil
	case 400:
		return nil, fmt.Errorf("duplicate application Id: %s", string(body))
	default:
		return nil, fmt.Errorf("unknown return code %d", response.StatusCode)
	}
}

func readJobsResponse(response *http.Response) (*DnbJobsResponse, error) {
	defer response.Body.Close()
	body, err := ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("cant read response body; %v", err)
	}
	switch response.StatusCode {
	case 200:
		application := DnbJobsResponse{}
		if err := ConvertToJobsResponse(body, &application); err != nil {
			return nil, err
		}
		return &application, nil
	case 400:
		return nil, fmt.Errorf("invalid uuid: %s", string(body))
	case 404:
		return nil, fmt.Errorf("uuid does not exist: %s", string(body))
	default:
		return nil, fmt.Errorf("unknown return code %c", response.StatusCode)
	}
}

func ConvertToApplicationsResponse(jsonBytes []byte, response *DnbApplicationsResponse) error {
	err := protojson.Unmarshal(jsonBytes, response)
	response.Status = strings.ToLower(response.Status)
	if err != nil {
		return fmt.Errorf("wrong response body format: %v; body: %s", err, string(jsonBytes))
	}
	return nil
}

func ConvertToJobsResponse(jsonBytes []byte, response *DnbJobsResponse) error {
	err := protojson.Unmarshal(jsonBytes, response)
	response.Status = strings.ToLower(response.Status)
	if err != nil {
		return fmt.Errorf("wrong response body format: %v; body: %s", err, string(jsonBytes))
	}
	return nil
}

func ReadAll(r io.Reader) ([]byte, error) {
	b := make([]byte, 0, 512)
	for {
		if len(b) == cap(b) {
			// Add more capacity (let append pick how much).
			b = append(b, 0)[:len(b)]
		}
		n, err := r.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return b, err
		}
	}
}
