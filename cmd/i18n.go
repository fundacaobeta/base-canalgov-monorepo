package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/fundacaobeta/base-canalgov-monorepo/internal/envelope"
	"github.com/knadh/go-i18n"
	"github.com/knadh/stuffbin"
	"github.com/zerodha/fastglue"
)

const (
	defLang = "pt-BR"
)

// handleGetI18nLang returns the JSON language pack for the given language code.
func handleGetI18nLang(r *fastglue.Request) error {
	var (
		app  = r.Context.(*App)
		lang = r.RequestCtx.UserValue("lang").(string)
	)
	i, err := loadI18nLang(lang, app.fs)
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendBytes(http.StatusOK, "application/json", i.JSON())
}

// initI18n initializes the i18n manager with the default language.
func initI18n(fs stuffbin.FileSystem) *i18n.I18n {
	lang := ko.String("app.lang")
	if lang == "" {
		lang = defLang
	}

	i, err := loadI18nLang(lang, fs)
	if err != nil {
		log.Fatalf("error initializing i18n manager: %v", err)
	}

	return i
}

// loadI18nLang loads the i18n language pack for the given language code.
func loadI18nLang(lang string, fs stuffbin.FileSystem) (*i18n.I18n, error) {
	// mergeFiles merges all JSON files in a directory into a single map.
	mergeFiles := func(langCode string) (map[string]interface{}, error) {
	        merged := make(map[string]interface{})
	        basePath := fmt.Sprintf("i18n/%s", langCode)

	        // Get all files in the directory.
	        files := fs.List()
	        found := false
	        for _, f := range files {
	                // Normalize path by removing leading slash for comparison.
	                cleanPath := strings.TrimPrefix(f, "/")
	                if strings.HasPrefix(cleanPath, basePath) && filepath.Ext(f) == ".json" {				b, err := fs.Read(f)
				if err != nil {
					continue
				}

				var data map[string]interface{}
				if err := json.Unmarshal(b, &data); err != nil {
					continue
				}

				for k, v := range data {
					merged[k] = v
				}
				found = true
			}
		}

		if !found {
			return nil, fmt.Errorf("no translation files found for %s", langCode)
		}

		return merged, nil
	}

	// Read and merge default language files.
	defaultData, err := mergeFiles(defLang)
	if err != nil {
		return nil, envelope.NewError(envelope.GeneralError, fmt.Sprintf("error reading default language `%s` : %v", defLang, err), nil)
	}

	b, _ := json.Marshal(defaultData)
	i, err := i18n.New(b)
	if err != nil {
		return nil, envelope.NewError(envelope.GeneralError, "error initializing i18n", nil)
	}

	// Load the selected language on top of it if it's different from default.
	if lang != defLang {
		langData, err := mergeFiles(lang)
		if err == nil {
			lb, _ := json.Marshal(langData)
			i.Load(lb)
		}
	}

	return i, nil
}
