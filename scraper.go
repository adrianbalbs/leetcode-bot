package main

import (
	"fmt"

	"github.com/go-rod/rod"
)

type LeetcodeLink struct {
	name string
	url  string
}

// TODO: Move this into a package
func problem(titleSlug string) {
}

func main() {
	page := rod.New().MustConnect().MustPage("https://neetcode.io/practice?tab=neetcode150")
	page.MustElement(`button.navbar-btn.is-rounded[data-tooltip="Show List View"]`).MustClick()
	tables := page.MustElements("table > tbody")

	// leetcodeLinks := []LeetcodeLink{}
	for _, table := range tables {
		for _, row := range table.MustElements("tr") {
			rowChild := row.MustElements("td")[2]
			problemTitle, err := rowChild.MustElement("a.table-text.text-color").Text()
			if err != nil {
				fmt.Println("Error: ", err)
				continue
			}
			link := rowChild.MustElement("a.has-tooltip-bottom").MustAttribute("href")
			fmt.Println(problemTitle, ":", *link)

		}
	}
	// TODO: for each link, call the graphql leetcode API for each problem and store in DB
}
