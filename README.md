# Dotlink - Configuration File Linking Tool

A Go utility that scans for `dotlink.yaml` configuration files and creates symbolic links for configuration management.

## Features

✨ **Smart Config Parsing**
- YAML-based configuration format
- Support for directory and file-level linking
- Environment variable expansion (e.g., `$HOME/.zsh`)

🔗 **Flexible Linking**
- Link entire directories to target locations
- Link individual files with custom target paths
- Automatic parent directory creation
- Safe overwriting of existing symlinks

📊 **Enhanced Logging**
- Emoji-based visual feedback
- Debug-level logging for troubleshooting
- Information messages for key operations
- Error handling with clear messages

🧪 **Comprehensive Testing**
- Full test coverage of core functionality
- Tests for config parsing, path expansion, and linking
- Dry-run mode support for safe testing

## Configuration Format

Create a `dotlink.yaml` file in your configuration directory:

```yaml
link:
  to: $HOME/.zsh

files:
  - file: .zshrc
    to: $HOME/.zshrc
  - file: .zsh_profile
    to: $HOME/.zsh_profile
```

### Configuration Options

- **link.to**: Directory where the config directory will be linked
- **files**: List of individual files to link
  - **file**: Path relative to the dotlink.yaml directory
  - **to**: Target path where the file will be linked (supports environment variables)

## Usage

### Basic Usage

```bash
./dotlink
```

This will:
1. Scan `$DOTLINK_ROOT` (defaults to `$HOME/.custom`) recursively
2. Find all `dotlink.yaml` files
3. Create symbolic links according to each configuration

### Environment Variables

- `DOTLINK_ROOT`: Root directory to scan (default: `$HOME/.custom`)

## Implementation Details

### Core Components

#### `config.go`
- **DotlinkConfig**: Main configuration structure
- **LoadConfig()**: Parses YAML configuration files
- **ExpandPath()**: Expands environment variables in paths

#### `linker.go`
- **Linker**: Main linking orchestrator
- **LinkDirectory()**: Creates directory-level symlinks
- **LinkFile()**: Creates file-level symlinks
- **ProcessConfig()**: Orchestrates the full linking process

#### `main.go`
- Scans for dotlink.yaml files
- Iterates and processes each configuration
- Provides user-friendly emoji-based logging

### Logging with Emojis

The tool uses emojis to provide clear visual feedback:

- 🔍 Scanning operations
- ✨ Found files
- 📂 File system operations
- 📁 Directory linking
- 📄 File operations
- 🔗 Successful links created
- 📖 Config processing
- ✅ Successful completion
- ❌ Errors
- ⚠️ Warnings

## Testing

Run the test suite:

```bash
go test -v
```

### Test Coverage

- `TestLoadConfig`: YAML parsing
- `TestGetConfigDir`: Configuration directory extraction
- `TestExpandPath`: Environment variable expansion
- `TestLinkerLinkFile`: File symlinking
- `TestLinkerLinkDirectory`: Directory symlinking
- `TestLinkerDryRun`: Dry-run mode verification
- `TestProcessConfig`: Full configuration processing
- `TestLoadConfigInvalidYAML`: Error handling

## Example Workflow

1. Create a configuration directory structure:
   ```
   ~/.custom/zsh/
   ├── dotlink.yaml
   ├── .zshrc
   └── .zsh_profile
   ```

2. Create `dotlink.yaml`:
   ```yaml
   link:
     to: $HOME/.zsh
   files:
     - file: .zshrc
       to: $HOME/.zshrc
     - file: .zsh_profile
       to: $HOME/.zsh_profile
   ```

3. Run dotlink:
   ```bash
   DOTLINK_ROOT=$HOME/.custom ./dotlink
   ```

4. Result:
   - `$HOME/.zsh` → symlink to `~/.custom/zsh`
   - `$HOME/.zshrc` → symlink to `~/.custom/zsh/.zshrc`
   - `$HOME/.zsh_profile` → symlink to `~/.custom/zsh/.zsh_profile`

## Dependencies

- `gopkg.in/yaml.v3`: YAML parsing
- `github.com/charmbracelet/log`: Enhanced logging with colors and emojis

## Building

```bash
go build -o dotlink
```

## Error Handling

The tool handles various error conditions gracefully:
- Missing source files
- Invalid YAML configuration
- File system permission issues
- Existing symlink conflicts

All errors are logged with context to help with troubleshooting.

