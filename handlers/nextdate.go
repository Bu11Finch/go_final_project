package handlers

import (
	"final/date"
	"net/http"
)

func (h *Handlers) GetnextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	_, err := date.CalcNextDate(nowStr, dateStr, repeat)
	if err != nil {
		http.Error(w, "Ошибка вычисления даты", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

}
