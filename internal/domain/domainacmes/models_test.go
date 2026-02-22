package domainacmes

import (
	"reflect"
	"testing"
	"time"
)

func TestAcme_ShouldKeepAssignedDomainState(t *testing.T) {
	now := time.Date(2026, 2, 22, 10, 30, 0, 0, time.UTC)
	prefix := HelperRandomAlphaPrefix(t, 8)
	expectedID := prefix + "-acme-001"
	expectedTitle := prefix + "-Servicio de Documentos"
	expectedDescription := prefix + "-Plantilla base"
	expectedCategory := prefix + "-backend"
	expectedContentB64 := prefix + "-Y29udGVudA=="
	expectedCoverImage := "https://example.com/" + prefix + "/cover.png"
	expectedTags := []string{prefix + "-go", prefix + "-grpc", prefix + "-postgres"}

	acme := Acme{
		ID:          expectedID,
		Title:       expectedTitle,
		Description: expectedDescription,
		Category:    expectedCategory,
		Tags:        expectedTags,
		UpdatedAt:   now,
		ContentB64:  expectedContentB64,
		CoverImage:  expectedCoverImage,
	}

	if acme.ID != expectedID {
		t.Fatalf("expected ID %s, got %s", expectedID, acme.ID)
	}
	if acme.Title != expectedTitle {
		t.Fatalf("expected Title %s, got %s", expectedTitle, acme.Title)
	}
	if acme.Description != expectedDescription {
		t.Fatalf("expected Description %s, got %s", expectedDescription, acme.Description)
	}
	if acme.Category != expectedCategory {
		t.Fatalf("expected Category %s, got %s", expectedCategory, acme.Category)
	}
	if !reflect.DeepEqual(acme.Tags, expectedTags) {
		t.Fatalf("expected Tags %v, got %v", expectedTags, acme.Tags)
	}
	if !acme.UpdatedAt.Equal(now) {
		t.Fatalf("expected UpdatedAt %v, got %v", now, acme.UpdatedAt)
	}
	if acme.ContentB64 != expectedContentB64 {
		t.Fatalf("expected ContentB64 %s, got %s", expectedContentB64, acme.ContentB64)
	}
	if acme.CoverImage != expectedCoverImage {
		t.Fatalf("expected CoverImage %s, got %s", expectedCoverImage, acme.CoverImage)
	}
}

func TestAcme_ShouldExposeZeroValueAsEmptyDomainState(t *testing.T) {
	var acme Acme

	if acme.ID != "" {
		t.Fatalf("expected empty ID, got %s", acme.ID)
	}
	if acme.Title != "" {
		t.Fatalf("expected empty Title, got %s", acme.Title)
	}
	if acme.Description != "" {
		t.Fatalf("expected empty Description, got %s", acme.Description)
	}
	if acme.Category != "" {
		t.Fatalf("expected empty Category, got %s", acme.Category)
	}
	if len(acme.Tags) != 0 {
		t.Fatalf("expected empty Tags, got %v", acme.Tags)
	}
	if !acme.UpdatedAt.IsZero() {
		t.Fatalf("expected zero UpdatedAt, got %v", acme.UpdatedAt)
	}
	if acme.ContentB64 != "" {
		t.Fatalf("expected empty ContentB64, got %s", acme.ContentB64)
	}
	if acme.CoverImage != "" {
		t.Fatalf("expected empty CoverImage, got %s", acme.CoverImage)
	}
}
