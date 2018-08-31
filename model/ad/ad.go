package ad

import (
	"github.com/shmy/dd-server/model"
)

var M = Model{&model.Model{model.Db.C("ad")}}

type Model struct {
	*model.Model
}
