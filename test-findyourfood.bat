@echo off
echo ========================================
echo Testing Find Your Food API Endpoints
echo ========================================
echo.

echo [1/5] GET /api/v1/public/categories
curl -X GET "http://localhost:8080/api/v1/public/categories"
echo.
echo.

echo [2/5] GET /api/v1/public/categories/MP/foods (Makanan Pokok)
curl -X GET "http://localhost:8080/api/v1/public/categories/MP/foods?limit=5"
echo.
echo.

echo [3/5] GET /api/v1/public/foods/search?q=nasi
curl -X GET "http://localhost:8080/api/v1/public/foods/search?q=nasi&limit=5"
echo.
echo.

echo [4/5] GET /api/v1/public/foods/search?q=mangga  
curl -X GET "http://localhost:8080/api/v1/public/foods/search?q=mangga&limit=5"
echo.
echo.

echo [5/5] GET /api/v1/public/foods/:id (Get detail with goroutines)
echo Getting food detail for first food...
curl -X GET "http://localhost:8080/api/v1/public/foods/uuid-pisang"
echo.
echo.

echo ========================================
echo Testing Complete!
echo ========================================
pause
