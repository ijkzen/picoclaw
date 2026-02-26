# Learnings - Theme Toggle Circular Reveal

- Project: PicoClaw web (Angular 21 + Angular Material)
- Theme is managed in `ThemeService` (signals) but **LayoutComponent** currently duplicates theme loading/toggling.
- Prefer centralizing theme state in ThemeService and let components consume it.
- CSS variables for Material theme are available via `mat.theme` generation in `styles.scss` (e.g., `--mat-sys-surface`).
- View Transitions API (Chrome 111+) is available and the preferred approach for a high-performance reveal animation.
- Fallback: if `document.startViewTransition` missing, perform direct toggle without animation.

Notes:
- Always pass click coordinates (clientX, clientY) from toolbar button to theme toggle orchestration.
- For precise coverage compute farthest-corner radius; simpler approach: use clip-path to 150%.

Change summary:
- Added async toggleThemeWithTransition(event?: MouseEvent) to ThemeService. It sets CSS variables --theme-toggle-x/y, uses document.startViewTransition when available to wrap the existing toggleTheme(), and falls back to toggleTheme() when the API is unavailable.
  Implementation preserves existing toggleTheme() behavior and awaits transition.ready when provided.
