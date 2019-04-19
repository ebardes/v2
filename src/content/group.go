package content

import "encoding/json"

type Slot interface {
	GetPreview() string
}

type Group struct {
	Slots map[int]Slot `json:"slots"`
}

type Image struct {
	Slot
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Video struct {
	Slot
	Name string `json:"name"`
	URL  string `json:"url"`
}

type HTML struct {
	Slot
	Name string `json:"name"`
	Body string `json:"body"`
}

func (img *Image) MarshalJSON() (b []byte, err error) {
	return json.Marshal(&struct {
		Type string `json:"type"`
		Name string `json:"name"`
		URL  string `json:"URL"`
	}{
		Type: "IMAGE",
		Name: img.Name,
		URL:  img.URL,
	})
}

func (vid *Video) MarshalJSON() (b []byte, err error) {
	return json.Marshal(&struct {
		Type string `json:"type"`
		Name string `json:"name"`
		URL  string `json:"URL"`
	}{
		Type: "VIDEO",
		Name: vid.Name,
		URL:  vid.URL,
	})
}

func (html *HTML) MarshalJSON() (b []byte, err error) {
	return json.Marshal(&struct {
		Type string `json:"type"`
		Name string `json:"name"`
		Body string `json:"body"`
	}{
		Type: "HTML",
		Name: html.Name,
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
