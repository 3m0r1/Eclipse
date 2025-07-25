
# Eclipse

An extensible Lua plugin manager in Go.

<img width="1531" height="1131" alt="eclipse" src="https://github.com/user-attachments/assets/674fac4d-f187-44e1-8132-2a77c4e32d22" />

## Background
Eclipse started out as a way to interpret lua plugins as efficiently as possible to allow for dynamic behavior for systems like C2s, where extensibility matters alot.

It focuses on 4 main components when it comes to plugin design:

- Events
- Commands
- Imports
- Exports

It was designed for the [Abyss C2 Framework](https://github.com/AbyssFramework).

## Features

#### Events
- Lifecycle events (``OnLoad``, ``OnUnload``, ``OnReady``)
- Custom events

#### Commands
- Variadiac & optional arguments
- Strict mode
- Variadiac returns

#### Imports
- Clean syntax using ``__call`` metamethod

#### Exports
- Ability to export commands with the ``Export`` field
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
    print('Loading greet plugin')
end

function OnUnload()
    print('Unloading greet plugin')
end

function OnReady(plugin)
    print('greet plugin is ready')
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
For more details, check out the [Plugins](https://github.com/3m0r1/Eclipse/tree/main/plugins/) folder.

## Usage
```
./eclipse [plugins]
```

## Resources
These resources were crucial during the development of Eclipse.

- [GopherLua](https://github.com/yuin/gopher-lua)
- [Metatables & Metamethods](https://www.lua.org/pil/13.html)
