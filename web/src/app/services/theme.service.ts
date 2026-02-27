import { Injectable, signal } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class ThemeService {
  private readonly THEME_KEY = 'picoclaw-theme';
  isDarkMode = signal<boolean>(false);

  // Highlight.js theme management
  private readonly HIGHLIGHT_LINK_ID = 'highlightjs-theme';
  // Use highlight theme bundles emitted by Angular build (non-injected)
  private readonly HIGHLIGHT_LIGHT = '/highlight-light.css';
  private readonly HIGHLIGHT_DARK = '/highlight-dark.css';

  constructor() {
    this.loadTheme();
  }

  private loadTheme(): void {
    const savedTheme = localStorage.getItem(this.THEME_KEY);
    if (savedTheme) {
      this.isDarkMode.set(savedTheme === 'dark');
    } else {
      // Check system preference
      this.isDarkMode.set(window.matchMedia('(prefers-color-scheme: dark)').matches);
    }
    this.applyTheme();
    this.updateHighlightJsTheme();
  }

  toggleTheme(): void {
    this.isDarkMode.update(current => !current);
    this.applyTheme();
    localStorage.setItem(this.THEME_KEY, this.isDarkMode() ? 'dark' : 'light');
  }

  async toggleThemeWithTransition(event?: MouseEvent): Promise<void> {
    // compute center fallback
    const x = event?.clientX ?? Math.round(window.innerWidth / 2);
    const y = event?.clientY ?? Math.round(window.innerHeight / 2);

    const html = document.documentElement;
    html.style.setProperty('--theme-toggle-x', `${x}px`);
    html.style.setProperty('--theme-toggle-y', `${y}px`);

    if (typeof (document as any).startViewTransition === 'function') {
      const transition = (document as any).startViewTransition(() => {
        this.toggleTheme();
      });

      // 手动触发动画，只作用于新视图
      const endRadius = Math.sqrt(
        Math.max(x, window.innerWidth - x) ** 2 + Math.max(y, window.innerHeight - y) ** 2
      );
      transition.ready
        .then(() => {
          document.documentElement.animate(
            {
              clipPath: [
                `circle(0px at ${x}px ${y}px)`,
                `circle(${endRadius}px at ${x}px ${y}px)`
              ]
            },
            {
              duration: 400,
              easing: 'ease-in-out',
              pseudoElement: '::view-transition-new(root)'
            }
          );
        })
        .catch(() => undefined);
    } else {
      this.toggleTheme();
    }
  }

  setDarkMode(isDark: boolean): void {
    this.isDarkMode.set(isDark);
    this.applyTheme();
    localStorage.setItem(this.THEME_KEY, isDark ? 'dark' : 'light');
    this.updateHighlightJsTheme();
  }

  private applyTheme(): void {
    const html = document.documentElement;
    if (this.isDarkMode()) {
      html.classList.add('dark');
      document.body.classList.add('mat-app-background');
    } else {
      html.classList.remove('dark');
      document.body.classList.remove('mat-app-background');
    }
    // keep highlight.js theme in sync with overall theme
    this.updateHighlightJsTheme();
  }

  private updateHighlightJsTheme(): void {
    const href = this.isDarkMode() ? this.HIGHLIGHT_DARK : this.HIGHLIGHT_LIGHT;
    const id = this.HIGHLIGHT_LINK_ID;
    let link = document.getElementById(id) as HTMLLinkElement | null;

    if (link) {
      if (link.getAttribute('href') !== href) {
        link.setAttribute('href', href);
      }
      return;
    }

    link = document.createElement('link');
    link.id = id;
    link.rel = 'stylesheet';
    link.href = href;
    document.head.appendChild(link);
  }
}
