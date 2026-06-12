package gemini

// Close releases the underlying gRPC connection. Callers (factory + routes)
// should defer this at app shutdown.
func (p *GeminiProvider) Close() error {
	if p == nil || p.client == nil {
		return nil
	}
	return p.client.Close()
}
