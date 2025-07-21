function OnLoad()
    print('Loading basic plugin')
end

function OnUnload()
    print('Unloading basic plugin')
end

function OnReady(plugin)
    print('Basic plugin is ready')

    -- you can now call imports safely
    local result = plugin.Imports.HelloMessage(plugin.Metadata.Name)
    local msg = result[1]

    print('[OnReady] HelloMessage:', msg)

    local result = plugin.Imports.Add(5, 5)
    local number = result[1]

    print('[OnReady] Add (5 + 5): ' .. number)
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

    Commands = { },

    Imports = {
        HelloMessage = {
            Plugin = 'Utils',
            Procedure = 'HelloMessage'
        },

        Add = {
            Plugin = 'Utils',
            Procedure = 'Add'
        }
    }
}