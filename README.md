# Dev Memory Assistant

Simple tooling to help short memory devs

## Running the paste tool (Windows)

Build the thing (you need go):

```
go build ./paste
```

or

```
go build -ldflags="-H windowsgui" ./paste
```

Adding the flags will make it so no console window pops up when running

Run the thing:

```
./paste.exe
```
