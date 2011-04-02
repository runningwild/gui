package gui


type Widget interface {
	Private__html() string
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

type HasPath interface {
	SetPath(string)
	GetPath() string
	OnPath(Hook)
	HandlePath() Refresh
}
