package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type SafeUser struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

var categories = []string{"facts", "rumors"}

func (srv *Server) apiPostsMasterHandler(w http.ResponseWriter, r *http.Request) {
	switch {

	case reApiPostsCategories.MatchString(r.URL.Path):
		srv.postsCategoriesHandler(w, r)

	case reApiPostsCategoriesName.MatchString(r.URL.Path):
		srv.apiPostsCategoriesNameHandler(w, r)

	case reApiPostsId.MatchString(r.URL.Path):
		srv.postsIdHandler(w, r)

	case reApiPostsCreate.MatchString(r.URL.Path):
		srv.postsCreateHandler(w, r)

	case reApiPostsIdLike.MatchString(r.URL.Path):
		srv.postsIdLikeHandler(w, r)

	case reApiPostsIdDislike.MatchString(r.URL.Path):
		srv.postsIdDislikeHandler(w, r)

	case reApiPostsIdReaction.MatchString(r.URL.Path):
		srv.postsIdReactionHandler(w, r)

	case reApiPostsIdCommentIdReaction.MatchString(r.URL.Path):
		srv.postsIdCommentIdReactionHandler(w, r)

	case reApiPostsIdComment.MatchString(r.URL.Path):
		srv.postsIdCommentHandler(w, r)

	case reApiPostsIdComments.MatchString(r.URL.Path):
		srv.postsIdCommentsHandler(w, r)

	case reApiPostsIdCommentIdLike.MatchString(r.URL.Path):
		srv.postsIdCommentIdLikeHandler(w, r)

	case reApiPostsIdCommentIdDislike.MatchString(r.URL.Path):
		srv.postsIdCommentIdDislikeHandler(w, r)
	default:
		http.NotFound(w, r)
	}
}

// postsCategoriesHandler returns a json list of all categories from the database
func (srv *Server) postsCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	sendObject(w, categories)
}

func (srv *Server) apiPostsCategoriesNameHandler(w http.ResponseWriter, r *http.Request) {
	categoryName := strings.TrimPrefix(r.URL.Path, "/api/posts/categories/")
	// /api/posts/categories/name -> name

	categoryName = strings.ToLower(categoryName)
	// Name -> name

	isValid := false
	for _, cat := range categories {
		if cat == categoryName {
			isValid = true
			break
		}
	}

	if !isValid {
		errorResponse(w, http.StatusNotFound)
		return
	}

	posts := srv.DB.GetCategoryPosts(categoryName)
	type ResponsePost struct {
		Id            int      `json:"id"`
		Title         string   `json:"title"`
		Content       string   `json:"content"`
		Author        SafeUser `json:"author"`
		Date          string   `json:"date"`
		CommentsCount int      `json:"comments_count"`
		LikesCount    int      `json:"likes_count"`
		Category      string   `json:"category"`
	}

	response := make([]ResponsePost, 0)
	for _, post := range posts {
		postAuthor := srv.DB.GetUserById(post.AuthorId)
		response = append(response, ResponsePost{
			Id:            post.Id,
			Title:         post.Title,
			Content:       post.Content,
			Date:          post.Date,
			Author:        SafeUser{Id: postAuthor.Id, Name: postAuthor.Name},
			CommentsCount: post.CommentsCount,
			LikesCount:    post.LikesCount,
			Category:      post.Category,
		})
	}

	sendObject(w, response)
}

// postsHandler returns a json list of all posts from the database
func (srv *Server) postsHandler(w http.ResponseWriter, r *http.Request) {
	posts := srv.DB.GetAllPosts()
	type ResponsePost struct {
		Id            int      `json:"id"`
		Title         string   `json:"title"`
		Content       string   `json:"content"`
		Author        SafeUser `json:"author"`
		Date          string   `json:"date"`
		CommentsCount int      `json:"comments_count"`
		LikesCount    int      `json:"likes_count"`
		Category      string   `json:"category"`
	}

	response := make([]ResponsePost, 0)
	for _, post := range posts {
		postAuthor := srv.DB.GetUserById(post.AuthorId)
		response = append(response, ResponsePost{
			Id:            post.Id,
			Title:         post.Title,
			Content:       post.Content,
			Date:          post.Date,
			Author:        SafeUser{Id: postAuthor.Id, Name: postAuthor.Name},
			CommentsCount: post.CommentsCount,
			LikesCount:    post.LikesCount,
			Category:      post.Category,
		})
	}

	sendObject(w, response)
}

