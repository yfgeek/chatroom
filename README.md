# Chat Room

An AES-256 encrypted chatroom server & client based on UDP

[简体中文](README-CN.md)


# Language

* Golang ![golang](http://i.imgur.com/UEdZpr4.png)

## Features

- [x] UDP protocol
- [x] Communication with the Json file between C & S
- [x] Local configuration with Json format
- [x] AES-256 encryption
- [x] random KEY when running the server


## Client

The client configuration file located at `~\chat-config.json`. 

Sample
```json
{
	"listen": ":52915",
	"remote": "127.0.0.1",
	"key": "AzonXhdbWCYoAA52GTE9FnldZEN4KhEsInFJe1oHYAgzQTRsCyEdUlBOPzd3HxgFbTAudDZobiU8TQYbURBFWVdvMisNSn5UIw8kei0gcjl1cGkeFTV9U0tEY2YaCkdPYl9nZRQSBGsMQgFzVlxhL0hGAlV/O0A+OGoJfBwpE0w="
}
```
The user should set his ID and nickname when running the client.

``
./client
``

## Server 

The client configuration file located at`~\chat-config.json`.

Sample
```json
{
	"listen": ":52915",
	"remote": "", //can be empty
	"key": "AzonXhdbWCYoAA52GTE9FnldZEN4KhEsInFJe1oHYAgzQTRsCyEdUlBOPzd3HxgFbTAudDZobiU8TQYbURBFWVdvMisNSn5UIw8kei0gcjl1cGkeFTV9U0tEY2YaCkdPYl9nZRQSBGsMQgFzVlxhL0hGAlV/O0A+OGoJfBwpE0w="
}
```

``
./server
``


## Reference

[https://github.com/digitalis-io/golang-udp-chat](https://github.com/digitalis-io/golang-udp-chat)

[https://github.com/gwuhaolin/lightsocks](https://github.com/gwuhaolin/lightsocks)
