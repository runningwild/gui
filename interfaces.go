package gui

type extraCommand string

type Widget interface {
	Private__html() (html string, extra_commands []extraCommand)
	Private__getId() Id
	Private__getChildren() []Widget
}

type String interface {
	SetString(string)
	GetString() string
}

type OnlyText interface {
	SetText(string)
	GetText() string
}

type Changeable interface {
	OnChange(Hook)
	HandleChange() Refresh
}

type Bool interface {
	GetBool() bool
	SetBool(bool)
	Toggle()
}

type Clickable interface {
	OnClick(Hook)
	HandleClick() Refresh
}

type PathHandler interface {
	Widget
	SetWidget(Widget)
	SetPath(string) Refresh
	GetPath() string
	OnPath(Hook)
}
