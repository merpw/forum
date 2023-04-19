package database_test

import "testing"

func TestOpsComments(t *testing.T) {

	var authorId, postId int
	// add user to database, to add post authored by user
	t.Run("AddUser", func(t *testing.T) {
		authorId = srv.DB.AddUser("testuser", "user@email", "password")
	})
	// add post to database, to add comment to it
	t.Run("AddPost", func(t *testing.T) {
		postId = srv.DB.AddPost("TEST POST TITLE", "test post content", authorId, "dummyCategory testCategory")
	})

	t.Run("AddComment", func(t *testing.T) {
		srv.DB.AddComment("test comment", postId, 1)
	})

	t.Run("UpdatePostsCommentsCount", func(t *testing.T) {
		srv.DB.UpdatePostsCommentsCount(postId, 1)
	})

}
