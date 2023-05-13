# dummy-ws-backend

Simple dummy websocket server for development purposes written in TypeScript.

## How to run?

- `npm install` - install dependencies
- `npm run start` - start the server

## Concept

- Server accepts WebSocket connections
- We use JSON as a message format
- Each message has a `type` property that defines the message type, it can be `handshake`, `get`, `post` or `error`
- Each message has an `item` property that contains the payload of the message

## Postman Workspace

[![Run in Postman](https://run.pstmn.io/button.svg)](https://restless-zodiac-600487.postman.co/workspace/merpw~8e6f6f99-c3c2-4738-b609-a958ed3a626a/ws-raw-request/644a6692a783d705e1c749a3)

All basic requests are implemented in Postman, you can use it to test the server behavior.

## WebSocket schemas

Here's a draft of the WebSocket messages that the server accepts

### Handshake, initial message

First request to authenticate the client

#### Request

```json
{
  "type": "handshake",
  "item": {
    "token": "<string>"
  }
}
```

#### Response

```json
{
  "type": "handshake",
  "item": {
    "data": {
      "userId": 1
    }
  }
}
```

### Get

```json
{
  "type": "get",
  "item": {
    "url": "<string>"
  }
}
```

#### Example: get all chats (chat ids)

##### Request

```json
{
  "type": "get",
  "item": {
    "url": "/chat/all"
  }
}
```

##### Response

```json
{
  "type": "get",
  "item": {
    "url": "/chat/all",
    "data": [1, 2, 3]
  }
}
```

### Post

```json
{
  "type": "post",
  "item": {
    "url": "<string>",
    "data": {
      ...
    }
  }
}
```

#### Example: send message to chat

##### Request

```json
{
  "type": "post",
  "item": {
    "url": "/chat/1/message",
    "data": {
      "content": "Hello world!"
    }
  }
}
```

##### Response

```json
{
  "type": "post",
  "item": {
    "url": "/chat/1/message",
    "data": {
      "messageId": 1
    }
  }
}
```

### Errors

```json
{
  "type": "error",
  "item": {
    "message": "Error message"
  }
}
```

For more examples, please refer to the Postman workspace and the TypeScript types in the `src/types.d.ts` file

[//]: # "TODO: maybe include url in error messages?"
