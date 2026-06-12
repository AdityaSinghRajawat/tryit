package config

type parseConsts struct {
	// confidenceThreshold: below this, the panel forces the confirm step (D7).
	confidenceThreshold float64
	// aiRepairRetries: extra AI calls after a schema-validation failure (§7.4).
	aiRepairRetries int
	// systemPromptTemplate is the verbatim AI parsing instruction (IMPL §7.1)
	// with " SCHEMA:\n" trailing so the live schema bytes get appended at use.
	systemPromptTemplate string
}

var parseI = &parseConsts{
	confidenceThreshold: 0.75,
	aiRepairRetries:     1,
	systemPromptTemplate: `You convert one section of API documentation into a single JSON object that matches the
RequestSpec schema provided below. You output ONLY that JSON object — no prose, no markdown fences.

RULES:
1. Use ONLY information present in the supplied documentation and page URL. Never invent
   endpoints, parameters, or auth that are not indicated by the text.
2. baseUrl: take the absolute host from an example request if present; otherwise infer it from
   the page URL's host and any "Base URL" statement. Never include a trailing slash.
3. path: use {name} tokens for path parameters (e.g. /v1/users/{id}); list each in pathParams.
4. Distinguish path vs query vs body parameters. Mark each required:true only when the docs say so.
   For repeated/array query params, use "values". Put nested request bodies verbatim under body.json.
5. auth: detect the type (bearer|basic|apiKey|none) and location.
   - Prefer an explicit example (Authorization: Bearer …, -u user:pass, ?key=…, X-API-Key:…).
   - Otherwise infer from prose in the supplied Authentication section.
   - NEVER output a real secret. Use placeholders: bearer/apiKey → valueRef "{{secret:NAME}}";
     basic → username/password each "{{secret:NAME}}". Choose a stable UPPER_SNAKE NAME derived
     from the host (e.g. STRIPE_API_KEY, TWILIO_SID, TWILIO_AUTH_TOKEN).
   - bearer prefix defaults to "Bearer "; apiKey prefix defaults to "".
6. confidence (0–1): follow the rubric in section CONFIDENCE.
7. Output must validate against SCHEMA. If unsure of a value, omit it rather than guessing,
   and lower confidence accordingly.

CONFIDENCE:
  0.90–1.00  A complete, unambiguous example request is present (method+URL+auth+body).
  0.75–0.89  Clear prose plus a partial example; minor inference only.
  0.50–0.74  Inferred mostly from prose; some ambiguity in params or auth location.
  0.00–0.49  Largely guessed; key facts (host/auth/required params) are missing.

SCHEMA:
`,
}

func GetParseConfidenceThreshold() float64 { return parseI.confidenceThreshold }
func GetAIRepairRetries() int              { return parseI.aiRepairRetries }
func GetParseSystemPromptTemplate() string { return parseI.systemPromptTemplate }
