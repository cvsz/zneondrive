# Windows Build and Unreal Engine Transfer Tools

These scripts replace accidentally tracked files whose repository paths contained a Windows absolute path. Scripts must not contain account-specific absolute paths, private endpoints, tokens, or credentials.

## Security requirements

- Obtain Unreal Engine archives only from an authorized source.
- Verify SHA-256 against a value obtained through an independent trusted channel.
- HTTPS is required by default; an HTTP exception is permitted only for isolated transfer networks and still requires hash verification.
- Never overwrite or remove an existing Unreal Engine installation automatically.
- Do not commit access tokens, signed URLs, passwords, or private endpoints.
- Installer scripts require explicit operator authorization.

## Validation

Run Go checks from the repository root using `tools/windows/build-game-api.ps1`. Verify Visual Studio using `tools/windows/test-vs-installation.ps1`.
