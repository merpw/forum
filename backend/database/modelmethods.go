package database

// type Post struct {
// 	Id            int    `json:"id"`
// 	Title         string `json:"title"`
// 	Content       string `json:"content"`
// 	Author        Author `json:"author"`
// 	Date          string `json:"date"`
// 	Likes         int    `json:"likes"`
// 	Dislikes      int    `json:"dislikes"`
// 	UserReaction  int    `json:"user_reaction"`
// 	CommentsCount int    `json:"comments_count"`
// }

// Add a method Add to the Post struct, that adds a post to the database
func (p *Post) Add() error {
	// Open a connection to the database
	// db := Opendb()

	// // Prepare the INSERT statement
	// stmt, err := db.Prepare("INSERT INTO posts (id, title, content, author, date, likes, dislikes, user_reaction, comments_count) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	// if err != nil {
	// 	return err
	// }
	// defer stmt.Finalize()

	// // Bind the values of the fields of the Post struct to the placeholders in the INSERT statement
	// _, err = stmt.Bind(p.Id, p.Title, p.Content, p.Author, p.Date, p.Likes, p.Dislikes, p.UserReaction, p.CommentsCount)
	// if err != nil {
	// 	return err
	// }

	// // Execute the INSERT statement
	// _, err = stmt.Step()
	// if err != nil {
	// 	return err
	// }

	return nil
}

// Add a method Get to the Post struct, that gets a post from the database
func (p *Post) Get() error {
	//TODO: implement
	return nil
}
