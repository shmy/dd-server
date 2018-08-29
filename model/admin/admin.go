package admin

import (
"github.com/shmy/dd-server/model"
)

var M = Model{&model.Model{model.Db.C("admin")}}

type Model struct {
	*model.Model
}
