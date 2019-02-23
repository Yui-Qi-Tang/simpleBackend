package pbserver

// Need to rebuild pianoplayPb
import (
	"context"
	"time"

	"simpleBackend/ann-service/pianogame"

	pianoplayPb "simpleBackend/ann-service/pianogame/protocol-buffer/pianoplay"

	"github.com/mongodb/mongo-go-driver/bson"
)

type pianoplayService struct{}

func (s *pianoplayService) Save(ctx context.Context, in *pianoplayPb.UserData) (*pianoplayPb.Response, error) {
	collection := pianogame.Mongodb.Database("piano_game").Collection("user_log")
	newGameData := bson.M{
		"uuid":      in.UUID,
		"to":        in.To,
		"piano_key": in.PianoKey,
		"from":      in.From,
		"msg_text":  in.Text,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, newGameData)
	if err != nil {
		return &pianoplayPb.Response{Success: false, Msg: "DB error"}, err
	}

	return &pianoplayPb.Response{Success: true, Msg: "save ok!"}, nil
}
