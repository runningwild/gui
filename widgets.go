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
	return &text{<-newId, t}
}

func EditText(t string) HasChangingText {
	return &edittext{text{<-newId, t}, nil}
}

func Button(t string) ClickableWithText {
	return &button{text{<-newId, t}, nil}
}

func Checkbox() Bool {
	c := &checkbox{<- newId, false, nil}
	c.OnChange(func() Refresh {
		fmt.Println("I am toggling", c)
		return NeedsRefresh
	})
	return c
}

func Table(rows ...[]Widget) Widget {
	return &table{<-newId, rows}
}

func Column(widgets ...Widget) Widget {
	ws := make([][]Widget, len(widgets))
	for i := range ws {
		ws[i] = []Widget{widgets[i]}
	}
	return &table{<-newId, ws}
}

func Row(widgets ...Widget) Widget {
	return &table{<-newId, [][]Widget{widgets}}
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

///////////////////////////////////////
// Everything below this is private! //
///////////////////////////////////////

type table struct {
	Id
	ws [][]Widget
}
func (t *table) Private__html() string {
	out := "<table>\n"
	for _,r := range t.ws {
		out += "  <tr>\n"
		for _,w := range r {
			out += "    <td>" + w.Private__html() + "</td>\n"
		}
		out += "  </tr>\n"
	}
	out += "</table>\n"
	return out
}
func (t *table) Private__getChildren() []Widget {
	out := []Widget{}
	for _,ws := range t.ws {
		out = append(out, ws...)
	}
	return out
}

type text struct {
	Id
	string
}
func (dat *text) Private__html() string {
	return html.EscapeString(dat.string)
}
func (b *text) GetText() string {
	return b.string
}
func (b *text) SetText(newt string) {
	b.string = newt
}

type edittext struct {
	text
	ChangeHandler
}
func (dat *edittext) Private__html() string {
	h := `<input type="text" onchange="say('onchange:` + string(dat.Private__getId()) + ":" + dat.GetText() +
		`:' + this.value)" value="` + dat.text.Private__html() + `" />`
	fmt.Println(h)
	return h
	return `<input type="text" onchange="say('onchange:` + string(dat.Private__getId()) + ":" + dat.GetText() +
		`:' + this.value)" value="` + dat.text.Private__html() + `" />`
}

type button struct {
	text
	ClickHandler
}
func (dat *button) Private__html() string {
	return `<input type="submit" onclick="say('onclick:` + string(dat.Private__getId()) + ":" + dat.GetText() + `')" value="` +
		html.EscapeString(dat.GetText()) + `" />`
}

type checkbox struct {
	Id
	BoolValue
	ChangeHandler
}
func (dat *checkbox) Private__html() string {
	checked := ""
	if dat.GetBool() {
		checked = "checked='checked' "
	}
	h := `<input type="checkbox" onchange="say('onchange:` + string(dat.Id) + `')" ` + checked + `" />`
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