// postsIdHandler returns a single post from the database that matches the incoming id of the post in the url
func (srv *Server) postsIdHandler(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Id            int      `json:"id"`
		Title         string   `json:"title"`
		Content       string   `json:"content"`
		Author        SafeUser `json:"author"`
		Date          string   `json:"date"`
		CommentsCount int      `json:"comments_count"`
		Comments      []struct {
			Id      int      `json:"id"`
			Content string   `json:"content"`
			Author  SafeUser `json:"author"`
		} `json:"comments"`
		LikesCount int    `json:"likes_count"`
		Category   string `json:"category"`
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse(w, http.StatusNotFound)
		return
	}

	// Get the post from the database
	post := srv.DB.GetPostById(id)

	if post == nil {
		errorResponse(w, http.StatusNotFound)
		return
	}

	postAuthor := srv.DB.GetUserById(post.AuthorId)

	sendObject(w, Response{
		Id:            post.Id,
		Title:         post.Title,
		Content:       post.Content,
		Author:        SafeUser{Id: postAuthor.Id, Name: postAuthor.Name},
		Date:          post.Date,
		CommentsCount: post.CommentsCount,
		Comments: []struct {
			Id      int      `json:"id"`
			Content string   `json:"content"`
			Author  SafeUser `json:"author"`
		}{},
		LikesCount: post.LikesCount,
		Category:   post.Category,
	},
	)
}

// postsCreateHandler creates a new post in the database
func (srv *Server) postsCreateHandler(w http.ResponseWriter, r *http.Request) {
	userId := srv.getUserId(w, r)
	if userId == -1 {
		errorResponse(w, http.StatusUnauthorized)
		return
	}

	requestBody := struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		Category string `json:"category"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Body is not valid", http.StatusBadRequest)
		return
	}

	requestBody.Title = strings.TrimSpace(requestBody.Title)
	requestBody.Content = strings.TrimSpace(requestBody.Content)

	if len(requestBody.Title) < 1 {
		http.Error(w, "Title is too short", http.StatusBadRequest)
		return
	}
	if len(requestBody.Content) < 1 {
		http.Error(w, "Content is too short", http.StatusBadRequest)
		return
	}

	if len(requestBody.Title) > 25 {
		http.Error(w, "Title is too long, maximum length is 25", http.StatusBadRequest)
		return
	}

	requestBody.Category = strings.TrimSpace(requestBody.Category)
	requestBody.Category = strings.ToLower(requestBody.Category)

	isValid := false
	for _, cat := range categories {
		if cat == requestBody.Category {
			isValid = true
			break
		}
	}

	if !isValid {
		http.Error(w, "Category is not valid", http.StatusBadRequest)
		return
	}

	id := srv.DB.AddPost(requestBody.Title, requestBody.Content, userId, requestBody.Category)
	sendObject(w, id)
}

// postsIdLikeHandler likes a post in the database
func (srv *Server) postsIdLikeHandler(w http.ResponseWriter, r *http.Request) {
	userId := srv.getUserId(w, r)
	if userId == -1 {
		errorResponse(w, http.StatusUnauthorized)
		return
	}
	postIdStr := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	postIdStr = strings.TrimSuffix(postIdStr, "/like")
	// /api/posts/1/like -> 1

	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		errorResponse(w, http.StatusNotFound)
	}

	post := srv.DB.GetPostById(postId)
	if post == nil {
		errorResponse(w, http.StatusNotFound)
		return
	}

	reaction := srv.DB.GetPostReaction(postId, userId)

	switch reaction {
	case 0: // if not reacted, add like
		srv.DB.AddPostReaction(postId, userId, 1)
		srv.DB.UpdatePostLikesCount(postId, +1)

		sendObject(w, +1)

	case 1: // if already liked, unlike
		srv.DB.RemovePostReaction(postId, userId)
		srv.DB.UpdatePostLikesCount(postId, -1)

		sendObject(w, 0)

	case -1: // if disliked, remove dislike and add like
		srv.DB.RemovePostReaction(postId, userId)
		srv.DB.UpdatePostDislikeCount(postId, -1)

		srv.DB.AddPostReaction(postId, userId, 1)
		srv.DB.UpdatePostLikesCount(postId, +1)

		sendObject(w, 1)
	}
}

func (srv *Server) postsIdReactionHandler(w http.ResponseWriter, r *http.Request) {
	userId := srv.getUserId(w, r)
	if userId == -1 {
		errorResponse(w, http.StatusUnauthorized)
		return
	}

	postIdStr := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	postIdStr = strings.TrimSuffix(postIdStr, "/reaction")
	// /api/posts/1/reaction -> 1

	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		errorResponse(w, http.StatusNotFound)
		return
	}

	post := srv.DB.GetPostById(postId)
	if post == nil {
		errorResponse(w, http.StatusNotFound)
		return
	}

	reaction := srv.DB.GetPostReaction(postId, userId)
	if userId == post.AuthorId {
		sendObject(w, struct {
			Reaction      int `json:"reaction"`
			LikesCount    int `json:"likes_count"`
			DislikesCount int `json:"dislikes_count"`
		}{
			Reaction:      reaction,
			LikesCount:    post.LikesCount,
			DislikesCount: post.DislikesCount,
		})
	} else {
		sendObject(w, struct {
			Reaction   int `json:"reaction"`
			LikesCount int `json:"likes_count"`
		}{
			Reaction:   reaction,
			LikesCount: post.LikesCount,
		})
	}
}

// postsPostsIdDislikeHandler dislikes a post in the database
func (srv *Server) postsIdDislikeHandler(w http.ResponseWriter, r *http.Request) {
	userId := srv.getUserId(w, r)
	if userId == -1 {
		errorResponse(w, http.StatusUnauthorized)
		return
	}
	postIdStr := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	postIdStr = strings.TrimSuffix(postIdStr, "/dislike")
	// /api/posts/1/dislike -> 1

	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		errorResponse(w, http.StatusNotFound)
	}

	post := srv.DB.GetPostById(postId)
	if post == nil {
		errorResponse(w, http.StatusNotFound)
		return
	}

	reaction := srv.DB.GetPostReaction(postId, userId)

	switch reaction {
	case 0: // if not reacted, add dislike
		srv.DB.AddPostReaction(postId, userId, -1)
		srv.DB.UpdatePostDislikeCount(postId, +1)

		sendObject(w, -1)

	case -1: // if already disliked, remove dislike
		srv.DB.RemovePostReaction(postId, userId)
		srv.DB.UpdatePostDislikeCount(postId, -1)

		sendObject(w, 0)

	case 1: // if liked, remove like and add dislike
		srv.DB.RemovePostReaction(postId, userId)
		srv.DB.UpdatePostLikesCount(postId, -1)

		srv.DB.AddPostReaction(postId, userId, -1)
		srv.DB.UpdatePostDislikeCount(postId, +1)

		sendObject(w, -1)
	}
}

// getPostId returns -1 if post was not found in database
func (srv *Server) getPostId(w http.ResponseWriter, r *http.Request) int {
	postId, err := strconv.Atoi(strings.Split(r.URL.Path, "/")[3])
	if err != nil {
		return -1
	}
	// Get the post from the database
	post := srv.DB.GetPostById(postId)
	if post == nil {
		return -1
	}
	return postId
}

// postsIdCommentHandler comments on a post in the database
func (srv *Server) postsIdCommentHandler(w http.ResponseWriter, r *http.Request) {
	userId := srv.getUserId(w, r)
	if userId == -1 {
		errorResponse(w, http.StatusUnauthorized)
		return
	}

	postId := srv.getPostId(w, r)
	if postId == -1 {
		errorResponse(w, http.StatusNotFound)
		return
	}

	requestBody := struct {
		Content string `json:"content"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Body is not valid", http.StatusBadRequest)
		return
	}

	requestBody.Content = strings.TrimSpace(requestBody.Content)

	if len(requestBody.Content) < 1 {
		http.Error(w, "Content is too short", http.StatusBadRequest)
		return
	}

	id := srv.DB.AddComment(requestBody.Content, postId, userId)
	sendObject(w, id)
}

