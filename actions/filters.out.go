package actions

import ics "github.com/darmiel/golang-ical"

type FilterOutAction struct{}

func (foa *FilterOutAction) Identifier() string {
	return "filters/filter-out"
}

///

var DummyFilterOutMessage = FilterOutMessage(0)

func (foa *FilterOutAction) Execute(_ *ics.VEvent, _ map[string]interface{}) (ActionMessage, error) {
	return DummyFilterOutMessage, nil
}
