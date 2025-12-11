## Root Cause Summary
- Windows is returning "Insufficient system resources exist to complete the requested service" while `go.exe` tries to open a file inside the Go module cache (`C:\Users\USER\go\pkg\mod...`).
- Most common causes:
  1. Low disk space on `C:` (GOPATH/mod cache lives under `C:\Users\USER\go`).
  2. Memory or OS resource exhaustion during build (gopls compiles a large dependency graph and opens many files).
  3. OneDrive/AV interference with the mod cache directory (file virtualization or scanning blocking file hydration/open).
  4. Corrupted module cache paths causing file open errors.

## Quick Diagnostics (to confirm the cause)
1. Check free space on `C:` and ensure at least 5–10 GB free.
2. Inspect current Go env:
   - `go env GOPATH GOMODCACHE GOCACHE GOPROXY`
3. Verify whether `C:\Users\USER\go` is under OneDrive; if so, note it.
4. Check RAM availability and background processes consuming memory.
5. Review Defender/AV real-time scanning status for `C:\Users\USER\go` and `C:\Users\USER\AppData\Local\go-build`.

## Fix Plan
1. Free up disk space on `C:` (target ≥10 GB free).
2. Clean Go caches and module cache:
   - `go clean -cache -modcache -fuzzcache`.
3. Move module cache off `C:` to a non-OneDrive local path with space:
   - Create `E:\Go\modcache` and `E:\Go\gocache`.
   - `go env -w GOMODCACHE=E:\Go\modcache`
   - `go env -w GOCACHE=E:\Go\gocache`
4. Add exclusions in Windows Defender/AV for the new `GOMODCACHE` and `GOCACHE` directories (to avoid file-open hooks).
5. Ensure a stable proxy and fresh install:
   - `go env -w GOPROXY=https://proxy.golang.org,direct`
   - Re-run `go install golang.org/x/tools/gopls@latest`
6. If still failing, pin a known-compatible gopls version:
   - `go install golang.org/x/tools/gopls@v0.15.3`
7. As a fallback, download the official gopls Windows binary from the gopls releases and place it in PATH (`C:\Users\USER\bin` or `%GOBIN%`).

## Post-Validation
- Run `gopls version` and `gopls -rpc.trace -v check` on a small Go file to confirm operation.
- In Trae/VS Code, ensure Go extension picks up the installed `gopls` and the correct `GOROOT/GOPATH`.

## Why This Works
- Frees disk and system resources, avoids OneDrive/AV interference, relocates heavy caches to a safe local path, and resets corrupted cache states—addressing the OS-level file open failure.