export interface ModelConfig {
  model_name: string;
  model: string;
  api_key: string;
  api_base?: string;
  proxy?: string;
  rpm?: number;
}

export interface ChannelConfig {
  enabled: boolean;
  [key: string]: any;
}

export interface ChannelsConfig {
  telegram: ChannelConfig;
  discord: ChannelConfig;
  slack: ChannelConfig;
  line: ChannelConfig;
  qq: ChannelConfig;
  dingtalk: ChannelConfig;
  wecom: ChannelConfig;
  wecom_app: ChannelConfig;
  feishu: ChannelConfig;
  whatsapp: ChannelConfig;
  onebot: ChannelConfig;
  maixcam: ChannelConfig;
}

export interface WebToolConfig {
  enabled: boolean;
  api_key?: string;
  max_results: number;
}

export interface ToolsConfig {
  web: {
    brave: WebToolConfig;
    tavily: WebToolConfig;
    duckduckgo: WebToolConfig;
    perplexity: WebToolConfig;
    proxy?: string;
  };
  cron: {
    exec_timeout_minutes: number;
  };
}

export interface HeartbeatConfig {
  enabled: boolean;
  interval: number;
}

export interface AgentDefaults {
  workspace: string;
  restrict_to_workspace: boolean;
  model_name: string;
  max_tokens: number;
  temperature?: number;
  max_tool_iterations: number;
}

export interface Config {
  agents: {
    defaults: AgentDefaults;
  };
  model_list: ModelConfig[];
  channels: ChannelsConfig;
  tools: ToolsConfig;
  heartbeat: HeartbeatConfig;
  gateway: {
    host: string;
    port: number;
  };
}

export interface Message {
  id: string;
  role: 'user' | 'assistant' | 'system';
  content: string;
  timestamp: Date;
  showRawContent?: boolean;
  isComplete?: boolean;
}
