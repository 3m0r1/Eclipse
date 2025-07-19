package plugin

import lua "github.com/yuin/gopher-lua"

func GetField(state *lua.LState, obj lua.LValue, key string) lua.LValue {
	return state.GetField(obj, key)
}

func IsNil(value lua.LValue) bool {
	return value == lua.LNil
}

func GetStringOr(state *lua.LState, obj lua.LValue, key, def string) string {
	field := GetField(state, obj, key)

	if !IsNil(field) {
		return lua.LVAsString(field)
	} else {
		return def
	}
}

func GetStringOrPanic(state *lua.LState, obj lua.LValue, key, message string) string {
	field := GetField(state, obj, key)

	if !IsNil(field) {
		return lua.LVAsString(field)
	} else {
		panic(message)
	}
}

func GetBoolOr(state *lua.LState, obj lua.LValue, key string, def bool) bool {
	field := GetField(state, obj, key)

	if !IsNil(field) {
		return lua.LVAsBool(field)
	} else {
		return def
	}
}

func GetBoolOrPanic(state *lua.LState, obj lua.LValue, key, message string) bool {
	field := GetField(state, obj, key)

	if !IsNil(field) {
		return lua.LVAsBool(field)
	} else {
		panic(message)
	}
}

func GetFunctionOrPanic(state *lua.LState, obj lua.LValue, key, message string) *lua.LFunction {
	field := GetField(state, obj, key)

	if !IsNil(field) {
		return field.(*lua.LFunction)
	} else {
		panic(message)
	}
}

func GetTableOrPanic(state *lua.LState, obj lua.LValue, key, message string) *lua.LTable {
	field := GetField(state, obj, key)

	if !IsNil(field) {
		return field.(*lua.LTable)
	} else {
		panic(message)
	}
}

func MakeCallableTable(state *lua.LState, table *lua.LTable, callFn lua.LGFunction) *lua.LTable {
	mt := state.NewTable()
	state.SetField(mt, "__call", state.NewFunction(callFn))
	state.SetMetatable(table, mt)
	return mt
}
