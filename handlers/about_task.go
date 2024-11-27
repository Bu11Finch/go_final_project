package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	nextdate "final/date"

	"final/task"
)

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
		var task task.Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			writeJSONResponse(w, nil, err)
			return
		}

		err = task.Checktitle()
		if err != nil {
			writeJSONResponse(w, nil, err)
			return
		}

		taskmod, err := task.Checkdate()
		if err != nil {
			writeJSONResponse(w, nil, err)
			return
		}

		err = taskmod.Countdate()
		if err != nil {
			writeJSONResponse(w, nil, err)
			return
		}

		id, err := h.TaskStorage.AddTaskToDataBase(taskmod)
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
		if err != nil {
			writeJSONResponse(w, nil, err)
			return
		}

		writeJSONResponse(w, map[string]interface{}{"tasks": tasks}, nil)
	}
}

func (h *Handlers) GetTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.FormValue("id")
		if id == "" {
			writeJSONResponse(w, nil, errors.New("Не указан идентификатор"))
			return
		}

		task, err := h.TaskStorage.FindTask(id)
		if err != "" {
			writeJSONResponse(w, nil, errors.New(err))
			return
		}

		writeJSONResponse(w, task, nil)
	}
}

func (h *Handlers) EditTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task task.Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			writeJSONResponse(w, nil, err)
			return
		}

		errstr := task.CheckId()
		if errstr != "" {
			writeJSONResponse(w, nil, errors.New(errstr))
			return
		}

		err = task.Checktitle()
		if err != nil {
			writeJSONResponse(w, nil, err)
			return
		}

		task, err = task.Checkdate()
		if err != nil {
			writeJSONResponse(w, nil, err)
			return
		}

		errstr = task.CheckRepeat()
		if errstr != "" {
			writeJSONResponse(w, nil, errors.New(errstr))
			return
		}

		errstr = h.TaskStorage.UpdateTask(task)
		if errstr != "" {
			writeJSONResponse(w, nil, errors.New(errstr))
			return
		}

		writeJSONResponse(w, map[string]interface{}{}, nil)
	}
}

func (h *Handlers) MarkTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.FormValue("id")
		if id == "" {
			writeJSONResponse(w, nil, errors.New("Не указан идентификатор"))
			return
		}

		task, err := h.TaskStorage.FindTask(id)
		if err != "" {
			writeJSONResponse(w, nil, errors.New("Задача не найдена"))
			return
		}

		if task.Repeat == "" {
			err = h.TaskStorage.DeleteTask(id)
			if err != "" {
				writeJSONResponse(w, nil, errors.New("Ошибка удаления задачи"))
				return
			}
		} else {
			now := time.Now()
			timeNow := now.Format(nextdate.ParseDate)
			date, errnotstr := nextdate.CalcNextDate(timeNow, task.Date, task.Repeat)
			if errnotstr != nil {
				writeJSONResponse(w, nil, errors.New("Ошибка вычисления следующей даты"))
				return
			}

			err = h.TaskStorage.UpdateTaskDate(date, id)
			if err != "" {
				writeJSONResponse(w, nil, errors.New("Ошибка обновления задачи"))
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
			writeJSONResponse(w, nil, errors.New("Не указан идентификатор"))
			return
		}

		err := h.TaskStorage.DeleteTask(id)
		if err != "" {
			writeJSONResponse(w, nil, errors.New("Ошибка удаления задачи"))
			return
		}

		writeJSONResponse(w, map[string]interface{}{}, nil)
	}
}
