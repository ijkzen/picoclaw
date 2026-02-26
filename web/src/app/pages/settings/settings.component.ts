import { Component, signal, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MatTabsModule } from '@angular/material/tabs';
import { MatCardModule } from '@angular/material/card';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { MatSelectModule } from '@angular/material/select';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatDividerModule } from '@angular/material/divider';
import { MatExpansionModule } from '@angular/material/expansion';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatListModule } from '@angular/material/list';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { Config, ModelConfig } from '../../models/config.model';
import { ApiService } from '../../services/api.service';

@Component({
  selector: 'app-settings',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MatTabsModule,
    MatCardModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
    MatSlideToggleModule,
    MatSelectModule,
    MatSnackBarModule,
    MatDividerModule,
    MatExpansionModule,
    MatFormFieldModule,
    MatListModule,
    MatProgressSpinnerModule
  ],
  template: `
    <div class="settings-container">
      <mat-card class="header-card">
        <mat-card-header>
          <mat-card-title>Settings</mat-card-title>
          <mat-card-subtitle>Configure your PicoClaw instance</mat-card-subtitle>
        </mat-card-header>
        <mat-card-actions align="end">
          <button mat-raised-button color="primary" (click)="saveConfig()" [disabled]="isSaving()">
            <mat-icon>save</mat-icon>
            {{ isSaving() ? 'Saving...' : 'Save Changes' }}
          </button>
        </mat-card-actions>
      </mat-card>

      @if (isLoading()) {
        <div class="loading-container">
          <mat-spinner></mat-spinner>
        </div>
      } @else {
        <mat-tab-group animationDuration="0ms">
          <!-- Models Tab -->
          <mat-tab label="Models">
            <div class="tab-content">
              <mat-card>
                <mat-card-header>
                  <mat-icon mat-card-avatar>stars</mat-icon>
                  <mat-card-title>Default Model</mat-card-title>
                </mat-card-header>
                <mat-card-content>
                  <mat-form-field appearance="outline" style="width: 100%;">
                    <mat-select [(value)]="defaultModel">
                      @for (model of config()?.model_list; track model.model_name) {
                        <mat-option [value]="model.model_name">
                          {{ model.model_name }} ({{ model.model }})
                        </mat-option>
                      }
                    </mat-select>
                  </mat-form-field>
                </mat-card-content>
              </mat-card>

              <mat-divider style="margin: 16px 0;"></mat-divider>

              <mat-card>
                <mat-card-header>
                  <mat-card-title>Model Configuration</mat-card-title>
                </mat-card-header>
                <mat-card-actions align="end">
                  <button mat-stroked-button color="primary" (click)="addNewModel()">
                    <mat-icon>add</mat-icon>
                    Add Model
                  </button>
                </mat-card-actions>
              </mat-card>

              <mat-accordion>
                @for (model of config()?.model_list; track model.model_name; let i = $index) {
                  <mat-expansion-panel>
                    <mat-expansion-panel-header>
                      <mat-panel-title>
                        @if (model.model_name === defaultModel()) {
                          <mat-icon style="margin-right: 8px; color: var(--mat-sys-primary);">star</mat-icon>
                        }
                        {{ model.model_name }}
                      </mat-panel-title>
                      <mat-panel-description>
                        {{ model.model }}
                      </mat-panel-description>
                    </mat-expansion-panel-header>

                    <div style="display: flex; flex-direction: column; gap: 16px;">
                      <mat-form-field appearance="outline">
                        <mat-label>Model Name</mat-label>
                        <input matInput [(ngModel)]="model.model_name">
                      </mat-form-field>

                      <mat-form-field appearance="outline">
                        <mat-label>Model ID</mat-label>
                        <input matInput [(ngModel)]="model.model">
                      </mat-form-field>

                      <mat-form-field appearance="outline">
                        <mat-label>API Key</mat-label>
                        <input matInput [(ngModel)]="model.api_key" type="password">
                      </mat-form-field>

                      <mat-form-field appearance="outline">
                        <mat-label>API Base URL (Optional)</mat-label>
                        <input matInput [(ngModel)]="model.api_base">
                      </mat-form-field>
                    </div>

                    <mat-action-row>
                      @if (model.model_name !== defaultModel()) {
                        <button mat-button color="primary" (click)="setDefaultModel(model.model_name)">
                          Set as Default
                        </button>
                      }
                      <button mat-button color="warn" (click)="deleteModel(i)">
                        <mat-icon>delete</mat-icon>
                        Delete
                      </button>
                    </mat-action-row>
                  </mat-expansion-panel>
                }
              </mat-accordion>
            </div>
          </mat-tab>

          <!-- Channels Tab -->
          <mat-tab label="Channels">
            <div class="tab-content">
              <mat-card>
                <mat-card-header>
                  <mat-icon mat-card-avatar>chat</mat-icon>
                  <mat-card-title>Chat Channel Configuration</mat-card-title>
                  <mat-card-subtitle>Configure integrations with chat platforms</mat-card-subtitle>
                </mat-card-header>
              </mat-card>

              <@for (channel of channelConfigs(); track channel.key) {
                <mat-card style="margin-top: 16px;">
                  <mat-card-header>
                    <mat-icon mat-card-avatar>{{ channel.icon }}</mat-icon>
                    <mat-card-title>{{ channel.name }}</mat-card-title>
                    <mat-card-subtitle>{{ channel.description }}</mat-card-subtitle>
                  </mat-card-header>
                  
                  <mat-card-content>
                    <mat-slide-toggle
                      [(ngModel)]="channel.config.enabled"
                      color="primary">
                      {{ channel.config.enabled ? 'Enabled' : 'Disabled' }}
                    </mat-slide-toggle>

                    @if (channel.config.enabled) {
                      <mat-divider style="margin: 16px 0;"></mat-divider>
                      
                      <@for (field of channel.fields; track field.key) {
                        <mat-form-field appearance="outline" style="width: 100%; margin-bottom: 8px;">
                          <mat-label>{{ field.label }}</mat-label>
                          @if (field.type === 'password') {
                            <input matInput [(ngModel)]="channel.config[field.key]" type="password">
                          } @else if (field.type === 'number') {
                            <input matInput [(ngModel)]="channel.config[field.key]" type="number">
                          } @else {
                            <input matInput [(ngModel)]="channel.config[field.key]">
                          }
                          @if (field.hint) {
                            <mat-hint>{{ field.hint }}</mat-hint>
                          }
                        </mat-form-field>
                      }
                    }
                  </mat-card-content>
                </mat-card>
              }
            </div>
          </mat-tab>

          <!-- Tools Tab -->
          <mat-tab label="Tools">
            <div class="tab-content">
              <mat-card>
                <mat-card-header>
                  <mat-icon mat-card-avatar>search</mat-icon>
                  <mat-card-title>Web Search</mat-card-title>
                </mat-card-header>
                <mat-card-content>
                  <!-- Brave -->
                  <mat-card appearance="outlined" style="margin-bottom: 16px;">
                    <mat-card-header>
                      <mat-card-title>Brave</mat-card-title>
                    </mat-card-header>
                    <mat-card-content>
                      <mat-slide-toggle [(ngModel)]="webProviders.brave.enabled" color="primary">
                        Enabled
                      </mat-slide-toggle>

                      @if (webProviders.brave.enabled) {
                        <mat-form-field appearance="outline" style="width: 100%; margin-top: 16px;">
                          <mat-label>API Key</mat-label>
                          <input matInput [(ngModel)]="webProviders.brave.api_key" type="password">
                        </mat-form-field>
                        
                        <mat-form-field appearance="outline" style="width: 100%;">
                          <mat-label>Max Results</mat-label>
                          <input matInput type="number" [(ngModel)]="webProviders.brave.max_results">
                        </mat-form-field>
                      }
                    </mat-card-content>
                  </mat-card>

                  <!-- Tavily -->
                  <mat-card appearance="outlined" style="margin-bottom: 16px;">
                    <mat-card-header>
                      <mat-card-title>Tavily</mat-card-title>
                    </mat-card-header>
                    <mat-card-content>
                      <mat-slide-toggle [(ngModel)]="webProviders.tavily.enabled" color="primary">
                        Enabled
                      </mat-slide-toggle>

                      @if (webProviders.tavily.enabled) {
                        <mat-form-field appearance="outline" style="width: 100%; margin-top: 16px;">
                          <mat-label>API Key</mat-label>
                          <input matInput [(ngModel)]="webProviders.tavily.api_key" type="password">
                        </mat-form-field>
                        
                        <mat-form-field appearance="outline" style="width: 100%;">
                          <mat-label>Max Results</mat-label>
                          <input matInput type="number" [(ngModel)]="webProviders.tavily.max_results">
                        </mat-form-field>
                      }
                    </mat-card-content>
                  </mat-card>

                  <!-- DuckDuckGo -->
                  <mat-card appearance="outlined" style="margin-bottom: 16px;">
                    <mat-card-header>
                      <mat-card-title>DuckDuckGo</mat-card-title>
                    </mat-card-header>
                    <mat-card-content>
                      <mat-slide-toggle [(ngModel)]="webProviders.duckduckgo.enabled" color="primary">
                        Enabled
                      </mat-slide-toggle>

                      @if (webProviders.duckduckgo.enabled) {
                        <mat-form-field appearance="outline" style="width: 100%; margin-top: 16px;">
                          <mat-label>Max Results</mat-label>
                          <input matInput type="number" [(ngModel)]="webProviders.duckduckgo.max_results">
                        </mat-form-field>
                      }
                    </mat-card-content>
                  </mat-card>

                  <!-- Perplexity -->
                  <mat-card appearance="outlined" style="margin-bottom: 16px;">
                    <mat-card-header>
                      <mat-card-title>Perplexity</mat-card-title>
                    </mat-card-header>
                    <mat-card-content>
                      <mat-slide-toggle [(ngModel)]="webProviders.perplexity.enabled" color="primary">
                        Enabled
                      </mat-slide-toggle>

                      @if (webProviders.perplexity.enabled) {
                        <mat-form-field appearance="outline" style="width: 100%; margin-top: 16px;">
                          <mat-label>API Key</mat-label>
                          <input matInput [(ngModel)]="webProviders.perplexity.api_key" type="password">
                        </mat-form-field>
                        
                        <mat-form-field appearance="outline" style="width: 100%;">
                          <mat-label>Max Results</mat-label>
                          <input matInput type="number" [(ngModel)]="webProviders.perplexity.max_results">
                        </mat-form-field>
                      }
                    </mat-card-content>
                  </mat-card>

                  <mat-form-field appearance="outline" style="width: 100%; margin-top: 16px;">
                    <mat-label>Proxy (Optional)</mat-label>
                    <input matInput [(ngModel)]="webProxy" placeholder="http://proxy.example.com:8080">
                  </mat-form-field>
                </mat-card-content>
              </mat-card>

              <mat-card style="margin-top: 16px;">
                <mat-card-header>
                  <mat-icon mat-card-avatar>schedule</mat-icon>
                  <mat-card-title>Scheduled Tasks</mat-card-title>
                </mat-card-header>
                <mat-card-content>
                  <mat-form-field appearance="outline" style="width: 100%;">
                    <mat-label>Execution Timeout (minutes)</mat-label>
                    <input matInput type="number" [(ngModel)]="cronTimeout">
                    <mat-hint>Set to 0 for no timeout</mat-hint>
                  </mat-form-field>
                </mat-card-content>
              </mat-card>
            </div>
          </mat-tab>

          <!-- Heartbeat Tab -->
          <mat-tab label="Heartbeat">
            <div class="tab-content">
              <mat-card>
                <mat-card-header>
                  <mat-icon mat-card-avatar>favorite</mat-icon>
                  <mat-card-title>Heartbeat Configuration</mat-card-title>
                  <mat-card-subtitle>Configure periodic task execution</mat-card-subtitle>
                </mat-card-header>
                
                <mat-card-content>
                  <mat-slide-toggle [(ngModel)]="heartbeatEnabled" color="primary">
                    {{ heartbeatEnabled ? 'Enabled' : 'Disabled' }}
                  </mat-slide-toggle>

                  @if (heartbeatEnabled) {
                    <mat-divider style="margin: 16px 0;"></mat-divider>
                    
                    <mat-form-field appearance="outline" style="width: 100%;">
                      <mat-label>Interval (minutes)</mat-label>
                      <input matInput type="number" [(ngModel)]="heartbeatInterval" min="5">
                      <mat-hint>Minimum 5 minutes</mat-hint>
                    </mat-form-field>
                  }
                </mat-card-content>
              </mat-card>

              <mat-card style="margin-top: 16px;">
                <mat-card-content>
                  <p style="display: flex; align-items: center; gap: 8px;">
                    <mat-icon>info</mat-icon>
                    Heartbeat tasks are defined in the HEARTBEAT.md file in your workspace.
                  </p>
                </mat-card-content>
              </mat-card>
            </div>
          </mat-tab>
        </mat-tab-group>
      }
    </div>
  `,
  styles: [`
    :host {
      display: block;
    }

    .settings-container {
      max-width: 900px;
      margin: 0 auto;
      padding: 16px;
    }

    .header-card {
      margin-bottom: 16px;
    }

    .loading-container {
      display: flex;
      justify-content: center;
      padding: 48px;
    }

    .tab-content {
      padding: 16px 0;
    }
  `]
})
export class SettingsComponent implements OnInit {
  config = signal<Config | null>(null);
  isLoading = signal(true);
  isSaving = signal(false);
  defaultModel = signal('');

