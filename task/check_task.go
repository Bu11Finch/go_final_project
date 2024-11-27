package task

import (
	"errors"
	nextdate "final/date"
	"time"
)

func (t Task) Checktitle() error {
	if t.Title == "" {
		return errors.New("Пустой заголовок")
	}
	return nil
}

func (t Task) Checkdate() (Task, error) {

	now := time.Now()
	if t.Date == "" {
		t.Date = now.Format(nextdate.ParseDate)
		return t, nil
	} else {
		date, err := time.Parse(nextdate.ParseDate, t.Date)
		if err != nil {
			return t, errors.New("Неправильный формат даты")
		}
		if date.Before(now) {
			if t.Repeat == "" {
				t.Date = now.Format(nextdate.ParseDate)
				return t, nil
			} else {
				nowtime := now.Format(nextdate.ParseDate)
				if nowtime != t.Date {
					nextDate, err := nextdate.CalcNextDate(nowtime, t.Date, t.Repeat)
					if err != nil {
						return t, errors.New("Ошибка вычисления даты")
					}
					t.Date = nextDate
					return t, nil
				} else {
					t.Date = nowtime
					return t, nil
				}
			}
		}
	}
	return t, nil
}

func (t Task) Countdate() error {
	if t.Repeat != "" {
		now := time.Now()
		nowtime := now.Format(nextdate.ParseDate)
		nextDate, err := nextdate.CalcNextDate(nowtime, t.Date, t.Repeat)
		if err != nil {
			return errors.New("Ошибка вычисления даты")
		}
		t.Date = nextDate
	}
	return nil
}

func (t Task) CheckId() string {
	if t.ID == "" {
		return "Не указан индентификатор задачи"
	} else {
		return ""
	}
}

func (t Task) CheckRepeat() string {
	if t.Repeat != "" {
		_, err := nextdate.ParseRepeatRules(t.Repeat)
		if err != nil {
			return "Правило повторения указано в неправильном формате"
		}
	}
	return ""
}
