import { ChangeDetectionStrategy, Component, input, output } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';

@Component({
  selector: 'app-settings-header',
  imports: [MatCardModule, MatButtonModule, MatIconModule],
  templateUrl: './settings-header.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class SettingsHeaderComponent {
  isSaving = input(false);
  isRestarting = input(false);
  save = output<void>();
}
