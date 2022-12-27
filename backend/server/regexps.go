package server

import (
	"regexp"
	"strings"
)

// method get
var reApiUser = regexp.MustCompile(`^\/api\/user\/?$`)
var reApiUserLikedPosts = regexp.MustCompile(`^\/api\/user\/liked\/posts\/?$`)
var reApiUserId = regexp.MustCompile(`^\/api\/user\/[[:digit:]]+\/?$`)
var reApiUserIdPosts = regexp.MustCompile(`^\/api\/user\/[[:digit:]]+\/posts\/?$`)
var reApiPosts = regexp.MustCompile(`^\/api\/posts\/?$`)
var reApiPostsCategories = regexp.MustCompile(`^\/api\/posts\/categories\/?$`)
var reApiPostsCategoriesFacts = regexp.MustCompile(`^\/api\/posts\/categories\/facts\/?$`)
var reApiPostsCategoriesRumors = regexp.MustCompile(`^\/api\/posts\/categories\/rumors\/?$`)
var reApiPostsId = regexp.MustCompile(`^\/api\/posts\/[[:digit:]]+\/?$`)

// method post
var reApiPostsCreate = regexp.MustCompile(`^\/api\/posts\/create\/?$`)
var reApiPostsIdLike = regexp.MustCompile(`^\/api\/posts\/[[:digit:]]+\/like\/?$`)
var reApiPostsIdDislike = regexp.MustCompile(`^\/api\/posts\/[[:digit:]]+\/dislike\/?$`)
var reApiPostsIdComment = regexp.MustCompile(`^\/api\/posts\/[[:digit:]]+\/comment\/?$`)
var reApiPostsIdCommentIdLike = regexp.MustCompile(`^\/api\/posts\/[[:digit:]]+\/comment\/[[:digit:]]+\/like\/?$`)
var reApiPostsIdCommentIdDislike = regexp.MustCompile(`^\/api\/posts\/[[:digit:]]+\/comment\/[[:digit:]]+\/dislike\/?$`)
var reApiSignup = regexp.MustCompile(`^\/api\/signup\/?$`)
var reApiLogin = regexp.MustCompile(`^\/api\/login\/?$`)
var reApiLogout = regexp.MustCompile(`^\/api\/logout\/?$`)

var getRegexps = []string{
	reApiUser.String(), reApiUserLikedPosts.String(), reApiUserId.String(), reApiUserIdPosts.String(), reApiPosts.String(), reApiPostsCategories.String(), reApiPostsCategoriesFacts.String(), reApiPostsCategoriesRumors.String(), reApiPostsId.String()}
var GetRegexp = regexp.MustCompile(strings.Join(getRegexps, "|"))

var postRegexps = []string{reApiPostsCreate.String(), reApiPostsIdLike.String(), reApiPostsIdDislike.String(), reApiPostsIdComment.String(), reApiPostsIdCommentIdLike.String(), reApiPostsIdCommentIdDislike.String(), reApiSignup.String(), reApiLogin.String(), reApiLogout.String()}
var PostRegexp = regexp.MustCompile(strings.Join(postRegexps, "|"))
