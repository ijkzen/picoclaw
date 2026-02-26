import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: '',
    loadComponent: () => import('./components/layout/layout.component').then(m => m.LayoutComponent),
    children: [
      { path: '', redirectTo: 'chat', pathMatch: 'full' },
      { path: 'chat', loadComponent: () => import('./pages/chat/chat.component').then(m => m.ChatComponent) },
      { path: 'settings', loadComponent: () => import('./pages/settings/settings.component').then(m => m.SettingsComponent) }
    ]
  }
];
