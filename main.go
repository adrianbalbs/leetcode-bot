package main

import (
	"fmt"

	"github.com/go-rod/rod"
)

func main() {
	page := rod.New().MustConnect().MustPage("https://neetcode.io/practice?tab=neetcode150")
	page.MustElement(`button.navbar-btn.is-rounded[data-tooltip="Show List View"]`).MustClick()
	tables := page.MustElements("table > tbody")
	for _, table := range tables {
		rows := table.MustElements("tr")
		for _, row := range rows {
			rowChild := row.MustElements("td")[2]
			title, err := rowChild.MustElement("a.table-text.text-color").Text()
			if err != nil {
				fmt.Println("No Title")
				return
			}
			link := rowChild.MustElement("a.has-tooltip-bottom").MustAttribute("href")
			fmt.Println(title, ":", *link)
		}

	}
}
