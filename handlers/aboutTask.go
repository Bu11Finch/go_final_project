package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"final/calcdate"
	"final/task"
)

const ParseDate = "20060102"

func writeJSONResponse(w http.ResponseWriter, data interface{}, err error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
	} else {
		json.NewEncoder(w).Encode(data)
	}
}

func (h *Handlers) AddTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var taskInput task.Task
		if err := json.NewDecoder(r.Body).Decode(&taskInput); err != nil {
			writeJSONResponse(w, nil, err)
			return
		}

		if err := taskInput.Checktitle(); err != nil {
			writeJSONResponse(w, nil, err)
			return
		}

		taskMod, err := taskInput.Checkdate()
		if err != nil {
			writeJSONResponse(w, nil, err)
			return
		}

		if err := taskMod.Countdate(); err != nil {
			writeJSONResponse(w, nil, err)
			return
		}

		id, err := h.TaskStorage.AddTaskToDataBase(taskMod)
		if err != nil {
			writeJSONResponse(w, nil, err)
			return
		}

		writeJSONResponse(w, map[string]interface{}{"id": id}, nil)
	}
}

func (h *Handlers) GetTasks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tasks, err := h.TaskStorage.GetTasks()
		writeJSONResponse(w, map[string]interface{}{"tasks": tasks}, err)
	}
}

func (h *Handlers) GetTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.FormValue("id")
		if id == "" {
			writeJSONResponse(w, nil, errors.New("Не указан идентификатор"))
			return
		}

		taskData, err := h.TaskStorage.FindTask(id)
		if err != "" {
			writeJSONResponse(w, nil, errors.New(err))
			return
		}

		writeJSONResponse(w, taskData, nil)
	}
}

func (h *Handlers) EditTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var taskInput task.Task
		if err := json.NewDecoder(r.Body).Decode(&taskInput); err != nil {
			writeJSONResponse(w, nil, errors.New("Ошибка десериализации JSON"))
			return
		}

		if errMsg := taskInput.CheckId(); errMsg != "" {
			writeJSONResponse(w, nil, errors.New(errMsg))
			return
		}

		if err := taskInput.Checktitle(); err != nil {
			writeJSONResponse(w, nil, err)
			return
		}

		updatedTask, err := taskInput.Checkdate()
		if err != nil {
			writeJSONResponse(w, nil, err)
			return
		}

		if errMsg := updatedTask.CheckRepeate(); errMsg != "" {
			writeJSONResponse(w, nil, errors.New(errMsg))
			return
		}

		if errMsg := h.TaskStorage.UpdateTask(updatedTask); errMsg != "" {
			writeJSONResponse(w, nil, errors.New(errMsg))
			return
		}

		writeJSONResponse(w, map[string]interface{}{}, nil)
	}
}

func (h *Handlers) MarkTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.FormValue("id")
		if id == "" {
			http.Error(w, `{"error": "Не указан индентификатор"}`, http.StatusBadRequest)
			return
		}

		taskData, err := h.TaskStorage.FindTask(id)
		if err != "" {
			http.Error(w, `{"error": "Задача не найдена"}`, http.StatusInternalServerError)
			return
		}

		if taskData.Repeat == "" {
			if err := h.TaskStorage.DeleteTask(id); err != "" {
				http.Error(w, `{"error": "Ошибка удаления задачи"}`, http.StatusInternalServerError)
			}
		} else {
			now := time.Now().Format(ParseDate)
			nextDate, calcErr := calcdate.CalcNextDate(now, taskData.Date, taskData.Repeat)
			if calcErr != nil {
				http.Error(w, `{"error": "Ошибка вычисления следующей даты"}`, http.StatusInternalServerError)
				return
			}
			if err := h.TaskStorage.UpdateTaskDate(nextDate, id); err != "" {
				http.Error(w, `{"error": "Ошибка обновления задачи"}`, http.StatusInternalServerError)
				return
			}
		}
		writeJSONResponse(w, map[string]interface{}{}, nil)
	}
}

func (h *Handlers) DeleteTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, `{"error": "Не указан индентификатор"}`, http.StatusBadRequest)
			return
		}

		if err := h.TaskStorage.DeleteTask(id); err != "" {
			http.Error(w, `{"error": "Ошибка удаления задачи"}`, http.StatusInternalServerError)
			return
		}

		writeJSONResponse(w, map[string]interface{}{}, nil)
	}
}
