function OnLoad()
    print('Loading basic plugin')
end

function OnUnload()
    print('Unloading basic plugin')
end

function OnCustom()
    print('calling custom event')
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

        OnCustom = OnCustom
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

            -- ctx is the caller (itself or other plugins)
            Run = function(ctx, args)
                local res = ctx.Imports.HelloMessage(args[1])
                print('Returned Message: ' .. res[1])
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