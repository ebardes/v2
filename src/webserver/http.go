package webserver

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"v2/config"
	"v2/content"

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
	FS http.Handler
}

func (s *static) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/" {
		w.Header().Add("Location", "index.go")
		w.WriteHeader(http.StatusFound)
	} else {
		s.FS.ServeHTTP(w, r)
	}
}

// Run launches the webserver and runs "forever"
func Run(config *config.Config) {
	cfg = config
	http.Handle("/", &static{
		http.FileServer(http.Dir(cfg.StaticDir)),
	})

	http.Handle("/post/", &wrap{f: post})
	http.Handle("/media/", &wrap{f: media})
	http.Handle("/delete/", &wrap{f: remove})

	enumInterfaces(config)
	addr := fmt.Sprintf(":%d", cfg.WebPort)
	log.Info().Msg("Listening to " + addr)
	log.Error().Err(http.ListenAndServe(addr, nil)).Msg("Error setting up HTTP server")
}

func enumInterfaces(cfg *config.Config) (err error) {
	ifs, err := net.Interfaces()
	if err != nil {
		return
	}

	cfg.Networks = nil
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

				url := fmt.Sprintf("http://%s:%d/", ip4, cfg.WebPort)
				log.Info().Str("interface", iface.Name).Str("URL", url).Msg("Connect")

				n := config.Network{
					Name:      iface.Name,
					IPAddress: ip4.String(),
				}
				cfg.Networks = append(cfg.Networks, n)
			}

		}
	}
	return
}

type wrap struct {
	f func(http.ResponseWriter, *http.Request) error
}

func (x *wrap) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := x.f(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%v", err)))
	}
}

func post(w http.ResponseWriter, r *http.Request) (err error) {
	cfg := config.GlobalConfig

	uri := r.RequestURI
	n := strings.LastIndex(uri, "/")
	if n < 0 {
		return fmt.Errorf("Invalid URI - must include numeric group")
	}
	uri = uri[n+1:]
	group, _ := strconv.Atoi(uri)
	if group <= 0 {
		return fmt.Errorf("Invalid URI - must include numeric group")
	}

	f, fh, err := r.FormFile("file")
	if err != nil {
		return
	}
	defer f.Close()
	log.Info().Str("filename", fh.Filename).Int("group", group).Msg("Upload")

	targetdir := path.Join(cfg.ContentDir, uri)
	os.MkdirAll(targetdir, 0777)
	dest := path.Join(targetdir, fh.Filename)

	fio, err := os.Create(dest)
	if err != nil {
		log.Error().Err(err).Str("destination", dest).Msg("Can not create file")
		return
	}
	defer fio.Close()

	_, err = io.Copy(fio, f)
	if err != nil {
		return
	}

	grp := cfg.Content[group]
	for _, slot := range grp.Slots {
		if slot.GetName() == fh.Filename {
			slot.SetSize(uint64(fh.Size))
			cfg.Save()
			return
		}
	}

	// new item adding
	for i := 1; i <= 255; i++ {
		_, ok := grp.Slots[i]
		if ok {
			continue
		}

		mime := fh.Header.Get("Content-Type")
		if strings.HasPrefix(mime, "image/") {
			slot := &content.Image{}
			slot.Name = fh.Filename
			slot.Size = uint64(fh.Size)
			slot.Mime = mime
			slot.URL = path.Join("media", uri, fh.Filename)

			grp.Slots[i] = slot
			cfg.Save()
			return
		}
	}
	return
}

func media(w http.ResponseWriter, r *http.Request) (err error) {
	cfg := config.GlobalConfig
	uri := r.URL.RequestURI()
	uri, err = url.QueryUnescape(uri)
	if err != nil {
		return
	}

	if strings.HasPrefix(uri, "/media/") {
		uri = uri[7:]
	}

	target := path.Join(cfg.ContentDir, uri)
	http.ServeFile(w, r, target)
	return
}

func remove(w http.ResponseWriter, r *http.Request) (err error) {
	cfg := config.GlobalConfig
	uri := r.URL.RequestURI()
	uri, err = url.QueryUnescape(uri)
	if err != nil {
		return
	}

	if strings.HasPrefix(uri, "/delete/") {
		uri = uri[8:]
	}

	var group, slot int

	_, err = fmt.Sscanf(uri, "%d/%d", &group, &slot)
	if err != nil {
		return
	}

	g, ok := cfg.Content[group]
	if !ok {
		return
	}

	s, ok := g.Slots[slot]

	target := path.Join(cfg.ContentDir, strconv.Itoa(group), s.GetName())
	os.Remove(target)
	delete(g.Slots, slot)
	cfg.Save()

	return
}
