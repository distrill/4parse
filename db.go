package main

import (
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func isErrDuplicateEntry(err error) bool {
	me, ok := err.(*mysql.MySQLError)
	return ok && me.Number == 1062
}

func insertThread(db *sqlx.DB, t thread) error {
	columns := "id"
	columnValues := ":id"

	if t.Title != "" {
		columns += ", title"
		columnValues += ", :title"
	}

	qs := `
		INSERT INTO threads (` + columns + `)
		VALUES (` + columnValues + `)
	`

	if _, err := db.NamedExec(qs, &t); err != nil && !isErrDuplicateEntry(err) {
		fmt.Printf("ERR: t %v\n", t)
		return err
	}

	return nil
}

func insertComments(db *sqlx.DB, cs []comment) error {
	// fmt.Println("inserting comments")

	qs := `
		INSERT INTO comments (id, image_url, poster_id, is_op, thread_id)
		VALUES (:id, :image_url, :poster_id, :is_op, :thread_id)
	`

	for _, c := range cs {
		if _, err := db.NamedExec(qs, &c); err != nil && !isErrDuplicateEntry(err) {
			fmt.Printf("ERR: c %v\n", c)
			return err
		}
	}

	return nil
}

func insertCommentData(db *sqlx.DB, cds []commentData) error {
	// fmt.Print("\ninserting comment data\n")

	qs := `
		INSERT INTO comment_data (data_type, contents, comment_id)
		VALUES (:data_type, :contents, :comment_id)
	`

	for _, cd := range cds {
		if _, err := db.NamedExec(qs, &cd); err != nil && !isErrDuplicateEntry(err) {
			fmt.Printf("ERR: cd %v\n", cd)
			return err
		}
	}

	return nil
}

func insertCommentRefereneces(db *sqlx.DB, crs []commentReference) error {
	// fmt.Print("\ninserting comment references\n")

	qs := `
		INSERT INTO comment_references (parent_comment_id, child_comment_id)
		VALUES (:parent_comment_id, :child_comment_id)
	`

	for _, cr := range crs {
		if _, err := db.NamedExec(qs, &cr); err != nil && !isErrDuplicateEntry(err) {
			fmt.Printf("ERR: cr %v\n", cr)
			return err
		}
	}

	return nil
}
