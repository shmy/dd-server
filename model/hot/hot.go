package hot

import (
	"github.com/shmy/dd-server/model"
)

var M = Model{&model.Model{model.Db.C("hots")}}

type Model struct {
	*model.Model
}
