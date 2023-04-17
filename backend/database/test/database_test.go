package database_test

import (
	"database/sql"
	"forum/server"
	"log"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var srv *server.Server

func TestMain(m *testing.M) {
	tmpDB, err := os.CreateTemp(".", "test.db")
	if err != nil {
		log.Fatal(err)
	}
	db, err = sql.Open("sqlite3", tmpDB.Name()+"?_foreign_keys=true")
	if err != nil {
		log.Fatal(err)
	}

	srv = server.Connect(db)
	err = srv.DB.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()
	if code != 0 {
		os.Exit(code)
	}

	// Clean only if all tests are successful (code=0)
	db.Close()
	tmpDB.Close()
	os.Remove(tmpDB.Name())
}

func TestOps(t *testing.T) {
	var privateChatId, userId, oponentId, channelChatId int

	t.Run("AddChat", func(t *testing.T) {
		privateChatId = srv.DB.AddChat(2)
	})

	t.Run("AddEmptyChat", func(t *testing.T) {
		channelChatId = srv.DB.AddChat(0)
	})

	t.Run("AddUser", func(t *testing.T) {
		userId = srv.DB.AddUser("testuser", "user@email", "password")
	})

	t.Run("AddOponent", func(t *testing.T) {
		oponentId = srv.DB.AddUser("testoponent", "oponent@email", "password")
	})

	// add dummy user, without session, to check if it is ignored

	t.Run("AddDummyUser", func(t *testing.T) {
		srv.DB.AddUser("dummy", "dummy@email", "password")
	})

	// login user and oponent, to create sessions, to get online users later

	t.Run("AddUserSession", func(t *testing.T) {
		srv.DB.AddSession("userToken", 9999999999, userId)
	})

	t.Run("AddOponentSession", func(t *testing.T) {
		srv.DB.AddSession("oponentToken", 9999999999, oponentId)
	})

	// add memberships to chat

	t.Run("AddUserMembership", func(t *testing.T) {
		srv.DB.AddMembership(userId, privateChatId)
	})

	t.Run("AddOponentMembership", func(t *testing.T) {
		srv.DB.AddMembership(oponentId, privateChatId)
	})

	// get chats by user id

	t.Run("GetChatsByUserId", func(t *testing.T) {
		chats := srv.DB.GetChatsByUserId(userId)
		if len(chats) != 1 {
			t.Fatal("Expected 1 chat, got", len(chats))
		}

		oponentChats := srv.DB.GetChatsByUserId(oponentId)
		if len(oponentChats) != 1 {
			t.Fatal("Expected 1 chat, got", len(oponentChats))
		}

		// check if chats are the same
		if chats[0].Id != oponentChats[0].Id {
			t.Fatal("Expected same chat, got different")
		}
	})

	// get online users

	t.Run("GetOnlineUsers", func(t *testing.T) {
		userOponents := srv.DB.GetOnlineUsersIdsAndNames(userId)
		if len(userOponents) != 1 {
			t.Fatal("Expected 1 user, got", len(userOponents))
		}

		oponentOponents := srv.DB.GetOnlineUsersIdsAndNames(oponentId)
		if len(oponentOponents) != 1 {
			t.Fatal("Expected 1 user, got", len(oponentOponents))
		}

		// check if users are the expected
		if userOponents[0].Id != oponentId {
			t.Fatal("Expected oponentId, got different")
		}

		if oponentOponents[0].Id != userId {
			t.Fatal("Expected userId, got different")
		}
	})

	// get private chat opoents by user id

	t.Run("GetPrivateChatOponentsByUserId", func(t *testing.T) {
		oponents := srv.DB.GetPrivateChatOponentsByUserId(userId)
		if len(oponents) != 1 {
			t.Fatal("Expected 1 oponent, got", len(oponents))
		}

		oponentOponents := srv.DB.GetPrivateChatOponentsByUserId(oponentId)
		if len(oponentOponents) != 1 {
			t.Fatal("Expected 1 oponent, got", len(oponentOponents))
		}

		// check if users are the expected
		if oponents[0].Id != oponentId {
			t.Fatal("Expected oponentId, got different")
		}

		if oponentOponents[0].Id != userId {
			t.Fatal("Expected userId, got different")
		}
	})

	// add messages to chat

	t.Run("AddUserMessage", func(t *testing.T) {
		srv.DB.AddMessage(userId, privateChatId, "helloFromUser")
	})

	t.Run("AddOponentMessage", func(t *testing.T) {
		srv.DB.AddMessage(oponentId, privateChatId, "helloFromOponent")
	})

	// add messages to channel chat, but only from user, let is say he is owner

	t.Run("AddUserMessageToChannel", func(t *testing.T) {
		srv.DB.AddMessage(userId, channelChatId, "helloFromUserToChannel")
	})

}
