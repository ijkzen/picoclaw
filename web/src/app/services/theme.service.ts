import { Injectable, signal } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class ThemeService {
  private readonly THEME_KEY = 'picoclaw-theme';
  isDarkMode = signal<boolean>(false);

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

    // set CSS vars for circular reveal positioning
    const html = document.documentElement;
    html.style.setProperty('--theme-toggle-x', `${x}px`);
    html.style.setProperty('--theme-toggle-y', `${y}px`);

    // Use View Transitions API when available
    const startViewTransition = (document as any).startViewTransition as
      | ((callback: () => void) => { ready?: Promise<void> | undefined })
      | undefined;

    if (typeof startViewTransition === 'function') {
      const transition = startViewTransition(() => {
        // perform the actual theme toggle inside the transition
        this.toggleTheme();
      });

      // wait for transition readiness if provided
      try {
        await (transition?.ready as Promise<void> | undefined);
      } catch {
        // ignore readiness errors and continue
      }
    } else {
      // Fallback when View Transitions API not available
      this.toggleTheme();
    }
  }

  setDarkMode(isDark: boolean): void {
    this.isDarkMode.set(isDark);
    this.applyTheme();
    localStorage.setItem(this.THEME_KEY, isDark ? 'dark' : 'light');
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
  }
}
