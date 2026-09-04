-- Loaded by LazyVim before plugins. Options only — keymaps live in keymaps.lua.

local opt = vim.opt

opt.relativenumber = true -- counts for motions; absolute on the cursor line
opt.scrolloff = 8 -- never read the last line at the very bottom
opt.sidescrolloff = 8
opt.wrap = false
opt.linebreak = true -- if wrap is turned on, break at words
opt.cursorline = true
opt.colorcolumn = "" -- a ruler is noise; formatters already enforce width
opt.signcolumn = "yes" -- stop the text jumping when a sign appears
opt.splitkeep = "screen"
opt.confirm = true -- ask rather than fail on :q with unsaved changes
opt.undofile = true
opt.undolevels = 10000
opt.updatetime = 200
opt.timeoutlen = 400

-- Search
opt.ignorecase = true
opt.smartcase = true
opt.inccommand = "split" -- live preview of :s

-- Indentation
opt.expandtab = true
opt.shiftwidth = 2
opt.tabstop = 2
opt.smartindent = true

-- Whitespace made visible, quietly.
opt.list = true
opt.listchars = { tab = "→ ", trail = "·", nbsp = "␣", extends = "›", precedes = "‹" }
opt.fillchars = { eob = " ", fold = " ", foldopen = "▾", foldclose = "▸", diff = "╱" }

-- Splits open where the eye expects them.
opt.splitright = true
opt.splitbelow = true

opt.termguicolors = true
opt.laststatus = 3 -- one status line for the whole window, not per split
opt.pumheight = 12
opt.clipboard = "unnamedplus"
