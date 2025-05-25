# JavaCleaner

JavaCleaner is a Windows utility written in Go that detects, analyzes, and safely removes old or unused versions of Java installed on a system. It uses the Windows Registry, environment variables, and process inspection to determine which Java installations can be safely removed.

## Features

- Detects all Java installations via:
  - Windows registry
  - `JAVA_HOME` environment variable
  - Common install paths like `Program Files\Java`
- Evaluates usage by:
  - Scanning currently running processes
  - Checking Windows services
  - Inspecting scheduled tasks (excluding false positives like `JavaUpdateSched`)
- Attempts proper uninstallation using the registry's uninstall string
- Falls back to safe directory deletion if uninstall fails
- Detects if the tool is run with administrative privileges
- Interactive or automated operation
- JSON or human-readable output modes
- Logging to `javacleaner.log`

## Requirements

- Windows 10/11
- Run as Administrator for full functionality
- Built using Go 1.20+

## Usage

```sh
JavaCleaner.exe [--interactive] [--auto-clean] [--output] [--verbose]
```

### Options

| Flag           | Description                                         |
|----------------|-----------------------------------------------------|
| `--interactive`| Runs in menu-driven interactive mode                |
| `--auto-clean` | Automatically attempts to remove unused Java        |
| `--output`     | Outputs detailed JSON report                        |
| `--verbose`    | Enables verbose logging                             |

### Example (interactive)

```sh
JavaCleaner.exe --interactive
```

### Example (auto-clean with output report)

```sh
JavaCleaner.exe --auto-clean --output
```

## Logging

- Actions and errors are logged to `javacleaner.log` in the current directory.
- When uninstall fails, JavaCleaner shows the command used and suggests rerunning as Administrator if not elevated.

## Known Limitations

- Does not yet scan `.bat`, `.cmd`, `.lnk`, `.xml`, or `.conf` files for hardcoded Java references (planned feature).
- If the uninstall entry does not include a valid `InstallLocation`, matching may not occur.
- Uninstallation depends on third-party uninstaller behavior and may silently fail.

