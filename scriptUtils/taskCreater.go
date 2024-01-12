package scriptUtils

import (
	"github.com/google/uuid"
	"github.com/slainsama/msgr_server/globals"
	"github.com/slainsama/msgr_server/models"
)

func TaskCreate(scriptName string, params []string) (newTask models.Task) {
	newId := uuid.New().String()
	newTask = models.Task{Id: newId, ScriptName: scriptName, Params: params}
	globals.TaskList[newId] = newTask
	return newTask
}