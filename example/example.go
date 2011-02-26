package main

import (
	"http"
	"fmt"
	"github.com/droundy/widgets"
)

// FIXME: There is no way to make the following return something that
// is both a widgets.Bool and a widgets.HasText.
type labelledcheckbox struct {
	widgets.Widget
	widgets.String
	widgets.Changeable
	widgets.Bool
}
func LabelledCheckbox(l string) interface { widgets.Widget; widgets.String; widgets.Changeable; widgets.Bool } {
	cb := widgets.Checkbox()
	label := widgets.Text(l)
	table := widgets.Row(cb, label)
	out := labelledcheckbox{ table, label, cb, cb }
	return &out
}

func main() {
	http.HandleFunc("/style.css", styleServer)
	buttonA := widgets.Button("A")
	buttonB := widgets.Button("B")
	buttonA.OnClick(func() widgets.Refresh {
		fmt.Println("I clicked on button A")
		buttonB.SetString(buttonB.GetString() + buttonB.GetString())
		return widgets.StillClean
	})
	buttonB.OnClick(func() widgets.Refresh {
		fmt.Println("I clicked on button A")
		t := buttonB.GetString()
		buttonB.SetString(t[:len(t)/2+1])
		return widgets.StillClean
	})
	iscool := widgets.Checkbox()
	name := widgets.EditText("Enter name here")
	hello := widgets.Text("Hello world!")
	name.OnChange(func() widgets.Refresh {
		hello.SetString("Hello " + name.GetString() + "!")
		return widgets.StillClean
	})
	testing_checkbox := LabelledCheckbox("testing")
	testing_checkbox.OnChange(func() widgets.Refresh {
		fmt.Println("Hello world")
		if testing_checkbox.GetBool() {
			testing_checkbox.SetString("this test is currently true")
		} else {
			testing_checkbox.SetString("this test is now false")
		}
		return widgets.NeedsRefresh
	})
	err := widgets.Run(
		widgets.Column(
		iscool,
		testing_checkbox,
		widgets.Row(buttonA, buttonB),
		widgets.Row(widgets.Text("Name:"), name),
		hello,
		widgets.Text("Goodbye world!"),
		))
	if err != nil {
		panic("ListenAndServe: " + err.String())
	}
}

func styleServer(c http.ResponseWriter, req *http.Request) {
	c.SetHeader("Content-Type", "text/css")
	fmt.Fprint(c, `
html {
    margin: 0;
    padding: 0;
}

body {
    margin: 0;
    padding: 0;
    background: #ffffff;
    font-family: arial,helvetica,"sans serif";
    font-size: 12pt;
}
h1 {
font-family: verdana,helvetica,"sans serif";
font-weight: bold;
font-size: 16pt;
}
h2 { font-family: verdana,helvetica,"sans serif";
font-weight: bold;
font-size: 14pt;
}
p {
font-family: arial,helvetica,"sans serif";
font-size:12pt;
}
li {
  font-family: arial,helvetica,"sans serif";
  font-size: 12pt;
}
a {
  color: #555599;
}
`)
}
