package server

import (
	"regexp"
	"strings"
)

// addSlashes escapes all slashes in a string ( / -> \/ )
func addSlashes(s string) string {
	return strings.Replace(s, "/", "\\/", -1)
}

// method get
var reApiMe = regexp.MustCompile(addSlashes(`^/api/me/?$`))
var reApiMePosts = regexp.MustCompile(addSlashes(`^/api/me/posts/?$`))

var reApiUserId = regexp.MustCompile(addSlashes(`^/api/user/[[:digit:]]+/?$`))
var reApiUserIdPosts = regexp.MustCompile(addSlashes(`^/api/user/[[:digit:]]+/posts/?$`))
var reApiPosts = regexp.MustCompile(addSlashes(`^/api/posts/?$`))
var reApiPostsCategories = regexp.MustCompile(addSlashes(`^/api/posts/categories/?$`))
var reApiPostsCategoriesFacts = regexp.MustCompile(addSlashes(`^/api/posts/categories/facts/?$`))
var reApiPostsCategoriesRumors = regexp.MustCompile(addSlashes(`^/api/posts/categories/rumors/?$`))
var reApiPostsId = regexp.MustCompile(addSlashes(`^/api/posts/[[:digit:]]+/?$`))

// method POST endpoints
var reApiPostsCreate = regexp.MustCompile(addSlashes(`^/api/posts/create/?$`))

var reApiPostsIdLike = regexp.MustCompile(addSlashes(`^/api/posts/[[:digit:]]+/like/?$`))
var reApiPostsIdDislike = regexp.MustCompile(addSlashes(`^/api/posts/[[:digit:]]+/dislike/?$`))
var reApiPostsIdReaction = regexp.MustCompile(addSlashes(`^/api/posts/[[:digit:]]+/reaction/?$`))

var reApiPostsIdComment = regexp.MustCompile(addSlashes(`^/api/posts/[[:digit:]]+/comment/?$`))

var reApiPostsIdCommentIdLike = regexp.MustCompile(addSlashes(`^/api/posts/[[:digit:]]+/comment/[[:digit:]]+/like/?$`))
var reApiPostsIdCommentIdDislike = regexp.MustCompile(addSlashes(`^/api/posts/[[:digit:]]+/comment/[[:digit:]]+/dislike/?$`))

var reApiSignup = regexp.MustCompile(addSlashes(`^/api/signup/?$`))
var reApiLogin = regexp.MustCompile(addSlashes(`^/api/login/?$`))
var reApiLogout = regexp.MustCompile(addSlashes(`^/api/logout/?$`))

var getRegexps = []string{
	reApiPosts.String(),
	reApiPostsId.String(),

	reApiUserId.String(),
	reApiUserIdPosts.String(),

	reApiMe.String(),
	reApiMePosts.String(),
	reApiPostsIdReaction.String(),

	reApiPostsCategories.String(),
	reApiPostsCategoriesFacts.String(),
	reApiPostsCategoriesRumors.String(),
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
