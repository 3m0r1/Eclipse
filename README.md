
# Eclipse

Lua plugin manager in Go.

## Background
Eclipse started out as a way to interpret lua plugins as efficiently as possible to allow for dynamic behavior for systems like C2s, where extensibility matters alot.

Eclipse focuses on 4 main components when it comes to plugin design:

- Events
- Commands
- Imports
- Exports

Eclipse was designed for the [Abyss C2 Framework](https://github.com/AbyssFramework).


## Example

Manager:
```go
mgr := manager.NewPluginManager()

basicPlugin := plugin.NewPlugin("./plugins/basic.lua")
mgr.LoadPlugins(basicPlugin)

basicPlugin.InvokeCommand("Greet", lua.LString("World"))
```

Plugin:
```lua
function OnLoad()
    print('Loading basic plugin')
end

function OnUnload()
    print('Unloading basic plugin')
end

function OnReady()
    print('Basic plugin is ready')
end

return {
    Metadata = {
        Name = 'GreetPlugin',
        Description = 'Plugin that has commands related to greeting',
        Author = '3m0r1',
        Version = '1.0.0'
    },

    Events = {
        OnLoad = OnLoad,
        OnUnload = OnUnload,
        OnReady = OnReady
    },

    Commands = {
        Greet = {
            Description = 'Greets the user',
            Use = 'Greet [user]',
            Args = {
                {
                    Name = 'user',
                    Default = 'DefaultUser',
                    Optional = true,
                }
            },
            Strict = true, -- allow only the exact number of argsument(s)?
            Return = false, -- are there return(s)?
            Export = false, -- make it importable by other plugin(s)?

           -- ctx refers to the local plugin (most of the time) commands are invoked
            Run = function(ctx, args)
                print('Hello:', args[1])
            end
        }
    },
	Imports = { }
}

```

## Usage
```
./eclipse [plugins]
```
