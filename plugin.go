package ip_resolver

import (
	"net/http"

	"go.uber.org/zap"
)

const (
	name = "ip_resolver"
)

type Configurer interface {
	// UnmarshalKey takes a single key and unmarshal it into a Struct.
	UnmarshalKey(name string, out any) error
	// Has checks if config section exists.
	Has(name string) bool
}

type Logger interface {
	NamedLogger(name string) *zap.Logger
}

type Plugin struct {
	log *zap.Logger
	cfg Configurer
}

func (p *Plugin) Init(cfg Configurer) error {
	// construct a named logger for the middleware
	p.cfg = cfg
	return nil
}

func (p *Plugin) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := "127.0.0.1:12345"
		p.log.Info("RewriteAddress from " + r.RemoteAddr + " to " + ip)

		r.RemoteAddr = ip
		next.ServeHTTP(w, r)
	})
}

func (p *Plugin) Name() string {
	return name
}
