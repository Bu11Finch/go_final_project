package handlers

import (
	"final/calcdate"
	"fmt"
	"net/http"
)

func (h *Handlers) GetnextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	// Вычисление следующей даты с использованием функции из calcdate
	nextDate, err := calcdate.CalcNextDate(nowStr, dateStr, repeat)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка вычисления даты: %v", err), http.StatusBadRequest)
		return
	}

	// Устанавливаем заголовок ответа и отправляем строку даты напрямую
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, nextDate)
}
