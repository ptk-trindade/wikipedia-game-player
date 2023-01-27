package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

// Page is a wikipedia page
type Page struct {
	name     string
	linkWord string
	url      string // has this format "https://cc.wikipedia.org/wiki/{name}"
	foundBy  *Page
}

func main() {
	fmt.Println("--- start ---")

	// a := getPage("https://zx.wikipedia.org/wiki/Jadfsi")
	// x := linksInBody(a)
	// fmt.Println(len(x), x)

	list := findWikipediaPath("en", "Penguin", "Penguin_(album)")

	fmt.Println(len(list), list)
}

/*
-- input:
cc: country code (2 characters),
start: start page
end: goal page
-- output:
list of wiki pages from start to end [[2]string{linkWord, pageName}]
if you click in "linkWord" you will go to "pageName"
*/
func findWikipediaPath(cc, start, end string) [][2]string {
	cc = strings.ToLower(cc)
	start = strings.ToLower(start)
	end = strings.ToLower(end)

	// VALIDATION
	{
		// check if start and end are not the same
		if start == end {
			return [][2]string{}
		}

		// check if start and end exist
		page := getPage(fmt.Sprintf("https://%s.wikipedia.org/wiki/%s", cc, start))
		length := linksInBody(page)

		if len(length) == 0 {
			fmt.Println("Start page doesn't exist or is empty, remeber that page's name is case sensitive (the first letter is always capitalized)")
			return [][2]string{}
		}

		page = getPage(fmt.Sprintf("https://%s.wikipedia.org/wiki/%s", cc, end))
		length = linksInBody(page)

		if len(length) == 0 {
			fmt.Println("End page doesn't exist or is empty, remeber that page's name is case sensitive")
			return [][2]string{}
		}
	}

	// SEARCH (BFS)
	pages := make(map[string]*Page) // map of all pages found so far, key is pageName

	var firstPage, lastPage *Page
	firstPage = &Page{
		name:    start,
		url:     fmt.Sprintf("https://%s.wikipedia.org/wiki/%s", cc, start),
		foundBy: nil,
	}

	pages[start] = firstPage

	fifo := []Page{*firstPage}
	for len(fifo) > 0 && lastPage == nil {
		// fifo.pop()
		curPage := fifo[0]
		fifo = fifo[1:]

		pageBody := getPage(curPage.url)
		pageReferences := linksInBody(pageBody)
		for _, pageReference := range pageReferences {
			pageName := strings.ToLower(pageReference[0])
			if pages[pageName] == nil { // this page wasn't found yet
				wordLink := pageReference[1]

				foundPage := Page{
					name:     pageName,
					linkWord: wordLink,
					url:      fmt.Sprintf("https://%s.wikipedia.org/wiki/%s", cc, pageName),
					foundBy:  &curPage,
				}

				if pageName == end {
					lastPage = &foundPage
					break
				}

				pages[pageName] = &foundPage
				fifo = append(fifo, foundPage)
			}
		}

	}

	pagePath := [][2]string{}
	for curPage := lastPage; curPage != nil; curPage = curPage.foundBy {
		pagePath = append(pagePath, [2]string{curPage.linkWord, curPage.name})
	}

	reverseSlice(pagePath)

	return pagePath
}

/*
-- input:
text: the html text of a wikipedia page
-- output:
list of all pages referenced in the text
[2]string{pageName, linkWord}
*/
func linksInBody(text string) [][2]string {
	pages := [][2]string{}
	urlRegex := regexp.MustCompile(`href="/wiki/([a-zA-Z0-9./?=_-]+)".+?>(.+?)</a>`)
	for _, pageName := range urlRegex.FindAllString(text, -1) { // -1 means all matches (doesn't separate the match in paranteses)
		matches := urlRegex.FindStringSubmatch(pageName) // returns the match in paranteses separated {wholeMatch, pageName, linkWord}
		pages = append(pages, [2]string{matches[1], matches[2]})
	}
	return pages
}

/*
-- input:
url: An url
-- output:
the html text of the page
*/
func getPage(url string) string {
	a := ""
	fmt.Print("Will get page: ", url, " (press enter to continue) ")
	fmt.Scan(&a)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("error while getting page: ", url)
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error while reading page: ", url)
		log.Fatal(err)
	}

	return string(body)
}

// reverseSlice reverses a slice in place.
func reverseSlice[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
