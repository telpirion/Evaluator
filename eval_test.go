package main

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

var candidate = `
// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

'use strict';

// [START translate_sampleassistance_googlecloudtranslate_v8_v3_translationserviceclient_getglossary]
// Import the Translate library
const {TranslationServiceClient} = require('@google-cloud/translate').v3;

// Instantiate a reusable client object in the global scope
const client = new TranslationServiceClient();

/**
 * Get a Glossary.
 */
async function getGlossary(
  projectId,
  location = 'us-central1',
  glossaryId = 'my-glossary'
) {
  // Construct request
  const request = {
    parent: \` + "`projects/${projectId}/locations/${location}`," +
	"name: " + "`projects/${projectId}/locations/${location}/glossaries/${glossaryId}`," +
	`};

  // Get Glossary
  const response = await client.getGlossary(request);
  const glossary = response[0];

  console.log(\` + "`Got glossary: ${glossary.name}`);" +
	`}
// [END translate_sampleassistance_googlecloudtranslate_v8_v3_translationserviceclient_getglossary]

module.exports = {
  getGlossary,
};
`

func TestMain(m *testing.M) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	os.Exit(m.Run())
}

func TestImports(t *testing.T) {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	checks, err := NewChecks(projectID, "us-west1")
	if err != nil {
		t.Fatal(err)
	}
	defer checks.Close()

	evalReq := EvalRequest{
		Candidate: candidate,
		Test:      "",
		Library:   "@google-cloud/translate",
		Language:  "Node.js",
	}

	got, err := checks.Imports(evalReq)
	if err != nil {
		t.Fatal(err)
	}

	if !got.IsPass {
		t.Error("expected pass; Imports() failed")
	}
}

func TestCLI(t *testing.T) {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	checks, err := NewChecks(projectID, "us-west1")
	if err != nil {
		t.Fatal(err)
	}
	defer checks.Close()

	evalReq := EvalRequest{
		Candidate: candidate,
		Test:      "",
		Library:   "@google-cloud/translate",
		Language:  "Node.js",
	}

	got, err := checks.CLI(evalReq)
	if err != nil {
		t.Fatal(err)
	}

	if !got.IsPass {
		t.Error("expected pass; CLI() failed")
	}
}

func TestCasing(t *testing.T) {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	checks, err := NewChecks(projectID, "us-west1")
	if err != nil {
		t.Fatal(err)
	}
	defer checks.Close()

	evalReq := EvalRequest{
		Candidate: candidate,
		Test:      "",
		Library:   "@google-cloud/translate",
		Language:  "Node.js",
	}

	got, err := checks.Casing(evalReq, "camelCase")
	if err != nil {
		t.Fatal(err)
	}

	if !got.IsPass {
		t.Error("expected pass; Casing() failed")
	}
}
