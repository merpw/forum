package database_test

import "testing"

func TestOpsComments(t *testing.T) {

	var authorId, postId int
	// add user to database, to add post authored by user
	t.Run(`
		Refactored tests which depend on each other,
	  so there's no need to split them to separate t.Runs.
	  If something will go wrong,
		panic trace will show which function caused it.
	`, func(t *testing.T) {
		authorId = srv.DB.AddUser("testuser", "user@email", "password")
		postId = srv.DB.AddPost("TEST POST TITLE", "test post content", authorId, "dummyCategory testCategory")
	})

	t.Run("AddComment", func(t *testing.T) {
		srv.DB.AddComment("test comment", postId, 1)
	})

	t.Run("UpdatePostsCommentsCount", func(t *testing.T) {
		srv.DB.UpdatePostsCommentsCount(postId, 1)
	})

}
