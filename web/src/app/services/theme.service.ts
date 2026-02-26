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
