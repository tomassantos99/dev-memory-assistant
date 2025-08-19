# Dev Memory Assistant

Simple tooling to help short memory devs.

## Tools

### Paste

Paste allows you to:
- Stores last 100 clipboard entries (text only)
- Provides an interface to search those entries and paste whichever is selected

How to use paste:
- Copy text (as many times as you want - only last 100 are stored)
- Press Ctrl + Shift + V to open the interface
- Use arrow keys to navigate up and down the copied text (or use the scroll wheel)
- Press Enter to paste the text (or double click it)
- You can use Ctrl + S to focus the search bar (or click it) and filter clipboard entries

## Running tools (Windows)

### Requirements

- Go
- winres (check building section on how to install with go)

### Building

Run (PowerShell):

```
./build.ps1
```

This automatically compiles windows resources and builds the executables (installs winres if not found). This also creates an `exe` directory with the executables.
The script builds the tool so it runs in the background. If you want to see it run on a console you can edit the script and remove `-ldflags="-H windowsgui"` from the build command.

Optionally, you can build it manually (example for paste tool):

```
go install github.com/tc-hib/go-winres@latest // install winres
cd ./paste && go-winres make // compile windows resources (manifest and icon)
go build -ldflags="-H windowsgui" -o ../exe/paste.exe // build, can remove the flag to see it run on the console
```

### Running

Run:

```
./paste.exe // (or double click the thing)
```
