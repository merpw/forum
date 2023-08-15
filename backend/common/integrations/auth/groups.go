package auth

import (
	"net/http"
	"strconv"
)

type Group struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

// GetGroup gets group data from the auth service. Returns nil if the group was not found.
func GetGroup(groupId int) (group *Group) {
	InternalRequest(http.MethodGet, "/api/groups/"+strconv.Itoa(groupId), nil, &group)
	return group
}

// GetGroupMembers gets group members from the auth service. Returns nil if the group was not found.
func GetGroupMembers(groupId int) (members *[]int) {
	InternalRequest(http.MethodGet, "/api/groups/"+strconv.Itoa(groupId)+"/members", nil, &members)
	return members
}
