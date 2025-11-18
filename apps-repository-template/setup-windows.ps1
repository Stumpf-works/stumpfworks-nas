# ============================================================================
# StumpfWorks NAS Apps Repository Setup for Windows (PowerShell)
# ============================================================================

$ErrorActionPreference = "Stop"

Write-Host ""
Write-Host "========================================================================" -ForegroundColor Cyan
Write-Host "   StumpfWorks NAS Apps Repository Setup" -ForegroundColor Cyan
Write-Host "========================================================================" -ForegroundColor Cyan
Write-Host ""

# Check if Git is installed
try {
    $gitVersion = git --version
    Write-Host "[OK] Git is installed: $gitVersion" -ForegroundColor Green
} catch {
    Write-Host "[ERROR] Git is not installed!" -ForegroundColor Red
    Write-Host ""
    Write-Host "Please install Git from: https://git-scm.com/download/win"
    Write-Host ""
    Read-Host "Press Enter to exit"
    exit 1
}

# Check if Python is installed
try {
    $pythonVersion = python --version
    Write-Host "[OK] Python is installed: $pythonVersion" -ForegroundColor Green
} catch {
    Write-Host "[WARNING] Python is not installed!" -ForegroundColor Yellow
    Write-Host "Python is needed for validation scripts."
    Write-Host "Download from: https://www.python.org/downloads/"
    Write-Host ""
    $continue = Read-Host "Continue anyway? (y/N)"
    if ($continue -ne "y" -and $continue -ne "Y") {
        exit 1
    }
}

# Ask for target directory
Write-Host ""
Write-Host "Where do you want to create the apps repository?"
Write-Host "Default: ..\stumpfworks-nas-apps"
Write-Host ""
$targetDir = Read-Host "Enter path (or press Enter for default)"
if ([string]::IsNullOrWhiteSpace($targetDir)) {
    $targetDir = "..\stumpfworks-nas-apps"
}

# Resolve to absolute path
$targetDir = Resolve-Path $targetDir -ErrorAction SilentlyContinue
if (-not $targetDir) {
    $targetDir = Join-Path (Get-Location).Path "..\stumpfworks-nas-apps"
}

# Check if directory exists
if (Test-Path $targetDir) {
    Write-Host ""
    Write-Host "[WARNING] Directory $targetDir already exists!" -ForegroundColor Yellow
    $confirm = Read-Host "Delete and recreate? (y/N)"
    if ($confirm -eq "y" -or $confirm -eq "Y") {
        Remove-Item -Recurse -Force $targetDir
    } else {
        Write-Host "Aborted."
        Read-Host "Press Enter to exit"
        exit 1
    }
}

# Create directory and copy files
Write-Host ""
Write-Host "[1/5] Creating directory: $targetDir" -ForegroundColor Cyan
New-Item -ItemType Directory -Path $targetDir | Out-Null

Write-Host "[2/5] Copying template files..." -ForegroundColor Cyan
Copy-Item -Path ".\*" -Destination $targetDir -Recurse -Force

# Remove files that shouldn't be in the new repo
Write-Host "[3/5] Cleaning up..." -ForegroundColor Cyan
$filesToRemove = @(
    "README_TEMPLATE.md",
    "setup-apps-repo.sh",
    "setup-windows.bat",
    "setup-windows.ps1",
    "MOVE_TO_NEW_REPO.md"
)

foreach ($file in $filesToRemove) {
    $filePath = Join-Path $targetDir $file
    if (Test-Path $filePath) {
        Remove-Item $filePath -Force
    }
}

# Initialize Git
Write-Host "[4/5] Initializing Git repository..." -ForegroundColor Cyan
Set-Location $targetDir
git init
git add .
git commit -m "Initial commit: StumpfWorks NAS Apps repository"

Write-Host "[5/5] Repository created!" -ForegroundColor Green
Write-Host ""
Write-Host "========================================================================" -ForegroundColor Green
Write-Host "   SUCCESS! Repository is ready." -ForegroundColor Green
Write-Host "========================================================================" -ForegroundColor Green
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host ""
Write-Host "1. Create GitHub repository:"
Write-Host "   - Go to: https://github.com/organizations/Stumpf-works/repositories/new"
Write-Host "   - Name: stumpfworks-nas-apps"
Write-Host "   - Make it Public"
Write-Host "   - Click 'Create repository'"
Write-Host ""
Write-Host "2. Push to GitHub:"
Write-Host "   cd $targetDir" -ForegroundColor Cyan
Write-Host "   git remote add origin https://github.com/Stumpf-works/stumpfworks-nas-apps.git" -ForegroundColor Cyan
Write-Host "   git branch -M main" -ForegroundColor Cyan
Write-Host "   git push -u origin main" -ForegroundColor Cyan
Write-Host ""
Write-Host "3. Add plugins:"
Write-Host "   mkdir plugins" -ForegroundColor Cyan
Write-Host "   (copy your plugins to plugins/)" -ForegroundColor Cyan
Write-Host ""
Write-Host "4. Generate registry:"
Write-Host "   python scripts\generate-registry.py" -ForegroundColor Cyan
Write-Host ""
Write-Host "5. Read documentation:"
Write-Host "   - START_HERE.md"
Write-Host "   - HOW_TO_USE_PROMPTS.md"
Write-Host ""
Write-Host "========================================================================" -ForegroundColor Green
Write-Host ""
Read-Host "Press Enter to exit"
