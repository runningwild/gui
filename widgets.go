//target:github.com/droundy/widgets

package widgets

import (
	"fmt"
	"html"
)

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

func Checkbox() Bool {
	c := &checkbox{false, <- NewId, nil, nil}
	c.OnChange(func() Refresh {
		c.Toggle()
		fmt.Println("I am toggling", c)
		return NeedsRefresh
	})
	c.OnClick(func() Refresh {
		c.Toggle()
		fmt.Println("I am toggling", c)
		return NeedsRefresh
	})
	return c
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

/////////////////////////////////////////
// Here is the event-handling stuff... //
/////////////////////////////////////////


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

type boolthing bool
func (b *boolthing) GetBool() bool {
	return bool(*b)
}
func (b *boolthing) SetBool(x bool) {
	*b = boolthing(x)
}
func (b *boolthing) Toggle() {
	*b = ! *b
}

type checkbox struct {
	boolthing
	Id
	onchange
	onclick
}
func (dat *checkbox) html() string {
	checked := ""
	if dat.GetBool() {
		checked = "checked='checked' "
	}
	h := `<input type="checkbox" onclick="say('onchange:` + string(dat.Id) + `')" ` + checked + `" />`
	fmt.Println(h)
	return h
}
func (b *checkbox) locate(id Id) Widget {
	fmt.Println("Doing locate", id, "and I am", b.Id)
	if b.Id == id {
		return b
	}
	return nil
}

