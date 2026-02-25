package channel

import (
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
)

// InitFormForChannel initializes the form for a specific channel
func (m *Model) InitFormForChannel(channelName string) {
	m.FormFields = nil
	m.FormInputs = nil
	m.FormCursor = 0

	switch channelName {
	case "telegram":
		m.InitTelegramForm()
	case "discord":
		m.InitDiscordForm()
	case "slack":
		m.InitSlackForm()
	case "qq":
		m.InitQQForm()
	case "dingtalk":
		m.InitDingTalkForm()
	case "wecom":
		m.InitWeComForm()
	case "wecom_app":
		m.InitWeComAppForm()
	case "feishu":
		m.InitFeishuForm()
	case "line":
		m.InitLineForm()
	case "onebot":
		m.InitOneBotForm()
	case "maixcam":
		m.InitMaixCamForm()
	case "whatsapp":
		m.InitWhatsAppForm()
	}

	// Initialize text inputs for non-bool fields
	for i, field := range m.FormFields {
		textInput := textinput.New()
		textInput.Placeholder = field.Placeholder
		textInput.SetValue(field.Value)
		textInput.Width = 50
		if field.Type == FieldTypePassword {
			textInput.EchoMode = textinput.EchoPassword
			textInput.EchoCharacter = '•'
		}

		m.FormInputs = append(m.FormInputs, textInput)

		// Focus first non-bool field
		if field.Type != FieldTypeBool && m.FormCursor == 0 {
			m.FormCursor = i
			m.FormInputs[i].Focus()
		}
	}
}

// InitTelegramForm initializes the Telegram configuration form
func (m *Model) InitTelegramForm() {
	cfg := m.Config.Channels.Telegram
	m.FormFields = []FormField{
		{Name: "enabled", Label: "Enabled", Type: FieldTypeBool, BoolValue: cfg.Enabled},
		{Name: "token", Label: "Bot Token", Type: FieldTypePassword, Value: cfg.Token, Placeholder: "123456789:ABCdefGHIjklMNOpqrsTUVwxyz", Required: true},
		{Name: "proxy", Label: "Proxy URL", Type: FieldTypeText, Value: cfg.Proxy, Placeholder: "socks5://127.0.0.1:1080"},
		{Name: "allow_from", Label: "Allowed Users", Type: FieldTypeArray, Value: StringArrayToString(cfg.AllowFrom), Placeholder: "user1, user2, user3"},
	}
}

// InitDiscordForm initializes the Discord configuration form
func (m *Model) InitDiscordForm() {
	cfg := m.Config.Channels.Discord
	m.FormFields = []FormField{
		{Name: "enabled", Label: "Enabled", Type: FieldTypeBool, BoolValue: cfg.Enabled},
		{Name: "token", Label: "Bot Token", Type: FieldTypePassword, Value: cfg.Token, Placeholder: "MTAxMDE...", Required: true},
		{Name: "allow_from", Label: "Allowed Users", Type: FieldTypeArray, Value: StringArrayToString(cfg.AllowFrom), Placeholder: "user1, user2"},
		{Name: "mention_only", Label: "Mention Only", Type: FieldTypeBool, BoolValue: cfg.MentionOnly},
	}
}

// InitSlackForm initializes the Slack configuration form
func (m *Model) InitSlackForm() {
	cfg := m.Config.Channels.Slack
	m.FormFields = []FormField{
		{Name: "enabled", Label: "Enabled", Type: FieldTypeBool, BoolValue: cfg.Enabled},
		{Name: "bot_token", Label: "Bot Token", Type: FieldTypePassword, Value: cfg.BotToken, Placeholder: "xoxb-...", Required: true},
		{Name: "app_token", Label: "App Token", Type: FieldTypePassword, Value: cfg.AppToken, Placeholder: "xapp-..."},
		{Name: "allow_from", Label: "Allowed Users", Type: FieldTypeArray, Value: StringArrayToString(cfg.AllowFrom), Placeholder: "user1, user2"},
	}
}

// InitQQForm initializes the QQ configuration form
func (m *Model) InitQQForm() {
	cfg := m.Config.Channels.QQ
	m.FormFields = []FormField{
		{Name: "enabled", Label: "Enabled", Type: FieldTypeBool, BoolValue: cfg.Enabled},
		{Name: "app_id", Label: "App ID", Type: FieldTypeText, Value: cfg.AppID, Placeholder: "12345678", Required: true},
		{Name: "app_secret", Label: "App Secret", Type: FieldTypePassword, Value: cfg.AppSecret, Placeholder: "Your app secret", Required: true},
		{Name: "allow_from", Label: "Allowed Users", Type: FieldTypeArray, Value: StringArrayToString(cfg.AllowFrom), Placeholder: "user1, user2"},
	}
}

