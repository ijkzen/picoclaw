import { Component, signal, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterOutlet, RouterLink, RouterLinkActive } from '@angular/router';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatListModule } from '@angular/material/list';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatCardModule } from '@angular/material/card';
import { MatDividerModule } from '@angular/material/divider';
import { MatTooltipModule } from '@angular/material/tooltip';
import { BreakpointObserver, Breakpoints } from '@angular/cdk/layout';
import { ThemeService } from '../../services/theme.service';

@Component({
  selector: 'app-layout',
  standalone: true,
  imports: [
    CommonModule,
    RouterOutlet,
    RouterLink,
    RouterLinkActive,
    MatSidenavModule,
    MatListModule,
    MatIconModule,
    MatButtonModule,
    MatToolbarModule,
    MatCardModule,
    MatDividerModule,
    MatTooltipModule
  ],
  template: `
    <mat-sidenav-container class="sidenav-container">
      <mat-sidenav
        #drawer
        class="sidenav"
        fixedInViewport
        [attr.role]="isHandset() ? 'dialog' : 'navigation'"
        [mode]="isHandset() ? 'over' : 'side'"
        [opened]="!isHandset()">
        
        <!-- Sidebar Header -->
        <mat-toolbar class="sidenav-toolbar">
          <div class="logo-container">
            <div class="logo">P</div>
            <span class="logo-text">PicoClaw</span>
          </div>
        </mat-toolbar>
        
        <mat-divider></mat-divider>

        <!-- Navigation -->
        <mat-nav-list class="nav-list">
          <a mat-list-item
             routerLink="/chat"
             routerLinkActive="active-link"
             [routerLinkActiveOptions]="{exact: true}"
             (click)="isHandset() && drawer.close()">
            <mat-icon matListItemIcon> chat </mat-icon>
            <span matListItemTitle> Chat </span>
          </a>
          
          <a mat-list-item
             routerLink="/settings"
             routerLinkActive="active-link"
             (click)="isHandset() && drawer.close()">
            <mat-icon matListItemIcon> settings </mat-icon>
            <span matListItemTitle> Settings </span>
          </a>
        </mat-nav-list>


      </mat-sidenav>

      <mat-sidenav-content class="sidenav-content">
        <!-- Main Toolbar -->
        <mat-toolbar color="primary" class="main-toolbar">
          @if (isHandset()) {
            <button
              type="button"
              mat-icon-button
              (click)="drawer.toggle()"
              aria-label="Toggle sidenav">
              <mat-icon aria-label="Side nav toggle icon">menu</mat-icon>
            </button>
          }
          
          <span class="toolbar-spacer"></span>
          
          <button
            mat-icon-button
            (click)="onToggleTheme($event)"
            matTooltip="Toggle theme">
            <mat-icon>{{ themeService.isDarkMode() ? 'light_mode' : 'dark_mode' }}</mat-icon>
          </button>
        </mat-toolbar>

        <!-- Main Content Area -->
        <div class="content-area">
          <router-outlet></router-outlet>
        </div>
      </mat-sidenav-content>
    </mat-sidenav-container>
  `,
  styles: [`
    :host {
      display: block;
      height: 100vh;
    }

    .sidenav-container {
      height: 100%;
    }

    .sidenav {
      width: 280px;
      display: flex;
      flex-direction: column;
    }

    .sidenav-toolbar {
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 16px;
    }

    .logo-container {
      display: flex;
      align-items: center;
      gap: 12px;
    }

    .logo {
      width: 40px;
      height: 40px;
      border-radius: 50%;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      display: flex;
      align-items: center;
      justify-content: center;
      color: white;
      font-weight: bold;
      font-size: 20px;
    }

    .logo-text {
      font-size: 20px;
      font-weight: 500;
      color: var(--mat-sys-on-surface);
    }

    .nav-list {
      flex: 1;
      padding: 16px 0;
    }

    .active-link {
      background: var(--mat-sys-primary-container) !important;
      color: var(--mat-sys-on-primary-container) !important;
    }



    .sidenav-content {
      display: flex;
      flex-direction: column;
      height: 100%;
      background: var(--mat-sys-surface-container);
    }

    .main-toolbar {
      position: sticky;
      top: 0;
      z-index: 1000;
      box-shadow: var(--mat-sys-level2);
    }

    .toolbar-spacer {
      flex: 1;
    }

    .content-area {
      flex: 1;
      display: flex;
      flex-direction: column;
      overflow: hidden;
      background: var(--mat-sys-surface);
    }
  `]
})
export class LayoutComponent implements OnInit {
  isHandset = signal(false);

  constructor(
    private breakpointObserver: BreakpointObserver,
    public themeService: ThemeService,
  ) {}

  ngOnInit(): void {
    this.breakpointObserver.observe(Breakpoints.Handset)
      .subscribe(result => {
        this.isHandset.set(result.matches);
      });
  }

  onToggleTheme(event: MouseEvent): void {
    // Delegate to ThemeService which handles storage and transitions
    void this.themeService.toggleThemeWithTransition(event);
  }
}
