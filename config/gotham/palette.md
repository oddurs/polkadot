# Gotham

Andrea Leopardi's palette. Upstream files are vendored verbatim in
[`upstream/`](upstream/) so provenance is checkable, not remembered.

There are **two** official Gotham palettes, and they are not the same. Using
one where the other belongs is the subtle way to get this wrong.

## Terminal — [gotham-contrib](https://github.com/whatyouhide/gotham-contrib)

The sixteen ANSI slots, as shipped for alacritty, kitty, iTerm, Terminal.app and
the rest. Verified slot-for-slot against `upstream/alacritty-gotham.yml`.

| slot | hex | role |
|------|-----|------|
| bg   | `#0a0f14` | ground |
| fg   | `#98d1ce` | body text |
| 0    | `#0a0f14` | black |
| 1    | `#c33027` | red |
| 2    | `#26a98b` | green |
| 3    | `#edb54b` | yellow |
| 4    | `#195465` | blue |
| 5    | `#4e5165` | magenta |
| 6    | `#33859d` | cyan |
| 7    | `#98d1ce` | white |
| 8    | `#10151b` | bright black — surface |
| 9    | `#d26939` | bright red — orange |
| 10   | `#081f2d` | bright green — panel |
| 11   | `#245361` | bright yellow — border |
| 12   | `#093748` | bright blue — selection |
| 13   | `#888ba5` | bright magenta — comment |
| 14   | `#599caa` | bright cyan — secondary text |
| 15   | `#d3ebe9` | bright white — emphasis |

Slots 10–12 are the reason Gotham stays calm: where most schemes spend the
bright range on louder text, Gotham spends it on darker blues used as UI
shades. A port that "fixes" this by brightening them is no longer Gotham.

Used by: Ghostty, Neovim, fish, fzf, starship, lazygit, delta, bat, btop, herdr.

## Editor — [vim-gotham](https://github.com/whatyouhide/vim-gotham)

The colorscheme's own palette, from `upstream/vim-gotham-README.md`. Close to
the terminal set but not identical — every channel differs by one or two.

| base | hex | | named | hex |
|------|-----|-|-------|-----|
| 0 | `#0c1014` | | red | `#c23127` |
| 1 | `#11151c` | | orange | `#d26937` |
| 2 | `#091f2e` | | yellow | `#edb443` |
| 3 | `#0a3749` | | magenta | `#888ca6` |
| 4 | `#245361` | | violet | `#4e5166` |
| 5 | `#599cab` | | blue | `#195466` |
| 6 | `#99d1ce` | | cyan | `#33859e` |
| 7 | `#d3ebe9` | | green | `#2aa889` |

Not used any more. Neovim used to extend the `gotham256` colorscheme with a
file of overrides, and those overrides had to use *these* values or sit a
shade off the highlights the plugin drew itself.

Neovim now ships its own complete colorscheme, `config/nvim/colors/gotham.lua`,
with no plugin underneath it and nothing to sit flush with. It uses the
**terminal** set above, which means `Normal` is byte-identical to the ground
behind Ghostty's window padding and there is no seam around the editor. The
values live in `config/nvim/lua/gotham/palette.lua`, read by both the
colorscheme and the statusline.

## Everything else gets it for free

`bat`, `delta` and `herdr` are configured to render through the terminal's own
sixteen colours (`--theme=base16`, `syntax-theme = base16`, `theme.name =
"terminal"`). They are Gotham because Ghostty is, and follow automatically if
the terminal theme changes.
