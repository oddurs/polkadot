-- The parts of the screen that are not the buffer.

return {
  -- What the statusline actually says. LazyVim's default is good; this trims
  -- the noise and adds the two things it leaves out — the LSP that is
  -- attached, and how many characters are selected.
  {
    "nvim-lualine/lualine.nvim",
    opts = function(_, opts)
      local icons = LazyVim.config.icons

      opts.sections.lualine_c = {
        LazyVim.lualine.root_dir(),
        { "filetype", icon_only = true, separator = "", padding = { left = 1, right = 0 } },
        { LazyVim.lualine.pretty_path() },
      }

      opts.sections.lualine_x = {
        -- Selection size, only while selecting.
        {
          function()
            local m = vim.fn.mode()
            if m:match("[vV\22]") then
              local s, e = vim.fn.line("v"), vim.fn.line(".")
              local lines = math.abs(s - e) + 1
              return m == "V" and lines .. " lines" or lines .. "L " .. vim.fn.wordcount().visual_chars .. "c"
            end
            return ""
          end,
          color = { fg = "#d26937" },
        },
        {
          function() return require("noice").api.status.command.get() end,
          cond = function() return package.loaded["noice"] and require("noice").api.status.command.has() end,
          color = { fg = "#edb443" },
        },
        -- Which language server is actually attached.
        {
          function()
            local names = {}
            for _, c in pairs(vim.lsp.get_clients({ bufnr = 0 })) do
              if c.name ~= "copilot" then names[#names + 1] = c.name end
            end
            return #names > 0 and " " .. table.concat(names, " ") or ""
          end,
          color = { fg = "#4e5166" },
        },
        { "diff", symbols = { added = icons.git.added, modified = icons.git.modified, removed = icons.git.removed } },
      }

      opts.sections.lualine_y = { { "progress", padding = { left = 1, right = 0 } }, { "location", padding = 0 } }
      opts.sections.lualine_z = { { function() return os.date("%H:%M") end, padding = { left = 1, right = 1 } } }
      return opts
    end,
  },

  -- The dashboard is the first thing you see; make it say something true.
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
      -- Big files: turn off the expensive things rather than refusing to open.
      bigfile = { enabled = true, size = 1.5 * 1024 * 1024 },
      quickfile = { enabled = true },
      statuscolumn = { enabled = true },
      input = { enabled = true },
    },
  },

  -- Fewer, quieter messages; the command line moves to a floating box.
  {
    "folke/noice.nvim",
    opts = {
      presets = { bottom_search = true, command_palette = true, lsp_doc_border = true },
      routes = {
        -- "written" after every save is not news.
        { filter = { event = "msg_show", find = "written" }, opts = { skip = true } },
        { filter = { event = "msg_show", find = "lines" }, opts = { skip = true } },
        { filter = { event = "msg_show", kind = "search_count" }, opts = { skip = true } },
      },
    },
  },

  { "folke/which-key.nvim", opts = { preset = "helix", delay = 250 } },
}
