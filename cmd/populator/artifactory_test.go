package main

import (
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