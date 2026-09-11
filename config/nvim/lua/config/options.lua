-- Loaded by LazyVim before plugins. Options only — keymaps live in keymaps.lua.

local opt = vim.opt

-- ── the editor as an editor ────────────────────────────────────────────────
-- Absolute line numbers. Relative numbering exists to make `7j` countable,
-- and counting motions is not how this editor is driven.
opt.number = true
opt.relativenumber = false

-- Shift+arrows start and extend a selection, the way every other editor on
-- this machine behaves. `startsel` begins one, `stopsel` ends it when an
-- unshifted motion follows. Selection lands in visual mode rather than select
-- mode, so ordinary commands still operate on it.
opt.keymodel = "startsel,stopsel"
opt.selectmode = ""

-- Mouse everywhere, with a right-click menu that first moves the cursor to
-- where you clicked.
opt.mouse = "a"
opt.mousemodel = "popup_setpos"
opt.mousescroll = "ver:3,hor:4"

-- ── presentation ───────────────────────────────────────────────────────────
opt.cursorline = true
opt.wrap = false
opt.linebreak = true -- if wrap is turned on, break at words, not mid-token
opt.smoothscroll = true -- wrapped lines scroll by screen line, not by buffer line
opt.scrolloff = 8 -- never read the last line at the very bottom
opt.sidescrolloff = 8
opt.colorcolumn = "" -- a ruler is noise; formatters already enforce width
opt.signcolumn = "yes" -- stop the text jumping when a sign appears
opt.winborder = "rounded" -- one border style for every float, set once
opt.pumheight = 12
opt.pumblend = 0 -- a translucent menu over code is harder to read
opt.laststatus = 3 -- one status line for the whole window, not per split
opt.splitkeep = "screen"
opt.splitright = true
opt.splitbelow = true
opt.termguicolors = true
-- Gotham is a dark palette and has no light counterpart, so say so rather
-- than making Neovim ask the terminal and wait for an answer.
opt.background = "dark"

-- Whitespace made visible, quietly. Tabs and trailing spaces are worth
-- seeing; a dot in every space is not.
opt.list = true
opt.listchars = { tab = "→ ", trail = "·", nbsp = "␣", extends = "›", precedes = "‹" }
opt.fillchars = { eob = " ", fold = " ", foldopen = "▾", foldclose = "▸", diff = "╱" }

-- ── behaviour ──────────────────────────────────────────────────────────────
opt.confirm = true -- ask rather than fail on :q with unsaved changes
opt.undofile = true
opt.undolevels = 10000
opt.updatetime = 200
opt.timeoutlen = 400

opt.ignorecase = true
opt.smartcase = true
opt.inccommand = "split" -- live preview of :s

opt.expandtab = true
opt.shiftwidth = 2
opt.tabstop = 2
opt.smartindent = true

opt.clipboard = "unnamedplus" -- one clipboard, shared with the rest of macOS
