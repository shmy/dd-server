package user

import (
	"github.com/shmy/dd-server/model"
)

var M = Model{&model.Model{model.Db.C("users")}}

type Model struct {
	*model.Model
}
