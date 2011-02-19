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

func Column(widgets ...Widget) Widget {
	return &column{widgets}
}

func Row(widgets ...Widget) Widget {
	return &row{widgets}
}

func RunSeparate(w func() Widget) os.Error {
	return websocket.RunSeparate("/", func() websocket.Handler {
		return &widgetwrapper{w(), []func(string){}}
	})
}

func Run(w Widget) os.Error {
	return websocket.Run("/", &widgetwrapper{w, []func(string){}})
}

///////////////////////////////////////
// Everything below this is private! //
///////////////////////////////////////

type column struct {
	ws []Widget
}
func (c *column) html() string {
	out := "<table>\n"
	for _,w := range c.ws {
		out += "<tr><td>" + w.html() + "</td></tr>\n"
	}
	out += "\n</table>"
	return out
}

type row struct {
	ws []Widget
}
func (c *row) html() string {
	out := "<table><tr>\n"
	for _,w := range c.ws {
		out += "<td>" + w.html() + "</td>\n"
	}
	out += "\n</tr></table>"
	return out
}

type text struct {
	string
}
func (dat *text) html() string {
	return html.EscapeString(dat.string)
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
