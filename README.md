# 🦦 copiloutre.nvim

Easily use any model, local or distant, with Github Copilot plugin extensions on Neovim.

## What is this plugin?

[zbirenbaum/copilot.lua](https://github.com/zbirenbaum/copilot.lua) is the best AI autocomplete plugin on Neovim. Other plugins like [yetone/avante.nvim](https://github.com/yetone/avante.nvim) are good AI companions, but they struggle to make autocomplete as good as Github Copilot extension. This is mainly because Github Copilot leverages *FIM (Fill-In-Middle) Completion* and not Chat completions like other plugins like [yetone/avante.nvim](https://github.com/yetone/avante.nvim), therefore the model does not have to output in a specific format but can just output code.

But Github Copilot uses a lot of telemetry, train on your code and does not allow you to use your local models or third-party ones. This plugin is made to allow you to use their plugin but with your own models.

## How does it work?

[zbirenbaum/copilot.lua](https://github.com/zbirenbaum/copilot.lua) and the [official copilot.nvim made by Github](https://github.com/github/copilot.vim) embark a NodeJS server that acts as the brain of the extension. This server will call Github servers for authentication, model selection, telemetry and completion.

This extension launches a local proxy in background and patches the Copilot NodeJS server at runtime to interact with this proxy. This proxy will answer to the Copilot NodeJS server and forwards completion requests to the model of your choice. Currently, only one implementation is supported: *Mistral AI FIM Completions API*.

## Get started

This plugin will patch the NodeJS server installed by [zbirenbaum/copilot.lua](https://github.com/zbirenbaum/copilot.lua). You only need to add this plugin as a dependency to launch it right before, like this with `lazy.nvim`:

```lua
  {
    "zbirenbaum/copilot.lua",
    cmd = "Copilot",
    event = "InsertEnter",
    dependencies = {
      "fuegoio/copiloutre.nvim", -- The important part
    },
    config = function()
      require("copilot").setup {}
    end,
  },

```

You also with the current MistralAI implementation to supply via an environment variable a `MISTRAL_API_KEY`:

```bash
export MISTRAL_API_KEY=your-api-key
```
