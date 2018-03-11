package view

import (
	"fmt"
	"os"
	"testing"
)

func TestSystem(t *testing.T) {
	v := InitView(os.Stdout)
	v.Render("index.tpl", nil)

	// fmt.Println(AllTemplates.Tree)
	for _, t := range AllTemplates.Templates() {
		fmt.Println(t.Name())
	}
}
