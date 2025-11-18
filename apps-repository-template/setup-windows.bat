@echo off
REM ============================================================================
REM StumpfWorks NAS Apps Repository Setup for Windows
REM ============================================================================

echo.
echo ========================================================================
echo    StumpfWorks NAS Apps Repository Setup
echo ========================================================================
echo.

REM Check if Git is installed
git --version >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Git is not installed!
    echo.
    echo Please install Git from: https://git-scm.com/download/win
    echo.
    pause
    exit /b 1
)

REM Check if Python is installed (needed for scripts)
python --version >nul 2>&1
if errorlevel 1 (
    echo [WARNING] Python is not installed!
    echo Python is needed for validation scripts.
    echo Download from: https://www.python.org/downloads/
    echo.
    echo Continue anyway? (Press Ctrl+C to cancel)
    pause
)

REM Ask for target directory
echo Where do you want to create the apps repository?
echo Default: ..\stumpfworks-nas-apps
echo.
set /p TARGET_DIR="Enter path (or press Enter for default): "
if "%TARGET_DIR%"=="" set TARGET_DIR=..\stumpfworks-nas-apps

REM Check if directory exists
if exist "%TARGET_DIR%" (
    echo.
    echo [WARNING] Directory %TARGET_DIR% already exists!
    set /p CONFIRM="Delete and recreate? (y/N): "
    if /i "%CONFIRM%"=="y" (
        rmdir /s /q "%TARGET_DIR%"
    ) else (
        echo Aborted.
        pause
        exit /b 1
    )
)

echo.
echo [1/5] Creating directory: %TARGET_DIR%
mkdir "%TARGET_DIR%"

echo [2/5] Copying template files...
xcopy /E /I /Q /Y . "%TARGET_DIR%"

REM Remove files that shouldn't be in the new repo
echo [3/5] Cleaning up...
if exist "%TARGET_DIR%\README_TEMPLATE.md" del "%TARGET_DIR%\README_TEMPLATE.md"
if exist "%TARGET_DIR%\setup-apps-repo.sh" del "%TARGET_DIR%\setup-apps-repo.sh"
if exist "%TARGET_DIR%\setup-windows.bat" del "%TARGET_DIR%\setup-windows.bat"
if exist "%TARGET_DIR%\setup-windows.ps1" del "%TARGET_DIR%\setup-windows.ps1"
if exist "%TARGET_DIR%\MOVE_TO_NEW_REPO.md" del "%TARGET_DIR%\MOVE_TO_NEW_REPO.md"

REM Initialize Git
echo [4/5] Initializing Git repository...
cd /d "%TARGET_DIR%"
git init
git add .
git commit -m "Initial commit: StumpfWorks NAS Apps repository"

echo [5/5] Repository created!
echo.
echo ========================================================================
echo    SUCCESS! Repository is ready.
echo ========================================================================
echo.
echo Next steps:
echo.
echo 1. Create GitHub repository:
echo    - Go to: https://github.com/organizations/Stumpf-works/repositories/new
echo    - Name: stumpfworks-nas-apps
echo    - Make it Public
echo    - Click "Create repository"
echo.
echo 2. Push to GitHub (run these commands in Git Bash or PowerShell):
echo    cd %TARGET_DIR%
echo    git remote add origin https://github.com/Stumpf-works/stumpfworks-nas-apps.git
echo    git branch -M main
echo    git push -u origin main
echo.
echo 3. Add plugins:
echo    mkdir plugins
echo    (copy your plugins to plugins/)
echo.
echo 4. Generate registry:
echo    python scripts\generate-registry.py
echo.
echo 5. Read documentation:
echo    - START_HERE.md
echo    - HOW_TO_USE_PROMPTS.md
echo.
echo ========================================================================
echo.
pause
