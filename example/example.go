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
	label.OnClick(func() widgets.Refresh {
		cb.Toggle()
		return cb.HandleChange()
	})
	out := labelledcheckbox{ table, label, cb, cb }
	return &out
}

type radiobuttons struct {
	widgets.Widget
	widgets.String
	widgets.Changeable
}
	
func RadioButtons(vs... string) interface{ widgets.Widget; widgets.String; widgets.Changeable } {
	var bs []interface{ widgets.Changeable; widgets.Bool; widgets.String }
	var ws []widgets.Widget
	for _,v := range vs {
		b := widgets.RadioButton(v)
		bs = append(bs, b)
		ws = append(ws, b)
	}
	col := widgets.Column(ws...)
	grp := widgets.RadioGroup(bs...)
	return &radiobuttons{ col, grp, grp }
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

	// Now let's test out a set of radio buttons
	radio := RadioButtons("apples", "lemons", "oranges", "pears")
	radio_report := widgets.Text("I like to eat tasty fruit")
	menu := widgets.Menu("apples", "lemons", "oranges", "pears")
	radio.OnChange(func() widgets.Refresh {
		menu.SetString(radio.GetString());
		radio_report.SetString("I like to eat " + radio.GetString())
		return widgets.NeedsRefresh
	})

	menu.OnChange(func() widgets.Refresh {
		radio.SetString(menu.GetString())
		radio_report.SetString("I like to eat " + radio.GetString())
		return radio.HandleChange()
	})

	err := widgets.Run(
		widgets.Column(
		iscool,
		testing_checkbox,
		widgets.Row(buttonA, buttonB),
		widgets.Row(widgets.Text("Name:"), name),
		hello,
		widgets.Text("Goodbye world!"),
		radio,
		radio_report,
		menu,
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
