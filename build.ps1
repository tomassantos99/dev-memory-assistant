# build.ps1

# Config
$apps = @(
    @{ Name = "paste"; Path = "paste"; NeedsRes = $true }
)

# Ensure go-winres is installed if any app needs it
$needsWinRes = $apps | Where-Object { $_.NeedsRes } | Measure-Object | Select-Object -ExpandProperty Count
if ($needsWinRes -gt 0) {
    if (-not (Get-Command go-winres -ErrorAction SilentlyContinue)) {
        Write-Host "go-winres not found. Installing..."
        go install github.com/tc-hib/go-winres@latest
    }
}

foreach ($app in $apps) {
    Write-Host "Building $($app.Name)..."

    # Build resources only if needed
    if ($app.NeedsRes) {
        $resDir = Join-Path $app.Path "winres"
        if (Test-Path $resDir) {
            Write-Host "Generating resources for $($app.Name)..."
            Push-Location $app.Path
            go-winres make
            Pop-Location
        }
    }

    # Build the Go executable
    $exePath = "./exe/$($app.Name).exe"
    Write-Host "Compiling $($app.Name) to $exePath..."
    go build -o $exePath "./$($app.Path)"
}

Write-Host "All builds completed."
