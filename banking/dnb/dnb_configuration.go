package dnb

import "time"

type DnbConfiguration struct {
	requestTimeout time.Duration
	tcpConnectTimeout time.Duration
	tlsHandshakeTimeout time.Duration
	host string
	port string
	createSlug string
	checkStatusSlug string
	checkStatusParameterName string
	contentType string
}

func NewDnbDefaultConfiguration() DnbConfiguration {
	return DnbConfiguration{
		requestTimeout:           5 * time.Second,
		tcpConnectTimeout:        5 * time.Second,
		tlsHandshakeTimeout:      5 * time.Second,
		host:                     "localhost",
		port:                     "8000",
		createSlug:               "api/applications",
		checkStatusSlug:          "api/jobs",
		checkStatusParameterName: "application_id",
		contentType:              "application/json",
	}
}