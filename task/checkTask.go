package task

import (
	"errors"
	"time"

	nextdate "final/calcdate"
)

const ParseDate = "20060102"

func (t *Task) Checktitle() error {
	if t.Title == "" {
		return errors.New("Пустой заголовок")
	}
	return nil
}

func (t *Task) Checkdate() (Task, error) {
	now := time.Now()
	if t.Date == "" {
		t.Date = now.Format(ParseDate)
		return *t, nil
	}

	parsedDate, err := time.Parse(ParseDate, t.Date)
	if err != nil {
		return *t, errors.New("Неправильный формат даты")
	}

	if parsedDate.Before(now) {
		nowStr := now.Format(ParseDate)
		if t.Repeat == "" {
			t.Date = nowStr
		} else if nowStr != t.Date {
			nextDate, calcErr := nextdate.CalcNextDate(nowStr, t.Date, t.Repeat)
			if calcErr != nil {
				return *t, errors.New("Ошибка вычисления даты")
			}
			t.Date = nextDate
		} else {
			t.Date = nowStr
		}
	}
	return *t, nil
}

func (t *Task) Countdate() error {
	if t.Repeat != "" {
		now := time.Now().Format(ParseDate)
		nextDate, err := nextdate.CalcNextDate(now, t.Date, t.Repeat)
		if err != nil {
			return errors.New("Ошибка вычисления даты")
		}
		t.Date = nextDate
	}
	return nil
}

func (t *Task) CheckId() string {
	if t.ID == "" {
		return "Не указан индентификатор задачи"
	}
	return ""
}

func (t *Task) CheckRepeate() string {
	if t.Repeat != "" {
		if _, err := nextdate.ParseRepeatRules(t.Repeat); err != nil {
			return "Правило повторения указано в неправильном формате"
		}
	}
	return ""
}
