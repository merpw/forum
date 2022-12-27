package server

import (
	"regexp"
	"strings"
)

// TODO: maybe refactor, create regexp for all routes as a variables and use them in the following regexp arrays

const reApiMe = `^\/api\/me\/?$`
const reApiMeLikedPosts = `^\/api\/me\/liked\/posts\/?$`
const reApiUserId = `^\/api\/user\/[[:digit:]]+\/?$`
const reApiUserIdPosts = `^\/api\/user\/[[:digit:]]+\/posts\/?$`
const reApiPosts = `^\/api\/posts\/?$`
const reApiPostsCategories = `^\/api\/posts\/categories\/?$`
const reApiPostsCategoryFacts = `^\/api\/posts\/category\/facts\/?$`
const reApiPostsCategoryRumors = `^\/api\/posts\/category\/rumors\/?$`
const reApiPostsId = `^\/api\/posts\/[[:digit:]]+\/?$`

const reApiPostsCreate = `^\/api\/posts\/create\/?$`
const reApiPostsIdLike = `^\/api\/posts\/[[:digit:]]+\/like\/?$`
const reApiPostsIdDislike = `^\/api\/posts\/[[:digit:]]+\/dislike\/?$`
const reApiPostsIdComment = `^\/api\/posts\/[[:digit:]]+\/comment\/?$`
const reApiPostsIdCommentIdLike = `^\/api\/posts\/[[:digit:]]+\/comment\/[[:digit:]]+\/like\/?$`
const reApiPostsIdCommentIdDislike = `^\/api\/posts\/[[:digit:]]+\/comment\/[[:digit:]]+\/dislike\/?$`
const reApiSignup = `^\/api\/signup\/?$`
const reApiLogin = `^\/api\/login\/?$`
const reApiLogout = `^\/api\/logout\/?$`

// var getRegexps = []string{
// 	`^\/api\/(?:post(?:s\/category\/[a-zA-Z0-9_.-]+|\/[[:digit:]]+)\/?|posts\/categories\/?|me\/liked\/posts\/?|user\/[[:digit:]]+(?:\/(?:posts\/?)?)?|posts\/?|me\/?)$`,
// }

var getRegexps = []string{
	reApiMe, reApiMeLikedPosts, reApiUserId, reApiUserIdPosts, reApiPosts, reApiPostsCategories, reApiPostsCategoryFacts, reApiPostsCategoryRumors, reApiPostsId}

var GetRegexp = regexp.MustCompile(strings.Join(getRegexps, "|"))

//	var postRegexps = []string{
//		`^\/api\/post(?:\/(?:[[:digit:]]+\/(?:(?:comment\/[[:digit:]]+\/)?dislike\/?|(?:comment\/[[:digit:]]+\/)?like\/?|comment\/?))?)?$`,
//	}
var postRegexps = []string{reApiPostsCreate, reApiPostsIdLike, reApiPostsIdDislike, reApiPostsIdComment, reApiPostsIdCommentIdLike, reApiPostsIdCommentIdDislike, reApiSignup, reApiLogin, reApiLogout}

var PostRegexp = regexp.MustCompile(strings.Join(postRegexps, "|"))