// InitDingTalkForm initializes the DingTalk configuration form
func (m *Model) InitDingTalkForm() {
	cfg := m.Config.Channels.DingTalk
	m.FormFields = []FormField{
		{Name: "enabled", Label: "Enabled", Type: FieldTypeBool, BoolValue: cfg.Enabled},
		{Name: "client_id", Label: "Client ID", Type: FieldTypeText, Value: cfg.ClientID, Placeholder: "dingxxxxxxxxxxxxxxxx", Required: true},
		{Name: "client_secret", Label: "Client Secret", Type: FieldTypePassword, Value: cfg.ClientSecret, Placeholder: "Your client secret", Required: true},
		{Name: "allow_from", Label: "Allowed Users", Type: FieldTypeArray, Value: StringArrayToString(cfg.AllowFrom), Placeholder: "user1, user2"},
	}
}

// InitWeComForm initializes the WeCom Bot configuration form
func (m *Model) InitWeComForm() {
	cfg := m.Config.Channels.WeCom
	m.FormFields = []FormField{
		{Name: "enabled", Label: "Enabled", Type: FieldTypeBool, BoolValue: cfg.Enabled},
		{Name: "token", Label: "Token", Type: FieldTypePassword, Value: cfg.Token, Placeholder: "Webhook token", Required: true},
		{Name: "encoding_aes_key", Label: "Encoding AES Key", Type: FieldTypePassword, Value: cfg.EncodingAESKey, Placeholder: "Optional AES key"},
		{Name: "webhook_url", Label: "Webhook URL", Type: FieldTypeText, Value: cfg.WebhookURL, Placeholder: "https://qyapi.weixin.qq.com/..."},
		{Name: "webhook_host", Label: "Webhook Host", Type: FieldTypeText, Value: cfg.WebhookHost, Placeholder: "0.0.0.0"},
		{Name: "webhook_port", Label: "Webhook Port", Type: FieldTypeNumber, Value: strconv.Itoa(cfg.WebhookPort), Placeholder: "8080"},
		{Name: "webhook_path", Label: "Webhook Path", Type: FieldTypeText, Value: cfg.WebhookPath, Placeholder: "/webhook/wecom"},
		{Name: "allow_from", Label: "Allowed Users", Type: FieldTypeArray, Value: StringArrayToString(cfg.AllowFrom), Placeholder: "user1, user2"},
		{Name: "reply_timeout", Label: "Reply Timeout", Type: FieldTypeNumber, Value: strconv.Itoa(cfg.ReplyTimeout), Placeholder: "30"},
	}
}

// InitWeComAppForm initializes the WeCom App configuration form
func (m *Model) InitWeComAppForm() {
	cfg := m.Config.Channels.WeComApp
	m.FormFields = []FormField{
		{Name: "enabled", Label: "Enabled", Type: FieldTypeBool, BoolValue: cfg.Enabled},
		{Name: "corp_id", Label: "Corp ID", Type: FieldTypeText, Value: cfg.CorpID, Placeholder: "wwxxxxxxxxxxxxxxxx", Required: true},
		{Name: "corp_secret", Label: "Corp Secret", Type: FieldTypePassword, Value: cfg.CorpSecret, Placeholder: "Your corp secret", Required: true},
		{Name: "agent_id", Label: "Agent ID", Type: FieldTypeNumber, Value: strconv.FormatInt(cfg.AgentID, 10), Placeholder: "1000002", Required: true},
		{Name: "token", Label: "Token", Type: FieldTypeText, Value: cfg.Token, Placeholder: "Verification token"},
		{Name: "encoding_aes_key", Label: "Encoding AES Key", Type: FieldTypePassword, Value: cfg.EncodingAESKey, Placeholder: "Optional AES key"},
		{Name: "webhook_host", Label: "Webhook Host", Type: FieldTypeText, Value: cfg.WebhookHost, Placeholder: "0.0.0.0"},
		{Name: "webhook_port", Label: "Webhook Port", Type: FieldTypeNumber, Value: strconv.Itoa(cfg.WebhookPort), Placeholder: "8080"},
		{Name: "webhook_path", Label: "Webhook Path", Type: FieldTypeText, Value: cfg.WebhookPath, Placeholder: "/webhook/wecom"},
		{Name: "allow_from", Label: "Allowed Users", Type: FieldTypeArray, Value: StringArrayToString(cfg.AllowFrom), Placeholder: "user1, user2"},
		{Name: "reply_timeout", Label: "Reply Timeout", Type: FieldTypeNumber, Value: strconv.Itoa(cfg.ReplyTimeout), Placeholder: "30"},
	}
}

