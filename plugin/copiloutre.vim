" Title:        copiloutre.vim
" Author:       Alexis Tacnet <https://github.com/fuegoio>

" Prevents the plugin from being loaded multiple times. If the loaded
" variable exists, do nothing more. Otherwise, assign the loaded
" variable and continue running this instance of the plugin.
if exists("g:loaded_copiloutre")
    finish
endif
let g:loaded_copiloutre = 1

" Defines a package path for Lua. This facilitates importing the
" Lua modules from the plugin's dependency directory.
let s:lua_rocks_deps_loc =  expand("<sfile>:h:r") . "/../lua/copiloutre/deps"
exe "lua package.path = package.path .. ';" . s:lua_rocks_deps_loc . "/lua-?/init.lua'"

" Exposes the plugin's functions for use as commands in Neovim.
command! -nargs=0 CopiloutreStart lua require("copiloutre").start()
command! -nargs=0 CopiloutreStop lua require("copiloutre").stop()
command! -nargs=0 CopiloutreStatus lua require("copiloutre").status()