func (srv *Server) postsIdCommentsHandler(w http.ResponseWriter, r *http.Request) {
	// userIdStr := strings.TrimPrefix(r.URL.Path, "/api/user/")
	// userIdStr = strings.TrimSuffix(userIdStr, "/posts")
	// /api/user/1/posts -> 1

	postId, err := strconv.Atoi(strings.Split(r.URL.Path, "/")[3])
	if err != nil {
		errorResponse(w, http.StatusNotFound)
		return
	}

	post := srv.DB.GetPostById(postId)
	if post == nil {
		errorResponse(w, http.StatusNotFound)
		return
	}

	type ResponseComment struct {
		Id            int      `json:"id"`
		PostId        int      `json:"post_id"`
		AuthorId      int      `json:"author_id"`
		Content       string   `json:"content"`
		Author        SafeUser `json:"author"` // TODO: maybe remove
		Date          string   `json:"date"`
		LikesCount    int      `json:"likes_count"`
		DislikesCount int      `json:"dislikes_count"`
	}

	// posts := srv.DB.GetUserPosts(userId)
	comments := srv.DB.GetPostComments(postId)

	response := make([]ResponseComment, 0)
	for _, comment := range comments {
		user := srv.DB.GetUserById(comment.AuthorId)
		response = append(response, ResponseComment{
			Id:            comment.Id,
			PostId:        postId,
			AuthorId:      comment.AuthorId,
			Content:       comment.Content,
			Author:        SafeUser{user.Id, user.Name},
			Date:          comment.Date,
			LikesCount:    comment.LikesCount,
			DislikesCount: comment.DislikesCount,
		})
	}

	sendObject(w, response)
}

