//target:github.com/droundy/widgets

package widgets

import (
	"fmt"
	"html"
)

func Empty() Widget {
	return Text("")
}

func Text(t string) interface { Widget; String; Clickable } {
	return &text{<-newId, t, nil}
}

func EditText(t string) interface { Widget; Changeable; String } {
	return &edittext{text{<-newId, t, nil}, nil}
}

func Button(t string) interface { Widget; Clickable; String } {
	return &button{text{<-newId, t, nil}, nil}
}

func Checkbox() interface { Widget; Changeable; Bool } {
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
	ClickHandler
}
func (dat *text) Private__html() string {
	return `<span onclick="say('onclick:` + string(dat.Private__getId()) + ":" +
		dat.GetString() + `')">` + html.EscapeString(dat.string) + `</span>`
}
func (b *text) GetString() string {
	return b.string
}
func (b *text) SetString(newt string) {
	b.string = newt
}

type edittext struct {
	text
	ChangeHandler
}
func (dat *edittext) Private__html() string {
	h := `<input type="text" onchange="say('onchange:` + string(dat.Private__getId()) + ":" + dat.GetString() +
		`:' + this.value)" value="` + html.EscapeString(dat.text.GetString()) + `" />`
	//fmt.Println(h)
	return h
	return `<input type="text" onchange="say('onchange:` + string(dat.Private__getId()) + ":" + dat.GetString() +
		`:' + this.value)" value="` + dat.text.Private__html() + `" />`
}

type button struct {
	text
	ClickHandler
}
func (dat *button) Private__html() string {
	return `<input type="submit" onclick="say('onclick:` + string(dat.Private__getId()) + ":" + dat.GetString() + `')" value="` +
		html.EscapeString(dat.GetString()) + `" />`
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
	//fmt.Println(h)
	return h
}

func RadioButton(v string) interface { Widget; Changeable; Bool; String } {
	out := &radiobutton{ text{ <-newId, v, nil }, false, nil, nil }
	out.OnChange(func() Refresh {
		fmt.Println("I am toggling", out)
		return NeedsRefresh
	})
	out.OnClick(func() Refresh {
		out.Toggle()
		return out.HandleChange()
	})
	return out
}

type radiobutton struct {
	text
	BoolValue
	ChangeHandler
	ClickHandler
}
func (dat *radiobutton) Private__html() string {
	checked := ""
	if dat.GetBool() {
		checked = " checked='checked' "
	}
	return `<input type="radio" onchange="say('onchange:` + string(dat.Private__getId()) + `')" value="` +
		html.EscapeString(dat.GetString()) + `"` + checked +
		`/><span onclick="say('onclick:` + string(dat.Private__getId()) + ":" +
		dat.GetString() + `')">` + html.EscapeString(dat.string) + `</span>`
}

func RadioGroup(butts... interface{ Changeable; Bool; String }) interface { String; Changeable } {
	out := radiogroup{ butts, nil }
	numselected := 0
	for i := range butts {
		b := butts[i]
		if b.GetBool() {
			numselected += 1
		}
		b.OnChange(func () Refresh {
			bval := b.GetBool()
			for _,b2 := range out.buttons {
				if b2.GetString() != b.GetString() {
					b2.SetBool(!bval)
				}
			}
			out.HandleChange()
			return NeedsRefresh
		})
	}
	for _,b := range out.buttons {
		if b.GetBool() {
			numselected -= 1
		}
		if numselected == 0 {
			b.SetBool(true)
			numselected -= 1
		} else {
			b.SetBool(false)
		}
	}
	return &out
}

type radiogroup struct {
	buttons []interface{ Changeable; Bool; String }
	ChangeHandler
}

func (dat *radiogroup) GetString() string {
	for _,b := range dat.buttons {
		if b.GetBool() {
			return b.GetString();
		}
	}
	panic("Bug: radio group should always have one button selected!")
}

func (dat *radiogroup) SetString(v string) {
	foundstring := false
	for _,b := range dat.buttons {
		if b.GetString() == v {
			foundstring = true
		}
	}
	if foundstring == false {
		panic("Cannot SetString to "+v+" in radio group:  no such option!")
	}
	for _,b := range dat.buttons {
		if b.GetString() == v {
			b.SetBool(true)
		} else {
			b.SetBool(false)
		}
	}
}
