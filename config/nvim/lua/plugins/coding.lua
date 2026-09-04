-- Writing and changing code.

return {
  -- Treesitter: the parsers worth having on disk, plus textobjects so `vif`,
  -- `daf`, `cia` work on functions, arguments and classes.
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
      opts.incremental_selection = {
        enable = true,
        keymaps = {
          init_selection = "<C-space>",
          node_incremental = "<C-space>",
          node_decremental = "<BS>",
        },
      }
      return opts
    end,
  },
  {
    "nvim-treesitter/nvim-treesitter-textobjects",
    event = "VeryLazy",
    opts = {
      textobjects = {
        select = {
          enable = true,
          lookahead = true,
          keymaps = {
            ["af"] = "@function.outer", ["if"] = "@function.inner",
            ["ac"] = "@class.outer", ["ic"] = "@class.inner",
            ["aa"] = "@parameter.outer", ["ia"] = "@parameter.inner",
            ["al"] = "@loop.outer", ["il"] = "@loop.inner",
            ["ai"] = "@conditional.outer", ["ii"] = "@conditional.inner",
          },
        },
        move = {
          enable = true,
          set_jumps = true,
          goto_next_start = { ["]f"] = "@function.outer", ["]c"] = "@class.outer", ["]a"] = "@parameter.inner" },
          goto_previous_start = { ["[f"] = "@function.outer", ["[c"] = "@class.outer", ["[a"] = "@parameter.inner" },
        },
        swap = {
          enable = true,
          -- Reorder arguments without retyping them.
          swap_next = { ["<leader>ca"] = "@parameter.inner" },
          swap_previous = { ["<leader>cA"] = "@parameter.inner" },
        },
      },
    },
  },

  -- LSP. LazyVim wires the servers; these are the behaviours worth changing.
  {
    "neovim/nvim-lspconfig",
    opts = {
      inlay_hints = { enabled = true },
      codelens = { enabled = true },
      diagnostics = {
        virtual_text = { spacing = 4, prefix = "●", source = "if_many" },
        severity_sort = true,
        float = { border = "rounded", source = "if_many" },
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
        -- Inlay hints are the reason to use a typed language in an editor.
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

  -- Surround and comment come from LazyVim; this adds the pairs it leaves out.
  {
    "nvim-mini/mini.pairs",
    opts = { modes = { insert = true, command = true, terminal = false } },
  },

  -- Faster than reaching for the mouse when a file has thirty of the same word.
  {
    "folke/flash.nvim",
    opts = { modes = { char = { jump_labels = true }, search = { enabled = false } } },
  },

  -- Persist the session per directory, so reopening a project restores it.
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
}
