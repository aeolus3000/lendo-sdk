package banking

import "time"

type Configuration struct {
	RequestTimeout           time.Duration `default:"10s"`
	TcpConnectTimeout        time.Duration `default:"10s"`
	TlsHandshakeTimeout      time.Duration `default:"10s"`
	Host                     string        `required:"true"`
	Port                     string        `required:"true"`
	CreateSlug               string        `required:"true"`
	CheckStatusSlug          string        `required:"true"`
	CheckStatusParameterName string        `required:"true"`
	ContentType              string        `required:"true"`
}
