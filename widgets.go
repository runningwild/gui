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
	locate(id Id) Widget
}

func Empty() Widget {
	return Text("")
}

func Text(t string) HasText {
	tt := text(t)
	return &tt
}

func EditText(t string) HasChangingText {
	return &edittext{text(t), <- NewId, nil}
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

type HasChangingText interface {
	HasText
	OnChange(Hook)
	HandleChange() Refresh
}

type Changeable interface {
	Widget
	OnChange(Hook)
	HandleChange() Refresh
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

type Id string
var NewId <-chan Id
func init() {
	nid := make(chan Id, 5)
	go func() {
		i := 0
		for {
			i++
			nid <- Id(fmt.Sprint(i))
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
func (t *table) locate(id Id) Widget {
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
func (*text) locate(id Id) Widget {
	return nil
}
func (b *text) GetText() string {
	return string(*b)
}
func (b *text) SetText(newt string) {
	*b = text(newt)
}

type edittext struct {
	text
	Id
	onchange
}
func (dat *edittext) html() string {
	h := `<input type="text" onchange="say('onchange:` + string(dat.Id) + ":" + string(dat.text) +
		`:' + this.value)" value="` + dat.text.html() + `" />`
	fmt.Println(h)
	return h
	return `<input type="text" onchange="say('onchange:` + string(dat.Id) + ":" + string(dat.text) +
		`:' + this.value)" value="` + dat.text.html() + `" />`
}
func (w *edittext) locate(id Id) Widget {
	if w.Id == id {
		return w
	}
	return nil
}

type onchange Hook
func (o *onchange) OnChange(h Hook) {
	*o = onchange(h)
}
func (o *onchange) HandleChange() Refresh {
	if *o == nil {
		return StillClean
	}
	return (*o)()
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
	Id
	onclick
}
func (dat *button) html() string {
	return `<input type="submit" onclick="say('onclick:` + string(dat.Id) + ":" + string(dat.text) + `')" value="` +
		html.EscapeString(string(dat.text)) + `" />`
}
func (b *button) locate(id Id) Widget {
	if b.Id == id {
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
		clicked := w.w.locate(Id(evts[1]))
		if clicked != nil {
			if clicked, ok := clicked.(Clickable); ok {
				r := clicked.HandleClick()
				fmt.Println("HandleClick gave", r)
			}
		}
	case "onchange":
		if len(evts) == 4 {
			changed := w.w.locate(Id(evts[1]))
			if changed, ok := changed.(HasChangingText); ok {
				changed.SetText(evts[3])
				r := changed.HandleChange()
				fmt.Println("HandleChange gave", r)
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
