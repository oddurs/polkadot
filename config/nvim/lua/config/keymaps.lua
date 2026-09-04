-- Additions to LazyVim's defaults. Nothing here overrides a builtin.

local map = vim.keymap.set

-- Centre the view when jumping through search results and half-pages, so the
-- match is never at the very edge of the screen.
map("n", "n", "nzzzv")
map("n", "N", "Nzzzv")
map("n", "<C-d>", "<C-d>zz")
map("n", "<C-u>", "<C-u>zz")

-- Move the selection, keeping indentation correct.
map("v", "J", ":m '>+1<CR>gv=gv", { desc = "Move selection down" })
map("v", "K", ":m '<-2<CR>gv=gv", { desc = "Move selection up" })

-- Paste over a selection without losing the register.
map("x", "<leader>p", [["_dP]], { desc = "Paste without yanking" })

-- Yank to the system clipboard explicitly.
map({ "n", "v" }, "<leader>y", [["+y]], { desc = "Yank to clipboard" })

-- Keep the cursor in place when joining lines.
map("n", "J", "mzJ`z")

-- Escape from the terminal split.
map("t", "<C-\\><C-n>", "<C-\\><C-n>", { desc = "Terminal: normal mode" })

-- Quick file-relative work.
map("n", "<leader>fy", function()
  local path = vim.fn.expand("%:.")
  vim.fn.setreg("+", path)
  vim.notify("Copied " .. path)
end, { desc = "Copy relative path" })