// InitFeishuForm initializes the Feishu configuration form
func (m *Model) InitFeishuForm() {
	cfg := m.Config.Channels.Feishu
	m.FormFields = []FormField{
		{Name: "enabled", Label: "Enabled", Type: FieldTypeBool, BoolValue: cfg.Enabled},
		{Name: "app_id", Label: "App ID", Type: FieldTypeText, Value: cfg.AppID, Placeholder: "cli_xxxxxxxxxxxxxxxx", Required: true},
		{Name: "app_secret", Label: "App Secret", Type: FieldTypePassword, Value: cfg.AppSecret, Placeholder: "Your app secret", Required: true},
		{Name: "encrypt_key", Label: "Encrypt Key", Type: FieldTypePassword, Value: cfg.EncryptKey, Placeholder: "Optional encrypt key"},
		{Name: "verification_token", Label: "Verification Token", Type: FieldTypeText, Value: cfg.VerificationToken, Placeholder: "Verification token"},
		{Name: "allow_from", Label: "Allowed Users", Type: FieldTypeArray, Value: StringArrayToString(cfg.AllowFrom), Placeholder: "user1, user2"},
	}
}

// InitLineForm initializes the LINE configuration form
func (m *Model) InitLineForm() {
	cfg := m.Config.Channels.LINE
	m.FormFields = []FormField{
		{Name: "enabled", Label: "Enabled", Type: FieldTypeBool, BoolValue: cfg.Enabled},
		{Name: "channel_secret", Label: "Channel Secret", Type: FieldTypePassword, Value: cfg.ChannelSecret, Placeholder: "Your channel secret", Required: true},
		{Name: "channel_access_token", Label: "Channel Access Token", Type: FieldTypePassword, Value: cfg.ChannelAccessToken, Placeholder: "Your access token", Required: true},
		{Name: "webhook_host", Label: "Webhook Host", Type: FieldTypeText, Value: cfg.WebhookHost, Placeholder: "0.0.0.0"},
		{Name: "webhook_port", Label: "Webhook Port", Type: FieldTypeNumber, Value: strconv.Itoa(cfg.WebhookPort), Placeholder: "8080"},
		{Name: "webhook_path", Label: "Webhook Path", Type: FieldTypeText, Value: cfg.WebhookPath, Placeholder: "/webhook/line"},
		{Name: "allow_from", Label: "Allowed Users", Type: FieldTypeArray, Value: StringArrayToString(cfg.AllowFrom), Placeholder: "user1, user2"},
	}
}

// InitOneBotForm initializes the OneBot configuration form
func (m *Model) InitOneBotForm() {
	cfg := m.Config.Channels.OneBot
	m.FormFields = []FormField{
		{Name: "enabled", Label: "Enabled", Type: FieldTypeBool, BoolValue: cfg.Enabled},
		{Name: "ws_url", Label: "WebSocket URL", Type: FieldTypeText, Value: cfg.WSUrl, Placeholder: "ws://127.0.0.1:3001", Required: true},
		{Name: "access_token", Label: "Access Token", Type: FieldTypePassword, Value: cfg.AccessToken, Placeholder: "Optional access token"},
		{Name: "reconnect_interval", Label: "Reconnect Interval", Type: FieldTypeNumber, Value: strconv.Itoa(cfg.ReconnectInterval), Placeholder: "5"},
		{Name: "group_trigger_prefix", Label: "Group Trigger Prefix", Type: FieldTypeArray, Value: StringArrayToString(cfg.GroupTriggerPrefix), Placeholder: "!, /cmd"},
		{Name: "allow_from", Label: "Allowed Users", Type: FieldTypeArray, Value: StringArrayToString(cfg.AllowFrom), Placeholder: "user1, user2"},
	}
}

// InitMaixCamForm initializes the MaixCam configuration form
func (m *Model) InitMaixCamForm() {
	cfg := m.Config.Channels.MaixCam
	m.FormFields = []FormField{
		{Name: "enabled", Label: "Enabled", Type: FieldTypeBool, BoolValue: cfg.Enabled},
		{Name: "host", Label: "Host", Type: FieldTypeText, Value: cfg.Host, Placeholder: "192.168.1.100", Required: true},
		{Name: "port", Label: "Port", Type: FieldTypeNumber, Value: strconv.Itoa(cfg.Port), Placeholder: "8080"},
		{Name: "allow_from", Label: "Allowed Users", Type: FieldTypeArray, Value: StringArrayToString(cfg.AllowFrom), Placeholder: "user1, user2"},
	}
}

