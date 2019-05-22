package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path"
	"strconv"
	"v2/content"
)

const (
	// DefaultName is the name of the configuration file on the server
	DefaultName    = ".v2cfg.json"
	ProtocolSACN   = "sacn"
	ProtocolArtNet = "artnet"
)

type Layer struct {
	StartAddress uint   `json:"dmx_start"`
	Personality  string `json:"personality"`
}

// Display Describes a pane which is a single display
type Display struct {
	ID     uint    `json:"id"`
	Layers []Layer `json:"layers"`
}

type Network struct {
	Name      string
	IPAddress string
}

// Config The system configuration
type Config struct {
	DebugLevel  int                   `json:"debuglevel"`
	Universe    uint                  `json:"universe"`
	WebPort     uint                  `json:"port"`
	Interface   string                `json:"interface"`
	Protocol    string                `json:"protocol"`
	Displays    []Display             `json:"display"`
	ContentDir  string                `json:"content-directory"`
	Content     map[int]content.Group `json:"groups"`
	Networks    []Network             `json:"-"`
	Protocols   []string              `json:"-"`
	TemplateDir string                `json:"-"`
	StaticDir   string                `json:"-"`
}

var GlobalConfig Config

// Save Saves the config files in a default location
func (c *Config) Save() error {
	temp := getLocation() + ".tmp"

	c.Normalize()

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
func (c *Config) Load() (err error) {
	f, err := os.Open(getLocation())
	if err != nil {
		return err
	}
	defer f.Close()

	d := json.NewDecoder(f)
	return d.Decode(c)
}

// Normalize performs sanitiy checking of the state of the config
func (c *Config) Normalize() {
	for g, grp := range c.Content {
		if g == 0 {
			continue
		}
		r := make([]int, 0)
		for i, slt := range grp.Slots {
			target := path.Join(c.ContentDir, strconv.Itoa(g), slt.GetName())
			if _, err := os.Stat(target); err != nil {
				r = append(r, i)
			}
		}

		for _, x := range r {
			delete(grp.Slots, x)
		}
	}

	for i := 1; i <= 255; i++ {
		if group, ok := c.Content[i]; !ok {
			c.Content[i] = content.Group{
				Slots: make(map[int]content.Slot),
			}
			break
		} else {
			if len(group.Slots) == 0 {
				break
			}
		}
	}
}

func getLocation() string {
	u, _ := user.Current()
	return path.Join(u.HomeDir, DefaultName)
}

func (c *Config) GetDisplay(id uint) (d *Display, err error) {
	id--
	if id < 0 || id > uint(len(c.Displays)) {
		err = fmt.Errorf("Display ID out of range")
		return
	}
	d = &c.Displays[id]
	return
}
