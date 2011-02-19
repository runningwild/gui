//target:github.com/droundy/widgets

package widgets

import (
	"os"
	"fmt"
	"html"
	"strings"
	"github.com/droundy/widgets/websocket"
)

type Widget interface {
	html() string
	locate(id string) Widget
}

func Empty() Widget {
	return Text("")
}

func Text(t string) HasText {
	tt := text(t)
	return &tt
}

func Button(t string) ClickableWithText {
	return &button{text(t), <- NewId, nil}
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

type HasText interface {
	Widget
	SetText(string)
	GetText() string
}

type Clickable interface {
	Widget
	OnClick(Hook)
	HandleClick() Refresh
}

type ClickableWithText interface {
	Clickable
	SetText(string)
	GetText() string
}

type Refresh bool
const (
	NeedsRefresh Refresh = true
	StillClean Refresh = false
)
type Hook func() Refresh
func (r Refresh) String() string {
	if r {
		return "NeedsRefresh"
	}
	return "StillClean"
}

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
func (t *table) locate(id string) Widget {
	for _,r := range t.ws {
		for _,w := range r {
			if ans := w.locate(id); ans != nil {
				return ans
			}
		}
	}
	return nil
}

type text string
func (dat *text) html() string {
	return html.EscapeString(string(*dat))
}
func (*text) locate(id string) Widget {
	return nil
}
func (b *text) GetText() string {
	return string(*b)
}
func (b *text) SetText(newt string) {
	*b = text(newt)
}

type onclick Hook
func (o *onclick) OnClick(h Hook) {
	*o = onclick(h)
}
func (o *onclick) HandleClick() Refresh {
	if *o == nil {
		return StillClean
	}
	return (*o)()
}

type button struct {
	text
	id string
	onclick
}
func (dat *button) html() string {
	return `<input type="submit" onclick="say('onclick:` + dat.id + ":" + string(dat.text) + `')" value="` +
		html.EscapeString(string(dat.text)) + `" />`
}
func (b *button) locate(id string) Widget {
	if b.id == id {
		return b
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
		clicked := w.w.locate(evts[1])
		if clicked != nil {
			if clicked, ok := clicked.(Clickable); ok {
				r := clicked.HandleClick()
				fmt.Println("HandleClick gave", r)
			}
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
	out += w.w.html()
	for _,send := range w.sends {
		send(out)
	}
}
