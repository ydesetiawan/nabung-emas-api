# VS Code Go Navigation - Quick Fix Guide

## ✅ Setup Complete!

All Go tools have been installed and configured. Your VS Code should now support:
- ✅ Go to Definition (Cmd+Click or F12)
- ✅ Find All References (Shift+F12)
- ✅ Auto-completion
- ✅ Hover documentation
- ✅ Code formatting

## 🔧 How to Enable Go-to-Definition

### Step 1: Reload VS Code Window
1. Press `Cmd+Shift+P` (Command Palette)
2. Type: `Reload Window`
3. Press Enter

### Step 2: Wait for gopls to Index
- Look at the **bottom right** of VS Code
- You'll see "gopls: indexing..." 
- Wait until it says "gopls: ready" or disappears
- This may take 30-60 seconds for the first time

### Step 3: Test Navigation
Try these shortcuts on any function or struct:
- **Cmd+Click** - Go to definition
- **F12** - Go to definition
- **Shift+F12** - Find all references
- **Cmd+T** - Go to symbol in workspace
- **Cmd+Shift+O** - Go to symbol in file

## 🎯 Quick Test

Open any file and try:
1. Go to `internal/handlers/pocket_handler.go`
2. Click on `PocketService` (line 11)
3. Press **F12** or **Cmd+Click**
4. It should jump to `internal/services/pocket_service.go`

## 🔍 Troubleshooting

### If go-to-definition still doesn't work:

#### Option 1: Install/Update Go Tools
1. Press `Cmd+Shift+P`
2. Type: `Go: Install/Update Tools`
3. Select **ALL** tools
4. Click **OK**
5. Wait for installation to complete
6. Reload window

#### Option 2: Restart Go Language Server
1. Press `Cmd+Shift+P`
2. Type: `Go: Restart Language Server`
3. Press Enter
4. Wait for gopls to re-index

#### Option 3: Check Go Extension
1. Make sure **Go extension** is installed
2. Press `Cmd+Shift+X` (Extensions)
3. Search for "Go"
4. Install "Go" by Go Team at Google
5. Reload window

#### Option 4: Verify Go Installation
```bash
# Check Go is installed
go version

# Check gopls is installed
gopls version

# Check GOPATH
go env GOPATH
```

## 📚 Useful VS Code Shortcuts

| Shortcut | Action |
|----------|--------|
| `Cmd+Click` | Go to definition |
| `F12` | Go to definition |
| `Shift+F12` | Find all references |
| `Cmd+T` | Go to symbol in workspace |
| `Cmd+Shift+O` | Go to symbol in file |
| `Cmd+P` | Quick file open |
| `Cmd+Shift+F` | Search in all files |
| `F2` | Rename symbol |
| `Shift+F12` | Peek references |

## 🎨 VS Code Settings Applied

The following settings have been configured in `.vscode/settings.json`:

```json
{
  "go.useLanguageServer": true,
  "go.buildOnSave": "workspace",
  "go.lintOnSave": "workspace",
  "go.autocompleteUnimportedPackages": true,
  "[go]": {
    "editor.formatOnSave": true,
    "editor.codeActionsOnSave": {
      "source.organizeImports": "explicit"
    }
  },
  "gopls": {
    "ui.semanticTokens": true,
    "ui.completion.usePlaceholders": true
  }
}
```

## ✨ Additional Features Now Available

1. **Auto-formatting** - Code formats on save
2. **Auto-imports** - Imports organize automatically
3. **IntelliSense** - Smart code completion
4. **Hover docs** - Hover over any function to see documentation
5. **Error detection** - Real-time error highlighting
6. **Refactoring** - Rename symbols across files

## 🚀 Next Steps

1. **Reload VS Code** (Cmd+Shift+P → Reload Window)
2. **Wait for indexing** (check bottom right)
3. **Test navigation** (Cmd+Click on any function)
4. **Start coding!** 🎉

---

**Note:** If you still have issues after reloading, try closing and reopening VS Code completely.
