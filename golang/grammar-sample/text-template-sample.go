package main

import (
	"fmt"
	"os"
	"text/template"
)

func main() {
	/*type Inventory struct {
		Material string
		Count    int
	}*/
	maps := map[string]interface{}{
		"key1": "111",
		"key2": "222",
	}
	//sweaters := Inventory{"axe", 1}
	html := `{{eq len(.maxdm) 0}}.maxdm{{else}}4g{{end}}`
	tmpl, err := template.New("test").Option("missingkey=error").Parse(html)
	if err != nil {
		panic(err)
	}
	fmt.Println("-----------")
	err = tmpl.Execute(os.Stdout, maps)
	fmt.Println()
	fmt.Println("-----------")
	if err != nil {
		panic(err)
	}
}
