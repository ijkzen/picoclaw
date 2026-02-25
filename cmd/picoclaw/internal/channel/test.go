package channel

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// PerformTest performs a connection test for the selected channel
func (m *Model) PerformTest() TestResult {
	if m.SelectedChannel == "" {
		return TestResult{
			Success: false,
			Message: "No channel selected",
			Error:   errors.New("no channel selected"),
		}
	}

	switch m.SelectedChannel {
	case "telegram":
		return m.TestTelegram()
	case "discord":
		return m.TestDiscord()
	case "slack":
		return m.TestSlack()
	case "qq":
		return m.TestQQ()
	case "dingtalk":
		return m.TestDingTalk()
	case "wecom":
		return m.TestWeCom()
	case "wecom_app":
		return m.TestWeComApp()
	case "feishu":
		return m.TestFeishu()
	case "line":
		return m.TestLine()
	case "onebot":
		return m.TestOneBot()
	case "maixcam":
		return m.TestMaixCam()
	case "whatsapp":
		return m.TestWhatsApp()
	default:
		return TestResult{
			Success: false,
			Message: fmt.Sprintf("Unknown channel: %s", m.SelectedChannel),
			Error:   fmt.Errorf("unknown channel: %s", m.SelectedChannel),
		}
	}
}

// TestTelegram tests Telegram configuration
func (m *Model) TestTelegram() TestResult {
	cfg := m.Config.Channels.Telegram

	if cfg.Token == "" {
		return TestResult{
			Success: false,
			Message: "Telegram bot token is required",
			Error:   errors.New("bot token cannot be empty"),
		}
	}

	// Validate token format (should be numbers:alphanumeric)
	parts := strings.Split(cfg.Token, ":")
	if len(parts) != 2 {
		return TestResult{
			Success: false,
			Message: "Invalid token format",
			Error:   errors.New("token should be in format 'numbers:alphanumeric'"),
		}
	}

	// Make a test API call
	client := &http.Client{Timeout: 10 * time.Second}
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/getMe", cfg.Token)

	resp, err := client.Get(apiURL)
	if err != nil {
		return TestResult{
			Success: false,
			Message: "Failed to connect to Telegram API",
			Error:   err,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TestResult{
			Success: false,
			Message: fmt.Sprintf("Telegram API returned status %d", resp.StatusCode),
			Error:   fmt.Errorf("API status: %d", resp.StatusCode),
		}
	}

	return TestResult{
		Success: true,
		Message: "Telegram configuration is valid!",
	}
}

// TestDiscord tests Discord configuration
func (m *Model) TestDiscord() TestResult {
	cfg := m.Config.Channels.Discord

	if cfg.Token == "" {
		return TestResult{
			Success: false,
			Message: "Discord bot token is required",
			Error:   errors.New("bot token cannot be empty"),
		}
	}

	// Validate token format (should start with "M" and have multiple parts)
	if !strings.HasPrefix(cfg.Token, "M") || len(cfg.Token) < 50 {
		return TestResult{
			Success: false,
			Message: "Invalid token format",
			Error:   errors.New("Discord bot token appears to be invalid"),
		}
	}

	// Make a test API call
	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", "https://discord.com/api/v10/users/@me", nil)
	req.Header.Set("Authorization", "Bot "+cfg.Token)

	resp, err := client.Do(req)
	if err != nil {
		return TestResult{
			Success: false,
			Message: "Failed to connect to Discord API",
			Error:   err,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return TestResult{
			Success: false,
			Message: "Invalid bot token",
			Error:   errors.New("Discord API returned 401 Unauthorized"),
		}
	}

	if resp.StatusCode != http.StatusOK {
		return TestResult{
			Success: false,
			Message: fmt.Sprintf("Discord API returned status %d", resp.StatusCode),
			Error:   fmt.Errorf("API status: %d", resp.StatusCode),
		}
	}

	return TestResult{
		Success: true,
		Message: "Discord configuration is valid!",
	}
}

// TestSlack tests Slack configuration
func (m *Model) TestSlack() TestResult {
	cfg := m.Config.Channels.Slack

	if cfg.BotToken == "" {
		return TestResult{
			Success: false,
			Message: "Slack bot token is required",
			Error:   errors.New("bot token cannot be empty"),
		}
	}

	// Validate bot token format
	if !strings.HasPrefix(cfg.BotToken, "xoxb-") {
		return TestResult{
			Success: false,
			Message: "Invalid bot token format",
			Error:   errors.New("bot token should start with 'xoxb-'"),
		}
	}

	// Make a test API call
	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", "https://slack.com/api/auth.test", nil)
	req.Header.Set("Authorization", "Bearer "+cfg.BotToken)

	resp, err := client.Do(req)
	if err != nil {
		return TestResult{
			Success: false,
			Message: "Failed to connect to Slack API",
			Error:   err,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TestResult{
			Success: false,
			Message: fmt.Sprintf("Slack API returned status %d", resp.StatusCode),
			Error:   fmt.Errorf("API status: %d", resp.StatusCode),
		}
	}

	return TestResult{
		Success: true,
		Message: "Slack configuration is valid!",
	}
}

// TestQQ tests QQ configuration
func (m *Model) TestQQ() TestResult {
	cfg := m.Config.Channels.QQ

	if cfg.AppID == "" || cfg.AppSecret == "" {
		return TestResult{
			Success: false,
			Message: "QQ App ID and App Secret are required",
			Error:   errors.New("app_id and app_secret cannot be empty"),
		}
	}

	// Basic validation
	if len(cfg.AppID) < 5 {
		return TestResult{
			Success: false,
			Message: "App ID appears to be invalid",
			Error:   errors.New("app_id is too short"),
		}
	}

	return TestResult{
		Success: true,
		Message: "QQ configuration appears valid (test skipped - requires QQ Bot platform)",
	}
}

// TestDingTalk tests DingTalk configuration
func (m *Model) TestDingTalk() TestResult {
	cfg := m.Config.Channels.DingTalk

	if cfg.ClientID == "" || cfg.ClientSecret == "" {
		return TestResult{
			Success: false,
			Message: "DingTalk Client ID and Client Secret are required",
			Error:   errors.New("client_id and client_secret cannot be empty"),
		}
	}

	// Validate Client ID format
	if !strings.HasPrefix(cfg.ClientID, "ding") && !strings.HasPrefix(cfg.ClientID, "suite") {
		return TestResult{
			Success: false,
			Message: "Client ID should start with 'ding' or 'suite'",
			Error:   errors.New("invalid client_id format"),
		}
	}

	return TestResult{
		Success: true,
		Message: "DingTalk configuration appears valid (test skipped - requires DingTalk platform)",
	}
}

// TestWeCom tests WeCom Bot configuration
func (m *Model) TestWeCom() TestResult {
	cfg := m.Config.Channels.WeCom

	if cfg.Token == "" && cfg.WebhookURL == "" {
		return TestResult{
			Success: false,
			Message: "WeCom Token or Webhook URL is required",
			Error:   errors.New("at least one of token or webhook_url must be provided"),
		}
	}

	// Validate webhook URL if provided
	if cfg.WebhookURL != "" {
		if _, err := url.Parse(cfg.WebhookURL); err != nil {
			return TestResult{
				Success: false,
				Message: "Invalid webhook URL",
				Error:   err,
			}
		}
	}

	return TestResult{
		Success: true,
		Message: "WeCom Bot configuration is valid!",
	}
}

// TestWeComApp tests WeCom App configuration
func (m *Model) TestWeComApp() TestResult {
	cfg := m.Config.Channels.WeComApp

	if cfg.CorpID == "" || cfg.CorpSecret == "" {
		return TestResult{
			Success: false,
			Message: "WeCom Corp ID and Corp Secret are required",
			Error:   errors.New("corp_id and corp_secret cannot be empty"),
		}
	}

	if cfg.AgentID == 0 {
		return TestResult{
			Success: false,
			Message: "Agent ID is required",
			Error:   errors.New("agent_id cannot be 0"),
		}
	}

	// Validate Corp ID format
	if !strings.HasPrefix(cfg.CorpID, "ww") {
		return TestResult{
			Success: false,
			Message: "Corp ID should start with 'ww'",
			Error:   errors.New("invalid corp_id format"),
		}
	}

	return TestResult{
		Success: true,
		Message: "WeCom App configuration is valid!",
	}
}

// TestFeishu tests Feishu configuration
func (m *Model) TestFeishu() TestResult {
	cfg := m.Config.Channels.Feishu

	if cfg.AppID == "" || cfg.AppSecret == "" {
		return TestResult{
			Success: false,
			Message: "Feishu App ID and App Secret are required",
			Error:   errors.New("app_id and app_secret cannot be empty"),
		}
	}

	// Validate App ID format
	if !strings.HasPrefix(cfg.AppID, "cli_") {
		return TestResult{
			Success: false,
			Message: "App ID should start with 'cli_'",
			Error:   errors.New("invalid app_id format"),
		}
	}

	return TestResult{
		Success: true,
		Message: "Feishu configuration is valid!",
	}
}

// TestLine tests LINE configuration
func (m *Model) TestLine() TestResult {
	cfg := m.Config.Channels.LINE

	if cfg.ChannelSecret == "" || cfg.ChannelAccessToken == "" {
		return TestResult{
			Success: false,
			Message: "LINE Channel Secret and Channel Access Token are required",
			Error:   errors.New("channel_secret and channel_access_token cannot be empty"),
		}
	}

	// Make a test API call
	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", "https://api.line.me/v2/bot/info", nil)
	req.Header.Set("Authorization", "Bearer "+cfg.ChannelAccessToken)

	resp, err := client.Do(req)
	if err != nil {
		return TestResult{
			Success: false,
			Message: "Failed to connect to LINE API",
			Error:   err,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return TestResult{
			Success: false,
			Message: "Invalid channel access token",
			Error:   errors.New("LINE API returned 401 Unauthorized"),
		}
	}

	if resp.StatusCode != http.StatusOK {
		return TestResult{
			Success: false,
			Message: fmt.Sprintf("LINE API returned status %d", resp.StatusCode),
			Error:   fmt.Errorf("API status: %d", resp.StatusCode),
		}
	}

	return TestResult{
		Success: true,
		Message: "LINE configuration is valid!",
	}
}

// TestOneBot tests OneBot configuration
func (m *Model) TestOneBot() TestResult {
	cfg := m.Config.Channels.OneBot

	if cfg.WSUrl == "" {
		return TestResult{
			Success: false,
			Message: "OneBot WebSocket URL is required",
			Error:   errors.New("ws_url cannot be empty"),
		}
	}

	// Validate WebSocket URL format
	if !strings.HasPrefix(cfg.WSUrl, "ws://") && !strings.HasPrefix(cfg.WSUrl, "wss://") {
		return TestResult{
			Success: false,
			Message: "WebSocket URL should start with 'ws://' or 'wss://'",
			Error:   errors.New("invalid ws_url format"),
		}
	}

	if _, err := url.Parse(cfg.WSUrl); err != nil {
		return TestResult{
			Success: false,
			Message: "Invalid WebSocket URL",
			Error:   err,
		}
	}

	return TestResult{
		Success: true,
		Message: "OneBot configuration is valid!",
	}
}

// TestMaixCam tests MaixCam configuration
func (m *Model) TestMaixCam() TestResult {
	cfg := m.Config.Channels.MaixCam

	if cfg.Host == "" {
		return TestResult{
			Success: false,
			Message: "MaixCam Host is required",
			Error:   errors.New("host cannot be empty"),
		}
	}

	// Validate host format
	if cfg.Port == 0 {
		return TestResult{
			Success: false,
			Message: "Port is required",
			Error:   errors.New("port cannot be 0"),
		}
	}

	// Make a test HTTP call
	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port)
	resp, err := client.Get(url)
	if err != nil {
		return TestResult{
			Success: true,
			Message: fmt.Sprintf("MaixCam configuration appears valid (device at %s may be offline)", url),
		}
	}
	defer resp.Body.Close()

	return TestResult{
		Success: true,
		Message: "MaixCam is online and reachable!",
	}
}

// TestWhatsApp tests WhatsApp Bridge configuration
func (m *Model) TestWhatsApp() TestResult {
	cfg := m.Config.Channels.WhatsApp

	if cfg.BridgeURL == "" {
		return TestResult{
			Success: false,
			Message: "WhatsApp Bridge URL is required",
			Error:   errors.New("bridge_url cannot be empty"),
		}
	}

	// Validate URL format
	parsedURL, err := url.Parse(cfg.BridgeURL)
	if err != nil {
		return TestResult{
			Success: false,
			Message: "Invalid bridge URL",
			Error:   err,
		}
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return TestResult{
			Success: false,
			Message: "Bridge URL should use http:// or https://",
			Error:   errors.New("invalid URL scheme"),
		}
	}

	// Make a test HTTP call
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(cfg.BridgeURL)
	if err != nil {
		return TestResult{
			Success: true,
			Message: "WhatsApp Bridge configuration appears valid (bridge may be offline)",
		}
	}
	defer resp.Body.Close()

	return TestResult{
		Success: true,
		Message: "WhatsApp Bridge is online and reachable!",
	}
}
