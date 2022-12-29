// author by lipengfei5 @2022-03-24

package httpx

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/net/http2"
)

func Transport(config Config) http.RoundTripper {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(config.DialTimeout) * time.Millisecond,
			KeepAlive: time.Duration(config.DialKeepAlive) * time.Millisecond,
		}).DialContext,
		MaxIdleConns:        config.MaxIdleConnection,
		MaxIdleConnsPerHost: config.MaxIdleConnsPerHost,
		IdleConnTimeout:     time.Duration(config.IdleConnTimeout) * time.Millisecond,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: config.InsecureSkipVerify},
	}

	http2.ConfigureTransport(transport)
	return otelhttp.NewTransport(transport)
}
