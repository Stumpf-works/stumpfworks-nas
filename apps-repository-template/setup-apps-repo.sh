#!/bin/bash
set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   StumpfWorks NAS Apps Repository Setup                   ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""

# Check if we're in the template directory
if [ ! -f "START_HERE.md" ]; then
    echo -e "${YELLOW}⚠️  Please run this script from the apps-repository-template directory${NC}"
    exit 1
fi

# Ask for target directory
echo -e "${BLUE}Where do you want to create the apps repository?${NC}"
read -p "Path (default: ../stumpfworks-nas-apps): " TARGET_DIR
TARGET_DIR=${TARGET_DIR:-../stumpfworks-nas-apps}

# Create target directory
if [ -d "$TARGET_DIR" ]; then
    echo -e "${YELLOW}⚠️  Directory $TARGET_DIR already exists!${NC}"
    read -p "Delete and recreate? (y/N): " CONFIRM
    if [ "$CONFIRM" = "y" ] || [ "$CONFIRM" = "Y" ]; then
        rm -rf "$TARGET_DIR"
    else
        echo "Aborted."
        exit 1
    fi
fi

echo -e "${GREEN}📁 Creating directory: $TARGET_DIR${NC}"
mkdir -p "$TARGET_DIR"

# Copy template files
echo -e "${GREEN}📋 Copying template files...${NC}"
cp -r ./* "$TARGET_DIR/"

# Remove the README_TEMPLATE.md and this script from target
rm -f "$TARGET_DIR/README_TEMPLATE.md"
rm -f "$TARGET_DIR/setup-apps-repo.sh"

# Initialize git
echo -e "${GREEN}🔧 Initializing git repository...${NC}"
cd "$TARGET_DIR"
git init

# Create .gitignore if not exists
if [ ! -f ".gitignore" ]; then
    cat > .gitignore <<EOF
# Python
__pycache__/
*.py[cod]
*$py.class
*.so
.Python
env/
venv/

# OS
.DS_Store
Thumbs.db

# Editor
.vscode/
.idea/
*.swp
*.swo
*~

# Temporary
tmp/
temp/
*.tmp
EOF
fi

# Initial commit
echo -e "${GREEN}📝 Creating initial commit...${NC}"
git add .
git commit -m "Initial commit: StumpfWorks NAS Apps repository structure"

echo ""
echo -e "${GREEN}✅ Repository created successfully!${NC}"
echo ""
echo -e "${BLUE}Next steps:${NC}"
echo ""
echo "1. Create GitHub repository:"
echo "   https://github.com/organizations/Stumpf-works/repositories/new"
echo "   Name: stumpfworks-nas-apps"
echo ""
echo "2. Push to GitHub:"
echo "   cd $TARGET_DIR"
echo "   git remote add origin https://github.com/Stumpf-works/stumpfworks-nas-apps.git"
echo "   git branch -M main"
echo "   git push -u origin main"
echo ""
echo "3. Add plugins:"
echo "   mkdir -p plugins"
echo "   cp -r /path/to/plugin plugins/my-plugin"
echo ""
echo "4. Generate registry:"
echo "   python3 scripts/generate-registry.py"
echo ""
echo "5. Read documentation:"
echo "   cat START_HERE.md"
echo ""
echo -e "${GREEN}🎉 Happy plugin development!${NC}"
