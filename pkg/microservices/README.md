# Microservices Package

This package gives you all the functionality to interact with ioFog both via Local API and WebSockets:

 - send new message to ioFog with REST (PostMessage)
 - fetch next unread messages from ioFog (GetNextMessages)
 - fetch messages for time period and list of publishers (GetMessagesFromPublishersWithinTimeFrame)
 - get config options (GetConfig)
 - create IoMessage, encode(decode) to(from) raw bytes, encode(decode) data to(from) base64 string (IoMessage methods)
 - connect to ioFog Control Channel via WebSocket (EstablishControlWsConnection)
 - connect to ioFog Message Channel via WebSocket (EstablishMessageWsConnection) and publish new message via this channel (SendMessageViaSocket)

## Code snippets: 

Get sdk:
```go
go get github.com/eclipse-iofog/iofog-sdk-go
```

Import package:
```go
import (
	msvcs "github.com/eclipse-iofog/iofog-sdk-go/pkg/microservices"
)
```

Create IoFog client with default settings:
```go
client, err := msvcs.NewDefaultIoFogClient()
```

Or specify host, port, ssl and container id explicitly:
```go
client, err := msvcs.NewIoFogClient("IoFog", false, "containerId", 54321)
```


#### REST calls

Get list of next unread IoMessages:
```go
messages, err := client.GetNextMessages()
```

Post new IoMessage to ioFog via REST call:
```go
response, err := client.PostMessage(&msvcs.IoMessage{
	SequenceNumber:1,
	SequenceTotal:1,
	InfoType:"text",
	InfoFormat:"utf-8",
	ContentData: []byte("foo"),
	ContextData: []byte("bar"),
})
```

Get an array of IoMessages from specified publishers within given timeframe:
```go
messages, err := client.GetMessagesFromPublishersWithinTimeFrame(&msvcs.MessagesQueryParameters{
	TimeFrameStart: 1234567890123,
	TimeFrameEnd: 1234567892123,
	Publishers: []string{"sefhuiw4984twefsdoiuhsdf", "d895y459rwdsifuhSDFKukuewf", "SESD984wtsdidsiusidsufgsdfkh"},
})
```

Get container's config:
```go
config, err := client.GetConfig()
```

#### WebSocket calls

Establish connection with message ws. This call returns two channels, so
 you can listen to incoming messages and receipts:
```go
dataChannel, receiptChannel := client.EstablishMessageWsConnection()
for {
	select {
	case msg := <-dataChannel:
		// msg is IoMessage received
	case r := <-receiptChannel:
		// r is response with ID and Timestamp
}
```

After establishing this connection you can send your own message to IoFog:
```go
client.SendMessageViaSocket(&msvcs.IoMessage{
	Tag: "aaa",
	SequenceNumber: 127,
	ContentData: []byte("Here goes some test data"),
	ContextData: []byte("This one is test too"),
})
```


Establish connection with control ws and pass channel to listen to incoming config update signals:
```go
confChannel := client.EstablishControlWsConnection()
for {
	select {
	case <-confChannel:
		// signal received
		// we can fetch new config now
		config, err := client.GetConfig()
}
```
