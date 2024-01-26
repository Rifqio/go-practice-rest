package models

import (
	"database/sql"
	"errors"
	"time"
)

type Post struct {
	ID      int            `json:"id"`
	Title   string         `json:"title"`
	Content sql.NullString `json:"content"`
	Created time.Time      `json:"created"`
	Expires time.Time      `json:"expires"`
}

type PostModel struct {
	DB *sql.DB
}

/*
* Query Cheat Sheet
* DB.Query() -is used for SELECT queries which return multiple rows.
* DB.QueryRow() -is used for SELECT queries which return a single row.
* DB.Exec() -is used for INSERT, UPDATE and DELETE queries, and it does not return any rows.
 */

func (post *PostModel) Insert(title string, content string) (int, error) {
	query := `insert into posts (title, content, created, expires) 
 			  values (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL 30 DAY))`

	result, err := post.DB.Exec(query, title, content)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}

	return int(id), nil
}

func (post *PostModel) Get(id int) (*Post, error) {
	query := `select * from posts where expires > UTC_TIMESTAMP() and id = ?`
	row := post.DB.QueryRow(query, id)

	p := &Post{}

	// Use row.Scan() to copy the values from each field in sql.Row to the
	// corresponding field in the Snippet struct. Notice that the arguments
	// to row.Scan are *pointers* to the place you want to copy the data into,
	// and the number of arguments must be exactly the same as the number of
	// columns returned by your statement.
	err := row.Scan(&p.ID, &p.Title, &p.Content, &p.Created, &p.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return p, nil
}

func (post *PostModel) Latest() ([]*Post, error) {
	query := `select * from posts where expires > UTC_TIMESTAMP() order by created desc limit 10`
	rows, err := post.DB.Query(query)

	if err != nil {
		return nil, err
	}

	// We defer rows.Close() to ensure the sql.Rows resultset is
	// always properly closed before the Latest() method returns. This defers
	// statement should come *after* you check for an error from the Query()
	// method. Otherwise, if Query() returns an error, you'll get a panic
	// trying to close a nil resultset.
	defer rows.Close()

	// Initialize empty slice to hold the posts
	posts := []*Post{}

	// Use rows.Next to iterate through the rows in the resultset. This
	// prepares the first (and then each subsequent) row to be acted on by the
	// rows.Scan() method. If iteration over all the rows completes then the
	// resultset automatically closes itself and frees-up the underlying
	// database connection.

	for rows.Next() {
		// Create a pointer to a new zeroed Post struct.
		p := &Post{}
		// Use rows.Scan() to copy the values from each field in the row to the
		// new Post object that we created. Again, the arguments to row.Scan()
		// must be pointers to the place you want to copy the data into, and the
		// number of arguments must be exactly the same as the number of
		// columns returned by your statement.
		err = rows.Scan(&p.ID, &p.Title, &p.Content, &p.Created, &p.Expires)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of posts.
		posts = append(posts, p)
	}

	// When the rows.Next() loop has finished we call rows.Err() to retrieve any
	// error that was encountered during the iteration. It's important to
	// call this - don't assume that a successful iteration was completed
	// over the whole resultset.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything went OK then return the Posts slice.
	return posts, nil
}
