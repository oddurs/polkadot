-- Bootstrap lazy.nvim, then the spec.
--
-- Built on LazyVim 16 with blink.cmp for completion and snacks.nvim for the
-- picker, so no extra here brings in telescope, fzf-lua or nvim-cmp — mixing
-- those with the defaults is how a LazyVim config ends up fighting itself.

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
  -- Go, Rust, Python, Zig, plus the formats every project carries.
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

  -- Editing
  "coding.mini-surround", -- cs"' ds( ysiw)
  "coding.yanky", -- a yank ring, with a picker over it

  -- Moving around
  "editor.harpoon2", -- pin the four files you are actually working in
  "editor.aerial", -- symbol outline
  "editor.inc-rename", -- LSP rename with live preview
  "editor.refactoring", -- extract function / variable (needs nvim 0.12+)
  "editor.illuminate", -- highlight the other uses of the symbol under the cursor
  "editor.dial", -- ctrl-a/ctrl-x that understands dates, booleans, hex
  "editor.mini-diff", -- inline diff against the index

  -- Quality gates
  "formatting.prettier",
  "linting.eslint",
  "test.core",
  "dap.core",

  -- Utility
  "util.dot", -- syntax for the dotfiles this repo is made of
  "util.mini-hipatterns", -- render #0a0f14 in its own colour
  "util.octo", -- GitHub issues and PRs in the editor, via gh
  "util.project",

  -- UI
  "ui.treesitter-context", -- keep the enclosing function on screen
  "ui.indent-blankline",

  -- Copilot is here because there is a copilot-chat licence on this machine
  -- already. Run :Copilot auth once; delete these two lines to drop it.
  "ai.copilot",
  "ai.copilot-chat",
}

local spec = { { "LazyVim/LazyVim", import = "lazyvim.plugins" } }
for _, e in ipairs(extras) do
  table.insert(spec, { import = "lazyvim.plugins.extras." .. e })
end
table.insert(spec, { import = "plugins" })

require("lazy").setup({
  spec = spec,
  defaults = { lazy = false, version = false },
  install = { colorscheme = { "gotham256", "habamax" } },
  checker = { enabled = true, notify = false },
  change_detection = { notify = false },
  ui = { border = "rounded", backdrop = 100 },
  performance = {
    rtp = {
      disabled_plugins = {
        "gzip", "tarPlugin", "tohtml", "tutor", "zipPlugin", "rplugin", "netrwPlugin",
      },
    },
  },
})
