package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vietbui/chat-quality-agent/ai"
	"github.com/vietbui/chat-quality-agent/api/middleware"
	"github.com/vietbui/chat-quality-agent/config"
	"github.com/vietbui/chat-quality-agent/db"
	"github.com/vietbui/chat-quality-agent/db/models"
	"github.com/vietbui/chat-quality-agent/pkg"
	"golang.org/x/crypto/bcrypt"
)

// GetSettings returns all non-secret settings for the tenant
func GetSettings(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	var settings []models.AppSetting
	db.DB.Where("tenant_id = ?", tenantID).Find(&settings)

	result := make(map[string]string)
	for _, s := range settings {
		if s.ValuePlain != "" {
			result[s.SettingKey] = s.ValuePlain
		} else if len(s.ValueEncrypted) > 0 {
			// Return masked value for encrypted settings
			result[s.SettingKey] = "••••••••"
		}
	}

	// Also get tenant info
	var tenant models.Tenant
	db.DB.First(&tenant, "id = ?", tenantID)

	c.JSON(http.StatusOK, gin.H{
		"settings": result,
		"tenant": gin.H{
			"name":     tenant.Name,
			"timezone": getSettingValue(settings, "timezone", "Asia/Ho_Chi_Minh"),
			"language": getSettingValue(settings, "language", "vi"),
		},
	})
}

func getSettingValue(settings []models.AppSetting, key, defaultVal string) string {
	for _, s := range settings {
		if s.SettingKey == key && s.ValuePlain != "" {
			return s.ValuePlain
		}
	}
	return defaultVal
}

// SaveAISettings saves AI provider and API key
func SaveAISettings(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	var req struct {
		Provider  string `json:"provider" binding:"required,oneof=claude gemini openai"`
		APIKey    string `json:"api_key"`
		Model     string `json:"model"`
		BaseURL   string `json:"base_url"`
		BatchMode string `json:"batch_mode"`
		BatchSize string `json:"batch_size"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "details": err.Error()})
		return
	}

	cfg, _ := config.Load()

	// Save provider (plain)
	upsertSetting(tenantID, "ai_provider", req.Provider, nil)

	// Save model (plain)
	if req.Model != "" {
		upsertSetting(tenantID, "ai_model", req.Model, nil)
	}

	// Save base URL (plain, optional — empty string clears it)
	if req.BaseURL != "" {
		upsertSetting(tenantID, "ai_base_url", req.BaseURL, nil)
	} else {
		// Explicitly clear: delete the setting if empty
		db.DB.Where("tenant_id = ? AND setting_key = ?", tenantID, "ai_base_url").Delete(&models.AppSetting{})
	}

	// Save API key (encrypted), or keep the existing key when the UI sends the masked value.
	if req.APIKey != "" && req.APIKey != "••••••••" {
		encrypted, err := pkg.Encrypt([]byte(req.APIKey), cfg.EncryptionKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "encryption_failed"})
			return
		}
		upsertSetting(tenantID, "ai_api_key", "", encrypted)
	} else {
		var existing models.AppSetting
		if err := db.DB.Where("tenant_id = ? AND setting_key = ?", tenantID, "ai_api_key").First(&existing).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no_api_key_configured"})
			return
		}
	}

	// Save batch settings
	if req.BatchMode != "" {
		upsertSetting(tenantID, "ai_batch_mode", req.BatchMode, nil)
	}
	if req.BatchSize != "" {
		upsertSetting(tenantID, "ai_batch_size", req.BatchSize, nil)
	}

	c.JSON(http.StatusOK, gin.H{"message": "saved"})
}

// SaveAnalysisSettings saves batch mode and batch size settings
func SaveAnalysisSettings(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	var req struct {
		BatchMode string `json:"batch_mode" binding:"required"`
		BatchSize string `json:"batch_size"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "details": err.Error()})
		return
	}

	upsertSetting(tenantID, "ai_batch_mode", req.BatchMode, nil)
	if req.BatchSize != "" {
		upsertSetting(tenantID, "ai_batch_size", req.BatchSize, nil)
	}

	c.JSON(http.StatusOK, gin.H{"message": "saved"})
}

// TestAIKey tests the AI API key by making a simple request
func TestAIKey(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	cfg, _ := config.Load()

	// Get the encrypted API key
	var setting models.AppSetting
	if err := db.DB.Where("tenant_id = ? AND setting_key = ?", tenantID, "ai_api_key").First(&setting).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no_api_key_configured"})
		return
	}

	apiKey, err := pkg.Decrypt(setting.ValueEncrypted, cfg.EncryptionKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "decrypt_failed"})
		return
	}

	// Get provider
	var providerSetting models.AppSetting
	provider := "claude"
	if err := db.DB.Where("tenant_id = ? AND setting_key = ?", tenantID, "ai_provider").First(&providerSetting).Error; err == nil {
		provider = providerSetting.ValuePlain
	}

	var modelSetting models.AppSetting
	model := ""
	if err := db.DB.Where("tenant_id = ? AND setting_key = ?", tenantID, "ai_model").First(&modelSetting).Error; err == nil {
		model = modelSetting.ValuePlain
	}

	var baseURLSetting models.AppSetting
	baseURL := ""
	if err := db.DB.Where("tenant_id = ? AND setting_key = ?", tenantID, "ai_base_url").First(&baseURLSetting).Error; err == nil {
		baseURL = baseURLSetting.ValuePlain
	}

	var providerClient interface {
		AnalyzeChat(ctx context.Context, systemPrompt string, chatTranscript string) (ai.AIResponse, error)
	}
	switch provider {
	case "claude":
		providerClient = ai.NewClaudeProvider(string(apiKey), model, cfg.AIMaxTokens, baseURL)
	case "gemini":
		providerClient = ai.NewGeminiProvider(string(apiKey), model, baseURL)
	case "openai":
		providerClient = ai.NewOpenAIProvider(string(apiKey), model, cfg.AIMaxTokens, baseURL)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported_provider"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	resp, err := providerClient.AnalyzeChat(ctx, "Reply with exactly: OK", "Test connectivity.")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "provider": provider, "model": resp.Model, "message": "API key hoạt động"})
}

