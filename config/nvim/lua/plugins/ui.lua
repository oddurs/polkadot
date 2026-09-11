-- Everything that is not the buffer: theme, statusline, floats, file tree.
--
-- One spec per plugin. The previous layout configured lualine, snacks and
-- which-key from two files each, and whichever lazy.nvim merged last won.

local P = require("gotham.palette")

return {
  -- The colorscheme is a file in colors/, not a plugin, so there is nothing
  -- to install and nothing to load — LazyVim just names it.
  { "LazyVim/LazyVim", opts = { colorscheme = "gotham" } },

  -- The statusline. LazyVim's default is good; this trims the noise and adds
  -- the two things it leaves out — which language server is attached, and how
  -- much is selected.
  {
    "nvim-lualine/lualine.nvim",
    opts = function(_, opts)
      local icons = LazyVim.config.icons

      opts.options = vim.tbl_extend("force", opts.options or {}, {
        theme = {
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
        },
        section_separators = { left = "", right = "" },
        component_separators = { left = "", right = "" },
        globalstatus = true,
      })

      opts.sections.lualine_c = {
        LazyVim.lualine.root_dir(),
        { "filetype", icon_only = true, separator = "", padding = { left = 1, right = 0 } },
        { LazyVim.lualine.pretty_path() },
      }

      opts.sections.lualine_x = {
        -- How much is selected, only while selecting.
        {
          function()
            local m = vim.fn.mode()
            if m:match("[vV\22]") then
              local s, e = vim.fn.line("v"), vim.fn.line(".")
              local lines = math.abs(s - e) + 1
              return m == "V" and lines .. " lines"
                or lines .. "L " .. vim.fn.wordcount().visual_chars .. "c"
            end
            return ""
          end,
          color = { fg = P.orange },
        },
        {
          function() return require("noice").api.status.command.get() end,
          cond = function()
            return package.loaded["noice"] and require("noice").api.status.command.has()
          end,
          color = { fg = P.yellow },
        },
        {
          function()
            local names = {}
            for _, c in pairs(vim.lsp.get_clients({ bufnr = 0 })) do
              names[#names + 1] = c.name
            end
            return #names > 0 and " " .. table.concat(names, " ") or ""
          end,
          color = { fg = P.muted },
        },
        {
          "diff",
          symbols = {
            added = icons.git.added,
            modified = icons.git.modified,
            removed = icons.git.removed,
          },
        },
      }

      opts.sections.lualine_y = {
        { "progress", padding = { left = 1, right = 0 } },
        { "location", padding = 0 },
      }
      opts.sections.lualine_z = {
        { function() return os.date("%H:%M") end, padding = { left = 1, right = 1 } },
      }
      return opts
    end,
  },

  -- snacks carries the picker, the dashboard, notifications, indent guides
  -- and lazygit. One spec, so the settings are all visible together.
  {
    "folke/snacks.nvim",
    opts = {
      dashboard = {
        preset = {
          header = table.concat({
            "                      ",
            "   ▄▄▄  ▄▄▄▄▄ ▄▄▄▄▄   ",
            "  ▐▌ ▝▘ ▐▌ ▐▌   █     ",
            "  ▐▌▝▜▌ ▐▌ ▐▌   █     ",
            "   ▀▀▀  ▀▀ ▀▀   ▀     ",
            "                      ",
          }, "\n"),
          keys = {
            { icon = " ", key = "f", desc = "Find file", action = ":lua Snacks.dashboard.pick('files')" },
            { icon = " ", key = "n", desc = "New file", action = ":ene | startinsert" },
            { icon = " ", key = "g", desc = "Grep", action = ":lua Snacks.dashboard.pick('live_grep')" },
            { icon = " ", key = "r", desc = "Recent", action = ":lua Snacks.dashboard.pick('oldfiles')" },
            { icon = " ", key = "p", desc = "Projects", action = ":lua Snacks.picker.projects()" },
            { icon = " ", key = "d", desc = "Dotfiles", action = ":e ~/Code/polkadot" },
            { icon = "󰒲 ", key = "L", desc = "Lazy", action = ":Lazy" },
            { icon = " ", key = "q", desc = "Quit", action = ":qa" },
          },
        },
        sections = {
          { section = "header" },
          { section = "keys", gap = 1, padding = 1 },
          { section = "startup" },
        },
      },
      -- Indent guides, drawn once and never animated. This is the only
      -- source of them now; indent-blankline was doing the same job twice.
      indent = { enabled = true, animate = { enabled = false } },
      scroll = { enabled = false }, -- smooth scroll fights fast movement
      notifier = { enabled = true, timeout = 2500 },
      words = { enabled = true },
      statuscolumn = { enabled = true },
      input = { enabled = true },
      quickfile = { enabled = true },
      -- Big files: turn off the expensive things rather than refusing to open.
      bigfile = { enabled = true, size = 1.5 * 1024 * 1024 },
      lazygit = {
        configure = false, -- our own config.yml is already Gotham
        win = { style = "lazygit" },
      },
    },
    keys = {
      { "<leader>gg", function() Snacks.lazygit() end, desc = "Lazygit" },
      { "<leader>gl", function() Snacks.lazygit.log() end, desc = "Lazygit log (cwd)" },
      { "<leader>gf", function() Snacks.lazygit.log_file() end, desc = "Lazygit file history" },
      { "<leader>gb", function() Snacks.git.blame_line() end, desc = "Git blame line" },
    },
  },

  -- Fewer, quieter messages; the command line moves to a floating box.
  {
    "folke/noice.nvim",
    opts = {
      presets = { bottom_search = true, command_palette = true, lsp_doc_border = true },
      routes = {
        -- "written" after every save is not news, and Ctrl+S makes it often.
        { filter = { event = "msg_show", find = "written" }, opts = { skip = true } },
        { filter = { event = "msg_show", find = "lines" }, opts = { skip = true } },
        { filter = { event = "msg_show", kind = "search_count" }, opts = { skip = true } },
      },
    },
  },

  { "folke/which-key.nvim", opts = { preset = "helix", delay = 300 } },

  -- The file tree: on the left where the mouse expects it, a single click to
  -- open, and out of the way afterwards.
  {
    "nvim-neo-tree/neo-tree.nvim",
    opts = {
      close_if_last_window = true,
      window = {
        width = 32,
        mappings = {
          -- One click opens a file, rather than the default double click.
          ["<LeftRelease>"] = "open",
        },
      },
      filesystem = {
        follow_current_file = { enabled = true },
        use_libuv_file_watcher = true,
        filtered_items = { visible = true, hide_dotfiles = false, hide_gitignored = true },
      },
    },
  },
}