  // Tools state
  webProviders = {
    brave: { enabled: false, api_key: '', max_results: 5 },
    tavily: { enabled: false, api_key: '', max_results: 5 },
    duckduckgo: { enabled: true, max_results: 5 },
    perplexity: { enabled: false, api_key: '', max_results: 5 }
  };
  webProxy = '';
  cronTimeout = 5;
  heartbeatEnabled = true;
  heartbeatInterval = 30;

  channelConfigs = signal<Array<{
    key: string;
    name: string;
    icon: string;
    description: string;
    config: any;
    fields: Array<{ key: string; label: string; type?: string; hint?: string; fullWidth?: boolean }>;
  }>>([
    {
      key: 'telegram',
      name: 'Telegram',
      icon: 'send',
      description: 'Telegram Bot integration',
      config: { enabled: false, token: '', allow_from: [] },
      fields: [
        { key: 'token', label: 'Bot Token', type: 'password' },
        { key: 'allow_from', label: 'Allowed User IDs (comma separated)' }
      ]
    },
    {
      key: 'discord',
      name: 'Discord',
      icon: 'chat',
      description: 'Discord Bot integration',
      config: { enabled: false, token: '', allow_from: [], mention_only: false },
      fields: [
        { key: 'token', label: 'Bot Token', type: 'password' },
        { key: 'allow_from', label: 'Allowed User IDs (comma separated)' }
      ]
    },
    {
      key: 'slack',
      name: 'Slack',
      icon: 'forum',
      description: 'Slack Bot integration',
      config: { enabled: false, bot_token: '', app_token: '', allow_from: [] },
      fields: [
        { key: 'bot_token', label: 'Bot Token', type: 'password' },
        { key: 'app_token', label: 'App Token', type: 'password' },
        { key: 'allow_from', label: 'Allowed User IDs (comma separated)' }
      ]
    },
    {
      key: 'line',
      name: 'LINE',
      icon: 'message',
      description: 'LINE Messaging API integration',
      config: { enabled: false, channel_secret: '', channel_access_token: '', webhook_host: '0.0.0.0', webhook_port: 18791, webhook_path: '/webhook/line', allow_from: [] },
      fields: [
        { key: 'channel_secret', label: 'Channel Secret', type: 'password' },
        { key: 'channel_access_token', label: 'Channel Access Token', type: 'password', fullWidth: true },
        { key: 'webhook_host', label: 'Webhook Host' },
        { key: 'webhook_port', label: 'Webhook Port', type: 'number' },
        { key: 'webhook_path', label: 'Webhook Path' }
      ]
    }
  ]);

