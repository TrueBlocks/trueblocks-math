# Image Style Guide — Hidden Mathematics Series

All generated images (Mermaid diagrams, R plots, and future D3 visualizations) must follow these visual standards for a consistent, modern, colorful look across the book series.

## Color Palette

| Name       | Hex       | Usage                                      |
|------------|-----------|---------------------------------------------|
| Blue       | `#2563EB` | Primary data, axes, main concepts           |
| Purple     | `#7C3AED` | Secondary data, annotations, relationships  |
| Emerald    | `#059669` | Positive values, highlights, success states  |
| Red        | `#DC2626` | Warnings, negative values, key callouts      |
| Amber      | `#D97706` | Tertiary data, accents                       |
| Sky        | `#0EA5E9` | Lighter accent, background highlights        |
| Rose       | `#F43F5E` | Warm accent, emphasis                        |
| Dark Slate | `#1E293B` | Text, titles, labels                         |
| Slate      | `#64748B` | Gridlines, edges, secondary labels           |
| Light Gray | `#F8FAFC` | Plot area backgrounds                        |
| Border     | `#E2E8F0` | 1px image border                             |

## Typography

- **Title**: 13pt, Dark Slate (`#1E293B`), bold
- **Subtitle / Caption**: 10pt, Slate (`#64748B`)
- **Axis Labels**: 10pt, Dark Slate
- **Annotations**: 9pt, Slate
- **Font Family**: Helvetica Neue (fallback: sans-serif)

## Image Dimensions

- **Full-width images**: 6.5 inches wide (fits Word margins)
- **Aspect ratio**: Prefer 16:10 for charts, flexible for diagrams
- **Resolution**: 300 DPI for print-quality output
- **Border**: 1px `#E2E8F0` border on all images

## R Plots (ggplot2)

All R scripts must:
1. Source `common.R` at the top: `source(Sys.getenv("COMMON_R"))`
2. Use named colors from palette: `math_blue`, `math_purple`, `math_emerald`, `math_red`, `math_amber`, `math_sky`, `math_rose`
3. Use `math_theme()` for consistent styling
4. Call `save_chart(p)` to export with standard dimensions and border
5. Prefer minimal gridlines — major gridlines only in light gray
6. Never use default ggplot2 colors

## Mermaid Diagrams

All Mermaid diagrams must use the custom theme that matches the palette:

- Primary fill: `#2563EB` (blue) with `#1D4ED8` stroke, white text
- Secondary fill: `#7C3AED` (purple) with `#6D28D9` stroke, white text
- Accent fill: `#059669` (emerald) with `#047857` stroke, white text
- Warning fill: `#DC2626` (red) with `#B91C1C` stroke, white text
- Neutral fill: `#F8FAFC` with `#E2E8F0` stroke, dark text
- Edge/line color: `#64748B` (slate)
- Edge label color: `#1E293B`
- Background: white

## Captioning

- Every image must have a caption below it
- Captions are descriptive: explain what the image shows, not just label it
- Format: italic, 10pt, centered
- Example: *"The derivative chain extends through six orders — from position to pop — each describing a subtler aspect of motion."*

## Quality Checklist

Before including any image:
- [ ] Colors match the palette — no default theme colors
- [ ] Text is legible at print size (no text smaller than 8pt)
- [ ] Image has a clear explanatory purpose
- [ ] Border is present (1px, `#E2E8F0`)
- [ ] Caption is descriptive and meaningful
- [ ] Image renders without errors
