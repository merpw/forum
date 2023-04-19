package database_test

import "testing"

func TestOpsPosts(t *testing.T) {
	var userId int
	// Add user, to be able to add posts

	t.Run("AddUserToTestOpsPosts", func(t *testing.T) {
		userId = srv.DB.AddUser("OpsPostTestUser", "opsposttest@email", "password")
	})

	// Add two posts

	var postId1, postId2 int
	t.Run("AddPost1ToTestOpsPosts", func(t *testing.T) {
		postId1 = srv.DB.AddPost("Post Title 1", "post content 1", userId, "OpsPostCategory testCategory")
	})
	t.Run("AddPost2ToTestOpsPosts", func(t *testing.T) {
		postId2 = srv.DB.AddPost("Post Title 2", "post content 2", userId, "OpsPostCategory testCategory")
	})

	// Get all posts

	t.Run("GetAllPostsToTestOpsPosts", func(t *testing.T) {
		posts := srv.DB.GetAllPosts()
		if len(posts) < 2 { // 2 posts added at least, but maybe more from other tests
			t.Errorf("Expected 2 posts, got %d", len(posts))
		}
	})

	// Get post by id

	t.Run("GetPostByIdToTestOpsPosts", func(t *testing.T) {
		post := srv.DB.GetPostById(postId1)
		if post == nil {
			t.Errorf("Expected post with id %d, got nil", postId1)
		}
		if post.Id != postId1 {
			t.Errorf("Expected post with id %d, got %d", postId1, post.Id)
		}
		noPost := srv.DB.GetPostById(-1)
		if noPost != nil {
			t.Errorf("Expected nil, got post with id %d", noPost.Id)
		}
	})

	// Get user posts

	t.Run("GetUserPostsToTestOpsPosts", func(t *testing.T) {
		posts := srv.DB.GetUserPosts(userId)
		if len(posts) != 2 {
			t.Errorf("Expected 2 posts, got %d", len(posts))
		}
		// check user posts ids
		if posts[0].Id != postId1 && posts[0].Id != postId2 {
			t.Errorf("Expected post id %d or %d, got %d", postId1, postId2, posts[0].Id)
		}
		if posts[1].Id != postId1 && posts[1].Id != postId2 {
			t.Errorf("Expected post id %d or %d, got %d", postId1, postId2, posts[1].Id)
		}
	})

	// Like post to get it as liked by user

	t.Run("LikePostToTestOpsPosts", func(t *testing.T) {
		srv.DB.AddPostReaction(postId1, userId, 1)
	})

	// Get user liked posts

	t.Run("GetUserLikedPostsToTestOpsPosts", func(t *testing.T) {
		posts := srv.DB.GetUserPostsLiked(userId)
		if len(posts) != 1 {
			t.Errorf("Expected 1 post, got %d", len(posts))
		}
		if posts[0].Id != postId1 {
			t.Errorf("Expected post id %d, got %d", postId1, posts[0].Id)
		}
	})

	// Get category posts

	t.Run("GetCategoryPostsToTestOpsPosts", func(t *testing.T) {
		posts := srv.DB.GetCategoryPosts("OpsPostCategory testCategory")
		if len(posts) != 2 {
			t.Errorf("Expected 2 posts, got %d", len(posts))
		}
		// check category posts ids
		if posts[0].Id != postId1 && posts[0].Id != postId2 {
			t.Errorf("Expected post id %d or %d, got %d", postId1, postId2, posts[0].Id)
		}
		if posts[1].Id != postId1 && posts[1].Id != postId2 {
			t.Errorf("Expected post id %d or %d, got %d", postId1, postId2, posts[1].Id)
		}
	})

	// add two comment to post, to get them as comments of post

	var commentId1, commentId2 int
	t.Run("AddComment1ToTestOpsPosts", func(t *testing.T) {
		commentId1 = srv.DB.AddComment("Comment content 1", postId1, userId)
	})
	t.Run("AddComment2ToTestOpsPosts", func(t *testing.T) {
		commentId2 = srv.DB.AddComment("Comment content 2", postId1, userId)
	})

	// Get post comments

	t.Run("GetPostCommentsToTestOpsPosts", func(t *testing.T) {
		comments := srv.DB.GetPostComments(postId1)
		if len(comments) != 2 {
			t.Errorf("Expected 2 comments, got %d", len(comments))
		}
		// check comments ids
		if comments[0].Id != commentId1 && comments[0].Id != commentId2 {
			t.Errorf("Expected comment id %d or %d, got %d", commentId1, commentId2, comments[0].Id)
		}
		if comments[1].Id != commentId1 && comments[1].Id != commentId2 {
			t.Errorf("Expected comment id %d or %d, got %d", commentId1, commentId2, comments[1].Id)
		}
	})

	// Get comment by id

	t.Run("GetCommentByIdToTestOpsPosts", func(t *testing.T) {
		comment := srv.DB.GetCommentById(commentId1)
		if comment == nil {
			t.Errorf("Expected comment with id %d, got nil", commentId1)
		}
		if comment.Id != commentId1 {
			t.Errorf("Expected comment with id %d, got %d", commentId1, comment.Id)
		}
		noComment := srv.DB.GetCommentById(-1)
		if noComment != nil {
			t.Errorf("Expected nil, got comment with id %d", noComment.Id)
		}
	})
}
