#!/bin/bash

# Stumpf.Works NAS Update Script
# This script pulls the latest changes from git and restarts the server

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== Stumpf.Works NAS Update Script ===${NC}"
echo ""

# Check if we're in a git repository
if [ ! -d ".git" ]; then
    echo -e "${RED}Error: Not in a git repository${NC}"
    exit 1
fi

# Get current branch
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
echo -e "Current branch: ${YELLOW}${CURRENT_BRANCH}${NC}"

# Get current commit
CURRENT_COMMIT=$(git rev-parse --short HEAD)
echo -e "Current commit: ${YELLOW}${CURRENT_COMMIT}${NC}"
echo ""

# Fetch latest changes
echo -e "${GREEN}Fetching latest changes...${NC}"
git fetch origin

# Check if there are updates
LATEST_COMMIT=$(git rev-parse --short origin/${CURRENT_BRANCH})
if [ "$CURRENT_COMMIT" = "$LATEST_COMMIT" ]; then
    echo -e "${GREEN}Already up to date!${NC}"
    exit 0
fi

# Show how many commits behind
COMMITS_BEHIND=$(git rev-list --count HEAD..origin/${CURRENT_BRANCH})
echo -e "${YELLOW}${COMMITS_BEHIND} commit(s) behind${NC}"
echo ""

# Show changelog
echo -e "${GREEN}Changes:${NC}"
git log --oneline HEAD..origin/${CURRENT_BRANCH}
echo ""

# Confirm update
read -p "Do you want to update? (y/N) " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}Update cancelled${NC}"
    exit 0
fi

# Pull changes
echo -e "${GREEN}Pulling changes...${NC}"
git pull origin ${CURRENT_BRANCH}

# Check if backend was modified
if git diff --name-only ${CURRENT_COMMIT}..HEAD | grep -q "^backend/"; then
    echo -e "${YELLOW}Backend files changed. Rebuilding backend...${NC}"
    cd backend
    go build -o ../bin/stumpfworks-server ./cmd/stumpfworks-server
    cd ..
    echo -e "${GREEN}Backend rebuilt successfully${NC}"
fi

# Check if frontend was modified
if git diff --name-only ${CURRENT_COMMIT}..HEAD | grep -q "^frontend/"; then
    echo -e "${YELLOW}Frontend files changed. Rebuilding frontend...${NC}"
    cd frontend
    npm install
    npm run build
    cd ..
    echo -e "${GREEN}Frontend rebuilt successfully${NC}"
fi

# Ask to restart service
echo ""
echo -e "${GREEN}Update completed successfully!${NC}"
echo -e "New commit: ${YELLOW}$(git rev-parse --short HEAD)${NC}"
echo ""

if systemctl is-active --quiet stumpfworks-nas 2>/dev/null; then
    read -p "Restart stumpfworks-nas service? (y/N) " -n 1 -r
    echo ""
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${GREEN}Restarting service...${NC}"
        sudo systemctl restart stumpfworks-nas
        echo -e "${GREEN}Service restarted${NC}"
    fi
else
    echo -e "${YELLOW}Service not running. Please start it manually.${NC}"
fi

echo ""
echo -e "${GREEN}Update complete!${NC}"
