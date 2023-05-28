package ip_resolver_middleware

import (
	"github.com/netinternet/remoteaddr"
	"net/http"

	"github.com/roadrunner-server/api/v2/plugins/config"
	"github.com/roadrunner-server/errors"
	"go.uber.org/zap"
)

const name = "ip_resolver"

type Plugin struct {
	log *zap.Logger
	cfg *Config
}

func (p *Plugin) Init(cfg config.Configurer, log *zap.Logger) error {
	// check if we need to init this middleware
	if !cfg.Has(name) {
		return errors.E(errors.Disabled)
	}

	// populate configuration
	p.cfg = &Config{}
	err := cfg.UnmarshalKey(name, p.cfg)
	if err != nil {
		return err
	}

	return nil
}

// Middleware is our actual http middleware
func (p *Plugin) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, port := remoteaddr.Parse().IP(r)
		if port != "" {
			ip += ":" + port
		}

		r.RemoteAddr = ip
		next.ServeHTTP(w, r)
	})
}

func (p *Plugin) Name() string {
	return name
}
