package tao_ce

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"time"
)

type HealthcheckReport struct {
	AppsReports     map[string]int
	HealthcheckTime time.Time
}

type TAOCEService struct {
	BaseURL string
	Client  *http.Client
}

func New(baseURL string) *TAOCEService {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	return &TAOCEService{
		BaseURL: baseURL,
		Client:  client,
	}
}

func (s *TAOCEService) Healthcheck() (HealthcheckReport, error) {
	var gossReport map[string]interface{}
	var report map[string]int = make(map[string]int)

	response, err := s.Client.Get(s.BaseURL + "/health")
	if err != nil {
		// fmt.Fprintln(os.Stderr, "could not get healthcheck:", err)
		return HealthcheckReport{HealthcheckTime: time.Now()}, err
	}

	err = json.NewDecoder(response.Body).Decode(&gossReport)

	if err != nil {
		// fmt.Fprintln(os.Stderr, "could not decode healthcheck:", err)
		return HealthcheckReport{HealthcheckTime: time.Now()}, err
	}
	for _, result := range gossReport["results"].([]interface{}) {
		var app = result.(map[string]interface{})["meta"].(map[string]interface{})["app"].(string)
		if len(app) == 0 {
			continue
		}

		var status = int(result.(map[string]interface{})["result"].(float64))
		if _, ok := report[app]; ok {
			report[app] += status
		} else {
			report[app] = status
		}
	}

	return HealthcheckReport{AppsReports: report, HealthcheckTime: time.Now()}, nil
}
