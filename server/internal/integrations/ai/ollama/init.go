package ollama

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	"github.com/ollama/ollama/api"
)

type OllamaProvider struct {
	client *api.Client
	model  string
	host   string
}

// NewOllamaProvider constructs the client; connectivity is verified lazily
// on first Complete call. Surfaces a clear error when OLLAMA_HOST is empty
// or unparseable.
func NewOllamaProvider() (*OllamaProvider, error) {
	host := config.GetOllamaHost()
	if host == "" {
		return nil, errors.New("OLLAMA_HOST is empty")
	}

	u, err := url.Parse(host)
	if err != nil {
		return nil, fmt.Errorf("OLLAMA_HOST: invalid URL %q: %w", host, err)
	}

	return &OllamaProvider{
		client: api.NewClient(u, http.DefaultClient),
		model:  config.GetOllamaModel(),
		host:   host,
	}, nil
}
