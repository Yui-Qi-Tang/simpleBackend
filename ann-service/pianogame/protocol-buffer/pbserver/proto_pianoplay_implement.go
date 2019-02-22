package pbserver

// Need to rebuild pianoplayPb
import (
	"context"
	"log"

	//	"database/sql"
	//	"errors"
	//	"log"
	//	"simpleBackend/ann-service/pianogame"
	pianoplayPb "simpleBackend/ann-service/pianogame/protocol-buffer/pianoplay"
	//	"golang.org/x/crypto/bcrypt"
)

type pianoplayService struct{}

func (s *pianoplayService) Save(ctx context.Context, in *pianoplayPb.UserData) (*pianoplayPb.Response, error) {
	log.Println("I got", in.UUID, in.To, in.PianoKey, in.From)
	return &pianoplayPb.Response{Success: true, Msg: "test ok!"}, nil
}
