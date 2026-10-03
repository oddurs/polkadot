-- Bootstrap lazy.nvim, then the spec.
--
-- Built on LazyVim with blink.cmp for completion and snacks.nvim for the
-- picker, so no extra here brings in telescope, fzf-lua or nvim-cmp — mixing
-- those with the defaults is how a LazyVim config ends up fighting itself.
--
-- What is deliberately absent: the vim-operator tooling. Surround, yank
-- rings, jump-by-label, text objects and harpoon are excellent if you drive
-- the editor by motions, and dead weight if you drive it by mouse and arrow
-- keys. This config is the second kind. See lua/config/keymaps.lua.

local lazypath = vim.fn.stdpath("data") .. "/lazy/lazy.nvim"
if not (vim.uv or vim.loop).fs_stat(lazypath) then
  local out = vim.fn.system({
    "git", "clone", "--filter=blob:none", "--branch=stable",
    "https://github.com/folke/lazy.nvim.git", lazypath,
  })
  if vim.v.shell_error ~= 0 then
    vim.api.nvim_echo({
      { "Failed to clone lazy.nvim:\n", "ErrorMsg" },
      { out, "WarningMsg" },
      { "\nPress any key to exit..." },
    }, true, {})
    vim.fn.getchar()
    os.exit(1)
  end
end
vim.opt.rtp:prepend(lazypath)

local extras = {
  -- Languages, chosen from what is actually in ~/Code: Next.js and TypeScript,
  -- Go, Rust, Python, Zig, plus the formats every project carries. These are
  -- server and parser configuration; they cost nothing until a matching file
  -- is opened.
  "lang.typescript",
  "lang.tailwind",
  "lang.json",
  "lang.go",
  "lang.rust",
  "lang.python",
  "lang.zig",
  "lang.markdown",
  "lang.yaml",
  "lang.toml",
  "lang.docker",
  "lang.sql",
  "lang.git",

  -- Quality gates
  "formatting.prettier",
  "linting.eslint",

  -- Utility
  "editor.inc-rename", -- rename a symbol with a live preview of the edit
  "util.dot", -- syntax for the dotfiles this repo is made of
  "util.mini-hipatterns", -- render #0a0f14 in its own colour
  "util.project", -- project detection for the picker

  -- UI
  "ui.treesitter-context", -- keep the enclosing function on screen
}

local spec = { { "LazyVim/LazyVim", import = "lazyvim.plugins" } }
for _, e in ipairs(extras) do
  table.insert(spec, { import = "lazyvim.plugins.extras." .. e })
end
table.insert(spec, { import = "plugins" })

require("lazy").setup({
  spec = spec,
  -- `lazy = true` is the whole point of lazy.nvim. The previous value here was
  -- `false`, which forced every plugin that had not declared its own trigger
  -- to load during startup whether the session needed it or not.
  defaults = { lazy = true, version = false },
  install = { colorscheme = { "paris-guimard", "habamax" } },
  -- The update check is a network round trip on every start, to tell you
  -- something you can ask for with :Lazy whenever you feel like it.
  checker = { enabled = false },
  change_detection = { notify = false },
  ui = { border = "rounded", backdrop = 100 },
  performance = {
    rtp = {
      disabled_plugins = {
        "gzip",
        "matchit",
        "netrwPlugin",
        "rplugin",
        "tarPlugin",
        "tohtml",
        "tutor",
        "zipPlugin",
      },
    },
  },
})
