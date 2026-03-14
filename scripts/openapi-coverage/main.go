package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

var routePattern = regexp.MustCompile(`g\.(GET|POST|PUT|PATCH|DELETE)\("(/api/v1[^"]*)"`)

type openAPISpec struct {
	Paths map[string]map[string]any `yaml:"paths"`
}

func main() {
	handlerBytes, err := os.ReadFile("cmd/handlers.go")
	if err != nil {
		log.Fatalf("erro lendo handlers: %v", err)
	}

	specBytes, err := os.ReadFile("static/public/openapi.yaml")
	if err != nil {
		log.Fatalf("erro lendo openapi.yaml: %v", err)
	}

	actual := collectRoutes(string(handlerBytes))
	documented, err := collectDocumentedRoutes(specBytes)
	if err != nil {
		log.Fatalf("erro lendo spec OpenAPI: %v", err)
	}

	covered := 0
	missing := make([]string, 0)
	for _, route := range actual {
		if documented[route] {
			covered++
			continue
		}
		missing = append(missing, route)
	}

	sort.Strings(missing)
	coverage := float64(covered) / float64(len(actual)) * 100

	fmt.Printf("rotas_api_total=%d\n", len(actual))
	fmt.Printf("rotas_documentadas=%d\n", covered)
	fmt.Printf("cobertura_openapi=%.2f%%\n", coverage)

	if len(missing) > 0 {
		fmt.Println("rotas_sem_documentacao:")
		limit := len(missing)
		if limit > 25 {
			limit = 25
		}
		for i := 0; i < limit; i++ {
			fmt.Printf("- %s\n", missing[i])
		}
		if len(missing) > limit {
			fmt.Printf("- ... e mais %d rota(s)\n", len(missing)-limit)
		}
	}
}

func collectRoutes(src string) []string {
	matches := routePattern.FindAllStringSubmatch(src, -1)
	out := make([]string, 0, len(matches))
	for _, m := range matches {
		out = append(out, strings.ToUpper(m[1])+" "+m[2])
	}
	sort.Strings(out)
	return out
}

func collectDocumentedRoutes(specBytes []byte) (map[string]bool, error) {
	var spec openAPISpec
	if err := yaml.Unmarshal(specBytes, &spec); err != nil {
		return nil, err
	}

	out := make(map[string]bool)
	for path, methods := range spec.Paths {
		if !strings.HasPrefix(path, "/api/v1") {
			continue
		}
		for method := range methods {
			out[strings.ToUpper(method)+" "+path] = true
		}
	}
	return out, nil
}
