package database_test

import "testing"

func TestOpsComments(t *testing.T) {

	var authorId, postId, commentId int
	// add user to database, to add post authored by user
	t.Run("AddPost", func(t *testing.T) {
		authorId = DB.AddUser("testuser", "user@email", "password")
		postId = DB.AddPost("TEST POST TITLE", "test post content", authorId, "dummyCategory testCategory")
	})

	t.Run("AddComment", func(t *testing.T) {
		commentId = DB.AddComment("test comment", postId, 1)
		DB.UpdatePostsCommentsCount(postId, 1)
		comments := DB.GetPostComments(postId)
		if len(comments) != 1 {
			t.Fatalf("Expected 1 comment, got %d", len(comments))
		}
		if comments[0].Id != commentId {
			t.Fatalf("Expected comment id %d, got %d", commentId, comments[0].Id)
		}
	})

}
