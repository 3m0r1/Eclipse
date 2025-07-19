function OnLoad()
    print('Loading basic plugin')
end

function OnUnload()
    print('Unloading basic plugin')
end

function GetLocalPlugin()
    return _G.Plugin
end

function OnReady()
    print('Basic plugin is ready')
    local plugin = GetLocalPlugin()

    -- you can now call imports safely
    local res = plugin.Imports.HelloMessage(plugin.Metadata.Name)
    local msg = res[1]

    print('[OnReady] HelloMessage:', msg)
end

return {
    Metadata = {
        Name = 'BasicPlugin',
        Description = 'This is a basic plugin',
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
                local res = ctx.Imports.HelloMessage(args[1])
                local msg = res[1]
                print('[Greet] HelloMessage:', msg)
            end
        }
    },

    Imports = {
        HelloMessage = {
            Plugin = 'Utils',
            Procedure = 'HelloMessage'
        }
    }
}