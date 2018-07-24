package favorite

import (
	"github.com/shmy/dd-server/model"
)

var M = Model{&model.Model{model.Db.C("favorites")}}

type Model struct {
	*model.Model
}