  constructor(
    private apiService: ApiService,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {
    this.loadConfig();
  }

  loadConfig(): void {
    this.apiService.getConfig().subscribe({
      next: (config) => {
        this.config.set(config);
        this.defaultModel.set(config.agents?.defaults?.model_name || '');

        // Update tools state
        if (config.tools?.web) {
          this.webProviders.brave = { ...this.webProviders.brave, ...config.tools.web.brave };
          this.webProviders.tavily = { ...this.webProviders.tavily, ...config.tools.web.tavily };
          this.webProviders.duckduckgo = { ...this.webProviders.duckduckgo, ...config.tools.web.duckduckgo };
          this.webProviders.perplexity = { ...this.webProviders.perplexity, ...config.tools.web.perplexity };
          this.webProxy = config.tools.web.proxy || '';
        }
        if (config.tools?.cron) {
          this.cronTimeout = config.tools.cron.exec_timeout_minutes;
        }
        if (config.heartbeat) {
          this.heartbeatEnabled = config.heartbeat.enabled;
          this.heartbeatInterval = config.heartbeat.interval;
        }

        // Update channel configs with actual values
        this.channelConfigs.update(channels =>
          channels.map(ch => ({
            ...ch,
            config: config.channels?.[ch.key as keyof typeof config.channels] || ch.config
          }))
        );

        this.isLoading.set(false);
      },
      error: (error) => {
        console.error('Failed to load config:', error);
        this.snackBar.open('Failed to load configuration', 'Close', { duration: 3000 });
        this.isLoading.set(false);
      }
    });
  }

  saveConfig(): void {
    const config = this.config();
    if (!config) return;

    // Update default model
    if (config.agents?.defaults) {
      config.agents.defaults.model_name = this.defaultModel();
    }

    // Update web tools
    config.tools.web = {
      brave: this.webProviders.brave,
      tavily: this.webProviders.tavily,
      duckduckgo: this.webProviders.duckduckgo,
      perplexity: this.webProviders.perplexity,
      proxy: this.webProxy
    };

    // Update cron
    config.tools.cron = { exec_timeout_minutes: this.cronTimeout };

    // Update heartbeat
    config.heartbeat = { enabled: this.heartbeatEnabled, interval: this.heartbeatInterval };

    // Update channels from channelConfigs
    this.channelConfigs().forEach(ch => {
      if (config.channels) {
        (config.channels as any)[ch.key] = ch.config;
      }
    });

    this.isSaving.set(true);
    this.apiService.saveConfig(config).subscribe({
      next: () => {
        this.snackBar.open('Configuration saved successfully', 'Close', { duration: 3000 });
        this.isSaving.set(false);
      },
      error: (error) => {
        console.error('Failed to save config:', error);
        this.snackBar.open('Failed to save configuration', 'Close', { duration: 3000 });
        this.isSaving.set(false);
      }
    });
  }

  addNewModel(): void {
    const newModel: ModelConfig = {
      model_name: 'new-model',
      model: '',
      api_key: ''
    };

    this.config.update(cfg => {
      if (cfg) {
        return {
          ...cfg,
          model_list: [...cfg.model_list, newModel]
        };
      }
      return cfg;
    });
  }

  setDefaultModel(modelName: string): void {
    this.defaultModel.set(modelName);
  }

  deleteModel(index: number): void {
    this.config.update(cfg => {
      if (cfg) {
        const newList = [...cfg.model_list];
        newList.splice(index, 1);
        return { ...cfg, model_list: newList };
      }
      return cfg;
    });
  }
}
