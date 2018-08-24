// obtain workday wsdls
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

const wwsdirectory = "https://community.workday.com/sites/default/files/file-hosting/productionapi/index.html"
const wdsldirectory = "../../wsdl"

func main() {

	baseURL, err := url.Parse(wwsdirectory)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Get(wwsdirectory)
	if err != nil {
		log.Fatalln(err)
	}
	// get all links that end in .xsd or .wdsl
	for _, h := range getHrefs(resp.Body, []string{".xsd", ".wsdl"}) {
		fmt.Println(h)
		err := getAndWrite(baseURL, h, wdsldirectory)
		if err != nil {
			fmt.Println("\t", err)
		}
	}

}

func getAndWrite(baseURL *url.URL, fragment string, dir string) error {
	urlFragment, err := url.Parse(fragment)
	if err != nil {
		return err
	}
	url := baseURL.ResolveReference(urlFragment)

	resp, err := http.Get(url.String())
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	parts := strings.Split(fragment, "/")
	last := parts[len(parts)-1]
	err = ioutil.WriteFile(fmt.Sprintf("%s/%s", dir, last), body, 0644)
	if err != nil {
		return err
	}
	return nil
}

// getHrefs gets all the HTML anchor href attributes
// that contain any of the filter strings
func getHrefs(body io.Reader, filter []string) []string {
	var refs []string
	t := html.NewTokenizer(body)
	for {
		tt := t.Next()
		switch tt {
		case html.ErrorToken:
			return refs
		case html.StartTagToken, html.EndTagToken:
			token := t.Token()
			if "a" == token.Data {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						if len(filter) > 0 {
							for _, f := range filter {
								if strings.LastIndex(attr.Val, f) > 0 {
									refs = append(refs, attr.Val)
								}
							}
						} else {
							refs = append(refs, attr.Val)
						}
					}
				}
			}
		}
	}
}
