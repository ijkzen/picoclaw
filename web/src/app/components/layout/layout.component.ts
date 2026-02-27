import { BreakpointObserver, Breakpoints } from '@angular/cdk/layout';
import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, DestroyRef, OnInit, inject, signal } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatListModule } from '@angular/material/list';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatDividerModule } from '@angular/material/divider';
import { MatTooltipModule } from '@angular/material/tooltip';
import { NavigationEnd, Router, RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';
import { filter } from 'rxjs';
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
    MatDividerModule,
    MatTooltipModule
  ],
  templateUrl: './layout.component.html',
  host: { style: 'display: block; height: 100dvh;' },
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class LayoutComponent implements OnInit {
  private readonly destroyRef = inject(DestroyRef);
  private readonly router = inject(Router);

  isHandset = signal(false);
  pageTitle = signal('Chat');
  constructor(
    private breakpointObserver: BreakpointObserver,
    public themeService: ThemeService,
  ) {}

  ngOnInit(): void {
    this.setPageTitle(this.router.url);
    this.breakpointObserver.observe(Breakpoints.Handset)
      .pipe(takeUntilDestroyed(this.destroyRef))
      .subscribe(result => {
        this.isHandset.set(result.matches);
      });

    this.router.events
      .pipe(
        filter((event): event is NavigationEnd => event instanceof NavigationEnd),
        takeUntilDestroyed(this.destroyRef)
      )
      .subscribe((event) => {
        this.setPageTitle(event.urlAfterRedirects);
      });
  }

  onToggleTheme(event: MouseEvent): void {
    // Delegate to ThemeService which handles storage and transitions
    void this.themeService.toggleThemeWithTransition(event);
  }

  private setPageTitle(url: string): void {
    this.pageTitle.set(url.startsWith('/settings') ? 'Settings' : 'Chat');
  }
}
