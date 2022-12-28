package database

import (
	"encoding/json"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// InitDatabase initializes the database
func InitDatabase() error {
	db := OpenDB()

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		name TEXT,
		email TEXT,
		password TEXT,
		posts TEXT,
		comments TEXT
	)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS author (
		id INTEGER,
		name TEXT,
		FOREIGN KEY(id) REFERENCES users(id)
	)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS posts 
					(id INTEGER PRIMARY KEY,
					title TEXT, 
					content TEXT, 
					author INTEGER, 
					date TEXT, 
					likes INTEGER, 
					dislikes INTEGER, 
					user_reactions TEXT, 
					comments_count INTEGER, 
					comments_ids TEXT,
					categories TEXT,
					FOREIGN KEY(author) REFERENCES users(id))`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS comments 
					(id INTEGER PRIMARY KEY, 
					post INTEGER, 
					author INTEGER, 
					text TEXT, 
					date TEXT,
					likes INTEGER,
					dislikes INTEGER,
					user_reactions TEXT, 
					FOREIGN KEY(post) REFERENCES posts(id), 
					FOREIGN KEY(author) REFERENCES users(id))`)
	if err != nil {
		return err
	}

	err = insertSampleData()
	if err != nil {
		return err
	}

	return nil
}

// insert smaple data
func insertSampleData() error {

	db := OpenDB()

	// clean whole database
	_, err := db.Exec(`DELETE FROM users`)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = db.Exec(`DELETE FROM posts`)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = db.Exec(`DELETE FROM comments`)
	if err != nil {
		log.Println(err)
		return err
	}

	// Insert sample users
	_, err = db.Exec(`INSERT INTO users (id, name, email, password, posts, comments) VALUES (?, ?, ?, ?, ?, ?)`, 1, "John Smith", "john@example.com", "password123", "[1, 2, 3]", "[1, 2, 3]")
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = db.Exec(`INSERT INTO users (id, name, email, password, posts, comments) VALUES (?, ?, ?, ?, ?, ?)`, 2, "Jane Doe", "jane@example.com", "password456", "[4, 5, 6]", "[4, 5, 6]")
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = db.Exec(`INSERT INTO users (id, name, email, password, posts, comments) VALUES (?, ?, ?, ?, ?, ?)`, 3, "Bob Johnson", "bob@example.com", "password789", "[7, 8, 9]", "[7, 8, 9]")
	if err != nil {
		log.Println(err)
		return err
	}

	// Insert sample authors
	_, err = db.Exec(`INSERT INTO author (id, name) VALUES (?, ?)`, 1, "John Smith")
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = db.Exec(`INSERT INTO author (id, name) VALUES (?, ?)`, 2, "Jane Doe")
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = db.Exec(`INSERT INTO author (id, name) VALUES (?, ?)`, 3, "Bob Johnson")
	if err != nil {
		log.Println(err)
		return err
	}

	// make an empty map of user reactions
	uRs := make(map[int]int)
	// Convert the map to a JSON string.
	jsonData, err := json.Marshal(uRs)
	if err != nil {
		log.Println(err)
		return err
	}
	uRsStr := string(jsonData)

	// Insert sample comments
	usrReact := make(map[int]int)
	_, err = db.Exec(`INSERT INTO comments (id, post, author, text, date, likes, dislikes, user_reactions) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, 1, 1, 1, "This is a comment on the 1st post", "2022-01-01", 0, 0, vomitJSON(usrReact))
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = db.Exec(`INSERT INTO comments (id, post, author, text, date, likes, dislikes, user_reactions) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, 2, 1, 1, "This is a comment on the 1st post", "2022-01-02", 0, 0, vomitJSON(usrReact))
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = db.Exec(`INSERT INTO comments (id, post, author, text, date, likes, dislikes, user_reactions) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, 3, 1, 1, "This is a comment on the 1st post", "2022-01-03", 0, 0, vomitJSON(usrReact))
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = db.Exec(`INSERT INTO comments (id, post, author, text, date, likes, dislikes, user_reactions) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, 4, 2, 2, "This is a comment on the 2nd post", "2022-01-04", 0, 0, vomitJSON(usrReact))
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = db.Exec(`INSERT INTO comments (id, post, author, text, date, likes, dislikes, user_reactions) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, 5, 2, 2, "This is a comment on the 2nd post", "2022-01-05", 0, 0, vomitJSON(usrReact))
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = db.Exec(`INSERT INTO comments (id, post, author, text, date, likes, dislikes, user_reactions) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, 6, 2, 2, "This is a comment on the 2nd post", "2022-01-06", 0, 0, vomitJSON(usrReact))
	if err != nil {
		log.Println(err)
		return err
	}

	// Insert sample posts
	_, err = db.Exec(`INSERT INTO posts (id, title, content, author, date, likes, dislikes, user_reactions, comments_count, comments_ids, categories) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, 1, "First Post", "This is the first post", 1, "2022-01-01", 0, 0, uRsStr, 3, vomitJSON([]int{1, 2, 3}), vomitJSON([]string{"Facts", "Rumors"}))
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = db.Exec(`INSERT INTO posts (id, title, content, author, date, likes, dislikes, user_reactions, comments_count, comments_ids, categories) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, 2, "Second Post", "This is the second post", 1, "2022-01-02", 0, 0, uRsStr, 3, vomitJSON([]int{4, 5, 6}), vomitJSON([]string{"Facts", "Rumors"}))
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = db.Exec(`INSERT INTO posts (id, title, content, author, date, likes, dislikes, user_reactions, comments_count, comments_ids, categories) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, 3, "Third Post", "This is the third post", 1, "2022-01-03", 0, 0, uRsStr, 0, vomitJSON([]int{}), vomitJSON([]string{"Facts", "Rumors"}))
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = db.Exec(`INSERT INTO posts (id, title, content, author, date, likes, dislikes, user_reactions, comments_count, comments_ids, categories) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, 4, "Fourth Post", "This is the fourth post", 2, "2022-01-04", 0, 0, uRsStr, 0, vomitJSON([]int{}), vomitJSON([]string{"Facts", "Rumors"}))
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = db.Exec(`INSERT INTO posts (id, title, content, author, date, likes, dislikes, user_reactions, comments_count, comments_ids, categories) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, 5, "Fifth Post", "This is the fifth post", 2, "2022-01-05", 0, 0, uRsStr, 0, vomitJSON([]int{}), vomitJSON([]string{"Facts", "Rumors"}))
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = db.Exec(`INSERT INTO posts (id, title, content, author, date, likes, dislikes, user_reactions, comments_count, comments_ids, categories) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, 6, "Sixth Post", "This is the sixth post", 2, "2022-01-06", 0, 0, uRsStr, 0, vomitJSON([]int{}), vomitJSON([]string{"Facts", "Rumors"}))
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = db.Exec(`INSERT INTO posts (id, title, content, author, date, likes, dislikes, user_reactions, comments_count, comments_ids, categories) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, 7, "Seventh Post", "This is the seventh post", 3, "2022-01-07", 0, 0, uRsStr, 0, vomitJSON([]int{}), vomitJSON([]string{"Facts", "Rumors"}))
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = db.Exec(`INSERT INTO posts (id, title, content, author, date, likes, dislikes, user_reactions, comments_count, comments_ids, categories) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, 8, "Eighth Post", "This is the eighth post", 3, "2022-01-08", 0, 0, uRsStr, 0, vomitJSON([]int{}), vomitJSON([]string{"Facts", "Rumors"}))
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = db.Exec(`INSERT INTO posts (id, title, content, author, date, likes, dislikes, user_reactions, comments_count, comments_ids, categories) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, 9, "Ninth Post", "This is the ninth post", 3, "2022-01-09", 0, 0, uRsStr, 0, vomitJSON([]int{}), vomitJSON([]string{"Facts", "Rumors"}))
	_ = uRsStr
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// takes argument of any type and gives out the json string
func vomitJSON(myGod any) string {
	// convert to json
	json, err := json.Marshal(myGod)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(json)
}
