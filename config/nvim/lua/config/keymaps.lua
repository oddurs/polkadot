-- Keys for driving this editor with a mouse and the shortcuts every other
-- macOS application uses.
--
-- Each action is bound twice: once to the Cmd chord and once to the Ctrl
-- chord. Ghostty forwards Cmd combinations as CSI-u sequences carrying the
-- super modifier, which Neovim decodes to `<D-s>` and friends; see the keys
-- section of config/ghostty/config. The Ctrl form is what works in any other
-- terminal, and over ssh.
--
-- Cmd+C, Cmd+V and Cmd+X are deliberately *not* forwarded. Ghostty already
-- copies the selection and pastes through bracketed paste, which Neovim
-- handles correctly in every mode, and taking those keys would break copy
-- and paste everywhere else in the terminal.
--
-- Shift+arrow selection is not here because it is an option rather than a
-- mapping: `keymodel=startsel,stopsel` in options.lua.

local map = vim.keymap.set

--- Bind the same action to the Cmd chord and the Ctrl chord.
---@param modes string|string[]
---@param key string  the bare key, e.g. "s" for Cmd+S / Ctrl+S
---@param rhs string|function
---@param desc string
local function both(modes, key, rhs, desc)
  map(modes, "<D-" .. key .. ">", rhs, { desc = desc })
  map(modes, "<C-" .. key .. ">", rhs, { desc = desc })
end

-- ── the universal six ──────────────────────────────────────────────────────
-- <Cmd> mappings do not change mode, so saving from insert leaves you in
-- insert with the cursor where it was.
both({ "n", "i", "v", "s" }, "s", "<Cmd>silent! write<CR>", "Save file")
both({ "n", "i", "v", "s" }, "z", "<Cmd>silent! undo<CR>", "Undo")

-- Redo is Cmd+Shift+Z, which arrives as the capital letter.
map({ "n", "i", "v", "s" }, "<D-Z>", "<Cmd>silent! redo<CR>", { desc = "Redo" })
map({ "n", "i", "v", "s" }, "<C-y>", "<Cmd>silent! redo<CR>", { desc = "Redo" })

both("n", "a", "<Cmd>keepjumps normal! ggVG<CR>", "Select all")
both("i", "a", "<Esc><Cmd>keepjumps normal! ggVG<CR>", "Select all")

-- Copy and cut from a selection. Pasting is the system clipboard already,
-- because `clipboard=unnamedplus`; this only spells it the usual way.
map("v", "<C-c>", '"+y', { desc = "Copy" })
map("v", "<C-x>", '"+d', { desc = "Cut" })
map("i", "<C-v>", "<C-r><C-o>+", { desc = "Paste" })
map("c", "<C-v>", "<C-r>+", { desc = "Paste" })

-- ── finding things ─────────────────────────────────────────────────────────
both("n", "p", function() Snacks.picker.files() end, "Find file")
both("n", "f", function() Snacks.picker.lines() end, "Find in this file")
both("n", "b", "<Cmd>Neotree toggle<CR>", "Toggle file tree")

map("n", "<D-F>", function() Snacks.picker.grep() end, { desc = "Find in project" })
map("n", "<C-g>", function() Snacks.picker.grep() end, { desc = "Find in project" })

-- Comment. Cmd+/ arrives as <D-/>; terminals disagree about what Ctrl+/ is,
-- so both of its encodings are bound too.
for _, lhs in ipairs({ "<D-/>", "<C-_>", "<C-/>" }) do
  map("n", lhs, "gcc", { desc = "Toggle comment", remap = true })
  map("i", lhs, "<Esc>gcca", { desc = "Toggle comment", remap = true })
  map("v", lhs, "gc", { desc = "Toggle comment", remap = true })
end

-- ── editing, the way a normal editor does it ───────────────────────────────
-- Tab indents the selection and keeps it selected, instead of replacing it.
map("v", "<Tab>", ">gv", { desc = "Indent" })
map("v", "<S-Tab>", "<gv", { desc = "Outdent" })

-- Alt+arrows move the current line or selection, reindenting as they land.
map("n", "<A-Up>", "<Cmd>silent! move .-2<CR>==", { desc = "Move line up" })
map("n", "<A-Down>", "<Cmd>silent! move .+1<CR>==", { desc = "Move line down" })
map("i", "<A-Up>", "<Esc><Cmd>silent! move .-2<CR>==gi", { desc = "Move line up" })
map("i", "<A-Down>", "<Esc><Cmd>silent! move .+1<CR>==gi", { desc = "Move line down" })
map("v", "<A-Up>", ":move '<-2<CR>gv=gv", { silent = true, desc = "Move selection up" })
map("v", "<A-Down>", ":move '>+1<CR>gv=gv", { silent = true, desc = "Move selection down" })

-- Ctrl+Backspace deletes the word behind the cursor; most terminals send
-- that as Ctrl+H.
map("i", "<C-h>", "<C-w>", { desc = "Delete word back" })
map("i", "<C-BS>", "<C-w>", { desc = "Delete word back" })

-- Word-wise motion in insert mode. Normal mode already has these.
map("i", "<C-Left>", "<C-o>b", { desc = "Word left" })
map("i", "<C-Right>", "<C-o>w", { desc = "Word right" })

-- ── odds and ends ──────────────────────────────────────────────────────────
map("t", "<C-\\><C-n>", "<C-\\><C-n>", { desc = "Terminal: normal mode" })

map("n", "<leader>fy", function()
  local path = vim.fn.expand("%:.")
  vim.fn.setreg("+", path)
  vim.notify("Copied " .. path)
end, { desc = "Copy relative path" })

-- ── right-click menu ───────────────────────────────────────────────────────
-- `mousemodel=popup_setpos` moves the cursor to the click before opening
-- this, so every entry acts on whatever is under the pointer.
pcall(vim.cmd, "aunmenu PopUp")
vim.cmd([[
  anoremenu PopUp.Go\ to\ definition   <Cmd>lua vim.lsp.buf.definition()<CR>
  anoremenu PopUp.Find\ references     <Cmd>lua Snacks.picker.lsp_references()<CR>
  anoremenu PopUp.Rename\ symbol       <Cmd>lua vim.lsp.buf.rename()<CR>
  anoremenu PopUp.Code\ action         <Cmd>lua vim.lsp.buf.code_action()<CR>
  amenu     PopUp.-1-                  <Nop>
  vnoremenu PopUp.Copy                 "+y
  vnoremenu PopUp.Cut                  "+d
  anoremenu PopUp.Paste                "+p
  amenu     PopUp.-2-                  <Nop>
  anoremenu PopUp.Format\ file         <Cmd>lua require("conform").format({ lsp_format = "fallback" })<CR>
]])
