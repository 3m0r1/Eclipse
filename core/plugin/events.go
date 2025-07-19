package plugin

import lua "github.com/yuin/gopher-lua"

type Event struct {
	Name string
	Fn   *lua.LFunction
}

func MakeEvent(name string, fn *lua.LFunction) *Event {
	event := Event{
		Name: name,
		Fn:   fn,
	}

	return &event
}

func (event *Event) Fire(state *lua.LState, args ...lua.LValue) {
	state.CallByParam(
		lua.P{
			Fn:      event.Fn,
			NRet:    0,
			Protect: false,
		},
		args...,
	)
}
