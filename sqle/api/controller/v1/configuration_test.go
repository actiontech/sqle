package v1

import (
	"testing"
)

func TestNewCloneOptionsWithoutBranch(t *testing.T) {
	cloneOpts := newCloneOptions("https://example.com/test.git", "", nil)

	if cloneOpts.URL != "https://example.com/test.git" {
		t.Fatalf("unexpected clone url: %s", cloneOpts.URL)
	}
	if cloneOpts.Depth != shallowCloneDepth {
		t.Fatalf("unexpected clone depth: %d", cloneOpts.Depth)
	}
	if !cloneOpts.InsecureSkipTLS {
		t.Fatal("expected InsecureSkipTLS to be true")
	}
	if cloneOpts.SingleBranch {
		t.Fatal("expected SingleBranch to be false when branch is empty")
	}
	if cloneOpts.ReferenceName.String() != "" {
		t.Fatalf("unexpected reference name: %s", cloneOpts.ReferenceName.String())
	}
}

func TestNewCloneOptionsWithBranch(t *testing.T) {
	cloneOpts := newCloneOptions("https://example.com/test.git", "refs/heads/main", nil)

	if cloneOpts.URL != "https://example.com/test.git" {
		t.Fatalf("unexpected clone url: %s", cloneOpts.URL)
	}
	if cloneOpts.Depth != shallowCloneDepth {
		t.Fatalf("unexpected clone depth: %d", cloneOpts.Depth)
	}
	if !cloneOpts.InsecureSkipTLS {
		t.Fatal("expected InsecureSkipTLS to be true")
	}
	if !cloneOpts.SingleBranch {
		t.Fatal("expected SingleBranch to be true when branch is provided")
	}
	if cloneOpts.ReferenceName.String() != "refs/heads/main" {
		t.Fatalf("unexpected reference name: %s", cloneOpts.ReferenceName.String())
	}
}
