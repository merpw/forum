package database_test

import (
	"forum/database"
	"testing"
)

func TestOpsChats(t *testing.T) {
	testUser := createTestUser(1)
	testOpponent := createTestUser(2)
	testDummy := createTestUser(3)

	var privateChatId, userId, opponentId, channelChatId int

	t.Run("AddChat", func(t *testing.T) {
		privateChatId = DB.AddChat(database.PrivateChat)
	})

	t.Run("AddEmptyChat", func(t *testing.T) {
		channelChatId = DB.AddChat(database.ChannelChat)
	})

	t.Run("AddUser", func(t *testing.T) {
		userId = DB.AddUser(
			testUser.Name, testUser.Email, testUser.Password,
			testUser.FirstName, testUser.LastName, testUser.DoB, testUser.Gender,
		)
	})

	t.Run("AddOpponent", func(t *testing.T) {
		opponentId = DB.AddUser(
			testOpponent.Name, testOpponent.Email, testOpponent.Password,
			testOpponent.FirstName, testOpponent.LastName, testOpponent.DoB, testOpponent.Gender,
		)
	})

	// add dummy user, without session, to check if it is ignored

	t.Run("AddDummyUser", func(t *testing.T) {
		DB.AddUser(
			testDummy.Name, testDummy.Email, testDummy.Password,
			testDummy.FirstName, testDummy.LastName, testDummy.DoB, testDummy.Gender,
		)
	})

	// login user and opponent, to create sessions, to get online users later

	t.Run("AddUserSession", func(t *testing.T) {
		DB.AddSession("userToken", 9999999999, userId)
	})

	t.Run("AddOpponentSession", func(t *testing.T) {
		DB.AddSession("opponentToken", 9999999999, opponentId)
	})

	// add memberships to chat

	t.Run("AddUserMembershipToPrivateChat", func(t *testing.T) {
		DB.AddMembership(userId, privateChatId)
	})

	t.Run("AddOpponentMembershipToPrivateChat", func(t *testing.T) {
		DB.AddMembership(opponentId, privateChatId)
	})

	// get chats by user id

	t.Run("GetChatsByUserId", func(t *testing.T) {
		chats := DB.GetChats(userId, database.PrivateChat)
		if len(chats) != 1 {
			t.Fatal("Expected 1 chat, got", len(chats))
		}

		opponentChats := DB.GetChats(opponentId, database.PrivateChat)
		if len(opponentChats) != 1 {
			t.Fatal("Expected 1 chat, got", len(opponentChats))
		}

		// check if chats are the same
		if chats[0].Id != opponentChats[0].Id {
			t.Fatal("Expected same chat, got different")
		}
	})

	// add membership to channel chat , to check if it is included in GetChatsByUserId
	t.Run("AddUserMembershipToChannelChat", func(t *testing.T) {
		DB.AddMembership(userId, channelChatId)
	})

	// get chats by user id, including channel chat

	t.Run("GetChatsByUserIdIncludingChannelChat", func(t *testing.T) {
		chats := DB.GetChats(userId, database.AnyChat)
		if len(chats) != 2 {
			t.Fatal("Expected 2 chats, got", len(chats))
		}
	})

	t.Run("GetPrivateChatsByUserId", func(t *testing.T) {
		chats := DB.GetChats(userId, database.PrivateChat)
		if len(chats) != 1 {
			t.Fatal("Expected 1 chat, got", len(chats))
		}

		opponentChats := DB.GetChats(opponentId, database.PrivateChat)
		if len(opponentChats) != 1 {
			t.Fatal("Expected 1 chat, got", len(opponentChats))
		}

		// check if chats are the same
		if chats[0].Id != opponentChats[0].Id {
			t.Fatal("Expected same chat, got different")
		}
	})

	// get online users

	t.Run("GetOnlineUsers", func(t *testing.T) {
		userOpponents := DB.GetOnlineUsers(userId)
		if len(userOpponents) != 1 {
			t.Fatal("Expected 1 user, got", len(userOpponents))
		}

		opponentOpponents := DB.GetOnlineUsers(opponentId)
		if len(opponentOpponents) != 1 {
			t.Fatal("Expected 1 user, got", len(opponentOpponents))
		}

		// check if users are the expected
		if userOpponents[0].Id != opponentId {
			t.Fatal("Expected opponentId, got different")
		}

		if opponentOpponents[0].Id != userId {
			t.Fatal("Expected userId, got different")
		}
	})

	// get private chat opponents by user id

	t.Run("GetContacts", func(t *testing.T) {
		opponents := DB.GetContacts(userId)
		if len(opponents) != 2 {
			t.Fatal("Expected 2 opponents, got", len(opponents))
		}

		opponentOpponents := DB.GetContacts(opponentId)
		if len(opponentOpponents) != 2 {
			t.Fatal("Expected 2 opponents, got", len(opponentOpponents))
		}

		// check if users are the expected
		if opponents[0].Id != opponentId && opponents[1].Id != opponentId {
			t.Fatal("Expected opponentId, got different")
		}

		if opponentOpponents[0].Id != userId && opponentOpponents[1].Id != userId {
			t.Fatal("Expected userId, got different")
		}
	})

	// TODO: uncomment this and extend to check GetChat, when logic will be approved
	// add messages to chat

	/*
		t.Run("AddUserMessage", func(t *testing.T) {
			DB.AddMessage(userId, privateChatId, "helloFromUser")
		})

		t.Run("AddOpponentMessage", func(t *testing.T) {
			DB.AddMessage(opponentId, privateChatId, "helloFromOpponent")
		})

		// add messages to channel chat, but only from user, let is say he is owner

		t.Run("AddUserMessageToChannel", func(t *testing.T) {
			DB.AddMessage(userId, channelChatId, "helloFromUserToChannel")
		})
	*/

}
