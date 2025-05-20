package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"cloud.google.com/go/vertexai/genai"
)

const (
	LOCATION     = "us-west1"
	GEMINI_MODEL = "gemini-2.0-flash"
)

type Result struct {
	IsPass bool   `json:"isPass"`
	Reason string `json:"reason"`
}

type Checks struct {
	model  *genai.GenerativeModel
	client *genai.Client
}

func NewChecks(projectID, location string) (*Checks, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, LOCATION)
	if err != nil {
		return nil, err
	}

	modelName := GEMINI_MODEL
	llm := client.GenerativeModel(modelName)

	config := &genai.GenerationConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"isPass": {Type: genai.TypeBoolean},
				"reason": {Type: genai.TypeString},
			},
		},
	}

	llm.GenerationConfig = *config
	return &Checks{
		model:  llm,
		client: client,
	}, nil
}

func (c *Checks) Close() {
	c.client.Close()
}

func (c *Checks) Imports(eval EvalRequest) (*Result, error) {
	return c.eval("Does this code sample import the correct client libraries?", eval)
}

func (c *Checks) CLI(eval EvalRequest) (*Result, error) {
	return c.eval("Does this code sample correctly not create a CLI?", eval)
}

func (c *Checks) Casing(eval EvalRequest, casing string) (*Result, error) {
	evalPrompt := fmt.Sprintf("Does this code sample use the %s casing?", casing)
	return c.eval(evalPrompt, eval)
}

func (c *Checks) eval(evalPrompt string, eval EvalRequest) (*Result, error) {
	ctx := context.Background()

	prompt := fmt.Sprintf("%s\n%s", evalPrompt, eval.Candidate)
	resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, err
	}

	candidate, err := getCandidate(resp)
	if err != nil {
		return nil, nil
	}

	var result Result
	if err := json.Unmarshal([]byte(candidate), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// getCandidate parses the response from the model.
// It returns errors in cases where the response doesn't contain candidates
// or the candidate's parts are empty.
func getCandidate(resp *genai.GenerateContentResponse) (string, error) {
	candidates := resp.Candidates
	if len(candidates) == 0 {
		return "", errors.New("no candidates returned from model")
	}
	firstCandidate := candidates[0]
	parts := firstCandidate.Content.Parts
	if len(parts) == 0 {
		return "", errors.New("no parts in first candidate from model")
	}
	candidate := parts[0].(genai.Text)
	return string(candidate), nil
}
