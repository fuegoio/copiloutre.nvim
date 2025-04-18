local M = {}

local server_running = false
local server_process = nil

local function start(opts)
  opts = opts or {}
  if server_running then
    print("🦦 copiloutre server is already running.")
    return
  end

  local server_path = vim.api.nvim_get_runtime_file("proxy/proxy", false)[1]
  server_process = vim.fn.jobstart(server_path, {
    on_stdout = function(_, data)
      if data and opts.debug then
        print(table.concat(data, "\n"))
      end
    end,
    on_stderr = function(_, data)
      if data and opts.debug then
        print(table.concat(data, "\n"))
      end
    end,
    on_exit = function(_, _data)
      server_running = false
      print("🦦 copiloutre server stopped.")
    end,
  })

  if server_process ~= -1 then
    server_running = true
    print("🦦 copiloutre server started.")
  else
    print("Failed to start copiloutre server.")
  end
end

local function stop()
  if not server_running then
    print("🦦 copiloutre server is not running.")
    return
  end

  vim.fn.jobstop(server_process)
  server_running = false
  print("🦦 copiloutre server stopped.")
end

local function status()
  if server_running then
    print("🦦 copiloutre server is running.")
  else
    print("🦦 copiloutre server is not running.")
  end
end

local function setup(opts)
  start(opts)
end

M.start = start
M.stop = stop
M.status = status

M.setup = setup

return M
