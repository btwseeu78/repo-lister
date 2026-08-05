# Push Redesign: Local Docker Source to Remote Registry

Date: 2026-08-05
Status: Approved for planning
Scope: Push redesign only (pull and copy deferred to next layers)

## 1. Problem and Goal

Current push behavior expects a tar file source. That workflow is not needed for this project stage.

Goal: push a local Docker daemon image to a destination registry using pure Go and optional Kubernetes secret auth, with clear progress feedback and fail-fast validation.

## 2. Non-Goals

- No pull redesign in this layer.
- No copy redesign in this layer.
- No compatibility layer for old push semantics.
- No Docker CLI invocation.

## 3. Command Contract

Push uses explicit image flags:

- Required: `--source-image`
- Required: `--destination-image`
- Optional: `--secret`
- Optional: `--namespace` (default `default`)

Behavior rules:

- `source-image` is always a local Docker daemon image reference.
- `destination-image` is always a remote registry image reference.
- Old push behavior is removed immediately (breaking change accepted).
- If local image is missing, command fails immediately with an actionable message.

## 4. Architecture and Boundaries

### 4.1 Command Layer

File: `cmd/push.go`

Responsibilities:

- Define new push flags and help text.
- Enforce required args.
- Pass validated values to utility entrypoint.

### 4.2 Utility Layer

File: `utility/push.go`

Responsibilities:

- Parse and validate source/destination references.
- Resolve source image from local daemon.
- Resolve auth keychain (optional secret path supported).
- Push to destination registry.
- Render and finalize progress output.
- Return context-rich errors.

### 4.3 Shared Auth and Error Layers

Files: `utility/auth.go`, `utility/errors.go`

Responsibilities:

- Keep auth behavior stable for all commands, especially list.
- Keep registry/auth/network error mapping consistent.

Constraint:

- Do not change shared auth semantics used by list.

## 5. Data Flow

1. Parse `source-image` and `destination-image`.
2. Validate both references.
3. Load local source image from daemon.
4. Build keychain:
   - Secret provided -> Kubernetes keychain
   - Secret omitted -> default/anonymous behavior
5. Push image to destination registry.
6. Emit progress bar updates.
7. Emit success summary.

## 6. Progress UX

Progress style: character progress bar.

Target format:

- `[=====>     ] 52% pushing layers`

Rules:

- Animated redraw in interactive terminal mode.
- Plain periodic status lines when not in TTY/interactive output.
- Clean finalization on both success and error.
- If exact byte progress is unavailable from transport hooks, use phase-based percentage with deterministic step mapping.

## 7. Error Handling

### Validation errors

- Invalid source or destination reference -> fail before network operations.

### Local source errors

- Local image not found -> fail-fast with clear recovery hint.

### Auth errors

- Secret lookup/keychain creation failures surfaced with operation context.

### Registry/network errors

- Reuse existing registry error normalization path.

### Progress lifecycle errors

- Stop progress rendering cleanly before printing final error message.

## 8. Testing Strategy

### Unit tests

- Flag contract and required arg behavior.
- Reference validation failures.
- Missing local image fail-fast path.
- Optional secret flow still functional.
- Non-TTY progress fallback output path.

### Integration-lite tests

- Happy path with controlled stubs/mocks for remote write.
- Auth failure mapping and message consistency.

### Regression focus

- No behavior change in shared auth path used by list.

## 9. Rollout and Compatibility

- Immediate switch to new push semantics (breaking change accepted).
- No backward-compatible alias for removed tar-source behavior.
- Update command help and docs for new contract.

## 10. Next Implementation Layers

Layer 2 (next): pull redesign for external registry to local daemon direction.

Layer 3 (later): revisit copy for shared source/target abstractions if implementation pressure justifies it.
