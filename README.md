# schef (schedule format)

The lib helps build native rules for shedulers.

### format
> [#{tag}] [limitters] [duration or modifier]+[duration or modifier]-[duration or modifier]...

If the first tag not setted (string starts from rule), the rule will have `default` tag. 

`[]` - optional param


### tag
Starts from # string, identify a another one rule in string.

Example: `@month #onerror @day`

Execute every month, if got an error reqest `onerror` tag rule.

```golang
sf, err := Parse("@month #onerror @day")
if err != nil {
    ...
}

now := time.Now()
nd, err := sf.NextDate(now)
if err != nil || nd.Compare(now.AddDate(0, 1, 0)) != 0 {
    ...
}

nd, err = sf.NextDateTag("onerror", now)
if err != nil || nd.Compare(now.AddDate(0, 0, 1)) != 0 {
    ...
}

```

<!-- ### limitters
Alow to apply limits to rule

> {YYYY-MM-DD[:HH:MM[:SS]]}@from, {YYYY-MM-DD[:HH:MM[:SS]]}@to, {number}@times -->

### duration

> {number or empty}@second, {number or empty}@minute, {number or empty}@hour, {number or empty}@day, {number or empty}@week, {number or empty}@month, {number or empty}@year

<!-- 
### modifiers
> {number or L(last) or F(first) or R(random)}@day_of_week, {number or L(last) or F(first) or R(random)}@day_of_month, {number or L(last) or F(first) or R(random)}@day_of_year, {hh:mm}@time -->

## examples:
`@day` - every day

`2025-02-03@from 5@times @month` - starts from 2025-02-03 every month repeat 5 times, that's means 2025-07-03 will the latest execution

`30@day_of_month` - every month 30 the day (the february will skipped)

`2@weeks` - the same as `@month - 14@days` the same `14@days`

# Limitations
coming soon

## Usage

```golang
package main

import (
    "log"
    "strings"
    "time"

    "github.com/x1arch/schef"
)

func main() {
    str := "@month #onerror @day"
    sf, err := schef.Parse(str)
    log.Println(sf.NextDateFromNow())
    log.Println(sf.NextDateTagFromNow("onerror"))
    log.Println(sf.NextDate(time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)))
    log.Println(sf.NextDateTag("onerror", time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)))    
}
```