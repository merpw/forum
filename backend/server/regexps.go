package server

import (
	"regexp"
	"strings"
)

// TODO: maybe refactor, create regexp for all routes as a variables and use them in the following regexp arrays

var getRegexps = []string{
	`^\/api\/(?:post(?:s\/category\/[a-zA-Z0-9_.-]+|\/[[:digit:]]+)\/?|posts\/categories\/?|me\/liked\/posts\/?|user\/[[:digit:]]+(?:\/(?:posts\/?)?)?|posts\/?|me\/?)$`,
}

var GetRegexp = regexp.MustCompile(strings.Join(getRegexps, "|"))

var postRegexps = []string{
	`^\/api\/post(?:\/(?:[[:digit:]]+\/(?:(?:comment\/[[:digit:]]+\/)?dislike\/?|(?:comment\/[[:digit:]]+\/)?like\/?|comment\/?))?)?$`,
}

var PostRegexp = regexp.MustCompile(strings.Join(postRegexps, "|"))
