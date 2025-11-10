package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp file %s: %v", name, err)
	}
	return path
}

func Test_DryRunOutputsTitlesToFile(t *testing.T) {
	// Arrange temporary workspace
	tmp := t.TempDir()
	conversationsJSON := `[
	  {
	    "title": "Conversation 1",
	    "mapping": {
	      "a": {"message": {"create_time": 1634000000, "author": {"role": "user"}, "content": {"parts": ["Hello!"]}}}
	    }
	  },
	  {
	    "title": "Conversation 2",
	    "mapping": {
	      "b": {"message": {"create_time": 1634000100, "author": {"role": "assistant"}, "content": {"parts": ["Hi!"]}}}
	    }
	  }
	]`
	source := writeTempFile(t, tmp, "conversations.json", conversationsJSON)
	output := filepath.Join(tmp, "titles.txt")

	// Act: dry run should list titles
	var stdout, stderr strings.Builder
	err := run([]string{"-dry", source, output}, &stdout, &stderr)

	// Assert
	if err != nil {
		t.Fatalf("run returned error: %v (stderr=%s)", err, stderr.String())
	}
	data, readErr := os.ReadFile(output)
	if readErr != nil {
		t.Fatalf("failed reading output: %v", readErr)
	}
	expectedLines := []string{"Conversation 1", "Conversation 2"}
	got := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(got) != len(expectedLines) {
		t.Fatalf("unexpected number of lines. got=%d want=%d\ncontent=\n%s", len(got), len(expectedLines), string(data))
	}
	for i, line := range expectedLines {
		if got[i] != line {
			t.Fatalf("line %d mismatch: got=%q want=%q", i, got[i], line)
		}
	}
	// also stdout should include titles
	outStr := stdout.String()
	for _, line := range expectedLines {
		if !strings.Contains(outStr, line) {
			t.Fatalf("stdout missing %q. got: %s", line, outStr)
		}
	}
}

func Test_RunRequiresIncludeWhenMerging(t *testing.T) {
	tmp := t.TempDir()
	source := writeTempFile(t, tmp, "conversations.json", `[]`)
	csvOut := filepath.Join(tmp, "out.csv")
	var stdout, stderr strings.Builder
	err := run([]string{source, csvOut}, &stdout, &stderr)
	if err == nil {
		t.Fatalf("expected error when -include is missing in merge mode")
	}
	if !strings.Contains(err.Error(), "include") {
		t.Fatalf("error should mention include, got: %v", err)
	}
}