// InitWhatsAppForm initializes the WhatsApp configuration form
func (m *Model) InitWhatsAppForm() {
	cfg := m.Config.Channels.WhatsApp
	m.FormFields = []FormField{
		{Name: "enabled", Label: "Enabled", Type: FieldTypeBool, BoolValue: cfg.Enabled},
		{Name: "bridge_url", Label: "Bridge URL", Type: FieldTypeText, Value: cfg.BridgeURL, Placeholder: "http://localhost:3000/api", Required: true},
		{Name: "allow_from", Label: "Allowed Users", Type: FieldTypeArray, Value: StringArrayToString(cfg.AllowFrom), Placeholder: "user1, user2"},
	}
}

// SaveFormToConfig saves the form values to the config and enables the channel
func (m *Model) SaveFormToConfig() {
	switch m.SelectedChannel {
	case "telegram":
		m.SaveTelegramConfig()
		m.Config.Channels.Telegram.Enabled = true
	case "discord":
		m.SaveDiscordConfig()
		m.Config.Channels.Discord.Enabled = true
	case "slack":
		m.SaveSlackConfig()
		m.Config.Channels.Slack.Enabled = true
	case "qq":
		m.SaveQQConfig()
		m.Config.Channels.QQ.Enabled = true
	case "dingtalk":
		m.SaveDingTalkConfig()
		m.Config.Channels.DingTalk.Enabled = true
	case "wecom":
		m.SaveWeComConfig()
		m.Config.Channels.WeCom.Enabled = true
	case "wecom_app":
		m.SaveWeComAppConfig()
		m.Config.Channels.WeComApp.Enabled = true
	case "feishu":
		m.SaveFeishuConfig()
		m.Config.Channels.Feishu.Enabled = true
	case "line":
		m.SaveLineConfig()
		m.Config.Channels.LINE.Enabled = true
	case "onebot":
		m.SaveOneBotConfig()
		m.Config.Channels.OneBot.Enabled = true
	case "maixcam":
		m.SaveMaixCamConfig()
		m.Config.Channels.MaixCam.Enabled = true
	case "whatsapp":
		m.SaveWhatsAppConfig()
		m.Config.Channels.WhatsApp.Enabled = true
	}
}

// SaveTelegramConfig saves Telegram config from form
func (m *Model) SaveTelegramConfig() {
	cfg := &m.Config.Channels.Telegram
	cfg.Enabled = m.GetFieldValueBool("enabled")
	cfg.Token = m.GetFieldValue("token")
	cfg.Proxy = m.GetFieldValue("proxy")
	cfg.AllowFrom = m.GetFieldValueArray("allow_from")
}

// SaveDiscordConfig saves Discord config from form
func (m *Model) SaveDiscordConfig() {
	cfg := &m.Config.Channels.Discord
	cfg.Enabled = m.GetFieldValueBool("enabled")
	cfg.Token = m.GetFieldValue("token")
	cfg.AllowFrom = m.GetFieldValueArray("allow_from")
	cfg.MentionOnly = m.GetFieldValueBool("mention_only")
}

// SaveSlackConfig saves Slack config from form
func (m *Model) SaveSlackConfig() {
	cfg := &m.Config.Channels.Slack
	cfg.Enabled = m.GetFieldValueBool("enabled")
	cfg.BotToken = m.GetFieldValue("bot_token")
	cfg.AppToken = m.GetFieldValue("app_token")
	cfg.AllowFrom = m.GetFieldValueArray("allow_from")
}

// SaveQQConfig saves QQ config from form
func (m *Model) SaveQQConfig() {
	cfg := &m.Config.Channels.QQ
	cfg.Enabled = m.GetFieldValueBool("enabled")
	cfg.AppID = m.GetFieldValue("app_id")
	cfg.AppSecret = m.GetFieldValue("app_secret")
	cfg.AllowFrom = m.GetFieldValueArray("allow_from")
}

// SaveDingTalkConfig saves DingTalk config from form
func (m *Model) SaveDingTalkConfig() {
	cfg := &m.Config.Channels.DingTalk
	cfg.Enabled = m.GetFieldValueBool("enabled")
	cfg.ClientID = m.GetFieldValue("client_id")
	cfg.ClientSecret = m.GetFieldValue("client_secret")
	cfg.AllowFrom = m.GetFieldValueArray("allow_from")
}

