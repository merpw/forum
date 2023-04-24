package database_test

import (
	"database/sql"
	"fmt"
	"forum/database"
	"testing"
)

// AddUser(name, email, password string, firstName, lastName, dob, gender sql.NullString)
type TestUserSqlNullString struct {
	Name      string
	Email     string
	Password  string
	FirstName sql.NullString
	LastName  sql.NullString
	DoB       sql.NullString
	Gender    sql.NullString
}

func createTestUser(userIndex int) TestUserSqlNullString {
	return TestUserSqlNullString{
		Name:      "testUser " + fmt.Sprint(userIndex),
		Email:     "testUser" + fmt.Sprint(userIndex) + "@test.com",
		Password:  "testPassword",
		FirstName: sql.NullString{String: "testFirstName" + fmt.Sprint(userIndex), Valid: true},
		LastName:  sql.NullString{String: "testLastName" + fmt.Sprint(userIndex), Valid: true},
		DoB:       sql.NullString{String: "testDoB" + fmt.Sprint(userIndex), Valid: true},
		Gender:    sql.NullString{String: "testGender" + fmt.Sprint(userIndex), Valid: true},
	}
}

func TestOpsUsers(t *testing.T) {
	// add user
	testUser := createTestUser(1)
	var userId int
	t.Run("AddUser", func(t *testing.T) {
		userId = DB.AddUser(
			testUser.Name, testUser.Email, testUser.Password,
			testUser.FirstName, testUser.LastName, testUser.DoB, testUser.Gender,
		)
	})

	// get user by id
	var user *database.User
	t.Run("GetUserById", func(t *testing.T) {
		user = DB.GetUserById(userId)
		if user == nil {
			t.Fatalf("Expected user with id %d, got nil", userId)
		}
		if user.Id != userId {
			t.Fatalf("Expected user with id %d, got %d", userId, user.Id)
		}
		noUser := DB.GetUserById(-1)
		if noUser != nil {
			t.Fatalf("Expected nil, got user with id %d", noUser.Id)
		}
	})

	// get user by login
	t.Run("GetUserByLogin", func(t *testing.T) {
		userByName := DB.GetUserByLogin(user.Name)
		if userByName == nil {
			t.Fatalf("Expected user with login %s, got nil", user.Name)
		}
		if userByName.Name != user.Name {
			t.Fatalf("Expected user with name %s, got %s", user.Name, userByName.Name)
		}

		userByEmail := DB.GetUserByLogin(user.Email)
		if userByEmail == nil {
			t.Fatalf("Expected user with email %s, got nil", user.Email)
		}

		noUser := DB.GetUserByLogin("noUser")
		if noUser != nil {
			t.Fatalf("Expected nil, got user with Name %s", noUser.Name)
		}
	})

	// check is email or name taken
	t.Run("IsEmailOrNameTaken", func(t *testing.T) {
		if !DB.IsNameTaken(user.Name) {
			t.Fatalf("Expected name %s to be taken", user.Name)
		}
		if !DB.IsEmailTaken(user.Email) {
			t.Fatalf("Expected email %s to be taken", user.Email)
		}
	})

}
