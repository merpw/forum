package database_test

import "testing"

func TestOpsReactions(t *testing.T) {
	var userId, postId, commentId int
	t.Run("TestPostReaction", func(t *testing.T) {
		userId = DB.AddUser("OpsReactionTestUser", "opsreactiontest@email", "password")
		postId = DB.AddPost("Post Title", "post content", userId, "OpsReactionCategory testCategory")
		commentId = DB.AddComment("Comment content", postId, userId)
		DB.AddPostReaction(postId, userId, 1)
		reaction := DB.GetPostReaction(postId, userId)
		if reaction != 1 {
			t.Fatalf("Expected reaction 1, got %d", reaction)
		}
		DB.RemovePostReaction(postId, userId)
		reaction = DB.GetPostReaction(postId, userId)
		if reaction != 0 {
			t.Fatalf("Expected reaction 0, got %d", reaction)
		}
	})

	t.Run("TestCommentReaction", func(t *testing.T) {
		DB.AddCommentReaction(commentId, userId, -1)
		reaction := DB.GetCommentReaction(commentId, userId)
		if reaction != -1 {
			t.Fatalf("Expected reaction -1, got %d", reaction)
		}
		DB.RemoveCommentReaction(commentId, userId)
		reaction = DB.GetCommentReaction(commentId, userId)
		if reaction != 0 {
			t.Fatalf("Expected reaction 0, got %d", reaction)
		}
	})

	// update post likes count
	t.Run("UpdatePostLikesCount", func(t *testing.T) {
		DB.UpdatePostLikesCount(postId, 1)
	})

	// update post dislikes count
	t.Run("UpdatePostDislikeCount", func(t *testing.T) {
		DB.UpdatePostDislikeCount(postId, 1)
	})

	// update comment likes count
	t.Run("UpdateCommentLikesCount", func(t *testing.T) {
		DB.UpdateCommentLikesCount(commentId, 1)
	})

	// update comment dislikes count
	t.Run("UpdateCommentDislikeCount", func(t *testing.T) {
		DB.UpdateCommentDislikeCount(commentId, 1)
	})

}