// SaveWeComConfig saves WeCom config from form
func (m *Model) SaveWeComConfig() {
	cfg := &m.Config.Channels.WeCom
	cfg.Enabled = m.GetFieldValueBool("enabled")
	cfg.Token = m.GetFieldValue("token")
	cfg.EncodingAESKey = m.GetFieldValue("encoding_aes_key")
	cfg.WebhookURL = m.GetFieldValue("webhook_url")
	cfg.WebhookHost = m.GetFieldValue("webhook_host")
	cfg.WebhookPort = m.GetFieldValueInt("webhook_port")
	cfg.WebhookPath = m.GetFieldValue("webhook_path")
	cfg.AllowFrom = m.GetFieldValueArray("allow_from")
	cfg.ReplyTimeout = m.GetFieldValueInt("reply_timeout")
}

// SaveWeComAppConfig saves WeCom App config from form
func (m *Model) SaveWeComAppConfig() {
	cfg := &m.Config.Channels.WeComApp
	cfg.Enabled = m.GetFieldValueBool("enabled")
	cfg.CorpID = m.GetFieldValue("corp_id")
	cfg.CorpSecret = m.GetFieldValue("corp_secret")
	cfg.AgentID = m.GetFieldValueInt64("agent_id")
	cfg.Token = m.GetFieldValue("token")
	cfg.EncodingAESKey = m.GetFieldValue("encoding_aes_key")
	cfg.WebhookHost = m.GetFieldValue("webhook_host")
	cfg.WebhookPort = m.GetFieldValueInt("webhook_port")
	cfg.WebhookPath = m.GetFieldValue("webhook_path")
	cfg.AllowFrom = m.GetFieldValueArray("allow_from")
	cfg.ReplyTimeout = m.GetFieldValueInt("reply_timeout")
}

// SaveFeishuConfig saves Feishu config from form
func (m *Model) SaveFeishuConfig() {
	cfg := &m.Config.Channels.Feishu
	cfg.Enabled = m.GetFieldValueBool("enabled")
	cfg.AppID = m.GetFieldValue("app_id")
	cfg.AppSecret = m.GetFieldValue("app_secret")
	cfg.EncryptKey = m.GetFieldValue("encrypt_key")
	cfg.VerificationToken = m.GetFieldValue("verification_token")
	cfg.AllowFrom = m.GetFieldValueArray("allow_from")
}

// SaveLineConfig saves LINE config from form
func (m *Model) SaveLineConfig() {
	cfg := &m.Config.Channels.LINE
	cfg.Enabled = m.GetFieldValueBool("enabled")
	cfg.ChannelSecret = m.GetFieldValue("channel_secret")
	cfg.ChannelAccessToken = m.GetFieldValue("channel_access_token")
	cfg.WebhookHost = m.GetFieldValue("webhook_host")
	cfg.WebhookPort = m.GetFieldValueInt("webhook_port")
	cfg.WebhookPath = m.GetFieldValue("webhook_path")
	cfg.AllowFrom = m.GetFieldValueArray("allow_from")
}

// SaveOneBotConfig saves OneBot config from form
func (m *Model) SaveOneBotConfig() {
	cfg := &m.Config.Channels.OneBot
	cfg.Enabled = m.GetFieldValueBool("enabled")
	cfg.WSUrl = m.GetFieldValue("ws_url")
	cfg.AccessToken = m.GetFieldValue("access_token")
	cfg.ReconnectInterval = m.GetFieldValueInt("reconnect_interval")
	cfg.GroupTriggerPrefix = m.GetFieldValueArray("group_trigger_prefix")
	cfg.AllowFrom = m.GetFieldValueArray("allow_from")
}

// SaveMaixCamConfig saves MaixCam config from form
func (m *Model) SaveMaixCamConfig() {
	cfg := &m.Config.Channels.MaixCam
	cfg.Enabled = m.GetFieldValueBool("enabled")
	cfg.Host = m.GetFieldValue("host")
	cfg.Port = m.GetFieldValueInt("port")
	cfg.AllowFrom = m.GetFieldValueArray("allow_from")
}

// SaveWhatsAppConfig saves WhatsApp config from form
func (m *Model) SaveWhatsAppConfig() {
	cfg := &m.Config.Channels.WhatsApp
	cfg.Enabled = m.GetFieldValueBool("enabled")
	cfg.BridgeURL = m.GetFieldValue("bridge_url")
	cfg.AllowFrom = m.GetFieldValueArray("allow_from")
}
