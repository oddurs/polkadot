-- Gotham.
--
-- vim-gotham predates treesitter, LSP semantic tokens, and every floating
-- window in a modern config, so the plugin alone leaves most of the screen
-- unstyled. This carries the palette forward across all of it.
--
-- The palette is vim-gotham's own — NOT the terminal port. The two differ by a
-- channel or two everywhere, and mixing them leaves added highlights sitting
-- fractionally beside the ones the plugin draws itself.
-- Source: ../../gotham/upstream/vim-gotham-README.md

local P = {
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

-- Diff backgrounds, mixed from the palette rather than picked, so they sit on
-- the same ground as everything else.
local D = {
  add = "#0d2b26",
  add_emph = "#14453c",
  del = "#2a1416",
  del_emph = "#4d1f22",
  chg = "#101f2e",
  chg_emph = "#183349",
}

local function paint()
  local hl = function(g, o) vim.api.nvim_set_hl(0, g, o) end

  local groups = {
    -- ── treesitter ────────────────────────────────────────────────────
    ["@comment"] = { fg = P.muted, italic = true },
    ["@comment.documentation"] = { fg = P.comment, italic = true },
    ["@keyword"] = { fg = P.orange },
    ["@keyword.function"] = { fg = P.orange },
    ["@keyword.return"] = { fg = P.orange, bold = true },
    ["@keyword.import"] = { fg = P.orange },
    ["@keyword.exception"] = { fg = P.red },
    ["@keyword.conditional"] = { fg = P.orange },
    ["@keyword.repeat"] = { fg = P.orange },
    ["@keyword.operator"] = { fg = P.orange },
    ["@function"] = { fg = P.subtle },
    ["@function.call"] = { fg = P.subtle },
    ["@function.builtin"] = { fg = P.cyan },
    ["@function.method"] = { fg = P.subtle },
    ["@function.method.call"] = { fg = P.subtle },
    ["@constructor"] = { fg = P.cyan },
    ["@variable"] = { fg = P.fg },
    ["@variable.builtin"] = { fg = P.orange, italic = true },
    ["@variable.parameter"] = { fg = P.fg },
    ["@variable.member"] = { fg = P.cyan },
    ["@property"] = { fg = P.cyan },
    ["@string"] = { fg = P.yellow },
    ["@string.escape"] = { fg = P.green, bold = true },
    ["@string.special"] = { fg = P.green },
    ["@string.regexp"] = { fg = P.green },
    ["@character"] = { fg = P.yellow },
    ["@number"] = { fg = P.green },
    ["@boolean"] = { fg = P.green, bold = true },
    ["@constant"] = { fg = P.green },
    ["@constant.builtin"] = { fg = P.green, bold = true },
    ["@type"] = { fg = P.cyan },
    ["@type.builtin"] = { fg = P.cyan, italic = true },
    ["@type.definition"] = { fg = P.cyan },
    ["@attribute"] = { fg = P.subtle },
    ["@operator"] = { fg = P.green },
    ["@punctuation.delimiter"] = { fg = P.comment },
    ["@punctuation.bracket"] = { fg = P.comment },
    ["@punctuation.special"] = { fg = P.orange },
    ["@tag"] = { fg = P.orange },
    ["@tag.builtin"] = { fg = P.orange },
    ["@tag.attribute"] = { fg = P.cyan, italic = true },
    ["@tag.delimiter"] = { fg = P.comment },
    ["@module"] = { fg = P.subtle },
    ["@label"] = { fg = P.orange },
    ["@markup.heading"] = { fg = P.subtle, bold = true },
    ["@markup.strong"] = { fg = P.bright, bold = true },
    ["@markup.italic"] = { italic = true },
    ["@markup.link"] = { fg = P.cyan, underline = true },
    ["@markup.raw"] = { fg = P.yellow },
    ["@markup.list"] = { fg = P.orange },
    ["@markup.quote"] = { fg = P.comment, italic = true },
    ["@diff.plus"] = { fg = P.green },
    ["@diff.minus"] = { fg = P.red },

    -- ── LSP semantic tokens ───────────────────────────────────────────
    ["@lsp.type.class"] = { link = "@type" },
    ["@lsp.type.interface"] = { link = "@type" },
    ["@lsp.type.enum"] = { link = "@type" },
    ["@lsp.type.enumMember"] = { link = "@constant" },
    ["@lsp.type.namespace"] = { link = "@module" },
    ["@lsp.type.parameter"] = { link = "@variable.parameter" },
    ["@lsp.type.property"] = { link = "@property" },
    ["@lsp.type.typeParameter"] = { fg = P.cyan, italic = true },
    ["@lsp.typemod.variable.readonly"] = { fg = P.green },
    ["@lsp.typemod.function.defaultLibrary"] = { link = "@function.builtin" },
    ["@lsp.mod.deprecated"] = { strikethrough = true, fg = P.muted },

    -- ── editor chrome ─────────────────────────────────────────────────
    Normal = { fg = P.fg, bg = P.bg },
    NormalNC = { fg = P.fg, bg = P.bg },
    NormalFloat = { fg = P.fg, bg = P.surface },
    FloatBorder = { fg = P.border, bg = P.surface },
    FloatTitle = { fg = P.subtle, bg = P.surface, bold = true },
    WinSeparator = { fg = P.panel, bg = P.bg },
    CursorLine = { bg = P.surface },
    CursorColumn = { bg = P.surface },
    CursorLineNr = { fg = P.yellow, bold = true },
    LineNr = { fg = P.border },
    LineNrAbove = { fg = P.border },
    LineNrBelow = { fg = P.border },
    SignColumn = { bg = "NONE" },
    Visual = { bg = P.selection },
    VisualNOS = { bg = P.selection },
    Search = { fg = P.bg, bg = P.subtle },
    IncSearch = { fg = P.bg, bg = P.yellow },
    CurSearch = { fg = P.bg, bg = P.orange, bold = true },
    Substitute = { fg = P.bg, bg = P.red },
    MatchParen = { fg = P.yellow, bold = true, underline = true },
    Folded = { fg = P.comment, bg = P.surface },
    FoldColumn = { fg = P.border, bg = "NONE" },
    Whitespace = { fg = P.panel },
    NonText = { fg = P.panel },
    SpecialKey = { fg = P.panel },
    Conceal = { fg = P.muted },
    Directory = { fg = P.subtle },
    Title = { fg = P.bright, bold = true },
    ErrorMsg = { fg = P.red },
    WarningMsg = { fg = P.yellow },
    ModeMsg = { fg = P.comment },
    MsgArea = { fg = P.fg },
    Question = { fg = P.green },
    QuickFixLine = { bg = P.selection },
    ColorColumn = { bg = P.surface },
    Cursor = { fg = P.bg, bg = P.fg },
    lCursor = { fg = P.bg, bg = P.fg },
    TermCursor = { fg = P.bg, bg = P.fg },
    Winbar = { fg = P.subtle, bg = "NONE" },
    WinbarNC = { fg = P.muted, bg = "NONE" },

    -- completion menu
    Pmenu = { fg = P.fg, bg = P.surface },
    PmenuSel = { fg = P.bright, bg = P.selection, bold = true },
    PmenuSbar = { bg = P.surface },
    PmenuThumb = { bg = P.border },
    PmenuKind = { fg = P.cyan, bg = P.surface },
    PmenuExtra = { fg = P.muted, bg = P.surface },

    -- statusline / tabs
    StatusLine = { fg = P.subtle, bg = P.surface },
    StatusLineNC = { fg = P.muted, bg = P.surface },
    TabLine = { fg = P.muted, bg = P.surface },
    TabLineFill = { bg = P.bg },
    TabLineSel = { fg = P.bright, bg = P.panel, bold = true },

    -- ── diagnostics ───────────────────────────────────────────────────
    DiagnosticError = { fg = P.red },
    DiagnosticWarn = { fg = P.yellow },
    DiagnosticInfo = { fg = P.subtle },
    DiagnosticHint = { fg = P.green },
    DiagnosticOk = { fg = P.green },
    DiagnosticUnderlineError = { sp = P.red, undercurl = true },
    DiagnosticUnderlineWarn = { sp = P.yellow, undercurl = true },
    DiagnosticUnderlineInfo = { sp = P.subtle, undercurl = true },
    DiagnosticUnderlineHint = { sp = P.green, undercurl = true },
    DiagnosticVirtualTextError = { fg = P.red, bg = D.del },
    DiagnosticVirtualTextWarn = { fg = P.yellow, bg = P.surface },
    DiagnosticVirtualTextInfo = { fg = P.subtle, bg = P.surface },
    DiagnosticVirtualTextHint = { fg = P.green, bg = P.surface },

    -- ── diffs ─────────────────────────────────────────────────────────
    DiffAdd = { bg = D.add },
    DiffChange = { bg = D.chg },
    DiffDelete = { fg = P.red, bg = D.del },
    DiffText = { bg = D.chg_emph, bold = true },
    diffAdded = { fg = P.green },
    diffRemoved = { fg = P.red },
    diffChanged = { fg = P.yellow },
    diffFile = { fg = P.subtle, bold = true },
    diffLine = { fg = P.comment },

    -- ── git ───────────────────────────────────────────────────────────
    GitSignsAdd = { fg = P.green },
    GitSignsChange = { fg = P.yellow },
    GitSignsDelete = { fg = P.red },
    GitSignsAddInline = { bg = D.add_emph },
    GitSignsDeleteInline = { bg = D.del_emph },
    MiniDiffSignAdd = { fg = P.green },
    MiniDiffSignChange = { fg = P.yellow },
    MiniDiffSignDelete = { fg = P.red },

    -- ── plugins ───────────────────────────────────────────────────────
    -- snacks (picker, dashboard, notifier, indent)
    SnacksPickerDir = { fg = P.muted },
    SnacksPickerMatch = { fg = P.orange, bold = true },
    SnacksPickerBorder = { fg = P.border, bg = P.surface },
    SnacksPickerTitle = { fg = P.subtle, bg = P.surface, bold = true },
    SnacksPickerPrompt = { fg = P.green, bg = P.surface },
    SnacksPickerInputBorder = { fg = P.border, bg = P.surface },
    SnacksNormal = { fg = P.fg, bg = P.surface },
    SnacksIndent = { fg = P.panel },
    SnacksIndentScope = { fg = P.border },
    SnacksNotifierInfo = { fg = P.subtle },
    SnacksNotifierWarn = { fg = P.yellow },
    SnacksNotifierError = { fg = P.red },
    SnacksDashboardHeader = { fg = P.subtle },
    SnacksDashboardDesc = { fg = P.fg },
    SnacksDashboardKey = { fg = P.orange },
    SnacksDashboardIcon = { fg = P.cyan },
    SnacksDashboardFooter = { fg = P.muted },

    -- blink.cmp
    BlinkCmpMenu = { fg = P.fg, bg = P.surface },
    BlinkCmpMenuBorder = { fg = P.border, bg = P.surface },
    BlinkCmpMenuSelection = { fg = P.bright, bg = P.selection, bold = true },
    BlinkCmpLabelMatch = { fg = P.orange, bold = true },
    BlinkCmpLabelDeprecated = { fg = P.muted, strikethrough = true },
    BlinkCmpKind = { fg = P.cyan },
    BlinkCmpDoc = { fg = P.fg, bg = P.surface },
    BlinkCmpDocBorder = { fg = P.border, bg = P.surface },
    BlinkCmpSignatureHelp = { fg = P.fg, bg = P.surface },
    BlinkCmpGhostText = { fg = P.muted, italic = true },

    -- which-key
    WhichKey = { fg = P.orange },
    WhichKeyGroup = { fg = P.cyan },
    WhichKeyDesc = { fg = P.fg },
    WhichKeySeparator = { fg = P.muted },
    WhichKeyFloat = { bg = P.surface },
    WhichKeyBorder = { fg = P.border, bg = P.surface },

    -- neo-tree
    NeoTreeNormal = { fg = P.fg, bg = P.bg },
    NeoTreeNormalNC = { fg = P.fg, bg = P.bg },
    NeoTreeWinSeparator = { fg = P.panel, bg = P.bg },
    NeoTreeDirectoryName = { fg = P.subtle },
    NeoTreeDirectoryIcon = { fg = P.subtle },
    NeoTreeRootName = { fg = P.bright, bold = true },
    NeoTreeGitAdded = { fg = P.green },
    NeoTreeGitModified = { fg = P.yellow },
    NeoTreeGitDeleted = { fg = P.red },
    NeoTreeGitUntracked = { fg = P.comment },
    NeoTreeGitIgnored = { fg = P.muted },
    NeoTreeIndentMarker = { fg = P.panel },
    NeoTreeCursorLine = { bg = P.surface },

    -- treesitter-context, illuminate, flash, todo
    TreesitterContext = { bg = P.surface },
    TreesitterContextLineNumber = { fg = P.border, bg = P.surface },
    IlluminatedWordText = { bg = P.selection },
    IlluminatedWordRead = { bg = P.selection },
    IlluminatedWordWrite = { bg = P.selection, underline = true },
    FlashLabel = { fg = P.bg, bg = P.orange, bold = true },
    FlashMatch = { fg = P.bright, bg = P.blue },
    FlashCurrent = { fg = P.bg, bg = P.yellow },

    -- lazy / mason
    LazyNormal = { fg = P.fg, bg = P.surface },
    LazyButton = { bg = P.panel },
    LazyButtonActive = { bg = P.selection, bold = true },
    LazyH1 = { fg = P.bg, bg = P.subtle, bold = true },
    LazyProgressDone = { fg = P.green },
    LazyProgressTodo = { fg = P.border },
    MasonNormal = { fg = P.fg, bg = P.surface },
    MasonHeader = { fg = P.bg, bg = P.subtle, bold = true },
    MasonHighlight = { fg = P.cyan },
    MasonMuted = { fg = P.muted },

    -- noice / notify
    NoiceCmdlinePopup = { fg = P.fg, bg = P.surface },
    NoiceCmdlinePopupBorder = { fg = P.border, bg = P.surface },
    NoiceCmdlineIcon = { fg = P.green },
    NoiceConfirmBorder = { fg = P.border, bg = P.surface },
    NotifyINFOBorder = { fg = P.border },
    NotifyWARNBorder = { fg = P.yellow },
    NotifyERRORBorder = { fg = P.red },

    -- trouble, aerial, harpoon, hipatterns
    TroubleNormal = { fg = P.fg, bg = P.bg },
    TroubleText = { fg = P.fg },
    TroubleCount = { fg = P.orange, bg = P.panel },
    AerialLine = { bg = P.selection },
    HarpoonBorder = { fg = P.border },
    MiniHipatternsFixme = { fg = P.bg, bg = P.red, bold = true },
    MiniHipatternsHack = { fg = P.bg, bg = P.orange, bold = true },
    MiniHipatternsTodo = { fg = P.bg, bg = P.cyan, bold = true },
    MiniHipatternsNote = { fg = P.bg, bg = P.green, bold = true },
  }

  for g, o in pairs(groups) do
    hl(g, o)
  end

  -- The terminal inside nvim gets the terminal palette, since that is what it
  -- actually is.
  local term = {
    "#0a0f14", "#c33027", "#26a98b", "#edb54b", "#195465", "#4e5165", "#33859d", "#98d1ce",
    "#10151b", "#d26939", "#081f2d", "#245361", "#093748", "#888ba5", "#599caa", "#d3ebe9",
  }
  for i, c in ipairs(term) do
    vim.g["terminal_color_" .. (i - 1)] = c
  end
end

return {
  {
    "whatyouhide/vim-gotham",
    lazy = false,
    priority = 1000,
    config = function()
      vim.cmd.colorscheme("gotham256")
      paint()
      -- Re-apply after any later colorscheme load, so :colorscheme gotham256
      -- from the picker keeps the extensions.
      vim.api.nvim_create_autocmd("ColorScheme", {
        pattern = "gotham*",
        callback = paint,
      })
    end,
  },

  { "LazyVim/LazyVim", opts = { colorscheme = "gotham256" } },

  -- The statusline carries the palette too, or the bottom of the screen is the
  -- one thing still wearing the previous theme.
  {
    "nvim-lualine/lualine.nvim",
    opts = function(_, opts)
      local t = {
        normal = {
          a = { fg = P.bg, bg = P.subtle, gui = "bold" },
          b = { fg = P.subtle, bg = P.panel },
          c = { fg = P.comment, bg = "NONE" },
        },
        insert = { a = { fg = P.bg, bg = P.green, gui = "bold" } },
        visual = { a = { fg = P.bg, bg = P.orange, gui = "bold" } },
        replace = { a = { fg = P.bg, bg = P.red, gui = "bold" } },
        command = { a = { fg = P.bg, bg = P.yellow, gui = "bold" } },
        inactive = {
          a = { fg = P.muted, bg = "NONE" },
          b = { fg = P.muted, bg = "NONE" },
          c = { fg = P.muted, bg = "NONE" },
        },
      }
      opts.options = opts.options or {}
      opts.options.theme = t
      opts.options.section_separators = { left = "", right = "" }
      opts.options.component_separators = { left = "", right = "" }
      opts.options.globalstatus = true
      return opts
    end,
  },
}
