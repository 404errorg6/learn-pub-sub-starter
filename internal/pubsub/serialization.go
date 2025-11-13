package pubsub

import (
	"encoding/gob"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func Encode(gLog routing.GameLog) error {
	encoder := gob.Encoder{}
}
