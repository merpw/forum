package database_test

import "testing"

func TestOpsReactions(t *testing.T) {
	// add user, to be able to add post
	var userId int
	t.Run("AddUser", func(t *testing.T) {
		userId = srv.DB.AddUser("OpsReactionTestUser", "opsreactiontest@email", "password")
	})

	// add post, to be able to add reaction and comment
	var postId int
	t.Run("AddPost", func(t *testing.T) {
		postId = srv.DB.AddPost("Post Title", "post content", userId, "OpsReactionCategory testCategory")
	})

	// add comment, to be able to add reaction
	var commentId int
	t.Run("AddComment", func(t *testing.T) {
		commentId = srv.DB.AddComment("Comment content", postId, userId)
	})

	// add reaction to post
	t.Run("AddPostReaction", func(t *testing.T) {
		srv.DB.AddPostReaction(postId, userId, 1)
	})

	// update post likes count
	t.Run("UpdatePostLikesCount", func(t *testing.T) {
		srv.DB.UpdatePostLikesCount(postId, 1)
	})

	// update post dislikes count
	t.Run("UpdatePostDislikeCount", func(t *testing.T) {
		srv.DB.UpdatePostDislikeCount(postId, 1)
	})

	// get post reactions
	t.Run("GetPostReactions", func(t *testing.T) {
		reaction := srv.DB.GetPostReaction(postId, userId)
		if reaction != 1 {
			t.Fatalf("Expected reaction 1, got %d", reaction)
		}
	})

	// remove post reaction
	t.Run("RemovePostReaction", func(t *testing.T) {
		srv.DB.RemovePostReaction(postId, userId)
	})

	// add reaction to comment
	t.Run("AddCommentReaction", func(t *testing.T) {
		srv.DB.AddCommentReaction(commentId, userId, -1)
	})

	// update comment likes count
	t.Run("UpdateCommentLikesCount", func(t *testing.T) {
		srv.DB.UpdateCommentLikesCount(commentId, 1)
	})

	// update comment dislikes count
	t.Run("UpdateCommentDislikeCount", func(t *testing.T) {
		srv.DB.UpdateCommentDislikeCount(commentId, 1)
	})

	// get comment reactions
	t.Run("GetCommentReactions", func(t *testing.T) {
		reaction := srv.DB.GetCommentReaction(commentId, userId)
		if reaction != -1 {
			t.Fatalf("Expected reaction -1, got %d", reaction)
		}
	})

	// remove comment reaction
	t.Run("RemoveCommentReaction", func(t *testing.T) {
		srv.DB.RemoveCommentReaction(commentId, userId)
	})

}
