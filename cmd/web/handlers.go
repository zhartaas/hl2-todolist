package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"sort"
	"time"
)

// @summary Create Task
// @tags ToDo list
// @param request body TaskInput true "insert task title and date"
// @success 200 {object} TaskResponse
// @failure 404 {string} string "client error"
// @failure 500 {string} string "internal server error"
// @router /create [post]
func (app *application) createTask(w http.ResponseWriter, r *http.Request) {
	task, err := readAndValidateTask(r)
	if err != nil {
		app.errorLog.Println(err)
		app.handleError(w, err)
		return
	}

	taskID := uuid.New()
	task.ID = taskID.String()
	LocalMap.Store(taskID.String(), task)
	response, err := json.Marshal(&TaskResponse{ID: taskID.String()})
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

// @summary Update Task
// @tags ToDo list
// @param id query string true "insert task id"
// @param request body TaskInput true "insert task title and date"
// @success 200 {object} TaskResponse
// @failure 404 {string} string "client error"
// @failure 500 {string} string "internal server error"
// @router /update [put]
func (app *application) updateTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.FormValue("id")
	savedTask, exists := LocalMap.Load(taskID)
	if !exists {
		app.errorLog.Println("task doesn't exist")
		app.handleError(w, errors.New("task doesn't exists or incorrect id"))
	}
	taskFromMap := savedTask.(*Task)

	task, err := readAndValidateTask(r)
	if err != nil {
		app.errorLog.Println(err)
		app.handleError(w, err)
		return
	}
	taskFromMap.Title = task.Title
	taskFromMap.ActiveAt = task.ActiveAt

	w.WriteHeader(http.StatusNoContent)
}

func readAndValidateTask(r *http.Request) (*Task, error) {
	task := &Task{}

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, task)
	if err != nil {
		return nil, err
	}
	_, err = time.Parse(time.DateOnly, task.ActiveAt)
	titleIsValid := len(task.Title) <= 200
	if err != nil || !titleIsValid {
		return nil, errors.New("invalid input")
	}

	unique := true
	LocalMap.Range(func(key, value any) bool {
		currentTask := value.(*Task)
		unique = (currentTask.Title != task.Title) || (currentTask.ActiveAt != task.ActiveAt)
		return unique
	})
	if !unique {
		return nil, errors.New("task already exists")
	}

	return task, nil
}

// @summary Get Tasks
// @tags ToDo list
// @param status query string false "active or done"
// @success 200 {object} GetTasksResponse
// @failure 404 {string} string "client error"
// @failure 500 {string} string "internal server error"
// @router /getTasks [get]
func (app *application) getTasks(w http.ResponseWriter, r *http.Request) {
	status := r.FormValue("status")

	paramDone := false
	if status == "done" {
		paramDone = true
	}

	currentDate := time.Now()
	tasks := GetTasksResponse{}
	LocalMap.Range(func(key, value any) bool {
		task := value.(*Task)

		if task.Done != paramDone {
			return true
		}

		taskDate, _ := time.Parse(time.DateOnly, task.ActiveAt)
		if paramDone == false {
			if currentDate.After(taskDate) {
				return true
			}
		}
		if taskDate.Weekday() == time.Weekday(0) || taskDate.Weekday() == time.Weekday(6) {
			task.Title = fmt.Sprintf("ВЫХОДНОЙ - %s", task.Title)
		}

		tasks = append(tasks, task)
		return true
	})

	sort.Sort(tasks)
	response, err := json.MarshalIndent(tasks, "", "\t")
	if err != nil {
		app.errorLog.Println(err)
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// @summary Delete Tasks
// @tags ToDo list
// @param id query string true "insert task id"
// @success 204
// @failure 404 {string} string "client error"
// @failure 500 {string} string "internal server error"
// @router /delete [delete]
func (app *application) deleteTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.FormValue("id")

	_, exists := LocalMap.Load(taskID)
	if !exists {
		app.errorLog.Println("doesn't exist")

		app.handleError(w, errors.New("task doesn't exists or incorrect id"))
		return
	}

	LocalMap.Delete(taskID)

	w.WriteHeader(http.StatusNoContent)
}

// @summary Task Done
// @tags ToDo list
// @param id query string true "insert task id"
// @success 204
// @failure 404 {string} string "client error"
// @failure 500 {string} string "internal server error"
// @router /taskDone [put]
func (app *application) taskDone(w http.ResponseWriter, r *http.Request) {
	taskID := r.FormValue("id")

	taskFromMap, exists := LocalMap.Load(taskID)
	if !exists {
		app.errorLog.Println("doesn't exist")
		app.handleError(w, errors.New("task doesn't exists or incorrect id"))
		return
	}

	task := taskFromMap.(*Task)
	task.Done = true

	w.WriteHeader(http.StatusNoContent)
}
