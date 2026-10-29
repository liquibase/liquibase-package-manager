package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vifraa/gopom"
)

func TestGetCoreVersionFromPom(t *testing.T) {
	// Test with a different version
	pom := gopom.Project{
		Dependencies: &[]gopom.Dependency{
			{
				ArtifactID: &[]string{"liquibase-core"}[0],
				Version:    &[]string{"3.10.0"}[0],
			},
		},
	}
	got := GetCoreVersionFromPom(&pom)
	want := "3.10.0"
	if got != want {
		t.Errorf("GetCoreVersionFromPom() = %v, want %v", got, want)
	}

	// Test with no version specified
	pom = gopom.Project{
		Dependencies: &[]gopom.Dependency{
			{
				ArtifactID: &[]string{"liquibase-core"}[0],
			},
		},
		Properties: &gopom.Properties{
			Entries: map[string]string{
				"liquibase.version": "4.24.0",
			},
		},
	}
	got = GetCoreVersionFromPom(&pom)
	want = "4.24.0"
	if got != want {
		t.Errorf("GetCoreVersionFromPom() = %v, want %v", got, want)
	}

	// Test with no liquibase-core dependency
	pom = gopom.Project{
		Dependencies: &[]gopom.Dependency{
			{
				ArtifactID: &[]string{"not-liquibase-core"}[0],
				Version:    &[]string{"1.0.0"}[0],
			},
		},
		Properties: &gopom.Properties{
			Entries: map[string]string{
				"liquibase.version": "4.24.0",
			},
		},
	}
	got = GetCoreVersionFromPom(&pom)
	want = "4.24.0"
	if got != want {
		t.Errorf("GetCoreVersionFromPom() = %v, want %v", got, want)
	}
}
func TestGetCoreVersionFromPomWithProperty(t *testing.T) {
	// Create a sample POM object with a property
	pom := gopom.Project{
		Dependencies: &[]gopom.Dependency{
			{
				ArtifactID: &[]string{"liquibase-core"}[0],
				Version:    &[]string{"${liquibase.version}"}[0],
			},
		},
		Properties: &gopom.Properties{
			Entries: map[string]string{
				"liquibase.version": "4.24.0",
			},
		},
	}
	// Call the function and check the result
	got := GetCoreVersionFromPom(&pom)
	want := "4.24.0"
	if got != want {
		t.Errorf("GetCoreVersionFromPom() = %v, want %v", got, want)
	}
}

func TestGetCoreVersionFromPomWithNoDependencies(t *testing.T) {
	// Create a sample POM object with no dependencies
	pom := gopom.Project{
		Properties: &gopom.Properties{
			Entries: map[string]string{
				"liquibase.version": "4.24.0",
			},
		},
	}
	// Call the function and check the result
	got := GetCoreVersionFromPom(&pom)
	want := "4.24.0"
	if got != want {
		t.Errorf("GetCoreVersionFromPom() = %v, want %v", got, want)
	}
}

func TestGetCoreVersionFromPomWithNoVersion(t *testing.T) {
	// Create a sample POM object with no version specified
	pom := gopom.Project{
		Dependencies: &[]gopom.Dependency{
			{
				ArtifactID: &[]string{"liquibase-core"}[0],
			},
		},
		Properties: &gopom.Properties{
			Entries: map[string]string{
				"liquibase.version": "4.24.0",
			},
		},
	}
	// Call the function and check the result
	got := GetCoreVersionFromPom(&pom)
	want := "4.24.0"
	if got != want {
		t.Errorf("GetCoreVersionFromPom() = %v, want %v", got, want)
	}
}

func TestGetCoreVersionFromPomWithNoProperties(t *testing.T) {
	// Create a sample POM object with no properties
	pom := gopom.Project{
		Dependencies: &[]gopom.Dependency{
			{
				ArtifactID: &[]string{"liquibase-core"}[0],
				Version:    &[]string{"4.24.0"}[0],
			},
		},
	}
	// Call the function and check the result
	got := GetCoreVersionFromPom(&pom)
	want := "4.24.0"
	if got != want {
		t.Errorf("GetCoreVersionFromPom() = %v, want %v", got, want)
	}
}
func TestGetSHA1FromURL(t *testing.T) {
	const sha = "acff0ae38296af300d3da12e6e3fae3c166aa80f"
	tests := []struct {
		name   string
		status int
		body   string
		want   string
	}{
		{"checksum only", http.StatusOK, sha, sha},
		{"checksum with newline", http.StatusOK, sha + "\n", sha},
		{"checksum with filename", http.StatusOK, sha + "  osdt_core-21.23.0.0.jar\n", sha},
		{"uppercase checksum", http.StatusOK, "ACFF0AE38296AF300D3DA12E6E3FAE3C166AA80F", sha},
		{"plain text error message", http.StatusOK, "This tool has been identified as abusive and has been blocked", ""},
		{"html error page", http.StatusOK, "<html><body>Not Found</body></html>", ""},
		{"short body", http.StatusOK, "abc", ""},
		{"empty body", http.StatusOK, "", ""},
		{"non-200 status with checksum body", http.StatusForbidden, sha, ""},
		{"not found", http.StatusNotFound, "Not Found", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				fmt.Fprint(w, tt.body)
			}))
			defer server.Close()

			got := GetSHA1FromURL(server.URL + "/file.jar.sha1")
			if got != tt.want {
				t.Errorf("GetSHA1FromURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetSHA1FromURLUnreachable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	url := server.URL + "/file.jar.sha1"
	server.Close()

	got := GetSHA1FromURL(url)
	if got != "" {
		t.Errorf("GetSHA1FromURL() = %q, want empty string", got)
	}
}
