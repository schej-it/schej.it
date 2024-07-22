package main

import (
	"context"
	"time"

	"github.com/emersion/go-webdav"
	"github.com/emersion/go-webdav/caldav"
	"schej.it/server/utils"
)

func main() {
	httpClient := webdav.HTTPClientWithBasicAuth(nil, "liu.z.jonathan@gmail.com", "eypf-izki-chlg-cyzj")

	webdavClient, err := webdav.NewClient(httpClient, "https://caldav.icloud.com")
	if err != nil {
		panic(err)
	}

	caldavClient, err := caldav.NewClient(httpClient, "https://caldav.icloud.com")
	if err != nil {
		panic(err)
	}

	principal, err := webdavClient.FindCurrentUserPrincipal(context.Background())
	if err != nil {
		panic(err)
	}

	calendarHomeSet, err := caldavClient.FindCalendarHomeSet(context.Background(), principal)
	if err != nil {
		panic(err)
	}

	calendars, err := caldavClient.FindCalendars(context.Background(), calendarHomeSet)
	if err != nil {
		panic(err)
	}

	utils.PrintJson(calendars)

	for _, calendar := range calendars {
		if calendar.Name != "idk" {
			continue
		}

		events, err := caldavClient.QueryCalendar(context.Background(), calendar.Path, &caldav.CalendarQuery{
			CompRequest: caldav.CalendarCompRequest{
				Name: "VCALENDAR",
				Comps: []caldav.CalendarCompRequest{{
					Name: "VEVENT",
					Props: []string{
						"SUMMARY",
						"UID",
						"DTSTART",
						"DTEND",
						"DURATION",
					},
				}},
			},
			CompFilter: caldav.CompFilter{
				Name: "VCALENDAR",
				Comps: []caldav.CompFilter{{
					Name:  "VEVENT",
					Start: time.Now(),
					End:   time.Now().Add(time.Hour * 7 * 24),
					Props: []caldav.PropFilter{{
						Name: "DTSTART",
						ParamFilter: []caldav.ParamFilter{{
							Name: "VALUE",
							TextMatch: &caldav.TextMatch{
								Text:            "DATE",
								NegateCondition: true,
							},
						}},
					}},
				}},
			},
		})
		if err != nil {
			panic(err)
		}

		utils.PrintJson(events)
	}
}
