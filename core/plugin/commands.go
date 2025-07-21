package plugin

import (
	"errors"
	"fmt"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

type Argument struct {
	Name     string
	Default  lua.LValue
	Optional bool
}

type Command struct {
	Name             string
	Description      string
	Use              string
	Args             []*Argument
	ArgCount         int
	RequiredArgCount int
	Strict           bool
	Return           bool
	Export           bool
	RunFn            *lua.LFunction
}

type Export Command

func MakeCommand(state *lua.LState, name string, entry *lua.LTable) *Command {
	description := GetStringOr(
		state,
		entry,
		"Description",
		"N/A",
	)

	use := GetStringOr(
		state,
		entry,
		"Use",
		"N/A",
	)

	strict := GetBoolOr(
		state,
		entry,
		"Strict",
		false,
	)

	ret := GetBoolOr(
		state,
		entry,
		"Return",
		false,
	)

	export := GetBoolOr(
		state,
		entry,
		"Export",
		false,
	)

	runFn := GetFunctionOrPanic(
		state,
		entry,
		"Run",
		"Couldn't find run function",
	)

	cmd := Command{
		Name:        name,
		Description: description,
		Use:         use,
		Strict:      strict,
		Return:      ret,
		Export:      export,
		RunFn:       runFn,
		Args:        []*Argument{},
	}

	luaArgs := GetTableOrPanic(
		state,
		entry,
		"Args",
		"Couldn't find arguments table",
	)

	luaArgs.ForEach(func(key, value lua.LValue) {
		cmd.ArgCount++

		luaArg := value.(*lua.LTable)

		name := GetStringOr(
			state,
			luaArg,
			"Name",
			"N/A",
		)

		defaultVal := GetField(
			state,
			luaArg,
			"Default",
		)

		optional := GetBoolOr(
			state,
			luaArg,
			"Optional",
			false,
		)

		if !optional {
			cmd.RequiredArgCount++
		}

		arg := Argument{
			Name:     name,
			Default:  defaultVal,
			Optional: optional,
		}

		cmd.Args = append(cmd.Args, &arg)
	})

	return &cmd
}

func (cmd *Command) Invoke(state *lua.LState, ctx lua.LValue, args ...lua.LValue) (*lua.LTable, error) {

	if cmd.RequiredArgCount > len(args) {
		return nil, errors.New("amount of arguments less than required")
	}

	if cmd.Strict {
		if cmd.ArgCount != len(args) {
			return nil, errors.New("amount of arguments more than expected")
		}
	}

	argsTbl := state.NewTable()

	for index, arg := range args {
		if index < len(cmd.Args) {
			cmdArg := cmd.Args[index]

			if arg == lua.LNil {
				if cmdArg.Optional {
					argsTbl.Append(cmdArg.Default)
				} else {
					return nil, errors.New("required value is nil")
				}

			} else {
				argsTbl.Append(arg)
			}
		} else {
			argsTbl.Append(arg)
		}

	}

	if cmd.Return {
		err := state.CallByParam(
			lua.P{
				Fn:      cmd.RunFn,
				NRet:    1,
				Protect: true,
			},
			ctx,
			argsTbl,
		)
		ret := state.ToTable(-1)
		state.Pop(1)
		return ret, err
	} else {
		err := state.CallByParam(
			lua.P{
				Fn:      cmd.RunFn,
				NRet:    0,
				Protect: true,
			},
			ctx,
			argsTbl,
		)
		return nil, err
	}
}

func (cmd *Command) Help() {
	lines := []string{
		fmt.Sprintf("Name:        %s", cmd.Name),
		fmt.Sprintf("Description: %s", cmd.Description),
		fmt.Sprintf("Usage:       %s", cmd.Use),
		"\n",
		"Arguments:",
	}

	for _, arg := range cmd.Args {

		var argDefault string
		var argOpt string

		if arg.Default == lua.LNil {
			argDefault = "None"
		} else {
			argDefault = arg.Default.String()
		}

		if arg.Optional {
			argOpt = "Optional"
		} else {
			argOpt = "Required"
		}

		lines = append(lines, fmt.Sprintf("  %-10s (%s) Default: %s", arg.Name, argOpt, argDefault))
	}

	fmt.Println(strings.Join(lines, "\n"))
}
