library(ggplot2)

math_blue    <- "#2563EB"
math_purple  <- "#7C3AED"
math_emerald <- "#059669"
math_red     <- "#DC2626"
math_amber   <- "#D97706"
math_sky     <- "#0EA5E9"
math_rose    <- "#F43F5E"

math_theme <- function(base_size = 11) {
  theme_minimal(base_size = base_size) +
    theme(
      plot.title = element_text(face = "bold", size = base_size + 2),
      plot.subtitle = element_text(color = "#64748B"),
      axis.title = element_text(face = "bold"),
      panel.grid.minor = element_blank(),
      panel.grid.major = element_line(color = "#E2E8F0", linewidth = 0.4),
      plot.background = element_rect(fill = "white", color = NA),
      panel.background = element_rect(fill = "white", color = NA)
    )
}

save_chart <- function(plot, width = 4.5, height = 3.5, dpi = 300) {
  out <- Sys.getenv("IMAGERENDER_OUTPUT", "output.png")
  ggsave(out, plot = plot, width = width, height = height, dpi = dpi, bg = "white")
}
