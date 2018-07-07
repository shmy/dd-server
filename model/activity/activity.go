package activity

import (
	"github.com/shmy/dd-server/model"
)

var M = Model{&model.Model{model.Db.C("activitys")}}

type Model struct {
	*model.Model
}
