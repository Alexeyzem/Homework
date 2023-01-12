package processing_commands

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"strconv"
)

type task struct {
	Title string
	Id    int
}
type user struct {
	Name            string
	ChatId          int64
	CreatedTasks    task
	ExecutableTasks task
}
type MyUser struct {
	Name string
	Id   int64
}
type Out struct {
	Message string
	ChatId  int64
}

var Users = map[user]interface{}{}
var ActiveTask []task
var i = 0

func processingTemplate(fileName string, data interface{}) (string, error) {
	tmpl := template.New(fileName)
	tmpl, err := template.ParseFiles("processing_command/templates/" + fileName)
	if err != nil {
		log.Print(err)
		return "", errors.New("internal error")
	}
	wr := &bytes.Buffer{}
	err = tmpl.Execute(wr, data)
	if err != nil {
		log.Print(err)
		return "", err
	}
	return string(wr.Bytes()), nil
}

func isExecutor(userData MyUser, executableTask task) bool {
	for key := range Users {
		if userData.Id == key.ChatId && key.ExecutableTasks.Id == executableTask.Id {
			return true
		}
	}
	return false
}

func My(_ string, userData MyUser) ([]Out, error) {
	var users []user
	for key := range Users {
		if key.ChatId == userData.Id {
			users = append(users, key)
		}
	}
	var out []Out
	var supportive Out
	supportive.ChatId = userData.Id
	if len(users) == 0 {
		supportive.Message = "На вас не назначены задачи"
		out = append(out, supportive)
		return out, nil
	}
	for _, value := range users {
		if value.ExecutableTasks.Id == 0 {
			continue
		}
		text, err := processingTemplate("My", value)
		if err != nil {
			return nil, err
		}
		supportive.Message = supportive.Message + text + "\n"
		out = append(out, supportive)
	}
	return out, nil
}

func Assign(idS string, userData MyUser) ([]Out, error) {
	if len(idS) == 0 {
		return nil, errors.New("wrong ID")
	}
	id, err := strconv.Atoi(idS)
	if err != nil {
		return nil, err
	}
	var tasks = task{Id: 0}
	var out []Out
	var supportive Out
	for _, value := range ActiveTask {
		if value.Id == id {
			tasks = value
			break
		}
	}
	if tasks.Id == 0 {
		supportive.Message = "Данной задачи нет"
		supportive.ChatId = userData.Id
		out = append(out, supportive)
		return out, nil
	}

	var executor user
	for key := range Users {
		if key.ExecutableTasks == tasks {
			executor = key
			break
		}
	}
	if executor.ChatId == userData.Id {
		supportive.Message = `Задача "` + tasks.Title + `" уже назначена на вас`
		supportive.ChatId = executor.ChatId
		out = append(out, supportive)
		return out, nil
	} else {
		supportive.Message = `Задача "` + tasks.Title + `" назначена на вас`
		supportive.ChatId = userData.Id
		out = append(out, supportive)
	}
	var creators user
	for key := range Users {
		if key.CreatedTasks == tasks {
			creators = key
			break
		}
	}

	if creators.Name != userData.Name && executor.ChatId != 0 {
		supportive.Message = `Задача "` + tasks.Title + `" назначена на @` + userData.Name
		supportive.ChatId = creators.ChatId
		out = append(out, supportive)
		supportive.Message = `Задача "` + tasks.Title + `" назначена на @` + userData.Name
		supportive.ChatId = executor.ChatId
		out = append(out, supportive)
		delete(Users, user{Name: executor.Name, ChatId: executor.ChatId, ExecutableTasks: executor.ExecutableTasks})
	} else if executor.ChatId != 0 {
		supportive.Message = `Задача "` + tasks.Title + `" назначена на @` + userData.Name
		supportive.ChatId = executor.ChatId
		out = append(out, supportive)
		delete(Users, user{Name: executor.Name, ChatId: executor.ChatId, ExecutableTasks: executor.ExecutableTasks})
	}
	if creators.Name != userData.Name && executor.ChatId == 0 {
		supportive.Message = `Задача "` + tasks.Title + `" назначена на @` + userData.Name
		supportive.ChatId = creators.ChatId
		out = append(out, supportive)
	}
	Users[user{Name: userData.Name, ChatId: userData.Id, ExecutableTasks: tasks}] = nil
	return out, nil
}

