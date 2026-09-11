-- Turned off on purpose.
--
-- LazyVim's own spec and the language extras pull these in. Removing an
-- `import` line is not enough for anything in the core spec, so it is said
-- here instead, in one place, with the reason.
--
-- Every one is a single `enabled = true` away from coming back.

return {
  -- Motion and operator tooling. Excellent if you drive the editor by
  -- keystrokes; unused if you drive it by mouse and arrow keys.
  { "folke/flash.nvim", enabled = false }, -- jump to a two-character label
  { "nvim-mini/mini.ai", enabled = false }, -- around/inside text objects
  { "nvim-treesitter/nvim-treesitter-textobjects", enabled = false }, -- vif, daf, cia

  -- A second row of tabs, below the row of tabs Ghostty already draws, above
  -- a file tree that lists the same files. Ghostty's own tabs and splits do
  -- this job.
  { "akinsho/bufferline.nvim", enabled = false },

  -- The debugger. Four plugins and a UI for something this machine does from
  -- a terminal. `dap.core` in config/lazy.lua brings the whole set back.
  { "mfussenegger/nvim-dap", enabled = false },
  { "leoluz/nvim-dap-go", enabled = false },
  { "mfussenegger/nvim-dap-python", enabled = false },

  -- LazyVim ships these so a fresh install has a theme before you choose
  -- one. This config has chosen one: colors/gotham.lua.
  { "catppuccin/nvim", enabled = false },
  { "folke/tokyonight.nvim", enabled = false },
}
