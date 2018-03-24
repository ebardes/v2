package config

import (
	"encoding/json"
	"os"
	"os/user"
	"path"
	"v2/content"
)

const (
	// DefaultName is the name of the configuration file on the server
	DefaultName = ".v2cfg.json"
)

type Layer struct {
	StartAddress int    `json:"dmx_start"`
	Personality  string `json:"personality"`
}

// Pane Describes a pane which is a single display
type Display struct {
	ID     uint    `json:"id"`
	Layers []Layer `json:"layers"`
}

// Config The system configuration
type Config struct {
	DebugLevel int                   `json:"debuglevel"`
	Universe   uint                  `json:"universe"`
	WebPort    uint                  `json:"port"`
	Interface  string                `json:"interface"`
	Protocol   string                `json:"protocol"`
	Displays   []Display             `json:"display"`
	Content    map[int]content.Group `json:"groups"`
}

var cfg Config

func Get() *Config {
	return &cfg
}

// Save Saves the config files in a default location
func (c *Config) Save() error {
	temp := getLocation() + ".tmp"

	f, err := os.Create(temp)
	if err != nil {
		return err
	}
	defer f.Close()

	e := json.NewEncoder(f)
	e.SetIndent("", " ")
	err = e.Encode(c)
	if err != nil {
		return err
	}
	f.Close()

	os.Rename(temp, getLocation())
	return nil
}

// Load Reads the config file from the default location
func Load(cfg *Config) error {
	f, err := os.Open(getLocation())
	if err != nil {
		return err
	}
	defer f.Close()

	d := json.NewDecoder(f)
	d.Decode(cfg)
	return nil
}

func getLocation() string {
	u, _ := user.Current()
	return path.Join(u.HomeDir, DefaultName)
}
