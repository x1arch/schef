package schef

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// format: [#{tag}] [limitters] [duration or modifier]+[duration or modifier]-[duration or modifier]...
// tag: starts from # string - identify a another one rule
// limitters: {YYYY-MM-DD[:HH:MM[:SS]]}@from, {YYYY-MM-DD[:HH:MM[:SS]]}@to, {number}@times
// duration: {number or empty}@second, {number or empty}@minute, {number or empty}@hour, {number or empty}@day, {number or empty}@week, {number or empty}@month, {number or empty}@year
// modifiers {number or L(last) or F(first) or R(random)}@day_of_week, {number or L(last) or F(first) or R(random)}@day_of_month, {number or L(last) or F(first) or R(random)}@day_of_year, {hh:mm}@time, @date
// examples:
// `@day` - every day
// `2025-02-03@from 5@times @month` - starts from 2025-02-03 every month repeat 5 times, that's means 2025-07-03 will the latest execution
// `30@day_of_month` - every month 30 the day (the february will skipped)
// `2@weeks` - the same as `@month - 14@days` the same `14@days`
// `@month #onerror @day` - execute every month, on error execute every day
//
// Limitations

var (
	ErrWrongScheduleFormat = errors.New("wrong schedule format")
	ErrTagNotFound         = errors.New("tag not found")
)

type Schef struct {
	times        int
	modificators map[string][]func(time.Time) (time.Time, error)
}

func Parse(schedule string) (*Schef, error) {
	c := &Schef{modificators: map[string][]func(time.Time) (time.Time, error){}}
	schedule = strings.Trim(schedule, " ")

	tmp := ""
	tmpv := []string{}
	for i := 0; i < len(schedule); i++ {
		switch schedule[i : i+1] {
		case "#":
			if tmp != "" {
				tmpv = append(tmpv, tmp)
			}
			tmp = schedule[i : i+1]

		case "@":
			tmpv = append(tmpv, tmp)
			tmp = schedule[i : i+1]

		case " ":
			if tmp != "" {
				tmpv = append(tmpv, tmp)
				tmp = ""
			}

		case "+":
			if tmp != "" {
				tmpv = append(tmpv, tmp)
				tmp = ""
			}
			tmpv = append(tmpv, "+")

		case "-":
			if tmp != "" {
				tmpv = append(tmpv, tmp)
				tmp = ""
			}
			tmpv = append(tmpv, "-")

		default:
			tmp += schedule[i : i+1]
		}
	}
	tmpv = append(tmpv, tmp) // add last rule

	args := []string{}
	tag := "default"
	for _, v := range tmpv {
		if len(v) > 0 {
			if v[0:1] == "@" {
				_ = c.parts(tag, v, args...)
				args = []string{}
			} else if v[0:1] == "#" {
				tag = v[1:]
			} else {
				args = append(args, v)
			}
		} else {
			args = append(args, v)
		}
	}

	return c, nil
}

func (c *Schef) parts(tag, modifier string, args ...string) error {
	if _, ok := c.modificators[tag]; !ok {
		c.modificators[tag] = []func(time.Time) (time.Time, error){}
	}

	switch modifier {
	// periodic description
	case "@times":
		var err error
		if len(args) != 1 {
			return ErrWrongScheduleFormat
		}

		if c.times, err = strconv.Atoi(args[0]); err != nil {
			return ErrWrongScheduleFormat
		}

	// todo
	// case "@from":
	// case "@to":
	// case "@times":

	case "@second":
		var err error
		argJ := strings.Join(args, "")
		v := 0
		if argJ == "" {
			v = 1
		} else if v, err = strconv.Atoi(argJ); err != nil {
			return ErrWrongScheduleFormat
		}

		c.modificators[tag] = append(c.modificators[tag], func(t time.Time) (time.Time, error) {
			return t.Add(time.Duration(time.Duration(v) * time.Second)), nil
		})

	case "@minute":
		var err error
		argJ := strings.Join(args, "")
		v := 0
		if argJ == "" {
			v = 1
		} else if v, err = strconv.Atoi(argJ); err != nil {
			return ErrWrongScheduleFormat
		}

		c.modificators[tag] = append(c.modificators[tag], func(t time.Time) (time.Time, error) {
			return t.Add(time.Duration(time.Duration(v) * time.Minute)), nil
		})

	case "@hour":
		var err error
		argJ := strings.Join(args, "")
		v := 0
		if argJ == "" {
			v = 1
		} else if v, err = strconv.Atoi(argJ); err != nil {
			return ErrWrongScheduleFormat
		}

		c.modificators[tag] = append(c.modificators[tag], func(t time.Time) (time.Time, error) {
			return t.Add(time.Duration(time.Duration(v) * time.Hour)), nil
		})

	case "@day":
		var err error
		argJ := strings.Join(args, "")
		v := 0
		if argJ == "" {
			v = 1
		} else if v, err = strconv.Atoi(argJ); err != nil {
			return ErrWrongScheduleFormat
		}

		c.modificators[tag] = append(c.modificators[tag], func(t time.Time) (time.Time, error) {
			return t.AddDate(0, 0, v), nil
		})

	case "@week":
		var err error
		argJ := strings.Join(args, "")
		v := 0
		if argJ == "" {
			v = 1
		} else if v, err = strconv.Atoi(argJ); err != nil {
			return ErrWrongScheduleFormat
		}

		c.modificators[tag] = append(c.modificators[tag], func(t time.Time) (time.Time, error) {
			return t.AddDate(0, 0, v*7), nil
		})

	case "@month":
		var err error
		argJ := strings.Join(args, "")
		v := 0
		if argJ == "" {
			v = 1
		} else if v, err = strconv.Atoi(argJ); err != nil {
			return ErrWrongScheduleFormat
		}

		c.modificators[tag] = append(c.modificators[tag], func(t time.Time) (time.Time, error) {
			return t.AddDate(0, v, 0), nil
		})

	case "@year":
		var err error
		argJ := strings.Join(args, "")
		v := 0
		if argJ == "" {
			v = 1
		} else if v, err = strconv.Atoi(argJ); err != nil {
			return ErrWrongScheduleFormat
		}

		c.modificators[tag] = append(c.modificators[tag], func(t time.Time) (time.Time, error) {
			return t.AddDate(v, 0, 0), nil
		})

		// case "@day_of_week":
		// case "@day_of_month":
		// case "@day_of_year":
	}

	return nil
}

func (c *Schef) NextDateFromNow() (time.Time, error) {
	return c.NextDate(time.Now())
}

func (c *Schef) NextDateTagFromNow(tag string) (time.Time, error) {
	return c.NextDateTag(tag, time.Now())
}

func (c *Schef) NextDate(t time.Time) (time.Time, error) {
	return c.NextDateTag("default", t)
}

func (c *Schef) NextDateTag(tag string, t time.Time) (time.Time, error) {
	if _, ok := c.modificators[tag]; !ok {
		return time.Time{}, ErrTagNotFound
	}

	var err error
	for i, m := range c.modificators[tag] {
		if t, err = m(t); err != nil {
			return time.Time{}, errors.Join(fmt.Errorf("Broken on step %d", i), err)
		}
	}
	return t, nil
}
