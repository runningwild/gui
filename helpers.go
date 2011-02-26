package widgets

// This is a set of "helper" types that are designed to be embedded
// into struct types.  They are both used internally and are exported
// for users who may wish to define their own composite widgets.

import (
	"fmt"
)

type Id string
var newId <-chan Id
func init() {
	nid := make(chan Id, 5)
	go func() {
		i := 0
		for {
			i++
			nid <- Id(fmt.Sprint(i))
		}
	}()
	newId = nid
}
func (i *Id) Private__getId() Id {
	return *i
}
func (i *Id) Private__getChildren() []Widget {
	return []Widget{}
}

type CopyWidget struct {
	Widget
}
func (w *CopyWidget) Private__getId() Id {
	return w.Widget.Private__getId()
}
func (w *CopyWidget) Private__getChildren() []Widget {
	return w.Widget.Private__getChildren()
}
func (w *CopyWidget) Private__html() string {
	return w.Widget.Private__html()
}

type ChangeHandler Hook
func (o *ChangeHandler) OnChange(h Hook) {
	*o = ChangeHandler(h)
}
func (o *ChangeHandler) HandleChange() Refresh {
	if *o == nil {
		return StillClean
	}
	return (*o)()
}

type ClickHandler Hook
func (o *ClickHandler) OnClick(h Hook) {
	*o = ClickHandler(h)
}
func (o *ClickHandler) HandleClick() Refresh {
	if *o == nil {
		return StillClean
	}
	return (*o)()
}

type BoolValue bool
func (b *BoolValue) GetBool() bool {
	return bool(*b)
}
func (b *BoolValue) SetBool(x bool) {
	*b = BoolValue(x)
}
func (b *BoolValue) Toggle() {
	*b = ! *b
}

type BoolEcho struct {
	Bool Bool // this enables just the boolean portion
}
func (b *BoolEcho) GetBool() bool {
	return b.Bool.GetBool()
}
func (b *BoolEcho) SetBool(x bool) {
	b.Bool.SetBool(x)
}
func (b *BoolEcho) Toggle() {
	b.Bool.Toggle()
}
func (o *BoolEcho) OnChange(h Hook) {
	o.Bool.OnChange(h)
}
func (o *BoolEcho) HandleChange() Refresh {
	return o.Bool.HandleChange()
}


type TextEcho struct {
	Text HasText
}
func (t *TextEcho) SetText(s string) {
	t.Text.SetText(s)
}
func (t *TextEcho) GetText() string {
	return t.Text.GetText()
}
