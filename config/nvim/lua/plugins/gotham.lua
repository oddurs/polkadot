-- Gotham.
--
-- Andrea Leopardi's palette, wired into LazyVim as the default colorscheme.
-- The upstream plugin predates treesitter, so the highlight groups below carry
-- it forward to modern nvim — without them the syntax is flat, since old vim
-- groups no longer drive most of the buffer.
--
-- To go back to catppuccin, change `colorscheme` in lua/plugins/colorscheme.lua.

-- vim-gotham's own palette, not the terminal port. The two differ by a
-- channel or two everywhere; using the terminal values here would leave these
-- overrides a shade off the highlights the plugin draws itself.
-- Source: config/gotham/upstream/vim-gotham-README.md
local palette = {
  bg = "#0c1014", -- base0
  surface = "#11151c", -- base1
  panel = "#091f2e", -- base2
  selection = "#0a3749", -- base3
  border = "#245361", -- base4
  subtle = "#599cab", -- base5
  fg = "#99d1ce", -- base6
  bright = "#d3ebe9", -- base7
  red = "#c23127",
  orange = "#d26937",
  yellow = "#edb443",
  green = "#2aa889",
  blue = "#195466",
  cyan = "#33859e",
  comment = "#888ca6", -- magenta
  muted = "#4e5166", -- violet
}

return {
  {
    "whatyouhide/vim-gotham",
    lazy = false,
    priority = 1000,
    config = function()
      vim.cmd.colorscheme("gotham256")

      local p = palette
      local set = vim.api.nvim_set_hl

      -- Treesitter. The upstream scheme only fills the legacy groups.
      local ts = {
        ["@comment"] = { fg = p.muted, italic = true },
        ["@keyword"] = { fg = p.orange },
        ["@keyword.function"] = { fg = p.orange },
        ["@keyword.return"] = { fg = p.orange },
        ["@conditional"] = { fg = p.orange },
        ["@repeat"] = { fg = p.orange },
        ["@function"] = { fg = p.subtle },
        ["@function.call"] = { fg = p.subtle },
        ["@function.builtin"] = { fg = p.cyan },
        ["@method"] = { fg = p.subtle },
        ["@method.call"] = { fg = p.subtle },
        ["@constructor"] = { fg = p.cyan },
        ["@variable"] = { fg = p.fg },
        ["@variable.builtin"] = { fg = p.orange, italic = true },
        ["@parameter"] = { fg = p.fg },
        ["@field"] = { fg = p.cyan },
        ["@property"] = { fg = p.cyan },
        ["@string"] = { fg = p.yellow },
        ["@string.escape"] = { fg = p.green },
        ["@number"] = { fg = p.green },
        ["@boolean"] = { fg = p.green },
        ["@constant"] = { fg = p.green },
        ["@constant.builtin"] = { fg = p.green },
        ["@type"] = { fg = p.cyan },
        ["@type.builtin"] = { fg = p.cyan },
        ["@operator"] = { fg = p.green },
        ["@punctuation.delimiter"] = { fg = p.comment },
        ["@punctuation.bracket"] = { fg = p.comment },
        ["@tag"] = { fg = p.orange },
        ["@tag.attribute"] = { fg = p.cyan },
        ["@tag.delimiter"] = { fg = p.comment },
      }
      for group, opts in pairs(ts) do
        set(0, group, opts)
      end

      -- Editor chrome: quiet borders, a readable float, and a cursorline that
      -- is felt rather than seen.
      local ui = {
        Normal = { fg = p.fg, bg = p.bg },
        NormalFloat = { fg = p.fg, bg = p.surface },
        FloatBorder = { fg = p.border, bg = p.surface },
        FloatTitle = { fg = p.subtle, bg = p.surface, bold = true },
        WinSeparator = { fg = p.panel },
        CursorLine = { bg = p.surface },
        CursorLineNr = { fg = p.yellow, bold = true },
        LineNr = { fg = p.border },
        SignColumn = { bg = "NONE" },
        Visual = { bg = p.selection },
        Search = { fg = p.bg, bg = p.subtle },
        IncSearch = { fg = p.bg, bg = p.yellow },
        CurSearch = { fg = p.bg, bg = p.orange },
        MatchParen = { fg = p.yellow, bold = true },
        Pmenu = { fg = p.fg, bg = p.surface },
        PmenuSel = { fg = p.bright, bg = p.selection, bold = true },
        PmenuSbar = { bg = p.surface },
        PmenuThumb = { bg = p.border },
        StatusLine = { fg = p.subtle, bg = p.surface },
        WinBar = { fg = p.subtle, bg = "NONE" },
        WinBarNC = { fg = p.muted, bg = "NONE" },
        TabLineSel = { fg = p.bright, bg = p.panel },
        Folded = { fg = p.comment, bg = p.surface },
        Whitespace = { fg = p.panel },
        IndentBlanklineChar = { fg = p.panel },

        DiagnosticError = { fg = p.red },
        DiagnosticWarn = { fg = p.yellow },
        DiagnosticInfo = { fg = p.subtle },
        DiagnosticHint = { fg = p.green },
        DiagnosticUnderlineError = { sp = p.red, undercurl = true },
        DiagnosticUnderlineWarn = { sp = p.yellow, undercurl = true },

        DiffAdd = { fg = p.green, bg = "#0d2b26" },
        DiffChange = { fg = p.yellow, bg = "#1d2321" },
        DiffDelete = { fg = p.red, bg = "#2a1416" },
        DiffText = { fg = p.bright, bg = "#14453c" },
        GitSignsAdd = { fg = p.green },
        GitSignsChange = { fg = p.yellow },
        GitSignsDelete = { fg = p.red },
      }
      for group, opts in pairs(ui) do
        set(0, group, opts)
      end
    end,
  },

  -- LazyVim reads its colorscheme from here.
  {
    "LazyVim/LazyVim",
    opts = { colorscheme = "gotham256" },
  },
}
