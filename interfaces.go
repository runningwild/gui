package widgets


type Widget interface {
	Private__html() string
	Private__getId() Id
	Private__getChildren() []Widget
}

type HasText interface {
	Widget
	SetText(string)
	GetText() string
}

type OnlyText interface {
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

type Bool interface {
	Changeable
	GetBool() bool
	SetBool(bool)
	Toggle()
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
