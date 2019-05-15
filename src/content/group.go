package content

import "encoding/json"

type Slot interface {
	GetPreview() string
	GetURL() string

	GetName() string
	SetName(string)
	SetSize(uint64)
}

type SlotCommon struct {
	Slot
	Name string `json:"name"`
	Size uint64 `json:"size"`
	Mime string `json:"mime`
}

type Group struct {
	Slots map[int]Slot `json:"slots"`
}

type Image struct {
	SlotCommon
	URL string `json:"url"`
}

type Video struct {
	SlotCommon
	URL string `json:"url"`
}

type HTML struct {
	SlotCommon
	Body string `json:"body"`
}

func (me *SlotCommon) GetName() string     { return me.Name }
func (me *SlotCommon) SetSize(size uint64) { me.Size = size }
func (me *SlotCommon) SetName(name string) { me.Name = name }

func (me *Image) GetURL() string { return me.URL }

func (img *Image) MarshalJSON() (b []byte, err error) {
	return json.Marshal(&struct {
		Type string `json:"type"`
		Name string `json:"name"`
		Mime string `json:"mime"`
		Size uint64 `json:"size"`
		URL  string `json:"URL"`
	}{
		Type: "IMAGE",
		Name: img.Name,
		Mime: img.Mime,
		Size: img.Size,
		URL:  img.URL,
	})
}

func (vid *Video) MarshalJSON() (b []byte, err error) {
	return json.Marshal(&struct {
		Type string `json:"type"`
		Name string `json:"name"`
		Mime string `json:"mime"`
		Size uint64 `json:"size"`
		URL  string `json:"URL"`
	}{
		Type: "VIDEO",
		Name: vid.Name,
		Mime: vid.Mime,
		Size: vid.Size,
		URL:  vid.URL,
	})
}

func (html *HTML) MarshalJSON() (b []byte, err error) {
	return json.Marshal(&struct {
		Type string `json:"type"`
		Name string `json:"name"`
		Mime string `json:"mime"`
		Size uint64 `json:"size"`
		Body string `json:"body"`
	}{
		Type: "HTML",
		Name: html.Name,
		Mime: html.Mime,
		Size: html.Size,
		Body: html.Body,
	})
}

func (g *Group) UnmarshalJSON(b []byte) (err error) {
	var rawgroup struct {
		Slots map[int]*json.RawMessage
	}

	err = json.Unmarshal(b, &rawgroup)
	if err != nil {
		return
	}

	t := struct {
		Type string `json:"type"`
	}{}

	g.Slots = make(map[int]Slot)
	for k, v := range rawgroup.Slots {
		json.Unmarshal(*v, &t)
		switch t.Type {
		case "VIDEO":
			var vid Video
			err = json.Unmarshal(*v, &vid)
			g.Slots[k] = &vid
		case "IMAGE":
			var img Image
			err = json.Unmarshal(*v, &img)
			g.Slots[k] = &img
		case "HTML":
			var html HTML
			err = json.Unmarshal(*v, &html)
			g.Slots[k] = &html
		}
	}
	return
}
