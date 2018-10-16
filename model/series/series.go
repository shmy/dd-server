package _series

import (
	"github.com/shmy/dd-server/model"
)

var M = Model{&model.Model{model.Db.C("series")}}

type Model struct {
	*model.Model
}
