package database_test

import (
	"forum/database"
	"testing"
)

func TestOpsUsers(t *testing.T) {
	// add user
	var userId int
	t.Run("AddUser", func(t *testing.T) {
		userId = srv.DB.AddUser("OpsUserTestUser", "opsusertest@email", "password")
	})

	// get user by id
	var user *database.User
	t.Run("GetUserById", func(t *testing.T) {
		user = srv.DB.GetUserById(userId)
		if user == nil {
			t.Errorf("Expected user with id %d, got nil", userId)
		}
		if user.Id != userId {
			t.Errorf("Expected user with id %d, got %d", userId, user.Id)
		}
		noUser := srv.DB.GetUserById(-1)
		if noUser != nil {
			t.Errorf("Expected nil, got user with id %d", noUser.Id)
		}
	})

	// get user by login
	t.Run("GetUserByLogin", func(t *testing.T) {
		userByName := srv.DB.GetUserByLogin(user.Name)
		if userByName == nil {
			t.Errorf("Expected user with login %s, got nil", user.Name)
		}
		if userByName.Name != user.Name {
			t.Errorf("Expected user with name %s, got %s", user.Name, userByName.Name)
		}

		userByEmail := srv.DB.GetUserByLogin(user.Email)
		if userByEmail == nil {
			t.Errorf("Expected user with email %s, got nil", user.Email)
		}

		noUser := srv.DB.GetUserByLogin("noUser")
		if noUser != nil {
			t.Errorf("Expected nil, got user with Name %s", noUser.Name)
		}
	})

	// check is email or name taken
	t.Run("IsEmailOrNameTaken", func(t *testing.T) {
		if !srv.DB.IsNameTaken(user.Name) {
			t.Errorf("Expected name %s to be taken", user.Name)
		}
		if !srv.DB.IsEmailTaken(user.Email) {
			t.Errorf("Expected email %s to be taken", user.Email)
		}
	})

}
