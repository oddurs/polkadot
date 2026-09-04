-- On top of LazyVim's own. Each of these earns its keep daily.

local function aug(name)
  return vim.api.nvim_create_augroup("polkadot_" .. name, { clear = true })
end

-- Flash what was just yanked, so you can see the extent of it.
vim.api.nvim_create_autocmd("TextYankPost", {
  group = aug("yank"),
  callback = function() vim.highlight.on_yank({ timeout = 150 }) end,
})

-- Reopen a file where you left it.
vim.api.nvim_create_autocmd("BufReadPost", {
  group = aug("last_position"),
  callback = function(ev)
    if vim.bo[ev.buf].filetype:match("commit") then return end
    local mark = vim.api.nvim_buf_get_mark(ev.buf, '"')
    if mark[1] > 0 and mark[1] <= vim.api.nvim_buf_line_count(ev.buf) then
      pcall(vim.api.nvim_win_set_cursor, 0, mark)
    end
  end,
})

-- Absolute numbers in insert mode, relative in normal: counts where you need
-- them, the real line number where you are typing.
local numbers = aug("numbers")
vim.api.nvim_create_autocmd({ "InsertEnter", "WinLeave" }, {
  group = numbers,
  callback = function() if vim.wo.number then vim.wo.relativenumber = false end end,
})
vim.api.nvim_create_autocmd({ "InsertLeave", "WinEnter" }, {
  group = numbers,
  callback = function() if vim.wo.number then vim.wo.relativenumber = true end end,
})

-- Make the directory before writing to a path that does not exist yet.
vim.api.nvim_create_autocmd("BufWritePre", {
  group = aug("mkdir"),
  callback = function(ev)
    if ev.match:match("^%w%w+://") then return end
    vim.fn.mkdir(vim.fn.fnamemodify(vim.uv.fs_realpath(ev.match) or ev.match, ":p:h"), "p")
  end,
})

-- Terminal buffers have no use for line numbers or a sign column.
vim.api.nvim_create_autocmd("TermOpen", {
  group = aug("terminal"),
  callback = function()
    vim.opt_local.number = false
    vim.opt_local.relativenumber = false
    vim.opt_local.signcolumn = "no"
    vim.cmd.startinsert()
  end,
})

-- Close throwaway windows with plain `q`.
vim.api.nvim_create_autocmd("FileType", {
  group = aug("close_with_q"),
  pattern = {
    "help", "qf", "man", "lspinfo", "checkhealth", "notify", "startuptime",
    "gitsigns-blame", "dbout", "neotest-output", "neotest-summary", "fugitive",
  },
  callback = function(ev)
    vim.bo[ev.buf].buflisted = false
    vim.keymap.set("n", "q", "<cmd>close<cr>", { buffer = ev.buf, silent = true })
  end,
})

-- Prose gets soft wrap and spell check; code does not.
vim.api.nvim_create_autocmd("FileType", {
  group = aug("prose"),
  pattern = { "markdown", "gitcommit", "text", "tex" },
  callback = function()
    vim.opt_local.wrap = true
    vim.opt_local.linebreak = true
    vim.opt_local.spell = true
    vim.opt_local.conceallevel = 2
  end,
})

-- Re-balance splits when the terminal is resized.
vim.api.nvim_create_autocmd("VimResized", {
  group = aug("resize"),
  callback = function()
    local tab = vim.fn.tabpagenr()
    vim.cmd("tabdo wincmd =")
    vim.cmd("tabnext " .. tab)
  end,
})
