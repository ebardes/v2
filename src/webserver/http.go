package webserver

import (
	"fmt"
	"net"
	"net/http"
	"v2/config"

	"github.com/rs/zerolog/log"
)

// FN is a generic handler type
type FN func(http.ResponseWriter, *http.Request, *config.Config) error

var cfg *config.Config

// Register connections an HTTP path to a Function
func Register(pattern string, fn FN) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r, cfg)
		if err != nil {
			log.Error().Err(err).Msg("Server Error")
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("%v", err)))
		}
	})
}

type static struct {
}

var fs http.Handler

func (s *static) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/" {
		w.Header().Add("Location", "index.go")
		w.WriteHeader(http.StatusFound)
	} else {
		fs.ServeHTTP(w, r)
	}
}

// Run launches the webserver and runs "forever"
func Run(config *config.Config) {
	cfg = config
	fs = http.FileServer(http.Dir(cfg.Static))
	s := static{}
	http.Handle("/", &s)

	enumInterfaces()
	addr := fmt.Sprintf(":%d", cfg.WebPort)
	log.Info().Msg("Listening to " + addr)
	http.ListenAndServe(addr, nil)
}

func enumInterfaces() (err error) {
	ifs, err := net.Interfaces()
	if err != nil {
		return
	}

	for _, iface := range ifs {
		addrs, _ := iface.Addrs()
		if len(addrs) == 0 {
			continue
		}

		for _, addr := range addrs {
			switch addr.(type) {
			case *net.IPNet:
				ipnet := addr.(*net.IPNet)
				if ipnet.IP.IsLoopback() {
					continue
				}
				ip4 := ipnet.IP.To4()
				if ip4 == nil {
					continue // Skip IPv6 addresses
				}

				log.Info().Str("interface", iface.Name).Str("URL", fmt.Sprintf("http://%s:%d/", ip4, cfg.WebPort)).Msg("Connect")
			}

		}
	}
	return
}
