package collection

import (
	"github.com/shmy/dd-server/model"
)

var M = Model{&model.Model{model.Db.C("collections")}}

type Model struct {
	*model.Model
}
