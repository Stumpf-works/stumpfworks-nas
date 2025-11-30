#!/bin/bash

# Release script for StumpfWorks NAS
# Usage: ./release.sh [version]
# Example: ./release.sh 1.4.0

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if version is provided
if [ -z "$1" ]; then
    echo -e "${RED}Error: Version number required${NC}"
    echo "Usage: ./release.sh [version]"
    echo "Example: ./release.sh 1.4.0"
    exit 1
fi

VERSION="$1"
TAG="v${VERSION}"

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}  StumpfWorks NAS - Release Script${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${YELLOW}Creating release: ${TAG}${NC}"
echo ""

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo -e "${RED}Error: Not in a git repository${NC}"
    exit 1
fi

# Check for uncommitted changes
if ! git diff-index --quiet HEAD --; then
    echo -e "${RED}Error: You have uncommitted changes${NC}"
    echo "Please commit or stash your changes first"
    exit 1
fi

# Confirm with user
echo -e "${YELLOW}This will:${NC}"
echo "  1. Switch to main branch"
echo "  2. Pull latest changes"
echo "  3. Merge development branch"
echo "  4. Create tag ${TAG}"
echo "  5. Push to GitHub"
echo ""
read -p "Continue? (y/n) " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}Release cancelled${NC}"
    exit 0
fi

echo ""
echo -e "${YELLOW}Step 1/5: Switching to main branch...${NC}"
git checkout main
echo -e "${GREEN}✓ On main branch${NC}"
echo ""

echo -e "${YELLOW}Step 2/5: Pulling latest changes...${NC}"
git pull origin main
echo -e "${GREEN}✓ Main branch updated${NC}"
echo ""

echo -e "${YELLOW}Step 3/5: Merging development branch...${NC}"
git merge development --no-ff -m "Merge development for release ${TAG}"
echo -e "${GREEN}✓ Development merged${NC}"
echo ""

echo -e "${YELLOW}Step 4/5: Creating tag ${TAG}...${NC}"
git tag ${TAG} -m "Release ${VERSION}

Performance Optimizations:
- Gzip compression (60-80% smaller payloads)
- Code splitting (70% smaller initial bundle)
- Database indexes (5-10x faster queries)
- API caching (40-100x faster responses)
- React.memo (50% fewer re-renders)"

echo -e "${GREEN}✓ Tag created${NC}"
echo ""

echo -e "${YELLOW}Step 5/5: Pushing to GitHub...${NC}"
git push origin main
git push origin ${TAG}
echo -e "${GREEN}✓ Pushed to GitHub${NC}"
echo ""

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✅ Release ${TAG} completed successfully!${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${GREEN}Next steps:${NC}"
echo "  • Stable build will run automatically at 20:00"
echo "  • Or trigger manually: ssh root@46.4.25.15 'bash /root/stumpfworks-nas-build/scripts/auto-build-stable.sh'"
echo "  • Monitor Discord for build notifications"
echo ""
