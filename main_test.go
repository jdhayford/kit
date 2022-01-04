package main

import (
	"fmt"
	"os"
)

// \033[1A One Line Up?
// \033[2K Clear Line?
// \033[K ?

type UI struct {
	State string
	View
}

func (ui *UI) Render() {}

func (ui *UI) SetView(v View) {
	ui.View = v
}

func (ui *UI) Start() {
	fmt.Println("Starting")
	var b []byte = make([]byte, 1)
	for {
		os.Stdin.Read(b)
		fmt.Println(string(b))
		fmt.Println("I got the byte", b, "("+string(b)+")")
	}
}

func NewUI() UI {
	return UI{}
}

type View interface {
	Render() string
}

// func (v *View) Render() string {
// 	return ""
// }

type ListView struct {
	View
	items []string
	// selected string
}

func (lv ListView) Render() string {
	state := ""
	for item := range lv.items {
		state = state + fmt.Sprintf("\n%v", item)
	}
	return state
}

func NewListView() ListView {
	return ListView{items: []string{"hello", "world"}}
}

func meep() {
	lv := NewListView()
	ui := NewUI()
	ui.SetView(lv)
	ui.Start()
}
