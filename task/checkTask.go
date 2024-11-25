package task

import (
	"errors"
	"time"

	"final/date"
)

func (t *Task) Checktitle() error {
	if t.Title == "" {
		return errors.New("Пустой заголовок")
	}
	return nil
}

func (t *Task) Checkdate() (Task, error) {
	now := time.Now()
	if t.Date == "" {
		t.Date = now.Format(date.ParseDate)
		return *t, nil
	}

	parsedDate, err := time.Parse(date.ParseDate, t.Date)
	if err != nil {
		return *t, errors.New("Неправильный формат даты")
	}

	if parsedDate.Before(now) {
		nowStr := now.Format(date.ParseDate)
		if t.Repeat == "" {
			t.Date = nowStr
		} else if nowStr != t.Date {
			nextDate, calcErr := date.CalcNextDate(nowStr, t.Date, t.Repeat)
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
		now := time.Now().Format(date.ParseDate)
		nextDate, err := date.CalcNextDate(now, t.Date, t.Repeat)
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
		if _, err := date.ParseRepeatRules(t.Repeat); err != nil {
			return "Правило повторения указано в неправильном формате"
		}
	}
	return ""
}
