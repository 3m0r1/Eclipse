function OnLoad()
    print('Loading terminal plugin')
end

function OnUnload()
    print('unloading terminal plugin')
end

return {
    Metadata = {
        Name = 'Utils',
        Description = 'Plugin with helpful utilities',
        Author = '3m0r1',
        Version = '1.0.0'
    },

    Events = {
        OnLoad = OnLoad,
        OnUnload = OnUnload,
    },

    Commands = {
        HelloMessage = {
            Description = 'Returns a hello message',
            Use = 'HelloMessage [message]',
            Args = {
                {
                    Name = 'message',
                    Default = 'DefaultMessage',
                    Optional = false,
                }
            },

            Strict = true, -- allow only the exact number of argsument(s)?
            Return = true, -- are there return(s)?
            Export = true, -- make it importable by other plugin(s)?

            -- ctx is the caller (itself or other plugins)
            Run = function(ctx, args)
                local msg = 'Hello ' .. args[1]
                return { msg }
            end
        }
    },

    Imports = {}
}
