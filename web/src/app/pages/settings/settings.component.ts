import { CommonModule } from '@angular/common';
import { Component, OnDestroy, OnInit, signal } from '@angular/core';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatTabsModule } from '@angular/material/tabs';
import { Subscription, catchError, map, of, switchMap, timer } from 'rxjs';
import { Config, ModelConfig } from '../../models/config.model';
import { ApiService } from '../../services/api.service';
import { SettingsChannelsTabComponent } from './components/settings-channels-tab.component';
import { SettingsHeaderComponent } from './components/settings-header.component';
import { SettingsHeartbeatTabComponent } from './components/settings-heartbeat-tab.component';
import { SettingsModelsTabComponent } from './components/settings-models-tab.component';
import { SettingsToolsTabComponent } from './components/settings-tools-tab.component';
import { SettingsChannelItem, SettingsWebProviders } from './settings.types';

@Component({
  selector: 'app-settings',
  standalone: true,
  imports: [
    CommonModule,
    MatTabsModule,
    MatSnackBarModule,
    MatProgressSpinnerModule,
    SettingsHeaderComponent,
    SettingsModelsTabComponent,
    SettingsChannelsTabComponent,
    SettingsToolsTabComponent,
    SettingsHeartbeatTabComponent
  ],
  templateUrl: './settings.component.html',
  host: { style: 'display: block; height: 100%; min-height: 0;' }
})
export class SettingsComponent implements OnInit, OnDestroy {
  config = signal<Config | null>(null);
  isLoading = signal(true);
  isSaving = signal(false);
  isRestarting = signal(false);
  defaultModel = signal('');
  private restartPollSub?: Subscription;

  // Tools state
  webProviders: SettingsWebProviders = {
    brave: { enabled: false, api_key: '', max_results: 5 },
    tavily: { enabled: false, api_key: '', max_results: 5 },
    duckduckgo: { enabled: true, max_results: 5 },
    perplexity: { enabled: false, api_key: '', max_results: 5 }
  };
  webProxy = '';
  cronTimeout = 5;
  heartbeatEnabled = true;
  heartbeatInterval = 30;

  channelConfigs = signal<SettingsChannelItem[]>([
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
    },
    {
      key: 'dingtalk',
      name: 'DingTalk',
      icon: 'notifications',
      description: 'DingTalk Bot integration',
      config: { enabled: false, client_id: '', client_secret: '', allow_from: [] },
      fields: [
        { key: 'client_id', label: 'Client ID' },
        { key: 'client_secret', label: 'Client Secret', type: 'password' },
        { key: 'allow_from', label: 'Allowed User IDs (comma separated)' }
      ]
    },
    {
      key: 'qq',
      name: 'QQ',
      icon: 'chat_bubble',
      description: 'QQ Bot integration',
      config: { enabled: false, app_id: '', app_secret: '', allow_from: [] },
      fields: [
        { key: 'app_id', label: 'App ID' },
        { key: 'app_secret', label: 'App Secret', type: 'password' },
        { key: 'allow_from', label: 'Allowed User IDs (comma separated)' }
      ]
    },
    {
      key: 'feishu',
      name: 'Feishu',
      icon: 'work',
      description: 'Feishu Bot integration',
      config: { enabled: false, app_id: '', app_secret: '', encrypt_key: '', verification_token: '', allow_from: [] },
      fields: [
        { key: 'app_id', label: 'App ID' },
        { key: 'app_secret', label: 'App Secret', type: 'password' },
        { key: 'encrypt_key', label: 'Encrypt Key', type: 'password' },
        { key: 'verification_token', label: 'Verification Token' },
        { key: 'allow_from', label: 'Allowed User IDs (comma separated)' }
      ]
    },
    {
      key: 'wecom',
      name: 'WeCom',
      icon: 'business',
      description: 'WeCom Bot integration',
      config: { enabled: false, token: '', encoding_aes_key: '', webhook_url: '', webhook_host: '0.0.0.0', webhook_port: 18793, webhook_path: '/webhook/wecom', allow_from: [] },
      fields: [
        { key: 'token', label: 'Token' },
        { key: 'encoding_aes_key', label: 'Encoding AES Key', type: 'password' },
        { key: 'webhook_url', label: 'Webhook URL' },
        { key: 'webhook_host', label: 'Webhook Host' },
        { key: 'webhook_port', label: 'Webhook Port', type: 'number' },
        { key: 'webhook_path', label: 'Webhook Path' }
      ]
    },
    {
      key: 'onebot',
      name: 'OneBot',
      icon: 'developer_board',
      description: 'OneBot protocol (NapCat/Go-CQHTTP)',
      config: { enabled: false, ws_url: 'ws://127.0.0.1:3001', access_token: '', reconnect_interval: 30, group_trigger_prefix: [], allow_from: [] },
      fields: [
        { key: 'ws_url', label: 'WebSocket URL' },
        { key: 'access_token', label: 'Access Token', type: 'password' },
        { key: 'reconnect_interval', label: 'Reconnect Interval (seconds)', type: 'number' },
        { key: 'allow_from', label: 'Allowed User IDs (comma separated)' }
      ]
    },
    {
      key: 'maixcam',
      name: 'MaixCam',
      icon: 'camera_alt',
      description: 'MaixCam AI camera integration',
      config: { enabled: false, host: '0.0.0.0', port: 18794, allow_from: [] },
      fields: [
        { key: 'host', label: 'Host' },
        { key: 'port', label: 'Port', type: 'number' },
        { key: 'allow_from', label: 'Allowed Device IDs (comma separated)' }
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

  ngOnDestroy(): void {
    this.restartPollSub?.unsubscribe();
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
        const channels = config.channels as unknown as Record<string, unknown>;
        channels[ch.key] = ch.config;
      }
    });

    this.isSaving.set(true);
    this.apiService.saveConfig(config).subscribe({
      next: () => {
        this.isSaving.set(false);
        this.restartGatewayAndWait();
      },
      error: (error) => {
        console.error('Failed to save config:', error);
        this.snackBar.open('Failed to save configuration', 'Close', { duration: 3000 });
        this.isSaving.set(false);
      }
    });
  }

  private restartGatewayAndWait(): void {
    this.isRestarting.set(true);

    const startPolling = () => this.pollGatewayReady();
    this.apiService.restartGateway().subscribe({
      next: () => startPolling(),
      error: (error) => {
        // Restart may start before this request returns; continue by polling.
        console.warn('Gateway restart request interrupted, polling status:', error);
        startPolling();
      }
    });
  }

  private pollGatewayReady(): void {
    const maxAttempts = 80;
    let attempts = 0;

    this.restartPollSub?.unsubscribe();
    this.restartPollSub = timer(2000, 1500).pipe(
      map(() => ++attempts),
      switchMap((attempt) =>
        this.apiService.getStatus().pipe(
          map(() => ({ attempt, ok: true })),
          catchError(() => of({ attempt, ok: false }))
        )
      )
    ).subscribe(({ attempt, ok }) => {
      if (ok) {
        this.restartPollSub?.unsubscribe();
        this.isRestarting.set(false);
        this.snackBar.open('Configuration updated successfully', 'Close', { duration: 3000 });
        return;
      }

      if (attempt >= maxAttempts) {
        this.restartPollSub?.unsubscribe();
        this.isRestarting.set(false);
        this.snackBar.open('Config saved, but gateway restart timed out', 'Close', { duration: 5000 });
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
