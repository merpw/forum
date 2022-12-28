package database

import (
	_ "github.com/mattn/go-sqlite3"
)

// InitDatabase initializes the database
func InitDatabase() error {
	// Open a connection to the database
	db := OpenDB()

	// Create the users table
	// _, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (		 // type User struct {
	// 	id INTEGER PRIMARY KEY,									 // 	Id       int    `json:"user_id"`
	// 	name TEXT,												 // 	Name     string `json:"user_name"`
	// 	email TEXT,												 // 	Email    string `json:"user_email"`
	// 	password TEXT,											 // 	Password string `json:"user_password"`
	// 	posts TEXT,												 // 	Posts    []int  `json:"user_posts"`
	// 	comments TEXT											 // 	Comments []int  `json:"user_comments"`
	// )`)
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

	// Create the Author table
	// _, err = db.Exec(`CREATE TABLE IF NOT EXISTS author (	 // type Author struct {
	// 	id INTEGER PRIMARY KEY, 								 // 	Id   int    `json:"user_id"`
	// 	name TEXT 												 // 	Name string `json:"user_name"`
	// )`)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS author (
		id INTEGER PRIMARY KEY,
		name TEXT
	)`)
	if err != nil {
		return err
	}

	// Create the posts table
	// _, err = db.Exec(`CREATE TABLE IF NOT EXISTS posts       // type Post struct {
	// (id INTEGER PRIMARY KEY,                                 // 	Id            int    `json:"id"`
	// title TEXT,                                              // 	Title         string `json:"title"`
	// content TEXT, 											// 	Content       string `json:"content"`
	// author INTEGER, 											// 	Author        Author `json:"author"`
	// date TEXT, 												// 	Date          string `json:"date"`
	// likes INTEGER, 											// 	Likes         int    `json:"likes"`
	// dislikes INTEGER, 										// 	Dislikes      int    `json:"dislikes"`
	// user_reactions TEXT, 									// 	UserReactions  []UserReaction    `json:"user_reactions"`
	// comments_count INTEGER, 									// 	CommentsCount int    `json:"comments_count"`
	// FOREIGN KEY(author) REFERENCES users(id))`)				// to link to the users table
	// Create the posts table
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
					FOREIGN KEY(author) REFERENCES users(id))`)
	if err != nil {
		return err
	}

	// Create the comments table
	// _, err = db.Exec(`CREATE TABLE IF NOT EXISTS comments    // type Comment struct {
	// (id INTEGER PRIMARY KEY, 				 			    // 	Id      int    `json:"id"`
	// 	post INTEGER, 											// 	Post    Post   `json:"post"`
	// 	author INTEGER, 										// 	Author  Author `json:"author"`
	// 	text TEXT, 												// 	Text    string `json:"text"`
	// 	date TEXT, 												// 	Date    string `json:"date"`
	// 	FOREIGN KEY(post) REFERENCES posts(id), 				// to link to the posts table
	// 	FOREIGN KEY(author) REFERENCES users(id))`)				// to link to the users table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS comments 
					(id INTEGER PRIMARY KEY, 
					post INTEGER, 
					author INTEGER, 
					text TEXT, 
					date TEXT, 
					FOREIGN KEY(post) REFERENCES posts(id), 
					FOREIGN KEY(author) REFERENCES users(id))`)
	if err != nil {
		return err
	}

	// if the tables are empty, insert sample data
	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		err = insertSampleData()
		if err != nil {
			return err
		}
	}

	return nil
}

// insert smaple data
func insertSampleData() error {
	db := OpenDB()

	// Insert sample users
	_, err := db.Exec(`INSERT INTO users (id, name, email, password, posts, comments) VALUES (?, ?, ?, ?, ?, ?)`, 1, "John Smith", "john@example.com", "password123", "[1, 2, 3]", "[1, 2, 3]")
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO users (id, name, email, password, posts, comments) VALUES (?, ?, ?, ?, ?, ?)`, 2, "Jane Doe", "jane@example.com", "password456", "[4, 5, 6]", "[4, 5, 6]")
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO users (id, name, email, password, posts, comments) VALUES (?, ?, ?, ?, ?, ?)`, 3, "Bob Johnson", "bob@example.com", "password789", "[7, 8, 9]", "[7, 8, 9]")
	if err != nil {
		return err
	}

	// Insert sample authors
	_, err = db.Exec(`INSERT INTO author (id, name) VALUES (?, ?)`, 1, "John Smith")
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO author (id, name) VALUES (?, ?)`, 2, "Jane Doe")
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO author (id, name) VALUES (?, ?)`, 3, "Bob Johnson")
	if err != nil {
		return err
	}

	// Insert sample posts
	_, err = db.Exec(`INSERT INTO posts (id, title, content, author, date, likes, dislikes, user_reactions, comments_count) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, 1, "First Post", "This is the first post", 1, "2022-01-01", 0, 0, "[]", 1)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO posts (id, title, content, author, date, likes, dislikes, user_reactions, comments_count) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, 2, "Second Post", "This is the second post", 1, "2022-01-02", 0, 0, "[]", 1)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO posts (id, title, content, author, date, likes, dislikes, user_reactions, comments_count) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, 3, "Third Post", "This is the third post", 1, "2022-01-03", 0, 0, "[]", 1)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO posts (id, title, content, author, date, likes, dislikes, user_reactions, comments_count) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, 4, "Fourth Post", "This is the fourth post", 2, "2022-01-04", 0, 0, "[]", 1)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO posts (id, title, content, author, date, likes, dislikes, user_reactions, comments_count) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, 5, "Fifth Post", "This is the fifth post", 2, "2022-01-05", 0, 0, "[]", 1)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO posts (id, title, content, author, date, likes, dislikes, user_reactions, comments_count) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, 6, "Sixth Post", "This is the sixth post", 2, "2022-01-06", 0, 0, "[]", 1)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO posts (id, title, content, author, date, likes, dislikes, user_reactions, comments_count) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, 7, "Seventh Post", "This is the seventh post", 3, "2022-01-07", 0, 0, "[]", 1)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO posts (id, title, content, author, date, likes, dislikes, user_reactions, comments_count) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, 8, "Eighth Post", "This is the eighth post", 3, "2022-01-08", 0, 0, "[]", 1)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO posts (id, title, content, author, date, likes, dislikes, user_reactions, comments_count) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, 9, "Ninth Post", "This is the ninth post", 3, "2022-01-09", 0, 0, "[]", 1)
	if err != nil {
		return err
	}

	// Insert sample comments
	_, err = db.Exec(`INSERT INTO comments (id, content, author, date, likes, dislikes, user_reactions, post_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, 1, "This is a comment on the first post", 1, "2022-01-01", 0, 0, "[]", 1)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO comments (id, content, author, date, likes, dislikes, user_reactions, post_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, 2, "This is a comment on the second post", 1, "2022-01-02", 0, 0, "[]", 2)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO comments (id, content, author, date, likes, dislikes, user_reactions, post_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, 3, "This is a comment on the third post", 1, "2022-01-03", 0, 0, "[]", 3)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO comments (id, content, author, date, likes, dislikes, user_reactions, post_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, 4, "This is a comment on the fourth post", 2, "2022-01-04", 0, 0, "[]", 4)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO comments (id, content, author, date, likes, dislikes, user_reactions, post_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, 5, "This is a comment on the fifth post", 2, "2022-01-05", 0, 0, "[]", 5)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO comments (id, content, author, date, likes, dislikes, user_reactions, post_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, 6, "This is a comment on the sixth post", 2, "2022-01-06", 0, 0, "[]", 6)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO comments (id, content, author, date, likes, dislikes, user_reactions, post_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, 7, "This is a comment on the seventh post", 3, "2022-01-07", 0, 0, "[]", 7)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO comments (id, content, author, date, likes, dislikes, user_reactions, post_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, 8, "This is a comment on the eighth post", 3, "2022-01-08", 0, 0, "[]", 8)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO comments (id, content, author, date, likes, dislikes, user_reactions, post_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, 9, "This is a comment on the ninth post", 3, "2022-01-09", 0, 0, "[]", 9)
	if err != nil {
		return err
	}
	return nil
}
