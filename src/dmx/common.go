package dmx

import "log"

type DMXBuffer struct {
	Frame []byte
}

func (x *DMXBuffer) init() {
	log.Println("init dmx")
	x.Frame = make([]byte, 512)
}
