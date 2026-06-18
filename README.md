# MailTap

<picture>
   <img alt="Logo for Eloquent Salesforce Objects" src="art/screenshot.jpg">
</picture>

Application to capture and manage email messages for developers.

## Installation

Grab the latest build from the [Releases page](https://github.com/daikazu/mailtap/releases).

> **Heads up:** builds are **not code-signed** (this is a free, hobby project with no paid
> Apple/Microsoft developer account), so each OS shows a one-time warning the first time you
> open the app. This is expected and safe to bypass — instructions per platform below.

### macOS

1. Download `MailTap-<version>-universal.dmg` (runs on both Intel and Apple Silicon).
2. Open the `.dmg` and drag **MailTap** into your Applications folder.
3. On first launch macOS will say the app is from an "unidentified developer." Either:
   - **right-click** the app → **Open** → **Open** in the dialog, or
   - run `xattr -dr com.apple.quarantine /Applications/mailtap.app`

After the first launch it opens normally like any other app.

### Windows

1. Download `mailtap-<version>-windows-amd64.exe`.
2. Run it. Windows SmartScreen may show "Windows protected your PC."
3. Click **More info** → **Run anyway** (only needed once).

Requires the [WebView2 runtime](https://developer.microsoft.com/microsoft-edge/webview2/),
which is preinstalled on Windows 11 and most Windows 10 systems.

### Build from source

Requires [Go](https://go.dev) 1.24+, [Node.js](https://nodejs.org) 18+, and the
[Wails CLI](https://wails.io/docs/gettingstarted/installation):

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@v2.12.0
wails build
```

The compiled app lands in `build/bin/`.


