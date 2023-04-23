package database_test

import (
	"forum/database"
	"testing"
)

func TestOpsChats(t *testing.T) {
	var privateChatId, userId, oponentId, channelChatId int

	t.Run("AddChat", func(t *testing.T) {
		privateChatId = srv.DB.AddChat(database.PrivateChat)
	})

	t.Run("AddEmptyChat", func(t *testing.T) {
		channelChatId = srv.DB.AddChat(database.ChannelChat)
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

	t.Run("AddUserMembershipToPrivateChat", func(t *testing.T) {
		srv.DB.AddMembership(userId, privateChatId)
	})

	t.Run("AddOponentMembershipToPrivateChat", func(t *testing.T) {
		srv.DB.AddMembership(oponentId, privateChatId)
	})

	// get chats by user id

	t.Run("GetChatsByUserId", func(t *testing.T) {
		chats := srv.DB.GetChats(userId, database.PrivateChat)
		if len(chats) != 1 {
			t.Fatal("Expected 1 chat, got", len(chats))
		}

		oponentChats := srv.DB.GetChats(oponentId, database.PrivateChat)
		if len(oponentChats) != 1 {
			t.Fatal("Expected 1 chat, got", len(oponentChats))
		}

		// check if chats are the same
		if chats[0].Id != oponentChats[0].Id {
			t.Fatal("Expected same chat, got different")
		}
	})

	// add membership to channel chat , to check if it is included in GetChatsByUserId
	t.Run("AddUserMembershipToChannelChat", func(t *testing.T) {
		srv.DB.AddMembership(userId, channelChatId)
	})

	// get chats by user id, including channel chat

	t.Run("GetChatsByUserIdIncludingChannelChat", func(t *testing.T) {
		chats := srv.DB.GetChats(userId, database.AnyChat)
		if len(chats) != 2 {
			t.Fatal("Expected 2 chats, got", len(chats))
		}
	})

	t.Run("GetPrivateChatsByUserId", func(t *testing.T) {
		chats := srv.DB.GetChats(userId, database.PrivateChat)
		if len(chats) != 1 {
			t.Fatal("Expected 1 chat, got", len(chats))
		}

		oponentChats := srv.DB.GetChats(oponentId, database.PrivateChat)
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
		userOponents := srv.DB.GetOnlineUsers(userId)
		if len(userOponents) != 1 {
			t.Fatal("Expected 1 user, got", len(userOponents))
		}

		oponentOponents := srv.DB.GetOnlineUsers(oponentId)
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

	// get private chat oponents by user id

	t.Run("GetContacts", func(t *testing.T) {
		oponents := srv.DB.GetContacts(userId)
		if len(oponents) != 1 {
			t.Fatal("Expected 1 oponent, got", len(oponents))
		}

		oponentOponents := srv.DB.GetContacts(oponentId)
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

	// TODO: uncomment this and extend to check GetChat, when logic will be approved
	// add messages to chat

	/*
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
	*/

}
