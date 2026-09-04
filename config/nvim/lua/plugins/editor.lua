-- Editor behaviour on top of LazyVim's defaults.
--
-- The theme is one thing; the reason to open nvim is another. These are the
-- pieces that make the day-to-day feel finished: git in a keystroke, motion
-- that doesn't need counting, and a file tree that gets out of the way.

return {
  -- lazygit, one keystroke away and sharing the terminal's colours.
  {
    "folke/snacks.nvim",
    opts = {
      lazygit = {
        configure = false, -- our own config.yml is already Gotham
        win = { style = "lazygit", border = "rounded" },
      },
      indent = { enabled = true, animate = { enabled = false } },
      scroll = { enabled = false }, -- smooth scroll fights fast movement
      notifier = { enabled = true, timeout = 2500 },
      words = { enabled = true },
    },
    keys = {
      { "<leader>gg", function() require("snacks").lazygit() end, desc = "Lazygit" },
      { "<leader>gl", function() require("snacks").lazygit.log() end, desc = "Lazygit log (cwd)" },
      { "<leader>gf", function() require("snacks").lazygit.log_file() end, desc = "Lazygit file history" },
      { "<leader>gb", function() require("snacks").git.blame_line() end, desc = "Git blame line" },
    },
  },

  -- Jump anywhere with two characters.
  {
    "folke/flash.nvim",
    opts = { modes = { char = { jump_labels = true } } },
  },

  -- The file tree, positioned right and closing after you pick something.
  {
    "nvim-neo-tree/neo-tree.nvim",
    opts = {
      close_if_last_window = true,
      window = { width = 32 },
      filesystem = {
        follow_current_file = { enabled = true },
        use_libuv_file_watcher = true,
        filtered_items = { visible = true, hide_dotfiles = false, hide_gitignored = true },
      },
    },
  },

  -- Show the key hints sooner; the default delay is long enough to forget.
  { "folke/which-key.nvim", opts = { delay = 300 } },
}