func UnAssign(idS string, userData MyUser) ([]Out, error) {
	if len(idS) == 0 {
		return nil, errors.New("wrong ID")
	}
	id, err := strconv.Atoi(idS)
	if err != nil {
		return nil, err
	}
	var out []Out
	for key := range Users {
		if key.ExecutableTasks.Id == id && key.ChatId == userData.Id {
			out = append(out, Out{ChatId: userData.Id, Message: "Принято"})
			for creator := range Users {
				if creator.CreatedTasks.Id == id {
					out = append(out, Out{ChatId: userData.Id,
						Message: `Задача "` + creator.CreatedTasks.Title + `" осталась без исполнителя`})
					break
				}
			}
			var deleting = user{Name: userData.Name,
				ChatId:          userData.Id,
				ExecutableTasks: key.ExecutableTasks}
			delete(Users, deleting)
			return out, nil
		}
	}
	return []Out{{Message: "Задача не на вас", ChatId: userData.Id}}, nil
}

func Resolve(idS string, userData MyUser) ([]Out, error) {
	if len(idS) == 0 {
		return nil, errors.New("wrong ID")
	}
	id, err := strconv.Atoi(idS)
	if err != nil {
		return nil, err
	}
	var out []Out
	for key := range Users {
		if key.ExecutableTasks.Id == id && key.ChatId == userData.Id {
			out = append(out, Out{ChatId: userData.Id,
				Message: `Задача "` + key.ExecutableTasks.Title + `" выполнена`})
			for creator := range Users {
				if creator.CreatedTasks.Id == id {
					out = append(out, Out{ChatId: userData.Id,
						Message: `Задача "` + creator.CreatedTasks.Title + `" выполнена @` + key.Name})
					var deleting = user{Name: creator.Name,
						ChatId:       creator.ChatId,
						CreatedTasks: creator.CreatedTasks}
					delete(Users, deleting)
					break
				}
			}
			var deleting = user{Name: userData.Name,
				ChatId:          userData.Id,
				ExecutableTasks: key.ExecutableTasks}
			delete(Users, deleting)
			for i, value := range ActiveTask {
				if value.Id == id {
					ActiveTask = append(ActiveTask[:i], ActiveTask[i+1:]...)
				}
			}
			return out, nil
		}
	}
	return []Out{{Message: "Задача не на вас", ChatId: userData.Id}}, nil
}

func New(title string, userData MyUser) ([]Out, error) {
	if title == "" {
		return nil, errors.New("wrong task")
	}
	i++
	var out []Out
	var supportive Out
	var NewTask = task{Title: title, Id: i}
	ActiveTask = append(ActiveTask, NewTask)
	Users[user{Name: userData.Name, CreatedTasks: NewTask, ChatId: userData.Id}] = nil
	text, err := processingTemplate("New", NewTask)
	if err != nil {
		return nil, err
	}
	supportive.Message = text
	supportive.ChatId = userData.Id
	out = append(out, supportive)
	return out, nil
}

func Tasks(_ string, userData MyUser) ([]Out, error) {
	if len(ActiveTask) == 0 || ActiveTask[0].Id == 0 {
		return []Out{{Message: "Нет задач", ChatId: userData.Id}}, nil
	}
	var out []Out
	var creator user
	var msg string
	for _, value := range ActiveTask {
		for key := range Users {
			if value.Id == key.CreatedTasks.Id {
				creator = key
				break
			}
		}
		text, err := processingTemplate("Tasks", creator)
		if err != nil {
			return nil, err
		}
		if isExecutor(userData, value) {
			msg += text + "\nassignee: я\n/unassign_" + strconv.Itoa(value.Id) + " /resolve_" + strconv.Itoa(value.Id) + "\n"
		} else {
			var name string
			for key := range Users {
				if value.Id == key.ExecutableTasks.Id {
					name = key.Name
					break
				}
			}
			if len(name) != 0 {
				msg += text + "\nassignee: @" + name + "\n"
			} else {
				msg += text + "\n/assign_" + strconv.Itoa(value.Id) + "\n"
			}
		}
	}
	out = append(out, Out{ChatId: userData.Id, Message: msg})
	return out, nil
}

func Owner(_ string, userData MyUser) ([]Out, error) {
	var users []user
	var out []Out
	var supportive Out
	supportive.ChatId = userData.Id
	for key := range Users {
		if key.Name == userData.Name {
			users = append(users, key)
		}
	}
	if len(users) == 0 {
		supportive.Message = "У вас нет активных созданных задач"
		out = append(out, supportive)
		return out, nil
	}
	for _, value := range users {
		if value.CreatedTasks.Id == 0 {
			continue
		}
		text, err := processingTemplate("Owner", value)
		if err != nil {
			return nil, err
		}
		supportive.Message = supportive.Message + text + "\n"
	}
	out = append(out, supportive)
	return out, nil
}
