#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Find Your Food - Test Suite${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Check if API server is running
echo -e "${YELLOW}📡 Checking API server...${NC}"
if curl -s http://localhost:8080/health > /dev/null; then
    echo -e "${GREEN}✅ API server is running${NC}"
else
    echo -e "${RED}❌ API server is not running. Please start it with:${NC}"
    echo -e "   go run cmd/api/main.go"
    exit 1
fi

echo ""
echo -e "${YELLOW}🧪 Testing Public Endpoints...${NC}"
echo ""

# Test 1: Get all categories
echo -e "${BLUE}1. GET /api/v1/public/categories${NC}"
curl -s http://localhost:8080/api/v1/public/categories | jq '.status, .data | length'
echo ""

# Test 2: Search foods
echo -e "${BLUE}2. GET /api/v1/public/foods/search?q=nasi${NC}"
curl -s "http://localhost:8080/api/v1/public/foods/search?q=nasi&limit=5" | jq '.status, .count'
echo ""

# Test 3: Get foods by category (MP - Makanan Pokok)
echo -e "${BLUE}3. GET /api/v1/public/categories/MP/foods${NC}"
curl -s "http://localhost:8080/api/v1/public/categories/MP/foods?limit=5" | jq '.status, .count'
echo ""

# Test 4: Get food detail (first food from search)
FOOD_ID=$(curl -s "http://localhost:8080/api/v1/public/foods/search?q=nasi&limit=1" | jq -r '.data[0].id')
if [ "$FOOD_ID" != "null" ] && [ -n "$FOOD_ID" ]; then
    echo -e "${BLUE}4. GET /api/v1/public/foods/${FOOD_ID}${NC}"
    curl -s "http://localhost:8080/api/v1/public/foods/$FOOD_ID" | jq '.status, .data.code, .data.name, (.data.portion_photos | length)'
    echo ""
else
    echo -e "${RED}❌ No food found in search results${NC}"
fi

echo ""
echo -e "${GREEN}✅ Public API tests completed!${NC}"
echo ""

# Test WebSocket
echo -e "${YELLOW}🔌 Testing WebSocket Collaboration...${NC}"
echo ""

# Create test users
echo -e "${BLUE}Creating test users...${NC}"

USER1_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Researcher 1",
    "email": "test.researcher1@example.com",
    "password": "test12345",
    "role": "researcher"
  }')

USER2_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Researcher 2",
    "email": "test.researcher2@example.com",
    "password": "test12345",
    "role": "researcher"
  }')

# Login both users
echo -e "${BLUE}Logging in users...${NC}"

USER1_TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test.researcher1@example.com",
    "password": "test12345"
  }' | jq -r '.data.access_token // empty')

USER2_TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test.researcher2@example.com",
    "password": "test12345"
  }' | jq -r '.data.access_token // empty')

if [ -n "$USER1_TOKEN" ] && [ -n "$USER2_TOKEN" ]; then
    echo -e "${GREEN}✅ Users logged in successfully${NC}"
    echo ""
    echo -e "${YELLOW}📋 WebSocket Connection Details:${NC}"
    echo -e "   User 1 Token: ${USER1_TOKEN:0:50}..."
    echo -e "   User 2 Token: ${USER2_TOKEN:0:50}..."
    echo ""
    echo -e "${YELLOW}🧪 To test WebSocket collaboration, run in 2 terminals:${NC}"
    echo ""
    echo -e "${BLUE}Terminal 1:${NC}"
    echo -e "   wscat -c \"ws://localhost:8080/ws/collab?token=$USER1_TOKEN&room=test-room-001\""
    echo ""
    echo -e "${BLUE}Terminal 2:${NC}"
    echo -e "   wscat -c \"ws://localhost:8080/ws/collab?token=$USER2_TOKEN&room=test-room-001\""
    echo ""
    echo -e "${YELLOW}Then send messages like:${NC}"
    echo -e '   {"type":"food_search","payload":{"query":"nasi goreng"}}'
    echo -e '   {"type":"cursor_move","payload":{"x":500,"y":300,"page":"/find-food"}}'
else
    echo -e "${RED}❌ Failed to login users${NC}"
fi

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}✅ Test suite completed!${NC}"
echo -e "${BLUE}========================================${NC}"
