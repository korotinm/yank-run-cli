package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/korotinm/yank-run-cli/internal/api"
	"github.com/korotinm/yank-run-cli/internal/ids"
)

// LastSearch is the persisted last-search.json document.
type LastSearch struct {
	Query     string    `json:"query"`
	Timestamp int64     `json:"timestamp"`
	Hits      []api.Hit `json:"hits"`
}

// ErrNoRecentSearch — index used but no last-search.json exists.
var ErrNoRecentSearch = errors.New("no recent search; provide a full 64-char id")

// ErrBadRef — argument is neither a valid id nor a resolvable index.
var ErrBadRef = errors.New("not a valid id or index")

// cacheDir returns ~/.cache/yank, creating it if missing.
func cacheDir() (string, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(base, "yank")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

func filePath() (string, error) {
	dir, err := cacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "last-search.json"), nil
}

// Save overwrites last-search.json with the given query+hits.
func Save(query string, hits []api.Hit) error {
	path, err := filePath()
	if err != nil {
		return err
	}
	doc := LastSearch{
		Query:     query,
		Timestamp: time.Now().UnixMilli(),
		Hits:      hits,
	}
	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// Load reads last-search.json. Returns ErrNoRecentSearch if absent.
func Load() (*LastSearch, error) {
	path, err := filePath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrNoRecentSearch
		}
		return nil, err
	}
	var doc LastSearch
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("decode last-search.json: %w", err)
	}
	return &doc, nil
}

// ResolveRef turns a CLI argument into a snippet id.
//
//   - 64-char lowercase hex          -> id itself
//   - non-negative integer < 1000    -> index into last-search hits
//   - anything else                  -> ErrBadRef
func ResolveRef(arg string) (string, error) {
	if ids.IsValidID(arg) {
		return arg, nil
	}
	if !ids.IsNumericIndex(arg) {
		return "", ErrBadRef
	}
	idx, err := strconv.Atoi(arg)
	if err != nil {
		return "", ErrBadRef
	}
	doc, err := Load()
	if err != nil {
		return "", err
	}
	if idx < 0 || idx >= len(doc.Hits) {
		return "", fmt.Errorf("index %d out of range (0..%d)", idx, len(doc.Hits)-1)
	}
	return doc.Hits[idx].ID, nil
}
