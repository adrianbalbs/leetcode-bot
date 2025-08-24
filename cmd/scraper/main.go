package main

import (
	"adrainbalbs/leetcode-bot/leetcode"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/go-rod/rod"
)

const (
	maxWorkers = 10
	maxRetries = 4
)

func scrapeNeetcode() []string {
	page := rod.New().MustConnect().MustPage("https://neetcode.io/practice?tab=neetcode150")
	page.MustElement(`button.navbar-btn.is-rounded[data-tooltip="Show List View"]`).MustClick()
	tables := page.MustElements("table > tbody")

	neetcodeProblems := []string{}

	for _, table := range tables {
		for _, row := range table.MustElements("tr") {
			rowChild := row.MustElements("td")[2]

			link := rowChild.MustElement("a.has-tooltip-bottom").MustAttribute("href")
			neetcodeProblems = append(neetcodeProblems, strings.Replace(strings.Trim(*link, "/"),
				"https://leetcode.com/problems/", "", 1))
		}
	}
	return neetcodeProblems
}

func worker(ctx context.Context, client graphql.Client, jobs <-chan string, results chan<- *leetcode.GetProblemResponse) {
	for problem := range jobs {
		for attempt := 1; attempt <= maxRetries; attempt++ {
			response, err := leetcode.GetProblem(ctx, client, problem)
			if err == nil {
				results <- response
				break
			}

			// Use exponential backoff in between retry attempts
			backoff := time.Duration(attempt*attempt) * time.Second
			log.Printf("Retrying %s after error %v", problem, err)
			time.Sleep(backoff)
		}
	}
}

func main() {
	problems := scrapeNeetcode()
	jobs := make(chan string, len(problems))
	results := make(chan *leetcode.GetProblemResponse, len(problems))
	ctx := context.Background()

	client := graphql.NewClient(leetcode.LeetcodeURL,
		&http.Client{
			Transport: &leetcode.UserAgentTransport{
				Wrapped: http.DefaultTransport,
			},
			Timeout: 10 * time.Second,
		})

	for range maxWorkers {
		go worker(ctx, client, jobs, results)
	}

	for _, problemSlug := range problems {
		if problemSlug != "" {
			jobs <- problemSlug
		}
	}
	close(jobs)

	for range problems {
		res := <-results
		fmt.Printf("Id: %s, Title: %s, Title Slug %s\n", res.Question.QuestionId, res.Question.Title, res.Question.TitleSlug)
	}
}
