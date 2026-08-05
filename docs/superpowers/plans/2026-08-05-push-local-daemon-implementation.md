# Push Local Daemon Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace tar-based push with local-daemon-to-remote push using pure Go, optional Kubernetes secret auth, and terminal progress bar output.

**Architecture:** Keep command parsing in `cmd/push.go` and orchestration in `utility/push.go`. Load source images from the local daemon, then push to remote registry with `go-containerregistry` remote APIs and keychain auth from existing shared auth helpers. Add a focused progress reporter component to isolate TTY/non-TTY output behavior.

**Tech Stack:** Go 1.24, Cobra, go-containerregistry (`name`, `daemon`, `remote`, `authn`), existing Kubernetes keychain auth.

## Global Constraints

- Scope: Push redesign only (pull and copy deferred to next layers)
- No pull redesign in this layer.
- No copy redesign in this layer.
- No compatibility layer for old push semantics.
- No Docker CLI invocation.
- Required: `--source-image`
- Required: `--destination-image`
- Optional: `--secret`
- Optional: `--namespace` (default `default`)
- `source-image` is always a local Docker daemon image reference.
- `destination-image` is always a remote registry image reference.
- If local image is missing, command fails immediately with an actionable message.
- Do not change shared auth semantics used by list.
- Immediate switch to new push semantics (breaking change accepted).

---

## File Structure and Responsibilities

- Modify: `cmd/push.go`
- Modify: `utility/push.go`
- Modify: `utility/push_test.go`
- Create: `utility/push_progress.go`
- Create: `utility/push_progress_test.go`
- Modify: `docs/README.md`

Responsibilities:

- `cmd/push.go`: new flag contract and command help.
- `utility/push.go`: parse/validate refs, local image load, auth wiring, remote push orchestration.
- `utility/push_progress.go`: character progress bar + non-TTY fallback output.
- `utility/push_test.go`: push orchestration and fail-fast/error mapping tests.
- `utility/push_progress_test.go`: progress renderer behavior tests.
- `docs/README.md`: command usage examples and flag updates.

### Task 1: Switch Push CLI Contract

**Files:**
- Modify: `cmd/push.go`
- Test: `cmd/push_test.go` (create if missing)

**Interfaces:**
- Consumes: `utility.PushImage(sourceImage, destinationImage, secretName, namespace string) error`
- Produces: push command contract with flags:
  - `--source-image` (required)
  - `--destination-image` (required)
  - `--secret` (optional)
  - `--namespace` (optional, default `default`)

- [ ] **Step 1: Write the failing command flag tests**

```go
func TestPushCommandRequiresSourceAndDestinationImage(t *testing.T) {
    cmd := pushCmd
    cmd.SetArgs([]string{"--destination-image", "registry.io/team/app:1.0.0"})
    err := cmd.Execute()
    if err == nil {
        t.Fatalf("expected missing source-image error")
    }
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./cmd -run TestPushCommandRequiresSourceAndDestinationImage -v`
Expected: FAIL because flags still use old names and old required fields.

- [ ] **Step 3: Implement minimal flag contract changes in command layer**

```go
pushCmd.Flags().StringVar(&pushSourceImage, "source-image", "", "Local Docker source image reference (required)")
pushCmd.Flags().StringVar(&pushDestinationImage, "destination-image", "", "Destination registry image reference (required)")
_ = pushCmd.MarkFlagRequired("source-image")
_ = pushCmd.MarkFlagRequired("destination-image")
```

- [ ] **Step 4: Run targeted cmd tests**

Run: `go test ./cmd -run TestPush -v`
Expected: PASS for new contract tests.

- [ ] **Step 5: Commit**

```bash
git add cmd/push.go cmd/push_test.go
git commit -m "feat(push): switch CLI to source-image and destination-image"
```

### Task 2: Load Source Image from Local Daemon (Fail Fast)

**Files:**
- Modify: `utility/push.go`
- Test: `utility/push_test.go`

**Interfaces:**
- Consumes:
  - `CreateKeychain(namespace, secretName string) (authn.Keychain, error)`
  - `HandleRegistryError(err error, operation string, target string) error`
- Produces:
  - `func PushImage(sourceImageRef, destinationImageRef, secretName, namespace string) error`
  - internal helper: `func loadLocalImage(src name.Reference) (v1.Image, error)`