// postsIdCommentIdLikeHandler likes a comment on a post in the database
func (srv *Server) postsIdCommentIdLikeHandler(w http.ResponseWriter, r *http.Request) {

	userId := srv.getUserId(w, r)
	if userId == -1 {
		errorResponse(w, http.StatusUnauthorized)
		return
	}

	commentId, err := strconv.Atoi(strings.Split(r.URL.Path, "/")[5])
	// /api/posts/1/comment/2/like ->2
	if err != nil {
		errorResponse(w, http.StatusNotFound)
	}

	comment := srv.DB.GetCommentById(commentId)
	if comment == nil {
		errorResponse(w, http.StatusNotFound)
		return
	}

	reaction := srv.DB.GetCommentReaction(commentId, userId)

	switch reaction {
	case 0: // if not reacted, add like
		srv.DB.AddCommentReaction(commentId, userId, 1)
		srv.DB.UpdateCommentLikesCount(commentId, +1)

		sendObject(w, +1)

	case 1: // if already liked, unlike
		srv.DB.RemoveCommentReaction(commentId, userId)
		srv.DB.UpdateCommentLikesCount(commentId, -1)

		sendObject(w, 0)

	case -1: // if disliked, remove dislike and add like
		srv.DB.RemoveCommentReaction(commentId, userId)
		srv.DB.UpdateCommentDislikeCount(commentId, -1)

		srv.DB.AddCommentReaction(commentId, userId, 1)
		srv.DB.UpdateCommentLikesCount(commentId, +1)

		sendObject(w, 1)
	}
}

// postsIdCommentIdDislikeHandler dislikes a comment on a post in the database
func (srv *Server) postsIdCommentIdDislikeHandler(w http.ResponseWriter, r *http.Request) {
	userId := srv.getUserId(w, r)
	if userId == -1 {
		errorResponse(w, http.StatusUnauthorized)
		return
	}

	commentId, err := strconv.Atoi(strings.Split(r.URL.Path, "/")[5])
	if err != nil {
		errorResponse(w, http.StatusNotFound)
	}

	comment := srv.DB.GetCommentById(commentId)
	if comment == nil {
		errorResponse(w, http.StatusNotFound)
		return
	}

	reaction := srv.DB.GetCommentReaction(commentId, userId)

	switch reaction {
	case 0: // if not reacted, add dislike
		srv.DB.AddCommentReaction(commentId, userId, -1)
		srv.DB.UpdateCommentDislikeCount(commentId, +1)

		sendObject(w, -1)

	case -1: // if already disliked, remove dislike
		srv.DB.RemoveCommentReaction(commentId, userId)
		srv.DB.UpdateCommentDislikeCount(commentId, -1)

		sendObject(w, 0)

	case 1: // if liked, remove like and add dislike
		srv.DB.RemoveCommentReaction(commentId, userId)
		srv.DB.UpdateCommentLikesCount(commentId, -1)

		srv.DB.AddCommentReaction(commentId, userId, -1)
		srv.DB.UpdateCommentDislikeCount(commentId, +1)

		sendObject(w, -1)
	}
}

func (srv *Server) postsIdCommentIdReactionHandler(w http.ResponseWriter, r *http.Request) {
	userId := srv.getUserId(w, r)
	if userId == -1 {
		errorResponse(w, http.StatusUnauthorized)
		return
	}

	commentId, err := strconv.Atoi(strings.Split(r.URL.Path, "/")[5])
	if err != nil {
		errorResponse(w, http.StatusNotFound)
		return
	}

	comment := srv.DB.GetCommentById(commentId)
	if comment == nil {
		errorResponse(w, http.StatusNotFound)
		return
	}

	reaction := srv.DB.GetCommentReaction(commentId, userId)
	if userId == comment.AuthorId {
		sendObject(w, struct {
			Reaction      int `json:"reaction"`
			LikesCount    int `json:"likes_count"`
			DislikesCount int `json:"dislikes_count"`
		}{
			Reaction:      reaction,
			LikesCount:    comment.LikesCount,
			DislikesCount: comment.DislikesCount,
		})
	} else {
		sendObject(w, struct {
			Reaction   int `json:"reaction"`
			LikesCount int `json:"likes_count"`
		}{
			Reaction:   reaction,
			LikesCount: comment.LikesCount,
		})
	}
}
