import { ChangeDetectionStrategy, Component, input, output } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { SettingsWebProviders } from '../settings.types';

@Component({
  selector: 'app-settings-tools-tab',
  imports: [
    FormsModule,
    MatCardModule,
    MatIconModule,
    MatSlideToggleModule,
    MatFormFieldModule,
    MatInputModule
  ],
  templateUrl: './settings-tools-tab.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
  host: { class: 'block h-full' }
})
export class SettingsToolsTabComponent {
  webProviders = input.required<SettingsWebProviders>();
  webProxy = input('');
  cronTimeout = input(5);
  webProxyChange = output<string>();
  cronTimeoutChange = output<number>();

  onWebProxyChange(value: string): void {
    this.webProxyChange.emit(value);
  }

  onCronTimeoutChange(value: number): void {
    this.cronTimeoutChange.emit(value);
  }
}
