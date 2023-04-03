package server

import (
	"regexp"
	"strings"
)

// addSlashes escapes all slashes in a string ( / -> \/ )
func addSlashes(s string) string {
	return strings.ReplaceAll(s, "/", "\\/")
}

// pt function passes incoming pattern to regexp.MustCompile(addSlashes(pattern))
// The purpose of this function is to make it easier to read long regexps.
//
// Example:
//
//	pt(`^/api/posts/[[:digit:]]+/comment/[[:digit:]]+/like/?$`)
//
// is equivalent to
//
//	regexp.MustCompile(addSlashes(`^/api/posts/[[:digit:]]+/comment/[[:digit:]]+/like/?$`))
func pt(pattern string) *regexp.Regexp {
	return regexp.MustCompile(addSlashes(pattern))
}

// method GET endpoints

var reApiMe = pt(`^/api/me/?$`)
var reApiMePosts = pt(`^/api/me/posts/?$`)
var reApiMePostsLiked = pt(`^/api/me/posts/liked/?$`)

var reApiUserId = pt(`^/api/user/[[:digit:]]+/?$`)
var reApiUserIdPosts = pt(`^/api/user/[[:digit:]]+/posts/?$`)

var reApiPosts = pt(`^/api/posts/?$`)
var reApiPostsCategories = pt(`^/api/posts/categories/?$`)
var reApiPostsCategoriesName = pt(`^/api/posts/categories/[[:alnum:]]+/?$`)

var reApiPostsId = pt(`^/api/posts/[[:digit:]]+/?$`)

// method POST endpoints

var reApiPostsCreate = pt(`^/api/posts/create/?$`)

var reApiPostsIdLike = pt(`^/api/posts/[[:digit:]]+/like/?$`)
var reApiPostsIdDislike = pt(`^/api/posts/[[:digit:]]+/dislike/?$`)
var reApiPostsIdReaction = pt(`^/api/posts/[[:digit:]]+/reaction/?$`)

var reApiPostsIdCommentIdReaction = pt(`^/api/posts/[[:digit:]]+/comment/[[:digit:]]+/reaction/?$`)

var reApiPostsIdComment = pt(`^/api/posts/[[:digit:]]+/comment/?$`)

var reApiPostsIdComments = pt(`^/api/posts/[[:digit:]]+/comments/?$`)

var reApiPostsIdCommentIdLike = pt(`^/api/posts/[[:digit:]]+/comment/[[:digit:]]+/like/?$`)
var reApiPostsIdCommentIdDislike = pt(`^/api/posts/[[:digit:]]+/comment/[[:digit:]]+/dislike/?$`)

var reApiSignup = pt(`^/api/signup/?$`)
var reApiLogin = pt(`^/api/login/?$`)
var reApiLogout = pt(`^/api/logout/?$`)

var getRegexps = []string{
	reApiPosts.String(),
	reApiPostsId.String(),

	reApiUserId.String(),
	reApiUserIdPosts.String(),

	reApiMe.String(),
	reApiMePosts.String(),
	reApiMePostsLiked.String(),
	reApiPostsIdReaction.String(),

	reApiPostsCategories.String(),
	reApiPostsCategoriesName.String(),

	reApiPostsIdCommentIdReaction.String(),
	reApiPostsIdComments.String(),
}
var GetRegexp = regexp.MustCompile(strings.Join(getRegexps, "|"))

var postRegexps = []string{
	reApiPostsCreate.String(),

	reApiPostsIdLike.String(),
	reApiPostsIdDislike.String(),

	reApiPostsIdComment.String(),
	reApiPostsIdCommentIdLike.String(),
	reApiPostsIdCommentIdDislike.String(),

	reApiSignup.String(),
	reApiLogin.String(),
	reApiLogout.String(),
}
var PostRegexp = regexp.MustCompile(strings.Join(postRegexps, "|"))
