package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
)

var (
	prefix = "https://pdos.csail.mit.edu/6.828/2023/labs/"

	labFlag = flag.String("labs", "util,syscall,traps,pgtbl,lock,net,mmap,fs,cow,thread", "list of labs to generate")
)

func convertURL(lab string) {
	url := prefix + lab + ".html"

	rules := []md.Rule{
		{
			Filter: []string{"pre"},
			Replacement: func(content string, selec *goquery.Selection, options *md.Options) *string {
				return md.String("```\n" + content + "\n```")
			},
		},
		{
			Filter: []string{"tt"},
			Replacement: func(content string, selec *goquery.Selection, options *md.Options) *string {
				return md.String("`" + content + "`")
			},
		},
		{
			Filter: []string{"kbd"},
			Replacement: func(content string, selec *goquery.Selection, options *md.Options) *string {
				return md.String(content)
			},
		},
		{
			Filter: []string{"script"},
			Replacement: func(content string, selec *goquery.Selection, options *md.Options) *string {
				if !strings.HasPrefix(content, "g(") {
					return &content
				}
				level := content[3 : len(content)-2]

				return md.String(fmt.Sprintf("([%s](guidance.md))", level))
			},
		},
	}
	ccLicenseHook := func(markdown string) string {
		markdown = markdown + `
* * *
		
This is a derivative work based on [MIT 6.1810 Operating System Enginnering Fall 2023](` + fmt.Sprint(url) + `) 
used under [Creative Commons License](https://creativecommons.org/licenses/by/3.0/us/).

[![Creative Commons License](https://i.creativecommons.org/l/by/3.0/us/88x31.png)](https://creativecommons.org/licenses/by/3.0/us/)`

		return markdown
	}
	converter := md.NewConverter("", true, &md.Options{EscapeMode: "disabled"})
	converter.Keep("script")
	converter.AddRules(rules...)
	converter.After([]md.Afterhook{ccLicenseHook}...)

	if markdown, err := converter.ConvertURL(url); err != nil {
		log.Fatal(err)
	} else {
		err := os.WriteFile(lab+".md", []byte(markdown), 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

}
func main() {
	flag.Parse()

	labs := strings.Split(*labFlag, ",")
	for _, lab := range labs {
		convertURL(lab)
	}

}
