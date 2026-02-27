import { ChannelConfig } from '../../models/config.model';

export interface SettingsChannelField {
  key: string;
  label: string;
  type?: 'password' | 'number';
  hint?: string;
  fullWidth?: boolean;
}

export interface SettingsChannelItem {
  key: string;
  name: string;
  icon: string;
  description: string;
  config: ChannelConfig;
  fields: SettingsChannelField[];
}

export interface SettingsWebProvider {
  enabled: boolean;
  api_key?: string;
  max_results: number;
}

export interface SettingsWebProviders {
  brave: SettingsWebProvider;
  tavily: SettingsWebProvider;
  duckduckgo: SettingsWebProvider;
  perplexity: SettingsWebProvider;
}