- [ ] **Step 1: Write failing tests for local source resolution**

```go
func TestPushImageFailsWhenLocalSourceMissing(t *testing.T) {
    err := PushImage("local/not-found:dev", "registry.io/team/app:1.0.0", "", "default")
    if err == nil || !strings.Contains(err.Error(), "local image") {
        t.Fatalf("expected local image missing error, got %v", err)
    }
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./utility -run TestPushImageFailsWhenLocalSourceMissing -v`
Expected: FAIL because implementation still expects tar path and does not check local daemon.

- [ ] **Step 3: Implement minimal local-daemon loading path**

```go
srcRef, err := name.ParseReference(sourceImageRef)
if err != nil { return fmt.Errorf("failed to parse source image reference '%s': %w", sourceImageRef, err) }

img, err := daemon.Image(srcRef)
if err != nil {
    return fmt.Errorf("local image '%s' not found in Docker daemon. Tag the image first, then retry: %w", sourceImageRef, err)
}
```

- [ ] **Step 4: Run targeted utility tests**

Run: `go test ./utility -run TestPushImage.*Local -v`
Expected: PASS for local-source tests.

- [ ] **Step 5: Commit**

```bash
git add utility/push.go utility/push_test.go
git commit -m "feat(push): load source image from local daemon with fail-fast errors"
```

### Task 3: Push to Remote with Stable Shared Auth Behavior

**Files:**
- Modify: `utility/push.go`
- Test: `utility/push_test.go`

**Interfaces:**
- Consumes:
  - `CreateKeychain(namespace, secretName string) (authn.Keychain, error)`
  - `remote.Write(ref name.Reference, img v1.Image, options ...remote.Option) error`
- Produces:
  - push flow that honors optional secret while preserving shared auth semantics

- [ ] **Step 1: Write failing tests for auth and remote error mapping**

```go
func TestPushImageUsesDefaultKeychainWhenSecretEmpty(t *testing.T) {
    // inject fake keychain factory and verify empty secret path is used
}

func TestPushImageWrapsRegistryWriteError(t *testing.T) {
    // inject remote writer returning unauthorized and assert HandleRegistryError output
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./utility -run TestPushImage.*(Keychain|Registry) -v`
Expected: FAIL until injection points and auth/write path are updated.

- [ ] **Step 3: Implement minimal orchestration with injectable function vars**

```go
var createKeychainFn = CreateKeychain
var remoteWriteFn = remote.Write

kc, err := createKeychainFn(namespace, secretName)
if err != nil { return fmt.Errorf("failed to create keychain: %w", err) }

if err := remoteWriteFn(dstRef, img, remote.WithAuthFromKeychain(kc)); err != nil {
    return HandleRegistryError(err, "pushing image to", destinationImageRef)
}
```

- [ ] **Step 4: Run targeted utility tests**

Run: `go test ./utility -run TestPushImage -v`
Expected: PASS including secret-optional path and registry error mapping.

- [ ] **Step 5: Commit**

```bash
git add utility/push.go utility/push_test.go
git commit -m "feat(push): wire remote push with optional keychain auth"
```

### Task 4: Add Character Progress Bar and Non-TTY Fallback

**Files:**
- Create: `utility/push_progress.go`
- Create: `utility/push_progress_test.go`
- Modify: `utility/push.go`

**Interfaces:**
- Consumes:
  - remote progress events (or phase callbacks from push flow)
- Produces:
  - `type pushProgressReporter struct { ... }`
  - `func newPushProgressReporter(w io.Writer, isTTY bool) *pushProgressReporter`
  - methods:
    - `func (p *pushProgressReporter) Start(total int64)`
    - `func (p *pushProgressReporter) Update(done int64, phase string)`
    - `func (p *pushProgressReporter) Complete(msg string)`
    - `func (p *pushProgressReporter) Fail(err error)`

- [ ] **Step 1: Write failing progress renderer tests**

```go
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
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./utility -run TestProgressBar -v`
Expected: FAIL because progress reporter is not implemented.

- [ ] **Step 3: Implement minimal progress reporter and integrate into push flow**

```go
reporter := newPushProgressReporter(os.Stdout, isInteractiveTerminal())
reporter.Start(total)
// update in callbacks or deterministic phase points
reporter.Update(current, "pushing layers")
reporter.Complete("Successfully pushed image")
```

