//target:github.com/droundy/widgets

package widgets

import (
	"os"
	"fmt"
	"html"
	"github.com/droundy/widgets/websocket"
)

type Widget interface {
	html() string
}

func Empty() Widget {
	return &text{""}
}

func Text(t string) Widget {
	return &text{t}
}

func Button(t string) Widget {
	return &button{t, <- NewId}
}

func Table(rows ...[]Widget) Widget {
	return &table{rows}
}

func Column(widgets ...Widget) Widget {
	ws := make([][]Widget, len(widgets))
	for i := range ws {
		ws[i] = []Widget{widgets[i]}
	}
	return &table{ws}
}

func Row(widgets ...Widget) Widget {
	return &table{[][]Widget{widgets}}
}

func RunSeparate(w func() Widget) os.Error {
	return websocket.RunSeparate("/", func() websocket.Handler {
		return &widgetwrapper{w(), []func(string){}}
	})
}

func Run(w Widget) os.Error {
	return websocket.Run("/", &widgetwrapper{w, []func(string){}})
}

/////////////////////////////////////////
// Here is the event-handling stuff... //
/////////////////////////////////////////

var NewId <-chan string
func init() {
	nid := make(chan string, 5)
	go func() {
		i := 0
		for {
			i++
			nid <- fmt.Sprint(i)
		}
	}()
	NewId = nid
}


///////////////////////////////////////
// Everything below this is private! //
///////////////////////////////////////

type table struct {
	ws [][]Widget
}
func (t *table) html() string {
	out := "<table>\n"
	for _,r := range t.ws {
		out += "  <tr>\n"
		for _,w := range r {
			out += "    <td>" + w.html() + "</td>\n"
		}
		out += "  </tr>\n"
	}
	out += "</table>\n"
	return out
}

type text struct {
	string
}
func (dat *text) html() string {
	return html.EscapeString(dat.string)
}

type button struct {
	string
	id string
}
func (dat *button) html() string {
	return `<input type="submit" onclick="say('button-` + dat.string + dat.id + `')" value="` + html.EscapeString(dat.string) + `" />`
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
	out += w.w.html()
	for _,send := range w.sends {
		send(out)
	}
}
