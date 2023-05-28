package ipresolver

import (
	"github.com/netinternet/remoteaddr"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

const (
	name = "ipresolver"
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

func (p *Plugin) Init(logger Logger, cfg Configurer) error {
	// construct a named logger for the middleware
	p.log = logger.NamedLogger(name)
	p.cfg = cfg
	return nil
}

func (p *Plugin) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		originalAddr := r.RemoteAddr
		ip, port := remoteaddr.Parse().IP(r)

		if strings.ContainsRune(originalAddr, ':') && (port == "" || port == "-1") {
			parts := strings.Split(originalAddr, ":")
			port = parts[1]
		}

		if port != "" && port != "-1" {
			ip += ":" + port
		}

		p.log.Debug("rewrite address from " + originalAddr + " to " + ip)

		r.RemoteAddr = ip
		next.ServeHTTP(w, r)
	})
}

func (p *Plugin) Name() string {
	return name
}