- [ ] **Step 4: Run utility tests for progress paths**

Run: `go test ./utility -run TestProgress -v`
Expected: PASS for TTY and non-TTY fallback behavior.

- [ ] **Step 5: Commit**

```bash
git add utility/push_progress.go utility/push_progress_test.go utility/push.go
git commit -m "feat(push): add character progress bar with non-tty fallback"
```

### Task 5: Update and Expand Push Tests End-to-End at Utility Layer

**Files:**
- Modify: `utility/push_test.go`

**Interfaces:**
- Consumes:
  - `PushImage(sourceImageRef, destinationImageRef, secretName, namespace string) error`
  - injected test doubles for daemon load, keychain, and remote write
- Produces:
  - stable tests for validation, local missing, auth, remote error mapping, and success path

- [ ] **Step 1: Write failing table-driven tests for final behavior matrix**

```go
func TestPushImageBehaviorMatrix(t *testing.T) {
    cases := []struct {
        name string
        source string
        dest string
        secret string
        wantErrSubstr string
    }{
        {"missing source", "", "registry.io/team/app:1.0.0", "", "source image"},
        {"missing destination", "local/app:dev", "", "", "destination image"},
        {"invalid destination", "local/app:dev", ":::bad:::", "", "parse"},
    }
    // assert each case
}
```

- [ ] **Step 2: Run tests to verify expected failures**

Run: `go test ./utility -run TestPushImageBehaviorMatrix -v`
Expected: FAIL until final validation and messages are complete.

- [ ] **Step 3: Implement minimal validation/message improvements to satisfy tests**

```go
if strings.TrimSpace(sourceImageRef) == "" {
    return fmt.Errorf("source image is required")
}
if strings.TrimSpace(destinationImageRef) == "" {
    return fmt.Errorf("destination image is required")
}
```

- [ ] **Step 4: Run full utility tests**

Run: `go test ./utility -v`
Expected: PASS with deterministic results and no reliance on live registry.

- [ ] **Step 5: Commit**

```bash
git add utility/push_test.go utility/push.go
git commit -m "test(push): finalize behavior matrix for local-daemon push"
```

### Task 6: Update User Documentation and Final Verification

**Files:**
- Modify: `docs/README.md`

**Interfaces:**
- Consumes: final push CLI contract and behavior
- Produces: documented push examples using source-image and destination-image

- [ ] **Step 1: Write failing docs check script (grep assertions) locally**

```bash
grep -q "--source-image" docs/README.md
grep -q "--destination-image" docs/README.md
! grep -q "--source <path.tar>" docs/README.md
```

- [ ] **Step 2: Run docs checks to verify they fail first**

Run: `bash -c 'grep -q "--source-image" docs/README.md && grep -q "--destination-image" docs/README.md && ! grep -q "--source <path.tar>" docs/README.md'`
Expected: FAIL until docs are updated.

- [ ] **Step 3: Update push command docs to new semantics**

```markdown
repo-lister push \
  --source-image local/app:dev \
  --destination-image registry.io/team/app:1.0.0 \
  --secret regcred \
  --namespace default
```

- [ ] **Step 4: Run full verification suite**

Run: `go test ./... -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add docs/README.md
git commit -m "docs(push): describe local-daemon source and destination-image flow"
```

## Final Integration Checkpoint

- [ ] Run: `go test ./... -v`
- [ ] Run: `go build ./...`
- [ ] Manual smoke test in Docker-enabled environment:

```bash
repo-lister push \
  --source-image local/app:dev \
  --destination-image registry.io/team/app:1.0.0 \
  --secret regcred \
  --namespace default
```

- [ ] Confirm output includes progress bar and clean success line.

## Self-Review Results

Spec coverage check:

- Command contract covered in Task 1 and Task 6.
- Local daemon source and fail-fast behavior covered in Task 2 and Task 5.
- Optional secret auth with no shared-auth regression covered in Task 3 and Task 5.
- Character progress bar and non-TTY fallback covered in Task 4.
- Testing coverage and final verification covered in Task 5, Task 6, and Final Integration Checkpoint.

Placeholder scan:

- No unresolved placeholder markers remain.

Type and interface consistency check:

- `PushImage(sourceImageRef, destinationImageRef, secretName, namespace string) error` is used consistently across all tasks.
- Progress reporter method names are consistent across creation, test, and integration steps.
