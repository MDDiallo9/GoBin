package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	// SQL Statement
	statement := `INSERT INTO snippets (title, content, created, expires)
    VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(statement, title, content, expires)
	if err != nil {
		return 0, err
	}
	// Getting the last insert ID to be sure the INSERT worked
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Returned ID is int64 type , we convert it before returning
	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	// SQL Statement , worth trying in the msql shell before to see if it's correct
	statement := `SELECT id,title,content,created,expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// DB Query
	row := m.DB.QueryRow(statement, id)

	// New Snippet
	s := &Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// If the query returns no rows, then row.Scan() will return a
		// sql.ErrNoRows error. We use the errors.Is() function check for that
		// error specifically, and return our own ErrNoRecord error
		// instead.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecords
		} else {
			return nil, err
		}
	}

	/* LESS VERBOSE VERSION
	 s := &Snippet{}
	err := m.DB.QueryRow("SELECT ...", id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil */
	return s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	statement := `SELECT id,title,content,created,expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	// Query instead of QueryRow because the result will contain more than a single row
	rows,err := m.DB.Query(statement)
	if err != nil {
		return nil,err
	}
	// Important to close connection to DB once we're done
	defer rows.Close()

	// Empty slice of Snippet structs
	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}
		// Each row gets scanned then added to a snippet
		err = rows.Scan(&s.ID,&s.Title,&s.Content,&s.Created,&s.Expires)
		if err != nil {
			return nil,err
		}
		snippets = append(snippets, s)
	}
	// Important to know if any errors occured during the iteration on rows
	if err = rows.Err(); err != nil {
		return nil,err
	}

	return snippets, nil
}
