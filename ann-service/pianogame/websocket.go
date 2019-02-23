package pianogame

import (
	"log"
	"simpleBackend/ann-service/pianogame/datastructure"
	gameMsg "simpleBackend/ann-service/pianogame/msg"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"

	"google.golang.org/grpc"

	"context"
	pianoplayPb "simpleBackend/ann-service/pianogame/protocol-buffer/pianoplay"
)

func broadcastToClient(msg interface{}) {
	for client := range clients {
		client.SendMsg(msg)
	}
}

// close socket connection if error or leave
func connErrorOrExit(connErr error, u *datastructure.WebSocketUser) bool {
	if connErr == nil {
		return false
	}
	delete(clients, u) // remove socket user from shared model
	u.Close()          // close socket connectioin
	var userLeaveMsg interface{}
	if c := connErr.Error(); c == "client disconnects" {
		userLeaveMsg = &gameMsg.Exit{
			Base: gameMsg.Base{
				To:     "all",
				From:   u.GetID(),
				Action: gameMsg.ExitConn,
			},
			Text: strConcate("Good bye!", u.GetID()),
		}
	}
	broadcastToClient(userLeaveMsg)
	return true
}

func gameHandle(gamer *datastructure.WebSocketUser) {
	// for receivce message from websocket client
	var recMsg interface{}
	// receive msg
	for {
		err := gamer.GetConn().ReadJSON(&recMsg) // load data from client as inteface
		if connErrorOrExit(err, gamer) {
			return
		}
		recMsgMap := recMsg.(map[string]interface{}) // decode interface{}
		// check Action of receivced msg and do something
		switch act := recMsgMap["Action"]; act.(float64) { // decode map["Action"] interface as float64 and switch
		case gameMsg.SendPianoKey: // HINT: SendPianoKey is 'atoi' number, not a struct
			var pianoMsg gameMsg.PianoKey
			mapstructure.Decode(recMsgMap, &pianoMsg)
			pianoMsg.From = gamer.GetID()
			if pianoMsg.To == nil {
				pianoMsg.To = "all"
			}
			go saveMsg(
				gamer.GetID(),
				strconv.Itoa(int(pianoMsg.Key.(float64))),
				pianoMsg.To.(string),
				pianoMsg.From.(string),
				"no message",
			)
			broadcastToClient(pianoMsg)
		} // switch
	} // for
} // end of gameHandle()

func saveMsg(uuid, pianoKey, to, from, msgText string) {
	// grpc settings
	const (
		address = "localhost:9001" // gRPC server that is set in ann-servie/main.go now
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	/* Set connection */
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	grpcClient := pianoplayPb.NewPianoplayGreeterClient(conn)

	_, err = grpcClient.Save(
		ctx,
		&pianoplayPb.UserData{
			UUID:     uuid,
			PianoKey: pianoKey,
			To:       to,
			From:     from,
			Text:     msgText,
		},
	)
	if err != nil {
		log.Printf("could not greet: %v", err)
	}
} // saveMsg()
