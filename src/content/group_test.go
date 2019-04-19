package content

import (
	"encoding/json"
	"testing"
)

func TestJSONMarshal(t *testing.T) {
	g := Group{
		Slots: make(map[int]Slot),
	}
	img := Image{
		Name: "Test Bars",
		URL:  "TestBars.png",
	}
	g.Slots[1] = &img

	vid := Video{
		URL:  "Silly.mp4",
		Name: "Silly",
	}
	g.Slots[2] = &vid

	g.Slots[3] = &HTML{
		Name: "Test HTML",
		Body: `<html>
			<body>
			<h1>Hello World</h1>
			</body>
			</html>`,
	}

	bytes, err := json.MarshalIndent(g, "", " ")
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(string(bytes))
}

func TestJSONUnmarshal(t *testing.T) {
	blob := []byte(` {
		"slots": {
		 "1": {
		  "type": "IMAGE",
		  "name": "Test Bars",
		  "URL": "TestBars.png"
		 },
		 "2": {
		  "type": "VIDEO",
		  "name": "Silly",
		  "URL": "Silly.mp4"
		 },
		 "3": {
		  "type": "HTML",
		  "name": "Test HTML",
		  "body": "\u003chtml\u003e\n\t\t\t\u003cbody\u003e\n\t\t\t\u003ch1\u003eHello World\u003c/h1\u003e\n\t\t\t\u003c/body\u003e\n\t\t\t\u003c/html\u003e"
		 }
		}
	   }`)

	var g Group
	err := json.Unmarshal(blob, &g)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(g)
}
