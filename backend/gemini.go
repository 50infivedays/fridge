package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// GeminiRequest represents the request structure for Gemini API
type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

// GeminiContent represents the content part of a Gemini request
type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

// GeminiPart represents a part within the content
type GeminiPart struct {
	Text string `json:"text"`
}

// GeminiResponse represents the response from Gemini API
type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// processWithGemini sends the description to Gemini API and processes the response
func processWithGemini(description string, currentTime string) (*FridgeResponse, error) {
	// Get cached configuration
	config, err := GetConfig()
	if err != nil {
		return nil, fmt.Errorf("error getting configuration: %w", err)
	}

	// Check if API key is available
	if config.GeminiAPIKey == "" {
		return nil, fmt.Errorf("Gemini API key not found in config.yaml or environment variables")
	}

	// Construct the prompt for Gemini
	prompt := constructPrompt(description, currentTime)

	// Create the request body
	reqBody := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{
						Text: prompt,
					},
				},
			},
		},
	}

	// Convert request to JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", config.GeminiModel, config.GeminiAPIKey)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request to Gemini API: %w", err)
	}
	defer resp.Body.Close()

	// Read the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	// Check for successful status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Gemini API returned error: %s", string(respBody))
	}

	// Parse the response
	var geminiResp GeminiResponse
	if err := json.Unmarshal(respBody, &geminiResp); err != nil {
		return nil, fmt.Errorf("error parsing Gemini response: %w", err)
	}

	// Check if we have any candidates
	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no response from Gemini API")
	}

	// Extract the JSON string from the response
	jsonString := geminiResp.Candidates[0].Content.Parts[0].Text

	// Log the raw response for debugging
	log.Printf("Raw Gemini response text: %s", jsonString)

	// Clean up the JSON string (remove markdown code blocks if present)
	jsonString = cleanJSONString(jsonString)

	// Parse the JSON string into our response structure
	var fridgeResponse FridgeResponse
	if err := json.Unmarshal([]byte(jsonString), &fridgeResponse); err != nil {
		log.Printf("JSON parsing error: %v for string: %s", err, jsonString)

		// Try to create a default response if parsing fails
		return createDefaultResponse(jsonString), nil
	}

	return &fridgeResponse, nil
}

// constructPrompt creates the prompt for Gemini API
func constructPrompt(description string, clientTime string) string {
	// Use client-provided time or fall back to server time
	currentTime := clientTime
	if currentTime == "" {
		currentTime = time.Now().Format("2006-01-02 15:04:05")
	}
	
	// Parse the current time to calculate relative dates
	currentTimeObj, err := time.Parse("2006-01-02 15:04:05", currentTime)
	if err != nil {
		// If parsing fails, use current server time
		currentTimeObj = time.Now()
		currentTime = currentTimeObj.Format("2006-01-02 15:04:05")
	}
	
	// Calculate 5 days later for example
	fiveDaysLater := currentTimeObj.AddDate(0, 0, 5).Format("2006-01-02 15:04:05")
	return `你是一个专业的冰箱物品管理助手。请帮我解析以下关于放入冰箱的物品的描述，并将其转换为结构化的JSON格式。

描述中可能包含物品名称、数量、单位和过期时间。

重要说明：
1. 当前时间是：` + currentTime + `
2. 如果描述中提到相对时间（如"5天后"、"一周后"、"明天"等），请将其转换为绝对时间
3. 例如，如果当前时间是` + currentTime + `，"5天后"应该计算为` + fiveDaysLater + `
4. 如果没有明确提到过期时间，默认为7天后
5. 如果描述中没有明确提到数量或单位，请使用合理的默认值

描述: "` + description + `"

请将结果格式化为以下JSON结构:
{
  "items": [
    {
      "item": "物品名称",
      "quantity": 数量(整数),
      "unit": "单位(如个、斤、盒等)",
      "expireDate": "过期时间(格式为YYYY-MM-DD HH:MM:SS)"
    }
  ]
}

只返回JSON格式的结果，不要包含任何其他解释或文本。确保JSON格式正确且可解析。如果描述中包含多个物品，请在items数组中包含多个对象。`

}

// createDefaultResponse attempts to create a default response when JSON parsing fails
func createDefaultResponse(rawText string) *FridgeResponse {
	log.Printf("Creating default response from raw text")

	// Create a default response with a single item
	return &FridgeResponse{
		Items: []FridgeItem{
			{
				Item:       "未知物品",
				Quantity:   1,
				Unit:       "个",
				ExpireDate: time.Now().AddDate(0, 0, 7).Format("2006-01-02 15:04:05"),
			},
		},
	}
}

// cleanJSONString removes markdown code blocks and other non-JSON content
func cleanJSONString(input string) string {
	// Remove markdown code blocks if present
	jsonString := input

	// Find the first { character (start of JSON)
	jsonStartIndex := strings.Index(jsonString, "{")
	if jsonStartIndex > 0 {
		jsonString = jsonString[jsonStartIndex:]
	}

	// Find the last } character (end of JSON)
	jsonEndIndex := strings.LastIndex(jsonString, "}")
	if jsonEndIndex >= 0 && jsonEndIndex < len(jsonString)-1 {
		jsonString = jsonString[:jsonEndIndex+1]
	}

	// Remove ```json prefix if present
	if len(jsonString) > 7 && jsonString[:7] == "```json" {
		jsonString = jsonString[7:]
	} else if len(jsonString) > 3 && jsonString[:3] == "```" {
		jsonString = jsonString[3:]
	}

	// Remove trailing ``` if present
	if len(jsonString) > 3 && jsonString[len(jsonString)-3:] == "```" {
		jsonString = jsonString[:len(jsonString)-3]
	}

	// Trim whitespace
	jsonString = string(bytes.TrimSpace([]byte(jsonString)))

	// Debug logging
	log.Printf("Cleaned JSON string: %s", jsonString)

	return jsonString
}
