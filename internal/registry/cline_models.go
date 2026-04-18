package registry

// GetClineModels returns the Cline model definitions.
// Cline uses OpenRouter-compatible routing.
func GetClineModels() []*ModelInfo {
	now := int64(1732752000) // 2024-11-27

	return []*ModelInfo{
		{
			ID:                  "kwaipilot/kat-coder-pro",
			Object:              "model",
			Created:             now,
			OwnedBy:             "cline",
			Type:                "cline",
			DisplayName:         "Kat Coder Pro (Free)",
			Description:         "KwaiPilot Kat Coder Pro via Cline (Free)",
			ContextLength:       128000,
			MaxCompletionTokens: 32768,
		},
		{
			ID:                  "z-ai/glm-5",
			Object:              "model",
			Created:             now,
			OwnedBy:             "cline",
			Type:                "cline",
			DisplayName:         "GLM-5 (Free)",
			Description:         "Z-AI GLM-5 via Cline (Free)",
			ContextLength:       128000,
			MaxCompletionTokens: 32768,
		},
		{
			ID:                  "z-ai/glm-5.1",
			Object:              "model",
			Created:             now,
			OwnedBy:             "cline",
			Type:                "cline",
			DisplayName:         "GLM-5.1 (Free)",
			Description:         "Z-AI GLM-5.1 via Cline (Free)",
			ContextLength:       128000,
			MaxCompletionTokens: 32768,
		},
	}
}
