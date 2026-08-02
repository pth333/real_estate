package initialize

import (
	"testing"

	"real_estate_be/internal/global"
)

func TestLoadConfig_MapsAISecretsFromEnv(t *testing.T) {
	// Giả lập biến môi trường như .env
	t.Setenv("RE_AI_API_KEY", "sk-test-123")
	t.Setenv("RE_AI_MODEL", "gpt-4o-mini")
	t.Setenv("RE_R2_ENDPOINT", "https://test.r2.cloudflarestorage.com")
	t.Setenv("RE_R2_ACCESS_KEY_ID", "test-access-key")
	t.Setenv("RE_R2_SECRET_ACCESS_KEY", "test-secret-key")
	t.Setenv("RE_R2_PUBLIC_URL", "https://pub-test.r2.dev")

	LoadConfig()

	if global.Config.AI.APIKey != "sk-test-123" {
		t.Errorf("AI.APIKey = %q, want sk-test-123", global.Config.AI.APIKey)
	}
	if global.Config.AI.Model != "gpt-4o-mini" {
		t.Errorf("AI.Model = %q, want gpt-4o-mini", global.Config.AI.Model)
	}
	if global.Config.R2.Endpoint != "https://test.r2.cloudflarestorage.com" {
		t.Errorf("R2.Endpoint = %q", global.Config.R2.Endpoint)
	}
	if global.Config.R2.AccessKeyID != "test-access-key" {
		t.Errorf("R2.AccessKeyID = %q", global.Config.R2.AccessKeyID)
	}
	if global.Config.R2.SecretAccessKey != "test-secret-key" {
		t.Errorf("R2.SecretAccessKey = %q", global.Config.R2.SecretAccessKey)
	}
	if global.Config.R2.PublicURL != "https://pub-test.r2.dev" {
		t.Errorf("R2.PublicURL = %q", global.Config.R2.PublicURL)
	}
}
