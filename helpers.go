package gui

// This is a set of "helper" types that are designed to be embedded
// into struct types.  They are both used internally and are exported
// for users who may wish to define their own composite gui.

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

type PathHandler struct {
	Hook
	Path string
}
func (o *PathHandler) SetPath(p string) Refresh {
	o.Path = p
	if o.Hook == nil {
		return StillClean
	}
	return o.Hook()
}
func (o *PathHandler) GetPath() string {
	return o.Path
}
func (o *PathHandler) OnPath(h Hook) {
	o.Hook = h
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
