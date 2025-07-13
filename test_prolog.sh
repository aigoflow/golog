#!/bin/bash

# Prolog Server Test Suite
# Run this after starting the server with: go run main.go

BASE_URL="http://localhost:8080"
FAILED=0
PASSED=0

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🧪 Prolog Server Test Suite${NC}"
echo "==============================="

# Helper function to check if server is running
check_server() {
    if ! curl -s "$BASE_URL" > /dev/null 2>&1; then
        echo -e "${RED}❌ Server not running on $BASE_URL${NC}"
        echo "Start the server first: go run main.go"
        exit 1
    fi
}

# Helper function for test assertions
assert_response() {
    local test_name="$1"
    local expected_status="$2"
    local response="$3"
    local actual_status=$(echo "$response" | tail -n1)
    
    if [ "$actual_status" = "$expected_status" ]; then
        echo -e "${GREEN}✅ $test_name${NC}"
        ((PASSED++))
    else
        echo -e "${RED}❌ $test_name${NC}"
        echo "   Expected: $expected_status"
        echo "   Got: $actual_status"
        ((FAILED++))
    fi
}

# Helper function to test JSON response contains expected text
assert_contains() {
    local test_name="$1"
    local expected_text="$2"
    local response="$3"
    local body=$(echo "$response" | sed '$d')
    
    if echo "$body" | grep -q "$expected_text"; then
        echo -e "${GREEN}✅ $test_name${NC}"
        ((PASSED++))
    else
        echo -e "${RED}❌ $test_name${NC}"
        echo "   Expected to contain: $expected_text"
        echo "   Got: $body"
        ((FAILED++))
    fi
}

check_server

echo -e "\n${YELLOW}📝 Testing Basic Facts${NC}"
echo "----------------------"

# Test 1: Add simple atom fact
echo "Test 1: Add atom fact - likes(mary, pizza)"
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/facts" \
  -H "Content-Type: application/json" \
  -d '{
    "predicate": {
      "type": "compound",
      "value": "likes",
      "args": [
        {"type": "atom", "value": "mary"},
        {"type": "atom", "value": "pizza"}
      ]
    }
  }')
assert_response "Add likes fact" "200" "$response"

# Test 2: Add compound fact with variables
echo "Test 2: Add fact - parent(john, mary)"
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/facts" \
  -H "Content-Type: application/json" \
  -d '{
    "predicate": {
      "type": "compound", 
      "value": "parent",
      "args": [
        {"type": "atom", "value": "john"},
        {"type": "atom", "value": "mary"}
      ]
    }
  }')
assert_response "Add parent fact" "200" "$response"

# Test 3: Add another parent fact
echo "Test 3: Add fact - parent(mary, ann)"
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/facts" \
  -H "Content-Type: application/json" \
  -d '{
    "predicate": {
      "type": "compound",
      "value": "parent", 
      "args": [
        {"type": "atom", "value": "mary"},
        {"type": "atom", "value": "ann"}
      ]
    }
  }')
assert_response "Add second parent fact" "200" "$response"

echo -e "\n${YELLOW}🔍 Testing Basic Queries${NC}"
echo "------------------------"

# Test 4: Query with atom
echo "Test 4: Query - likes(mary, pizza)"
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/query" \
  -H "Content-Type: application/json" \
  -d '{
    "goals": [{
      "type": "compound",
      "value": "likes",
      "args": [
        {"type": "atom", "value": "mary"},
        {"type": "atom", "value": "pizza"}
      ]
    }]
  }')
assert_response "Query likes fact" "200" "$response"
assert_contains "Query returns success" '"success":true' "$response"

# Test 5: Query with variable
echo "Test 5: Query - parent(X, mary)"
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/query" \
  -H "Content-Type: application/json" \
  -d '{
    "goals": [{
      "type": "compound",
      "value": "parent", 
      "args": [
        {"type": "variable", "value": "X"},
        {"type": "atom", "value": "mary"}
      ]
    }]
  }')
assert_response "Query with variable" "200" "$response"
assert_contains "Variable binding found" '"X"' "$response"

