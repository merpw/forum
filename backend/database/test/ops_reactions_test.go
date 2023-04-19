package database_test

import "testing"

func TestOpsReactions(t *testing.T) {
	// add user, to be able to add post
	var userId int
	t.Run("AddUserToTestOpsReactions", func(t *testing.T) {
		userId = srv.DB.AddUser("OpsReactionTestUser", "opsreactiontest@email", "password")
	})

	// add post, to be able to add reaction and comment
	var postId int
	t.Run("AddPostToTestOpsReactions", func(t *testing.T) {
		postId = srv.DB.AddPost("Post Title", "post content", userId, "OpsReactionCategory testCategory")
	})

	// add comment, to be able to add reaction
	var commentId int
	t.Run("AddCommentToTestOpsReactions", func(t *testing.T) {
		commentId = srv.DB.AddComment("Comment content", postId, userId)
	})

	/*********************************************/
	/*********** Test reactions to post **********/
	/*********************************************/

	// add reaction to post
	t.Run("AddPostReactionToTestOpsReactions", func(t *testing.T) {
		srv.DB.AddPostReaction(postId, userId, 1)
	})

	// update post likes count
	t.Run("UpdatePostLikesCountToTestOpsReactions", func(t *testing.T) {
		srv.DB.UpdatePostLikesCount(postId, 1)
	})

	// update post dislikes count
	t.Run("UpdatePostDislikeCountToTestOpsReactions", func(t *testing.T) {
		srv.DB.UpdatePostDislikeCount(postId, 1)
	})

	// get post reactions
	t.Run("GetPostReactionsToTestOpsReactions", func(t *testing.T) {
		reaction := srv.DB.GetPostReaction(postId, userId)
		if reaction != 1 {
			t.Fatalf("Expected reaction 1, got %d", reaction)
		}
	})

	// remove post reaction
	t.Run("RemovePostReactionToTestOpsReactions", func(t *testing.T) {
		srv.DB.RemovePostReaction(postId, userId)
	})

	/***********************************************/
	/*********** Test reactions to comment *********/
	/***********************************************/

	// add reaction to comment
	t.Run("AddCommentReactionToTestOpsReactions", func(t *testing.T) {
		srv.DB.AddCommentReaction(commentId, userId, -1)
	})

	// update comment likes count
	t.Run("UpdateCommentLikesCountToTestOpsReactions", func(t *testing.T) {
		srv.DB.UpdateCommentLikesCount(commentId, 1)
	})

	// update comment dislikes count
	t.Run("UpdateCommentDislikeCountToTestOpsReactions", func(t *testing.T) {
		srv.DB.UpdateCommentDislikeCount(commentId, 1)
	})

	// get comment reactions
	t.Run("GetCommentReactionsToTestOpsReactions", func(t *testing.T) {
		reaction := srv.DB.GetCommentReaction(commentId, userId)
		if reaction != -1 {
			t.Fatalf("Expected reaction -1, got %d", reaction)
		}
	})

	// remove comment reaction
	t.Run("RemoveCommentReactionToTestOpsReactions", func(t *testing.T) {
		srv.DB.RemoveCommentReaction(commentId, userId)
	})

}
