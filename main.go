package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"text/template"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"gopkg.in/yaml.v2"
)

func main() {
	templateName := flag.String("t", "default.html", "template")
	style := flag.String("s", "monokai", "style")
	help := flag.Bool("h", false, "help")
	flag.Parse()

	if *help {
		fmt.Print("usage: press file.md\n\n")

		flag.PrintDefaults()
		return
	}

	mdFile := flag.Arg(0)
	md, err := os.ReadFile(mdFile)
	if err != nil {
		fatal(err)
	}

	tmpl, err := template.ParseFiles(*templateName)
	if err != nil {
		fatal(err)
	}

	sections := bytes.Split(md, []byte("\n---"))
	config, err := parseYaml(sections[0])
	if err != nil {
		fatal(err)
	}

	mdSections, err := parseMarkdown(sections, *style)
	if err != nil {
		fatal(err)
	}

	if err := createSlide(tmpl, mdSections, config, mdFile+".html"); err != nil {
		fatal(err)
	}

	fmt.Println("Done")
}

func createSlide(tmpl *template.Template, sections []string, config map[string]interface{}, outName string) error {
	fp, err := os.Create(outName)
	if err != nil {
		return err
	}

	data := struct {
		Sections []string
		Config   map[string]interface{}
		Len      int
	}{
		Sections: sections,
		Config:   config,
		Len:      len(sections),
	}
	tmpl.Execute(fp, data)

	fp.Close()

	return nil
}

func parseYaml(x []byte) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	if err := yaml.Unmarshal(x, m); err != nil {
		return nil, err
	}
	return m, nil
}

func parseMarkdown(sections [][]byte, style string) ([]string, error) {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle(style),
			),
		),
	)

	var out []string
	for i := 1; i < len(sections); i++ {
		var b bytes.Buffer
		if err := markdown.Convert(sections[i], &b); err != nil {
			return nil, err
		}
		out = append(out, b.String())
	}
	return out, nil
}

func fatal(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}
