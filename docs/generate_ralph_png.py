#!/usr/bin/env python3
"""Generate ralph.png from the ANSI pixel art sprite."""

from PIL import Image

# 256-color ANSI palette entries used in the sprite
# Mapped from the escape codes in test.dart / main.go
PALETTE = {
    0:   (0,   0,   0),     # black
    15:  (255, 255, 255),   # white
    37:  (0,   135, 135),   # teal
    38:  (0,   135, 175),   # dark cyan
    44:  (0,   215, 215),   # cyan
    58:  (95,  95,  0),     # dark olive
    94:  (135, 95,  0),     # brown (hair outline)
    178: (215, 175, 0),     # hair yellow (not used in grid but close)
    184: (215, 215, 0),     # yellow-green
    220: (255, 215, 0),     # yellow (skin)
    230: (255, 255, 215),   # light yellow
    232: (8,   8,   8),     # near-black
    233: (18,  18,  18),    # very dark grey
}

BG = None  # transparent

# The sprite as a pixel grid (top-half color, bottom-half color per cell)
# Each row pair from the half-block art = two pixel rows
# Decoded by hand from the escape sequences in test.dart

# Row format: each entry is (fg_color, bg_color) from the half-block char
# ▄ = bottom half block: fg = bottom pixel, bg = top pixel
# ▀ = top half block: fg = top pixel, bg = bottom pixel
# ' ' (space) = both pixels are bg color

# Sprite is 16 chars wide, 8 char rows = 16 pixel rows

# I'll define it as a 16x16 pixel grid directly
# BG = transparent

T = None  # transparent

grid = [
    # Row 0-1 (char row 0): "   ▄▄▄▄▄▄▄▄▄▄   "
    # top row
    [T,  T,  T,  94,  94,  94, 220, 94, 220, 94,  94,  94,  T,  T,  T, T],
    # bottom row
    [T,  T,  T,  94,  94,  94,  94,  94,  94,  94,  94,  94,  T,  T,  T, T],

    # Row 2-3 (char row 1): " ▄▄▄▄ ▄▄▄▄ ▄▄▄ "
    # top
    [T,  94, 94, 220, 94, 220,  94, 220, 94,  220, 220,  58, 94, 94,  T, T],
    # bottom
    [T, 220,220, 220,220, 220, 230,230, 220, 220, 220, 230,220, 94,  T, T],

    # Row 4-5 (char row 2): "  ▄ ▄  ▄▄▄ ▄▄▄▄ "
    # top
    [T,  94, 220, 94, 220, 220,  15,  0,  15, 220, 184,  15, 232,  15,  94, T],
    # bottom
    [T,  94, 220,220, 220, 220, 220, 15, 220, 220,  94, 220,  15, 220,  94, T],

    # Row 6-7 (char row 3): "▀▀     ▄▄▄▄▄  ▄▄"
    # top
    [94,  94, 220, 220, 220, 220, 220,  94,  94,   0,  94, 220, 220, 220,  94,  94],
    # bottom
    [ T,   T, 220, 220, 220, 220, 220, 220, 220, 220, 220,  94, 220, 220, 220,  T],

    # Row 8-9 (char row 4): "  ▄▄ ▄▄ ▄▄   ▀ "
    # top
    [T,  T,  94, 220, 220, 232,   0,   0, 220, 233, 232,  94,  94,  94,  94,  T],
    # bottom
    [T,  T,  37, 220, 220, 220, 220, 220, 220, 220, 220,  94,  94,  94,   T,  T],

    # Row 10-11 (char row 5): "▄▄▄▄▄▄▄▄▄ ▄▄▄▄  "
    # top
    [0,   0, 232,  38,  37,   0,  38,   0, 220, 220,  94,   0,  44, 232,  T,  T],
    # bottom
    [0,  44,  38,  38,  37,  233, 44,  38, 233, 220, 232, 220,   0,  44,  T,  T],

    # Row 12-13 (char row 6): "▄▄▄▄▄▄▄▄▄▄▄▄▄▄  "
    # top
    [0,  44,  44,   0,   0,  38,  44,  44,  37, 232,  44,  44,  38,   0,   0,  T],
    # bottom
    [0,   0,  44,  44,  44,  44,  44,  37,   0,   0,  44,  38,  44,  44,   0,  T],

    # Row 14-15 (char row 7): " ▀▄▄▄▄▄▄▄▄▄▄▄▄ "
    # top
    [T,   0,  38,  38,  38,  37,   0,   0, 232,  44,  44,  44,  44,  44,  38,   0],
    # bottom
    [T,   T,   0,   0, 232, 232,   0,  37,  38,  44,  44,  44,  44,  44,  44,   T],
]

SCALE = 16  # pixels per cell

width = 16
height = 16
img = Image.new('RGBA', (width * SCALE, height * SCALE), (0, 0, 0, 0))

for y, row in enumerate(grid):
    for x, color_idx in enumerate(row):
        if color_idx is None:
            continue  # transparent
        r, g, b = PALETTE[color_idx]
        for py in range(SCALE):
            for px in range(SCALE):
                img.putpixel((x * SCALE + px, y * SCALE + py), (r, g, b, 255))

img.save('docs/ralph.png')
print(f"Saved docs/ralph.png ({img.width}x{img.height})")
