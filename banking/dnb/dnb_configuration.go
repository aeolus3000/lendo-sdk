package dnb

import (
	"github.com/aeolus3000/lendo-sdk/banking"
	"time"
)

func NewDnbDefaultConfiguration() banking.Configuration {
	return banking.Configuration{
		RequestTimeout:           5 * time.Second,
		TcpConnectTimeout:        5 * time.Second,
		TlsHandshakeTimeout:      5 * time.Second,
		Host:                     "localhost",
		Port:                     "8000",
		CreateSlug:               "api/applications",
		CheckStatusSlug:          "api/jobs",
		CheckStatusParameterName: "application_id",
		ContentType:              "application/json",
	}
}

