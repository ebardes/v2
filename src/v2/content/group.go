package content

type Slot struct {
	Type string `json:"type"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Group struct {
	SlotID int          `json:"slot_id"`
	Slots  map[int]Slot `json:"slots"`
}
