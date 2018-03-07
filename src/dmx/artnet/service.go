package artnet

import (
	"config"
	"dmx"
)

type ArtNet struct {
	dmx.DMXBuffer
}

func NewService(*config.Config) (*ArtNet, error) {
	return nil, nil
}

func (x *ArtNet) Run() {

}

func (x *ArtNet) Stop() {

}
