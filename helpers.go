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
func (i *Id) getId() Id {
	return *i
}
func (i *Id) getChildren() []Widget {
	return []Widget{}
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