func Test_VersionFlagPrintsAndExits(t *testing.T) {
	var stdout, stderr strings.Builder
	err := run([]string{"-version"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("-version should not error: %v (stderr=%s)", err, stderr.String())
	}
	out := stdout.String()
	if !strings.HasPrefix(out, "chatgpt-merge ") {
		t.Fatalf("unexpected version output: %q", out)
	}
}

func Test_NoHeaderFlagOmitsHeader(t *testing.T) {
	tmp := t.TempDir()
	// Conversations with one selected title producing one row
	conversationsJSON := `[
  {
    "title": "Conversation 1",
    "mapping": {
      "a": {"message": {"create_time": 1634000000, "author": {"role": "user"}, "content": {"parts": ["Hello!"]}}}
    }
  },
  {
    "title": "Other",
    "mapping": {
      "b": {"message": {"create_time": 1634000100, "author": {"role": "assistant"}, "content": {"parts": ["Ignore"]}}}
    }
  }
]`
	source := writeTempFile(t, tmp, "conversations.json", conversationsJSON)
	include := writeTempFile(t, tmp, "include.txt", "Conversation 1\n")
	csvOut := filepath.Join(tmp, "out.csv")

	var stdout, stderr strings.Builder
	err := run([]string{"-no-header", "-include", include, source, csvOut}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run returned error: %v (stderr=%s)", err, stderr.String())
	}
	data, readErr := os.ReadFile(csvOut)
	if readErr != nil {
		t.Fatalf("failed reading csv: %v", readErr)
	}
	csv := string(data)
	if strings.HasPrefix(csv, "Timestamp,Role,Content\n") {
		t.Fatalf("expected no header, but header present: %q", csv)
	}
	expected := "2021-10-12T00:53:20Z,user,Hello!\n"
	if csv != expected {
		t.Fatalf("unexpected csv.\n got: %q\nwant: %q", csv, expected)
	}
}

func Test_BOMFlagWritesPrefix(t *testing.T) {
	tmp := t.TempDir()
	conversationsJSON := `[
  {
    "title": "Conversation 1",
    "mapping": {
      "a": {"message": {"create_time": 1634000000, "author": {"role": "user"}, "content": {"parts": ["Hello!"]}}}
    }
  }
]`
	source := writeTempFile(t, tmp, "conversations.json", conversationsJSON)
	include := writeTempFile(t, tmp, "include.txt", "Conversation 1\n")
	csvOut := filepath.Join(tmp, "out.csv")

	var stdout, stderr strings.Builder
	err := run([]string{"-bom", "-include", include, source, csvOut}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run returned error: %v (stderr=%s)", err, stderr.String())
	}
	bytesOut, readErr := os.ReadFile(csvOut)
	if readErr != nil {
		t.Fatalf("failed reading csv: %v", readErr)
	}
	if len(bytesOut) < 3 {
		t.Fatalf("output too short: %d bytes", len(bytesOut))
	}
	bom := []byte{0xEF, 0xBB, 0xBF}
	if bytesOut[0] != bom[0] || bytesOut[1] != bom[1] || bytesOut[2] != bom[2] {
		t.Fatalf("expected BOM prefix EF BB BF, got: % X", bytesOut[:3])
	}
	rest := string(bytesOut[3:])
	expectedPrefix := "Timestamp,Role,Content\n2021-10-12T00:53:20Z,user,Hello!\n"
	if rest != expectedPrefix {
		t.Fatalf("unexpected CSV after BOM.\n got: %q\nwant: %q", rest, expectedPrefix)
	}
}

// Ensure run signature remains stable when used by tests
var _ = func() any {
	var w io.Writer = os.Stdout
	_ = w
	return nil
}()

// --- Additional negative-path tests for error handling ---

type errWriter struct{ err error }

func (e errWriter) Write(_ []byte) (int, error) { return 0, e.err }

func Test_Version_WriterErrorIsPropagated(t *testing.T) {
	var stderr strings.Builder
	werr := errors.New("stdout write failed")
	stdout := errWriter{err: werr}
	if err := run([]string{"-version"}, stdout, &stderr); err == nil || !strings.Contains(err.Error(), werr.Error()) {
		t.Fatalf("expected stdout error to propagate, got: %v", err)
	}
}

func Test_DryRun_MirrorStdoutWriteError(t *testing.T) {
	tmp := t.TempDir()
	conversationsJSON := `[{"title":"A","mapping":{"k":{"message":{"create_time":1634000000,"author":{"role":"user"},"content":{"parts":["x"]}}}}}]`
	source := writeTempFile(t, tmp, "conversations.json", conversationsJSON)
	out := filepath.Join(tmp, "titles.txt")
	werr := errors.New("stdout write failed")
	stdout := errWriter{err: werr}
	var stderr strings.Builder
	if err := run([]string{"-dry", source, out}, stdout, &stderr); err == nil || !strings.Contains(err.Error(), werr.Error()) {
		t.Fatalf("expected stdout error to propagate in dry run, got: %v", err)
	}
}

func Test_Merge_InvalidOutputPathErrorWrapped(t *testing.T) {
	tmp := t.TempDir()
	// valid conversations and include file
	conversationsJSON := `[{"title":"A","mapping":{"k":{"message":{"create_time":1634000000,"author":{"role":"user"},"content":{"parts":["x"]}}}}}]`
	source := writeTempFile(t, tmp, "conversations.json", conversationsJSON)
	include := writeTempFile(t, tmp, "include.txt", "A\n")
	// Intentionally use directory path (tmp) as output file, which should fail to create
	var stdout, stderr strings.Builder
	err := run([]string{"-include", include, source, tmp}, &stdout, &stderr)
	if err == nil {
		t.Fatalf("expected error when output path is a directory")
	}
	if !strings.Contains(err.Error(), "error creating output file") {
		t.Fatalf("expected wrapped creation error, got: %v", err)
	}
}
