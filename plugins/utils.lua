function OnLoad()
    print('Loading utils plugin')
end

function OnUnload()
    print('Unloading utils plugin')
end

function OnReady()
    print('Utils plugin is ready')
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
        OnReady = OnReady
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

            -- ctx refers to the local plugin (most of the time) commands are invoked
            Run = function(ctx, args)
                local msg = 'Hello ' .. args[1]
                return { msg }
            end
        },

         Add = {
            Description = 'Adds two numbers together',
            Use = 'Add [number1] [number2]',
            Args = {
                {
                    Name = 'number1',
                    Optional = false,
                },
                {
                    Name = 'number2',
                    Optional = false,
                }
            },

            Strict = true, -- allow only the exact number of argsument(s)?
            Return = true, -- are there return(s)?
            Export = true, -- make it importable by other plugin(s)?

            -- ctx refers to the local plugin (most of the time) commands are invoked
            Run = function(ctx, args)
                local result = args[1] + args[2]
                return { result }
            end
        }
    },

    Imports = {}
}
