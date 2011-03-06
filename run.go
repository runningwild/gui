package gui

import (
	"os"
	"fmt"
	"strings"
	"github.com/droundy/gui/websocket"
)

// This file implements the code to actually interface with
// websockets.

func RunSeparate(w func() Widget) os.Error {
	return websocket.RunSeparate("/", func() websocket.Handler {
		return &widgetwrapper{w(), []func(string){}}
	})
}

func Run(w Widget) os.Error {
	return websocket.Run("/", &widgetwrapper{w, []func(string){}})
}


func Locate(id Id, w Widget) Widget {
	if id == w.Private__getId() {
		return w
	}
	for _,w = range w.Private__getChildren() {
		out := Locate(id, w)
		if out != nil {
			return out
		}
	}
	return nil
}


type widgetwrapper struct {
	w Widget
	sends []func(string)
}
func (w *widgetwrapper) Done(err os.Error) {
	fmt.Println("Done with error:", err)
}
func (w *widgetwrapper) AddSend(send func(string)) {
	w.sends = append(w.sends, send)
}
func (w *widgetwrapper) Handle(evt string) {
	fmt.Println("Got event:", evt)
	evts := strings.Split(evt, ":", -1)
	switch evts[0] {
	case "onclick":
		clicked := Locate(Id(evts[1]), w.w)
		if clicked != nil {
			if clicked, ok := clicked.(Clickable); ok {
				r := clicked.HandleClick()
				fmt.Println("HandleClick gave", r)
			}
		}
	case "onchange":
		if len(evts) < 1 {
			fmt.Println("A broken onchange!")
			break
		}
		changed := Locate(Id(evts[1]), w.w)
		switch changed := changed.(type) {
		case interface { Changeable; Bool }:
			if len(evts) == 2 {
				fmt.Println("I got a nice event to toggle")
				changed.Toggle()
				r := changed.HandleChange()
				fmt.Println("HandleChange gave", r)
			} else {
				fmt.Println("Ignoring a strange onchange text event with", len(evts), "events!")
			}
		case interface { Changeable; String }:
			if len(evts) == 4 {
				changed.SetString(evts[3])
				r := changed.HandleChange()
				fmt.Println("HandleChange gave", r)
			} else {
				fmt.Println("Ignoring a strange onchange text event with", len(evts), "events!")
			}
		case nil:
			fmt.Println("There is no event with id", evts[1])
		default:
			fmt.Printf("I don't understand this event of type %t", changed)
		}
	}
	// if evt == "First time" {
	// 	dat.Write("read-cookie")
	// 	return
	// } else if strings.HasPrefix(evt, "cookie is ") {
	// 	fmt.Println("got cookie:", evt)
	// 	dat.Cookie = readCookie(evt[len("cookie is "):])
	// 	dat.WriteCookie()
	// 	return
	// }
	out := `<h3> Event is ` + evt + "</h3>\n"
	out += w.w.Private__html()
	for _,send := range w.sends {
		send(out)
	}
}
