package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"nhooyr.io/websocket"
	"website-status-checker-in-go/states"
)

func main() {
	router := gin.Default()
	states.CurrState.Name = "Empty"
	fmt.Println("State before client Start up : ", states.CurrState)

	ctx, conn := subscribe()
	go func() {
		readSubscribedMessages := func() {
			_, bytes, err := conn.Read(ctx)
			if err != nil {
				fmt.Println("Error reading from websocket connection ! ", err.Error())
			}
			states.CurrState.Name = string(bytes)
			fmt.Println("Received an event : ", string(bytes))
		}

		for {
			readSubscribedMessages()
		}
	}()

	router.GET("/healthz", func(context *gin.Context) {
		fmt.Println("================================> Current State : ", states.CurrState)
		context.JSON(http.StatusOK, states.CurrState.Name)
	})
	http.ListenAndServe(":8000", router)
}

func subscribe() (context.Context, *websocket.Conn) {
	fmt.Println("Subscribing to PUSHPIN ")
	ctx := context.Background()
	conn, _, err := websocket.Dial(ctx, "ws://localhost:7999/subscribe", nil)
	if err != nil {
		fmt.Println("Error doing Dialing pushpin over websocket : " + err.Error())
	}
	fmt.Printf("Dialed and subscribed successfully !")
	return ctx, conn
}

