package utility

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestProgressBarTTYFormat(t *testing.T) {
	var b bytes.Buffer
	p := newPushProgressReporter(&b, true)

	p.Start(100)
	p.Update(52, "pushing layers")

	out := b.String()
	if !strings.Contains(out, "[====") || !strings.Contains(out, "52%") {
		t.Fatalf("unexpected tty progress output: %q", out)
	}
}

func TestProgressBarNonTTYFallback(t *testing.T) {
	var b bytes.Buffer
	p := newPushProgressReporter(&b, false)

	p.Start(100)
	p.Update(52, "pushing layers")

	out := b.String()
	if strings.Contains(out, "[") {
		t.Fatalf("unexpected bar output in non-tty mode: %q", out)
	}
	if !strings.Contains(out, "pushing layers") || !strings.Contains(out, "52%") {
		t.Fatalf("unexpected non-tty progress output: %q", out)
	}
}

func TestProgressBarFinalizeOnCompleteAndFail(t *testing.T) {
	t.Run("complete", func(t *testing.T) {
		var b bytes.Buffer
		p := newPushProgressReporter(&b, true)

		p.Start(100)
		p.Update(100, "done")
		p.Complete("Successfully pushed image")

		out := b.String()
		if !strings.Contains(out, "Successfully pushed image") {
			t.Fatalf("expected completion message, got: %q", out)
		}
		if !strings.HasSuffix(out, "\n") {
			t.Fatalf("expected final newline on completion, got: %q", out)
		}
	})

	t.Run("fail", func(t *testing.T) {
		var b bytes.Buffer
		p := newPushProgressReporter(&b, true)

		p.Start(100)
		p.Update(60, "pushing layers")
		p.Fail(errors.New("boom"))

		out := b.String()
		if !strings.Contains(out, "boom") {
			t.Fatalf("expected error in output, got: %q", out)
		}
		if !strings.HasSuffix(out, "\n") {
			t.Fatalf("expected final newline on failure, got: %q", out)
		}
	})
}