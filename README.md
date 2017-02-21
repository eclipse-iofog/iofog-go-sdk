# container-sdk-go

Golang SDK lets you easily build an ioElement with your favorite Go language. It gives you all the functionality to interact with ioFog both via Local API and WebSockets:

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
   go get github.com/iotracks/container-sdk-go
```

Import sdk:
```go
   import "github.com/iotracks/container-sdk-go"
```

Create IoFog client with default settings:
```go
  client := sdk.NewDefaultIoFogClient()
  	if client != nil {
  		// do smth...
  	}
```

Or specify host, port, ssl and container id explicitly:
```go
	client := sdk.NewIoFogClient("IoFog", 54321, false, "containerId")
  	if client != nil {
  		// do smth...
  	}
```

#### REST calls

Get list of next unread IoMessages. Pass true in parameters if you with to decode received data from base64 to raw bytes, pass false otherwise:
```go
	messages, err := client.GetNextMessages(true)
	if err != nil {
	    // handle bad request or error
		println(err.Error())
	} else {
	    // messages array contains received
	    // and parsed IoMessages
	    fmt.Printf("%+v", messages)
	}
```

Post new IoMessage to ioFog via REST call:
```go
	postMessageResponse, err := client.PostMessage(&sdk.IoMessage{
		SequenceNumber:1,
		SequenceTotal:1,
		InfoType:"text",
		InfoFormat:"utf-8",
		ContentData: []byte("foo"),
		ContextData: []byte("bar"),
	})
	if err != nil {
	    // handle bad request or error
		println(err.Error())
	} else {
	    // postMessageResponse contains sent message
	    // ID and generated Timestamp 
	    fmt.Printf("%+v", postMessageResponse)
	}
```

Get an array of IoMessages from specified publishers within given timeframe. Pass true in parameters if you with to decode received data from base64 to raw bytes, pass false otherwise:
```go
	timeFrameMessages, err := client.GetMessagesFromPublishersWithinTimeFrame(&sdk.MessagesQueryParameters{
		TimeFrameStart: 1234567890123,
		TimeFrameEnd: 1234567892123,
		Publishers: []string{"sefhuiw4984twefsdoiuhsdf", "d895y459rwdsifuhSDFKukuewf", "SESD984wtsdidsiusidsufgsdfkh"},

	}, true)
	if err != nil {
	    // handle bad request or error
		println(err.Error())
	} else {
	    // timeFrameMessages contains fields Messages,
	    // TimeFrameStart and TimeFrameEnd 
        fmt.Printf("%+v", timeFrameMessages)
    }
```

Get container's config:
```go
	config, err := client.GetConfig()
	if err != nil {
	    // handle bad request or error
		println(err.Error())
	} else {
	    // config is plain old go map with string keys
	    // and interface{} values
	    fmt.Println("Config: ", config)
	}
```

#### WebSocket calls
Establish connection with message ws and pass channels to listen to incoming messages and receipts or send your own message:
```go
		dataChannel := make(chan *sdk.IoMessage)
		receiptChannel := make(chan *sdk.PostMessageResponse)
		client.EstablishMessageWsConnection(dataChannel, receiptChannel)
		
		client.SendMessageViaSocket(&sdk.IoMessage{
        	Tag: "aaa",
        	SequenceNumber: 127,
        	ContentData: []byte("Here goes some test data"),
        	ContextData: []byte("This one is test too"),
        })

		for {
			select {
			case msg := <-dataChannel:
				// msg is IoMessage received
			case r := <-receiptChannel:
				// r is response with ID and Timestamp
		}
		

```

Establish connection with control ws and pass channel to listen to incoming config update signals:
```go
	confChannel := make(chan int)
	client.EstablishControlWsConnection(confChannel)
	for {
		select {
		case <-confChannel:
		    // signal received
		    // we can fetch new config now
			config, err := client.GetConfig()
			if err != nil {
				println(err.Error())
			} else {
			    fmt.Println("Config: ", config)
			}
	}
```







