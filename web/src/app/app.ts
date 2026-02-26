import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet],
  template: `
    <div class="h-screen w-screen overflow-hidden">
      <router-outlet></router-outlet>
    </div>
  `,
  styles: ``
})
export class App {
  title = 'PicoClaw Web';
}
