-- The Gotham palette, in one place.
--
-- These are the *terminal* Gotham values, the same sixteen Ghostty is
-- configured with, so Neovim's `Normal` background is byte-identical to the
-- ground behind the window padding and there is no seam around the editor.
--
-- vim-gotham's own set differs from this by a channel or two in every slot.
-- It was the right choice while the config extended the gotham256 colorscheme
-- and had to sit flush with highlights that plugin drew itself. Nothing does
-- that any more: colors/gotham.lua writes every group out directly.
--
-- Provenance: config/gotham/palette.md

return {
  bg = "#0a0f14", -- ground, == ghostty background
  surface = "#10151b", -- floats, statusline, cursorline
  panel = "#081f2d", -- the darker blue shade
  selection = "#093748", -- visual
  border = "#245361", -- borders, line numbers, indent guides
  subtle = "#599caa", -- secondary text, functions
  fg = "#98d1ce", -- body text
  bright = "#d3ebe9", -- emphasis
  red = "#c33027",
  orange = "#d26939",
  yellow = "#edb54b",
  green = "#26a98b",
  blue = "#195465",
  cyan = "#33859d",
  comment = "#888ba5",
  muted = "#4e5165",
}
