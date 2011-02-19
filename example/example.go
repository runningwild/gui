package main

import (
	"http"
	"fmt"
	"github.com/droundy/widgets"
)

func main() {
	http.HandleFunc("/style.css", styleServer)
	buttonA := widgets.Button("A")
	buttonB := widgets.Button("B")
	buttonA.OnClick(func() widgets.Refresh {
		fmt.Println("I clicked on button A")
		buttonB.SetText(buttonB.GetText() + buttonB.GetText())
		return widgets.StillClean
	})
	buttonB.OnClick(func() widgets.Refresh {
		fmt.Println("I clicked on button A")
		t := buttonB.GetText()
		buttonB.SetText(t[:len(t)/2+1])
		return widgets.StillClean
	})
	name := widgets.EditText("Enter name here")
	hello := widgets.Text("Hello world!")
	name.OnChange(func() widgets.Refresh {
		hello.SetText("Hello " + name.GetText() + "!")
		return widgets.StillClean
	})
	err := widgets.Run(
		widgets.Column(
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
