package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/seletskiy/tplutil"
)

var tpl = template.Must(template.New("index").Parse(tplutil.Strip(`
	<!DOCTYPE html>
	<html>
		<head>
			<script async src="https://www.googletagmanager.com/gtag/js?id=UA-147344009-1"></script>
			<script>
				window.dataLayer = window.dataLayer || [];
				function gtag(){dataLayer.push(arguments);}
				gtag('js', new Date());
				gtag('config', 'UA-147344009-1');
			</script>

			<style>
				.is_open {
					position: fixed;
					top: 50%;
					left: 50%;
					transform: translate(-50%, -50%);
				}

				.is_open span {
					display: inline-block;
					vertical-align: middle;
					line-height: normal;
					font-family: sans;
					font-weight: bold;
					font-size: 36px;
					padding: 10px;
					text-transform: uppercase;
				}

				.is_open span.open {
					color: forestgreen;
					background-color: springgreen;
				}

				.is_open span.closed {
					color: darkred;
					background-color: salmon;
				}
			</style>
			<title>Is NYSE open right now?</title>
		</head>

		<body>
			<div class="is_open">
				{{ if .open }}
					<span class="open">open</span>
				{{ else }}
					<span class="closed">closed</span>
				{{ end }}
			</div>
		</body>
	</html>
`)))

type (
	Status int

	Schedule       map[int]ScheduleYearly
	ScheduleYearly map[time.Month]ScheduleMontly
	ScheduleMontly map[int]Status
)

const (
	StatusOpen        Status = 1
	StatusClosedEarly        = 2
	StatusClosed             = 3
)

var (
	schedule = Schedule{
		2020: ScheduleYearly{
			time.January: ScheduleMontly{
				1:  StatusClosed,
				20: StatusClosed,
			},
			time.February: ScheduleMontly{
				17: StatusClosed,
			},
			time.April: ScheduleMontly{
				10: StatusClosed,
			},
			time.May: ScheduleMontly{
				25: StatusClosed,
			},
			time.July: ScheduleMontly{
				3: StatusClosed,
			},
			time.September: ScheduleMontly{
				7: StatusClosed,
			},
			time.November: ScheduleMontly{
				26: StatusClosedEarly,
			},
			time.December: ScheduleMontly{
				25: StatusClosedEarly,
			},
		},
		2019: ScheduleYearly{
			time.November: ScheduleMontly{
				28: StatusClosedEarly,
			},
			time.December: ScheduleMontly{
				25: StatusClosedEarly,
			},
		},
	}
)

func main() {
	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc(
		"/",
		func(writer http.ResponseWriter, request *http.Request) {
			err := tpl.Execute(
				writer,
				map[string]interface{}{
					"open": isOpen(time.Now().In(ny)),
				},
			)
			if err != nil {
				log.Fatal(err)
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}
		},
	)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func isOpen(now time.Time) bool {
	switch now.Weekday() {
	case time.Sunday, time.Saturday:
		return false
	}

	switch {
	case now.Hour() < 9 || now.Hour() > 16:
		return false
	case now.Hour() > 9 && now.Minute() < 30:
		return false
	}

	if yearly, ok := schedule[now.Year()]; ok {
		if monthly, ok := yearly[now.Month()]; ok {
			if daily, ok := monthly[now.Day()]; ok {
				switch {
				case daily == StatusClosedEarly:
					return now.Hour() < 13
				case daily == StatusClosed:
					return false
				}
			}
		}
	}

	return true
}
