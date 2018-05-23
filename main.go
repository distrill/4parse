package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type thread struct {
	ID    string `db:"id"`
	Title string `db:"title"`
}

type comment struct {
	ID       string `db:"id"`
	ImageURL string `db:"image_url"`
	PosterID string `db:"poster_id"`
	IsOP     bool   `db:"is_op"`
	ThreadID string `db:"thread_id"`
}

type commentData struct {
	ID        string `db:"id"`
	DataType  string `db:"data_type"`
	Contents  string `db:"contents"`
	CommentID string `db:"comment_id"`
}

type commentReference struct {
	ID              string `db:"id"`
	ParentCommentID string `db:"parent_comment_id"`
	ChildCommentID  string `db:"child_comment_id"`
}

func main() {
	db, err := sqlx.Connect("mysql", "four_parse:four_parse@/four_parse")
	if err != nil {
		log.Panic(err)
	}

	ts := getThreads()

	// TODO: these threads are empty "welcome to 4chan" threads, no comments, they error.
	// fix so they are handled properly and don't have to be skipped
	threads := map[string]bool{"904256": true, "4884770": true}

	for _, t := range ts {
		td := getThreadInfo(t)
		if threads[td.ID] == false {
			fmt.Printf("running for thread %v\n", td.ID)

			err = insertThread(db, td)
			if err != nil {
				log.Fatal(err)
			}

			cs := getCommentsInfo(t, td.ID)
			err = insertComments(db, cs)
			if err != nil {
				log.Fatal(err)
			}

			cd := getCommentData(t)
			err = insertCommentData(db, cd)
			if err != nil {
				log.Fatal(err)
			}

			cr := getCommentReferences(t)
			err = insertCommentRefereneces(db, cr)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
