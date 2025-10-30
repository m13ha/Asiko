package notifications

import (
    "bytes"
    "embed"
    "html/template"
    "sync"
)

//go:embed templates/*.html
var templatesFS embed.FS

var (
    tplCache = struct {
        mu sync.RWMutex
        m  map[string]*template.Template
    }{m: make(map[string]*template.Template)}
)

func getTemplate(path string) (*template.Template, error) {
    tplCache.mu.RLock()
    t, ok := tplCache.m[path]
    tplCache.mu.RUnlock()
    if ok {
        return t, nil
    }

    // Parse from embedded FS and cache
    parsed, err := template.ParseFS(templatesFS, path)
    if err != nil {
        return nil, err
    }
    tplCache.mu.Lock()
    tplCache.m[path] = parsed
    tplCache.mu.Unlock()
    return parsed, nil
}

func parseTemplate(templatePath string, data interface{}) (string, error) {
    t, err := getTemplate(templatePath)
    if err != nil {
        return "", err
    }
    var buf bytes.Buffer
    if err = t.Execute(&buf, data); err != nil {
        return "", err
    }
    return buf.String(), nil
}
