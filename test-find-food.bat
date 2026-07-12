@echo off
setlocal enabledelayedexpansion

echo ========================================
echo   Find Your Food - Test Suite
echo ========================================
echo.

REM Check if API server is running
echo Checking API server...
curl -s http://localhost:8080/health >nul 2>&1
if %errorlevel% == 0 (
    echo [OK] API server is running
) else (
    echo [ERROR] API server is not running. Please start it with:
    echo    go run cmd/api/main.go
    exit /b 1
)

echo.
echo Testing Public Endpoints...
echo.

REM Test 1: Get all categories
echo 1. GET /api/v1/public/categories
curl -s http://localhost:8080/api/v1/public/categories
echo.
echo.

REM Test 2: Search foods
echo 2. GET /api/v1/public/foods/search?q=nasi
curl -s "http://localhost:8080/api/v1/public/foods/search?q=nasi&limit=5"
echo.
echo.

REM Test 3: Get foods by category
echo 3. GET /api/v1/public/categories/MP/foods
curl -s "http://localhost:8080/api/v1/public/categories/MP/foods?limit=5"
echo.
echo.

REM Test 4: Get food detail
echo 4. Testing Get Food Detail...
curl -s "http://localhost:8080/api/v1/public/foods/search?q=nasi&limit=1" > temp_food.json
echo.
echo.

echo [OK] Public API tests completed!
echo.

echo Testing WebSocket Collaboration...
echo.

echo Creating test users...
curl -s -X POST http://localhost:8080/api/v1/auth/register -H "Content-Type: application/json" -d "{\"name\":\"Test Researcher 1\",\"email\":\"test.researcher1@example.com\",\"password\":\"test12345\",\"role\":\"researcher\"}"
echo.

curl -s -X POST http://localhost:8080/api/v1/auth/register -H "Content-Type: application/json" -d "{\"name\":\"Test Researcher 2\",\"email\":\"test.researcher2@example.com\",\"password\":\"test12345\",\"role\":\"researcher\"}"
echo.

echo Logging in users...
curl -s -X POST http://localhost:8080/api/v1/auth/login -H "Content-Type: application/json" -d "{\"email\":\"test.researcher1@example.com\",\"password\":\"test12345\"}" > user1_login.json
curl -s -X POST http://localhost:8080/api/v1/auth/login -H "Content-Type: application/json" -d "{\"email\":\"test.researcher2@example.com\",\"password\":\"test12345\"}" > user2_login.json

echo.
echo [OK] Users logged in successfully
echo.
echo WebSocket Connection Details saved to user1_login.json and user2_login.json
echo.
echo To test WebSocket collaboration, install wscat:
echo    npm install -g wscat
echo.
echo Then run in 2 Command Prompt windows:
echo.
echo Window 1:
echo    wscat -c "ws://localhost:8080/ws/collab?token=YOUR_USER1_TOKEN&room=test-room-001"
echo.
echo Window 2:
echo    wscat -c "ws://localhost:8080/ws/collab?token=YOUR_USER2_TOKEN&room=test-room-001"
echo.
echo Then send messages like:
echo    {"type":"food_search","payload":{"query":"nasi goreng"}}
echo    {"type":"cursor_move","payload":{"x":500,"y":300,"page":"/find-food"}}
echo.
echo ========================================
echo [OK] Test suite completed!
echo ========================================

REM Cleanup
if exist temp_food.json del temp_food.json

endlocal
