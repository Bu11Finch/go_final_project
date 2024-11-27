package handlers

import (
	"final/date"
	"net/http"
)

func (h *Handlers) GetnextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")
	if repeat == "" || nowStr == "" || dateStr == "" {
		http.Error(w, "Указаны некорректные данные в запросе", http.StatusBadRequest)
		return
	}
	nextdate, err := date.CalcNextDate(nowStr, dateStr, repeat)
	if err != nil {
		http.Error(w, "Ошибка вычисления даты", http.StatusBadRequest)
		return
	}

	w.Write([]byte(nextdate))

}
