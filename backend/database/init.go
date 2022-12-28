package database

import (
	_ "github.com/mattn/go-sqlite3"
)

// InitDatabase initializes the database
func InitDatabase() error {
	// Open a connection to the database
	db := Opendb()

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

	return nil
}
