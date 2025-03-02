package schef

import (
	"fmt"
	"testing"
	"time"
)

func ExampleParse() {
	cl, _ := Parse("45@day")

	fmt.Println(cl.NextDate(time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)))
	// Output: 2009-12-25 23:00:00 +0000 UTC <nil>
}

func TestTagRule(t *testing.T) {
	sf, err := Parse("@month #onerror @day")
	if err != nil {
		t.Error(err)
	}

	now := time.Now()
	nd, err := sf.NextDate(now)
	if err != nil || nd.Compare(now.AddDate(0, 1, 0)) != 0 {
		t.Error(err)
	}

	nd, err = sf.NextDateTag("onerror", now)
	if err != nil || nd.Compare(now.AddDate(0, 0, 1)) != 0 {
		t.Error(err)
	}

	t.Log(sf)
}

func TestParseDuration(t *testing.T) {
	now := time.Now()

	var tests = []struct {
		rule   string
		result time.Time
	}{
		{"", now},

		// raw
		{"@second", now.Add(time.Second)},
		{"@minute", now.Add(time.Minute)},
		{"@hour", now.Add(time.Hour)},
		{"@day", now.AddDate(0, 0, 1)},
		{"@week", now.AddDate(0, 0, 7)},
		{"@month", now.AddDate(0, 1, 0)},
		{"@year", now.AddDate(1, 0, 0)},

		// number
		{"45@second", now.Add(45 * time.Second)},
		{"45@minute", now.Add(45 * time.Minute)},
		{"45@hour", now.Add(45 * time.Hour)},
		{"45@day", now.AddDate(0, 0, 45)},
		{"45@week", now.AddDate(0, 0, 45*7)},
		{"45@month", now.AddDate(0, 45, 0)},
		{"45@year", now.AddDate(45, 0, 0)},

		// decrease
		{"10@second-5@second", now.Add(5 * time.Second)},
		{"10@minute-5@minute", now.Add(5 * time.Minute)},
		{"10@hour-5@hour", now.Add(5 * time.Hour)},
		{"10@day-5@day", now.AddDate(0, 0, 5)},
		{"10@week-5@week", now.AddDate(0, 0, 5*7)},
		{"10@month-5@month", now.AddDate(0, 5, 0)},
		{"10@year-5@year", now.AddDate(5, 0, 0)},

		// summ
		{"10@second+5@second", now.Add(15 * time.Second)},
		{"10@minute+5@minute", now.Add(15 * time.Minute)},
		{"10@hour+5@hour", now.Add(15 * time.Hour)},
		{"10@day+5@day", now.AddDate(0, 0, 15)},
		{"10@week+5@week", now.AddDate(0, 0, 15*7)},
		{"10@month+5@month", now.AddDate(0, 15, 0)},
		{"10@year+5@year", now.AddDate(15, 0, 0)},
	}

	for _, tt := range tests {
		t.Run(tt.rule, func(t *testing.T) {
			sf, err := Parse(tt.rule)
			if err != nil {
				t.Error(err)
			}

			if nd, err := sf.NextDate(now); err != nil || nd.Compare(tt.result) != 0 {
				t.Error(err)
			}

			t.Log(sf)
		})
	}
}

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Parse("45@day")
	}
}
