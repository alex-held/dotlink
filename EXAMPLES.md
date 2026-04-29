# Example dotlink.yaml Configuration

This directory demonstrates the `dotlink.yaml` configuration format and how to use dotlink.

## Directory Structure

```
my-dotfiles/
├── zsh/
│   ├── dotlink.yaml
│   ├── .zshrc
│   ├── .zsh_profile
│   └── .zshenv
├── vim/
│   ├── dotlink.yaml
│   └── .vimrc
└── git/
    ├── dotlink.yaml
    └── .gitconfig
```

## Zsh Configuration Example

**File: zsh/dotlink.yaml**

```yaml
link:
  to: $HOME/.zsh

files:
  - file: .zshrc
    to: $HOME/.zshrc
  - file: .zsh_profile
    to: $HOME/.zsh_profile
  - file: .zshenv
    to: $HOME/.zshenv
```

This configuration:
1. Creates a symbolic link from `$HOME/.zsh` to the `zsh` directory
2. Links `.zshrc` in the `zsh` directory to `$HOME/.zshrc`
3. Links `.zsh_profile` in the `zsh` directory to `$HOME/.zsh_profile`
4. Links `.zshenv` in the `zsh` directory to `$HOME/.zshenv`

## Vim Configuration Example

**File: vim/dotlink.yaml**

```yaml
files:
  - file: .vimrc
    to: $HOME/.vimrc
```

This configuration only links individual files without linking the directory itself.

## Git Configuration Example

**File: git/dotlink.yaml**

```yaml
link:
  to: $HOME/.git-config

files:
  - file: .gitconfig
    to: $HOME/.gitconfig
```

## Running Dotlink

### Standard Run

```bash
export DOTLINK_ROOT=/path/to/my-dotfiles
./dotlink
```

### Dry-Run Mode

To test without creating actual symlinks, modify `main.go` line:
```go
linker := NewLinker(true) // Set to true for dry-run mode
```

Then rebuild:
```bash
go build -o dotlink
./dotlink
```

## Environment Variables Used

- `$HOME`: User's home directory
- `$USER`: Current username
- Any other environment variables can be used in the configuration

## Output Example

```
🔍 Scanning root directory: /home/user/my-dotfiles
📋 Found 3 dotlink configuration(s)
📖 Processing config: /home/user/my-dotfiles/zsh/dotlink.yaml
📁 Linking directory to: $HOME/.zsh
🔗 Linked directory: /home/user/my-dotfiles/zsh -> /home/user/.zsh
📄 Linking 3 files
🔗 Linked file: /home/user/my-dotfiles/zsh/.zshrc -> /home/user/.zshrc
🔗 Linked file: /home/user/my-dotfiles/zsh/.zsh_profile -> /home/user/.zsh_profile
🔗 Linked file: /home/user/my-dotfiles/zsh/.zshenv -> /home/user/.zshenv
✅ Successfully processed config: /home/user/my-dotfiles/zsh/dotlink.yaml
...
✅ Dotlink completed successfully
```

## Notes

- Directory links (`link.to`) link the entire configuration directory to the target
- File links link individual files from the configuration directory
- All paths support environment variable expansion
- Parent directories are automatically created if they don't exist
- Existing symlinks are safely overwritten
- Use debug logging to troubleshoot issues:

```go
logger.SetLevel(log.DebugLevel)
```