// ListAIModels fetches available models from an OpenAI-compatible provider.
func ListAIModels(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	cfg, _ := config.Load()

	var req struct {
		Provider string `json:"provider" binding:"required"`
		APIKey   string `json:"api_key"`
		BaseURL  string `json:"base_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "details": err.Error()})
		return
	}
	if req.Provider != "openai" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "model_discovery_not_supported"})
		return
	}

	apiKey := req.APIKey
	if apiKey == "" || apiKey == "••••••••" {
		var setting models.AppSetting
		if err := db.DB.Where("tenant_id = ? AND setting_key = ?", tenantID, "ai_api_key").First(&setting).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no_api_key_configured"})
			return
		}
		decrypted, err := pkg.Decrypt(setting.ValueEncrypted, cfg.EncryptionKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "decrypt_failed"})
			return
		}
		apiKey = string(decrypted)
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	models, err := ai.FetchOpenAIModels(ctx, apiKey, req.BaseURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"models": models})
}

// SaveGeneralSettings saves general tenant settings
func SaveGeneralSettings(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	var req struct {
		CompanyName  string  `json:"company_name"`
		Timezone     string  `json:"timezone"`
		Language     string  `json:"language"`
		ExchangeRate float64 `json:"exchange_rate_vnd"`
		AppURL       string  `json:"app_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "details": err.Error()})
		return
	}

	// Update tenant name
	if req.CompanyName != "" {
		db.DB.Model(&models.Tenant{}).Where("id = ?", tenantID).Updates(map[string]interface{}{
			"name":       req.CompanyName,
			"updated_at": time.Now(),
		})
	}

	// Save timezone and language as settings
	if req.Timezone != "" {
		upsertSetting(tenantID, "timezone", req.Timezone, nil)
	}
	if req.Language != "" {
		upsertSetting(tenantID, "language", req.Language, nil)
	}
	if req.ExchangeRate > 0 {
		upsertSetting(tenantID, "exchange_rate_vnd", fmt.Sprintf("%.0f", req.ExchangeRate), nil)
	}

	// Strip trailing slash from app URL
	appURL := strings.TrimRight(req.AppURL, "/")
	upsertSetting(tenantID, "app_url", appURL, nil)

	c.JSON(http.StatusOK, gin.H{"message": "saved"})
}

// ChangePassword changes the user's password
func ChangePassword(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "details": err.Error()})
		return
	}
	if err := validatePasswordComplexity(req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "weak_password", "message": err.Error()})
		return
	}

	var user models.User
	if err := db.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user_not_found"})
		return
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong_current_password"})
		return
	}

	// Hash new password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "hash_failed"})
		return
	}

	if err := db.DB.Model(&user).Updates(map[string]interface{}{
		"password_hash": string(hash),
		"updated_at":    time.Now(),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update_failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password_changed"})
}

// allowedSettingKeys is a whitelist of keys that can be set via the SaveSetting API.
// Sensitive keys like ai_api_key must be set through dedicated endpoints.
var allowedSettingKeys = map[string]bool{
	"onboarding_dismissed": true,
	"language":             true,
	"timezone":             true,
	"date_format":          true,
	"notification_enabled": true,
	"sync_interval":        true,
	"default_ai_provider":  true,
	"default_ai_model":     true,
}

// SaveSetting saves a single key-value setting
func SaveSetting(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	var req struct {
		Key   string `json:"key" binding:"required"`
		Value string `json:"value" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
		return
	}
	if !allowedSettingKeys[req.Key] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "setting_key_not_allowed"})
		return
	}
	upsertSetting(tenantID, req.Key, req.Value, nil)
	c.JSON(http.StatusOK, gin.H{"message": "saved"})
}

func upsertSetting(tenantID, key, plainValue string, encryptedValue []byte) {
	var existing models.AppSetting
	result := db.DB.Where("tenant_id = ? AND setting_key = ?", tenantID, key).First(&existing)

	if result.Error == nil {
		// Update
		updates := map[string]interface{}{"updated_at": time.Now()}
		if plainValue != "" {
			updates["value_plain"] = plainValue
			updates["value_encrypted"] = nil
		}
		if encryptedValue != nil {
			updates["value_encrypted"] = encryptedValue
			updates["value_plain"] = ""
		}
		db.DB.Model(&existing).Updates(updates)
	} else {
		// Create
		setting := models.AppSetting{
			ID:             pkg.NewUUID(),
			TenantID:       tenantID,
			SettingKey:     key,
			ValuePlain:     plainValue,
			ValueEncrypted: encryptedValue,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		db.DB.Create(&setting)
	}
}
