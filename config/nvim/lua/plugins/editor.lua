-- The buffer itself: parsing, language servers, formatting, git.

return {
  -- Treesitter for highlighting and indentation. The textobject queries are
  -- deliberately absent: `vif` and `daf` are motions, and this config is not
  -- driven by motions.
  {
    "nvim-treesitter/nvim-treesitter",
    opts = function(_, opts)
      vim.list_extend(opts.ensure_installed or {}, {
        "bash", "css", "diff", "dockerfile", "fish", "gitcommit", "gitignore",
        "go", "gomod", "gosum", "graphql", "html", "javascript", "jsdoc", "json",
        "jsonc", "lua", "luadoc", "make", "markdown", "markdown_inline", "python",
        "regex", "rust", "scss", "sql", "toml", "tsx", "typescript", "vim",
        "vimdoc", "yaml", "zig",
      })
      return opts
    end,
  },

  -- LSP. LazyVim wires the servers; these are the behaviours worth changing.
  {
    "neovim/nvim-lspconfig",
    opts = {
      -- Inlay hints are off by default. They are genuinely useful and they
      -- also reflow the line you are reading, which is the opposite of what
      -- a calm screen wants. `<leader>uh` turns them on per buffer.
      inlay_hints = { enabled = false },
      codelens = { enabled = false },
      diagnostics = {
        virtual_text = { spacing = 4, prefix = "●", source = "if_many" },
        severity_sort = true,
        float = { source = "if_many" },
      },
      servers = {
        lua_ls = {
          settings = {
            Lua = {
              hint = { enable = true, arrayIndex = "Disable" },
              workspace = { checkThirdParty = false },
            },
          },
        },
        vtsls = {
          settings = {
            typescript = {
              inlayHints = {
                parameterNames = { enabled = "literals" },
                variableTypes = { enabled = false },
                propertyDeclarationTypes = { enabled = true },
                functionLikeReturnTypes = { enabled = true },
              },
              preferences = { importModuleSpecifier = "non-relative" },
            },
          },
        },
        gopls = {
          settings = {
            gopls = {
              hints = {
                assignVariableTypes = true,
                compositeLiteralFields = true,
                parameterNames = true,
                rangeVariableTypes = true,
              },
              analyses = { unusedparams = true, shadow = true },
              staticcheck = true,
              gofumpt = true,
            },
          },
        },
      },
    },
  },

  -- Formatting. Two rules: shell scripts keep two-space indents, and a file
  -- with no configured formatter falls back to the LSP rather than nothing.
  {
    "stevearc/conform.nvim",
    opts = {
      formatters_by_ft = {
        fish = { "fish_indent" },
        sh = { "shfmt" },
        bash = { "shfmt" },
        go = { "goimports", "gofumpt" },
        lua = { "stylua" },
        toml = { "taplo" },
        ["_"] = { "trim_whitespace" },
      },
      formatters = { shfmt = { prepend_args = { "-i", "2", "-ci" } } },
      default_format_opts = { lsp_format = "fallback" },
    },
  },

  -- Brackets and quotes close themselves.
  {
    "nvim-mini/mini.pairs",
    opts = { modes = { insert = true, command = true, terminal = false } },
  },

  -- Reopen a project where you left it.
  {
    "folke/persistence.nvim",
    event = "BufReadPre",
    opts = {},
    keys = {
      { "<leader>qs", function() require("persistence").load() end, desc = "Restore session" },
      { "<leader>ql", function() require("persistence").load({ last = true }) end, desc = "Restore last session" },
      { "<leader>qd", function() require("persistence").stop() end, desc = "Don't save current session" },
    },
  },

  -- ── git ──────────────────────────────────────────────────────────────────
  {
    "lewis6991/gitsigns.nvim",
    opts = {
      signs = {
        add = { text = "▎" },
        change = { text = "▎" },
        delete = { text = "" },
        topdelete = { text = "" },
        changedelete = { text = "▎" },
        untracked = { text = "▎" },
      },
      current_line_blame = false, -- on demand, via <leader>gt
      current_line_blame_opts = { delay = 300, virt_text_pos = "eol" },
    },
    keys = {
      { "<leader>gt", "<cmd>Gitsigns toggle_current_line_blame<cr>", desc = "Toggle line blame" },
    },
  },

  -- Side-by-side diffs and file history, in the editor rather than a pager.
  {
    "sindrets/diffview.nvim",
    cmd = { "DiffviewOpen", "DiffviewFileHistory", "DiffviewClose" },
    opts = {
      enhanced_diff_hl = true,
      view = {
        default = { winbar_info = true },
        merge_tool = { layout = "diff3_mixed", disable_diagnostics = true },
      },
    },
    keys = {
      { "<leader>gv", "<cmd>DiffviewOpen<cr>", desc = "Diffview: working tree" },
      { "<leader>gV", "<cmd>DiffviewClose<cr>", desc = "Diffview: close" },
      { "<leader>gh", "<cmd>DiffviewFileHistory %<cr>", desc = "Diffview: this file's history" },
      { "<leader>gH", "<cmd>DiffviewFileHistory<cr>", desc = "Diffview: branch history" },
    },
  },

  -- Conflict resolution that does not involve counting angle brackets.
  {
    "akinsho/git-conflict.nvim",
    version = "*",
    event = "BufReadPre",
    opts = { default_mappings = true, disable_diagnostics = true },
  },
}
