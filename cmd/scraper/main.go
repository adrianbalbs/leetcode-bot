package main

import (
	"adrainbalbs/leetcode-bot/leetcode"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/go-rod/rod"
	_ "github.com/jackc/pgx/v5/stdlib"
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

func insertProblem(db *sql.DB, problem *leetcode.GetProblemResponse, playlistId int64) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var difficultyId int64
	err = tx.QueryRow(`
        INSERT INTO difficulties (name)
        VALUES ($1)
        ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
        RETURNING id
    `, problem.Question.Difficulty).Scan(&difficultyId)
	if err != nil {
		return 0, err
	}

	var problemId int64
	err = tx.QueryRow(`
        INSERT INTO problems (slug, title, difficulty_id)
        VALUES ($1, $2, $3)
        ON CONFLICT (slug) DO UPDATE SET title = EXCLUDED.title
        RETURNING id
    `, problem.Question.TitleSlug, problem.Question.Title, difficultyId).Scan(&problemId)
	if err != nil {
		return 0, err
	}

	for _, tag := range problem.Question.TopicTags {
		var tagId int64
		err := tx.QueryRow(`
			INSERT INTO tags (name)
			VALUES ($1)
			ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
			RETURNING id
		`, tag.Name).Scan(&tagId)
		if err != nil {
			return 0, err
		}

		_, err = tx.Exec(`
            INSERT INTO problem_tags (problem_id, tag_id)
            VALUES ($1, $2)
            ON CONFLICT DO NOTHING
		`, problemId, tagId)
		if err != nil {
			return 0, err
		}
	}

	_, err = tx.Exec(`
		INSERT INTO playlist_entries (playlist_id, problem_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, playlistId, problemId)

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return problemId, nil
}

func main() {
	problems := scrapeNeetcode()
	jobs := make(chan string, len(problems))
	results := make(chan *leetcode.GetProblemResponse, len(problems))
	ctx := context.Background()

	dbConnStr := os.Getenv("DATABASE_URL")
	db, err := sql.Open("pgx", dbConnStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// First, create the neetcode150 playlist

	fmt.Println("Creating Neetcode150 Playlist")
	var playlistId int64
	err = db.QueryRow(`
		INSERT INTO playlists (name, creator)
		VALUES($1, $2)
		ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`, "Neetcode150", "Neetcode.io").Scan(&playlistId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating playlist %v\n", err)
		os.Exit(1)
	}

	client := graphql.NewClient(leetcode.LeetcodeURL,
		&http.Client{
			Transport: &leetcode.UserAgentTransport{
				Wrapped: http.DefaultTransport,
			},
			Timeout: 10 * time.Second,
		})

	fmt.Println("Fetching and inserting problems")

	var wg sync.WaitGroup
	for range maxWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(ctx, client, jobs, results)
		}()
	}

	for _, problemSlug := range problems {
		if problemSlug != "" {
			jobs <- problemSlug
		}
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		fmt.Printf("Inserting %s\n", res.Question.Title)
		_, err := insertProblem(db, res, playlistId)
		if err != nil {
			log.Printf("Failed inserting problem %s: %v", res.Question.TitleSlug, err)
			continue
		}
	}

	fmt.Println("Finished scraping neetcode problems")
}
