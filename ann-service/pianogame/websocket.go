package pianogame

import (
	"simpleBackend/ann-service/pianogame/datastructure"
	gameMsg "simpleBackend/ann-service/pianogame/msg"

	"github.com/mitchellh/mapstructure"
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
	var recMsg interface{}
	for {
		err := gamer.GetConn().ReadJSON(&recMsg) // load data from client as inteface
		if connErrorOrExit(err, gamer) {
			return
		}
		recMsgMap := recMsg.(map[string]interface{}) // decode interface{}
		// check Action of receivced msg and do something
		switch act := recMsgMap["Action"]; act.(float64) { // decode map["Action"] interface as float64 and switch
		case gameMsg.SendPianoKey:
			var pianoMsg gameMsg.PianoKey
			mapstructure.Decode(recMsgMap, &pianoMsg)
			pianoMsg.From = gamer.GetID()
			if pianoMsg.To == nil {
				pianoMsg.To = "all"
			}
			broadcastToClient(pianoMsg)
		} // switch
	} // for
}
