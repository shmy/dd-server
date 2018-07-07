package video

import (
	"github.com/shmy/dd-server/model"
)

var M = Model{&model.Model{model.Db.C("videos")}}

type Model struct {
	*model.Model
}
