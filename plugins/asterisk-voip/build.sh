#!/bin/bash
set -e

echo "🏗️  Building Asterisk VoIP Plugin..."
echo ""

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Change to backend directory
cd backend

echo -e "${BLUE}📦 Downloading Go dependencies...${NC}"
go mod download

echo ""
echo -e "${BLUE}🔨 Building binary...${NC}"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o asterisk-manager .

echo ""
echo -e "${BLUE}📝 Build info:${NC}"
ls -lh asterisk-manager
file asterisk-manager

# Move binary to plugin root
mv asterisk-manager ..

cd ..

echo ""
echo -e "${GREEN}✅ Build complete!${NC}"
echo ""
echo "Binary: ./asterisk-manager"
echo ""
echo "Next steps:"
echo "  1. Test locally:     ./asterisk-manager"
echo "  2. Start Docker:     docker-compose up -d"
echo "  3. View logs:        docker-compose logs -f"
echo "  4. Test API:         curl http://localhost:8090/health"
echo ""
