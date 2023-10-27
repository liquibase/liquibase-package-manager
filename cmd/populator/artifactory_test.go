package main

import (
	"testing"
	"github.com/vifraa/gopom"
)

func TestGetCoreVersionFromPom(t *testing.T) {
	// Create a sample POM object
	pom := gopom.Project{
		Properties: &gopom.Properties{
			Entries: map[string]string{
				"liquibase.version": "4.24.0",
			},
		},
	}
	// Call the function and check the result
	got := GetCoreVersionFromPom(pom)
	want := "4.24.0"
	if got != want {
		t.Errorf("GetCoreVersionFromPom() = %v, want %v", got, want)
	}
}