# Test 6: Query that should fail
echo "Test 6: Query - parent(bob, charlie) [should fail]"
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/query" \
  -H "Content-Type: application/json" \
  -d '{
    "goals": [{
      "type": "compound",
      "value": "parent",
      "args": [
        {"type": "atom", "value": "bob"},
        {"type": "atom", "value": "charlie"}
      ]
    }]
  }')
assert_response "Failed query" "200" "$response"
assert_contains "Query returns failure" '"success":false' "$response"

echo -e "\n${YELLOW}📏 Testing Rules${NC}"
echo "----------------"

# Test 7: Add a rule - grandparent(X,Z) :- parent(X,Y), parent(Y,Z)
echo "Test 7: Add rule - grandparent(X,Z) :- parent(X,Y), parent(Y,Z)"
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/rules" \
  -H "Content-Type: application/json" \
  -d '{
    "head": {
      "type": "compound",
      "value": "grandparent",
      "args": [
        {"type": "variable", "value": "X"},
        {"type": "variable", "value": "Z"}
      ]
    },
    "body": [
      {
        "type": "compound",
        "value": "parent",
        "args": [
          {"type": "variable", "value": "X"},
          {"type": "variable", "value": "Y"}
        ]
      },
      {
        "type": "compound", 
        "value": "parent",
        "args": [
          {"type": "variable", "value": "Y"},
          {"type": "variable", "value": "Z"}
        ]
      }
    ]
  }')
assert_response "Add grandparent rule" "200" "$response"

# Test 8: Query using the rule
echo "Test 8: Query - grandparent(john, ann)"
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/query" \
  -H "Content-Type: application/json" \
  -d '{
    "goals": [{
      "type": "compound",
      "value": "grandparent",
      "args": [
        {"type": "atom", "value": "john"},
        {"type": "atom", "value": "ann"}
      ]
    }]
  }')
assert_response "Query grandparent rule" "200" "$response"
assert_contains "Rule query succeeds" '"success":true' "$response"

echo -e "\n${YELLOW}🔧 Testing Built-ins${NC}"
echo "--------------------"

# Test 9: Unification built-in
echo "Test 9: Test unification - =(X, mary)"
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/query" \
  -H "Content-Type: application/json" \
  -d '{
    "goals": [{
      "type": "compound",
      "value": "=",
      "args": [
        {"type": "variable", "value": "X"},
        {"type": "atom", "value": "mary"}
      ]
    }]
  }')
assert_response "Unification builtin" "200" "$response"
assert_contains "Unification binds variable" '"X"' "$response"

# Test 10: Type checking built-in
echo "Test 10: Test type checking - atom(mary)"
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/query" \
  -H "Content-Type: application/json" \
  -d '{
    "goals": [{
      "type": "compound",
      "value": "atom",
      "args": [
        {"type": "atom", "value": "mary"}
      ]
    }]
  }')
assert_response "Type checking builtin" "200" "$response"
assert_contains "Atom check succeeds" '"success":true' "$response"

# Test 11: Variable type check
echo "Test 11: Test variable check - var(X)"  
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/query" \
  -H "Content-Type: application/json" \
  -d '{
    "goals": [{
      "type": "compound",
      "value": "var",
      "args": [
        {"type": "variable", "value": "X"}
      ]
    }]
  }')
assert_response "Variable check builtin" "200" "$response"
assert_contains "Variable check succeeds" '"success":true' "$response"

echo -e "\n${YELLOW}🚫 Testing Error Cases${NC}"
echo "----------------------"

# Test 12: Invalid JSON
echo "Test 12: Invalid JSON request"
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/facts" \
  -H "Content-Type: application/json" \
  -d '{invalid json}')
assert_response "Invalid JSON" "400" "$response"

# Test 13: Missing required fields
echo "Test 13: Missing predicate field"
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/facts" \
  -H "Content-Type: application/json" \
  -d '{}')
assert_response "Missing predicate" "500" "$response"

echo -e "\n${YELLOW}📊 Test Results${NC}"
echo "================"
echo -e "Passed: ${GREEN}$PASSED${NC}"
echo -e "Failed: ${RED}$FAILED${NC}"

if [ $FAILED -eq 0 ]; then
    echo -e "\n${GREEN}🎉 All tests passed!${NC}"
    exit 0
else
    echo -e "\n${RED}💥 Some tests failed!${NC}"
    exit 1
fi