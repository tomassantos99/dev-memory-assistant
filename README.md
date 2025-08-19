# Dev Memory Assistant

Simple tooling to help short memory devs

## Running the paste tool (Windows)

### Requirements

- Go
- winres (check building section on how to install with go)

### Building

Run:

```
./build.ps1
```

This automatically compiles windows resources and builds the executables. This also creates an `exe` directory with the executables.
The script builds the tool so it runs in the background. If you want to see it run on a console you can edit the script and remove `-ldflags="-H windowsgui"` from the build command.

Optionally, you can build it manually:

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
