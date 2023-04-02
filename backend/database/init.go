package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type DB struct {
	*sql.DB
}

// TODO: add column changing support

// InitDatabase creates all the necessary tables and columns in sql.DB
//
// If the table already exists, it will be checked for the columns. If the columns are missing, they will be added.
func (db DB) InitDatabase() error {
	err := db.InitTable("users", []Column{
		{Name: "id", Type: "INTEGER", PrimaryKey: true},
		{Name: "name", Type: "TEXT"},
		{Name: "email", Type: "TEXT"},
		{Name: "password", Type: "TEXT"},
	})
	if err != nil {
		return err
	}

	err = db.InitTable("posts", []Column{
		{Name: "id", Type: "INTEGER", PrimaryKey: true},
		{Name: "title", Type: "TEXT"},
		{Name: "content", Type: "TEXT"},
		{Name: "author", Type: "INTEGER"},
		{Name: "date", Type: "TEXT"},
		{Name: "likes_count", Type: "INTEGER"},
		{Name: "dislikes_count", Type: "INTEGER"},
		{Name: "comments_count", Type: "INTEGER"},
		{Name: "categories", Type: "TEXT"},
	}, "FOREIGN KEY(author) REFERENCES users(id)",
	)

	if err != nil {
		return err
	}

	err = db.InitTable("comments", []Column{
		{Name: "id", Type: "INTEGER", PrimaryKey: true},
		{Name: "post_id", Type: "INTEGER"},
		{Name: "author_id", Type: "INTEGER"},
		{Name: "content", Type: "TEXT"},
		{Name: "date", Type: "TEXT"},
		{Name: "likes_count", Type: "INTEGER"},
		{Name: "dislikes_count", Type: "INTEGER"},
	}, "FOREIGN KEY(post_id) REFERENCES posts(id)",
		"FOREIGN KEY(author_id) REFERENCES users(id)")
	if err != nil {
		return err
	}

	err = db.InitTable("sessions", []Column{
		{Name: "id", Type: "INTEGER", PrimaryKey: true},
		{Name: "token", Type: "TEXT"},
		{Name: "expire", Type: "INTEGER"},
		{Name: "user_id", Type: "INTEGER"},
	}, "FOREIGN KEY(user_id) REFERENCES users(id)")
	if err != nil {
		return err
	}

	err = db.InitTable("post_reactions", []Column{
		{Name: "id", Type: "INTEGER", PrimaryKey: true},
		{Name: "post_id", Type: "INTEGER"},
		{Name: "author_id", Type: "INTEGER"},
		{Name: "reaction", Type: "INTEGER"},
	}, "FOREIGN KEY(author_id) REFERENCES users(id)",
		"FOREIGN KEY(post_id) REFERENCES posts(id)")
	if err != nil {
		return err
	}

	err = db.InitTable("comment_reactions", []Column{
		{Name: "id", Type: "INTEGER", PrimaryKey: true},
		{Name: "comment_id", Type: "INTEGER"},
		{Name: "author_id", Type: "INTEGER"},
		{Name: "reaction", Type: "INTEGER"},
	}, "FOREIGN KEY(author_id) REFERENCES users(id)",
		"FOREIGN KEY(comment_id) REFERENCES comments(id)")
	if err != nil {
		return err
	}

	return nil
}

type Column struct {
	Name string
	Type string

	// optional, for PRAGMA table_info
	Id         int
	PrimaryKey bool
	NotNull    bool
	Default    interface{}
}

// InitTable creates a table with the given name and columns
//
// If the table already exists, it will be checked and updated if necessary.
// Additional arguments are added to the end of the query.
//
// Note: this function does not support changing existing columns
func (db DB) InitTable(table string, columns []Column, additional ...string) error {
	query, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s)", table))
	if err != nil {
		return err
	}
	var oldColumns []Column
	for query.Next() {
		var column Column
		err = query.Scan(&column.Id, &column.Name, &column.Type, &column.NotNull, &column.Default, &column.PrimaryKey)
		if err != nil {
			return err
		}
		oldColumns = append(oldColumns, column)
	}
	query.Close()

	if len(oldColumns) == 0 {
		var q string
		for _, column := range append(columns) {
			q += column.Name + " " + column.Type
			if column.PrimaryKey {
				q += " PRIMARY KEY"
			} else {
				if column.NotNull {
					q += " NOT NULL"
				}
				if column.Default != nil {
					q += fmt.Sprintf(" DEFAULT %v", column.Default)
				}
			}
			q += ", "
		}
		for _, a := range additional {
			q += a + ", "
		}
		q = q[:len(q)-2] // remove last ", "

		_, err = db.Exec(fmt.Sprintf("CREATE TABLE %s (%s)", table, q))
		if err != nil {
			return err
		}
		return nil
	}

	if len(oldColumns) > len(columns) {
		return fmt.Errorf("old database has more columns than expected in table '%s', got %d, expected %d", table, len(oldColumns), len(columns))
	}

	// check existing columns
	for i, column := range oldColumns {
		if column.Name != columns[i].Name {
			return fmt.Errorf("old column #%d in %s has name %s, but should be %s", i, table, column.Name, columns[i].Name)
		}
		if column.Type != columns[i].Type {
			return fmt.Errorf("old column %s in %s has type %s, but should be %s", column.Name, table, column.Type, columns[i].Type)
		}
		if column.PrimaryKey != columns[i].PrimaryKey {
			return fmt.Errorf("old column %s in %s has PrimaryKey %v, but should be %v", column.Name, table, column.PrimaryKey, columns[i].PrimaryKey)
		}
		if column.NotNull != columns[i].NotNull {
			return fmt.Errorf("old column %s in %s has NotNull %v, but should be %v", column.Name, table, column.NotNull, columns[i].NotNull)
		}
	}

	// add new columns
	for i := len(oldColumns); i < len(columns); i++ {
		log.Printf("creating a new column '%s' in table '%s'", columns[i].Name, table)
		_, err = db.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table, columns[i].Name, columns[i].Type))
		if err != nil {
			return err
		}
		if columns[i].Default != nil {
			_, err = db.Exec(fmt.Sprintf("UPDATE %s SET %s = %v", table, columns[i].Name, columns[i].Default))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
