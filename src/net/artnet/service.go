package artnet

import "config"

type ArtNet struct {
}

func NewService(*config.Config) (*ArtNet, error) {
	return nil, nil
}

func (x *ArtNet) Run() {

}

func (x *ArtNet) Stop() {

}
