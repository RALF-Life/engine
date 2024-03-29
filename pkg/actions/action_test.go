package actions

import (
	ics "github.com/darmiel/golang-ical"
	"testing"
)

type test struct {
	action  string
	event   func() *ics.VEvent
	error   bool
	message ActionMessage
	with    map[string]interface{}
	check   func(event *ics.VEvent) bool
}

func getAction(name string) (Action, bool) {
	for _, a := range Actions {
		if a.Identifier() == name {
			return a, true
		}
	}
	return nil, false
}

func TestMixed(t *testing.T) {
	cases := []*test{
		// filters/filter-in doesn't need an Event to work
		{
			action:  "filters/filter-in",
			message: new(FilterInActionMessage),
		},
		// filters/filter-out doesn't need an Event to work
		{
			action:  "filters/filter-out",
			message: new(FilterOutActionMessage),
		},
		{
			action:  "actions/regex-replace",
			message: nil,
			event: func() *ics.VEvent {
				event := ics.NewEvent("a")
				event.SetSummary("Hello World!")
				return event
			},
			with: map[string]interface{}{
				"match":   "Hello ",
				"replace": "",
				"in":      []interface{}{"summary"},
			},
			check: func(event *ics.VEvent) bool {
				prop := event.GetProperty(ics.ComponentPropertySummary)
				if prop == nil {
					return false
				}
				return prop.Value == "World!"
			},
		},
		{
			action: "actions/clear-alarms",
			event: func() *ics.VEvent {
				event := ics.NewEvent("a")
				event.AddAlarm()
				return event
			},
			check: func(event *ics.VEvent) bool {
				return len(event.Alarms()) == 0
			},
		},
		{
			action: "actions/add-alarm",
			with: map[string]interface{}{
				"action":  "display",
				"trigger": "1d",
			},
			event: func() *ics.VEvent {
				return ics.NewEvent("ok")
			},
			check: func(event *ics.VEvent) bool {
				return len(event.Alarms()) == 1
			},
		},
		{
			action: "actions/clear-attendees",
			event: func() *ics.VEvent {
				event := ics.NewEvent("a")
				event.AddAttendee("aaa")
				return event
			},
			check: func(event *ics.VEvent) bool {
				return len(event.Alarms()) == 0
			},
		},
		{
			action: "actions/add-attendee",
			with: map[string]interface{}{
				"mail": "me@example.com",
			},
			event: func() *ics.VEvent {
				return ics.NewEvent("ok")
			},
			check: func(event *ics.VEvent) bool {
				return len(event.Attendees()) == 1
			},
		},
		{
			action: "ctx/set",
			with: map[string]interface{}{
				"hello": "world",
			},
			error: false,
		},
	}
	for i, c := range cases {
		action, exists := getAction(c.action)
		if !exists {
			t.Fatalf("cannot find action %s", c.action)
		}
		var event *ics.VEvent
		if c.event != nil {
			event = c.event()
		}
		sharedContext := make(map[string]interface{})
		ctx := &Context{
			Event:         event,
			SharedContext: sharedContext,
			With:          c.with,
			Verbose:       false,
		}
		resp, err := action.Execute(ctx)
		if err == nil && c.error {
			t.Fatalf("expected error for test %d but no returned", i+1)
		} else if err != nil && !c.error {
			t.Fatalf("got error %v for test %d but no error expected", err, i+1)
		}
		if resp != c.message {
			t.Fatalf("expected return %v but got %v", c.message, resp)
		}
		if c.check != nil && !c.check(event) {
			t.Fatalf("check failed for test %d", i+1)
		}
	}
}
