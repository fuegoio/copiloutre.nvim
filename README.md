# Copiloutre

## How to install

With `lazy.nvim`:

```
  {
    "fuegoio/copiloutre.nvim",
    opts = {
      debug = false,
    },
  },
  {
    "zbirenbaum/copilot.lua",
    cmd = "Copilot",
    event = "InsertEnter",
    dependencies = {
      "fuegoio/copiloutre.nvim",
    },
    config = function()
      require("copilot").setup {
        panel = {
          enabled = false,
        },
        suggestion = {
          auto_trigger = true,
          keymap = {
            accept = "<C-/>",
          },
        },
      }
    end,
  },

```
