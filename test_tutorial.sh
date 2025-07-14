#!/bin/bash

# Test all tutorial examples automatically
# This ensures all parsing and functionality works correctly

BASE_URL="http://localhost:3000/api/v1"

echo "🧪 Testing Tutorial Examples"
echo "============================"

# Create a test session
echo "Creating test session..."
SESSION_RESPONSE=$(curl -s -X POST "$BASE_URL/sessions" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "tutorial-test-'$(date +%s)'",
    "description": "Tutorial test session"
  }')

SESSION_ID=$(echo "$SESSION_RESPONSE" | jq -r '.id')
echo "Created session ID: $SESSION_ID"

# Step 1: Add parent(tom, bob).
echo -e "\n1. Testing: parent(tom, bob)."
curl -s -X POST "$BASE_URL/sessions/$SESSION_ID/facts" \
  -H "Content-Type: application/json" \
  -d '{
    "predicate": {
      "type": "compound",
      "value": "parent",
      "args": [
        {"type": "atom", "value": "tom"},
        {"type": "atom", "value": "bob"}
      ]
    }
  }' | jq -r '.status'

# Step 2: Add parent(bob, alice).
echo "2. Testing: parent(bob, alice)."
curl -s -X POST "$BASE_URL/sessions/$SESSION_ID/facts" \
  -H "Content-Type: application/json" \
  -d '{
    "predicate": {
      "type": "compound",
      "value": "parent",
      "args": [
        {"type": "atom", "value": "bob"},
        {"type": "atom", "value": "alice"}
      ]
    }
  }' | jq -r '.status'

# Step 3: Query parent(tom, X)
echo "3. Testing: parent(tom, X)"
RESULT=$(curl -s -X POST "$BASE_URL/sessions/$SESSION_ID/query" \
  -H "Content-Type: application/json" \
  -d '{
    "goals": [{
      "type": "compound",
      "value": "parent",
      "args": [
        {"type": "atom", "value": "tom"},
        {"type": "variable", "value": "X"}
      ]
    }]
  }')
echo "$RESULT" | jq -r '.solutions[0].bindings.X.value' | xargs -I {} echo "   Result: X = {}"

# Step 4: Add grandparent rule
echo "4. Testing: grandparent(X, Z) :- parent(X, Y), parent(Y, Z)."
curl -s -X POST "$BASE_URL/sessions/$SESSION_ID/rules" \
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
  }' | jq -r '.status'

# Step 5: Query grandparent(tom, X)
echo "5. Testing: grandparent(tom, X)"
RESULT=$(curl -s -X POST "$BASE_URL/sessions/$SESSION_ID/query" \
  -H "Content-Type: application/json" \
  -d '{
    "goals": [{
      "type": "compound",
      "value": "grandparent",
      "args": [
        {"type": "atom", "value": "tom"},
        {"type": "variable", "value": "X"}
      ]
    }]
  }')
echo "$RESULT" | jq -r '.solutions[0].bindings.X.value' | xargs -I {} echo "   Result: X = {}"

# Step 6: Test help builtin
echo "6. Testing: help"
RESULT=$(curl -s -X POST "$BASE_URL/sessions/$SESSION_ID/query" \
  -H "Content-Type: application/json" \
  -d '{
    "goals": [{
      "type": "compound",
      "value": "help",
      "args": []
    }]
  }')
echo "$RESULT" | jq -r 'if .solutions[0].success then "   Result: Yes" else "   Result: No" end'

# Step 7: Add score(alice, 95).
echo "7. Testing: score(alice, 95)."
curl -s -X POST "$BASE_URL/sessions/$SESSION_ID/facts" \
  -H "Content-Type: application/json" \
  -d '{
    "predicate": {
      "type": "compound",
      "value": "score",
      "args": [
        {"type": "atom", "value": "alice"},
        {"type": "number", "value": 95}
      ]
    }
  }' | jq -r '.status'

# Step 8: Add score(bob, 87).
echo "8. Testing: score(bob, 87)."
curl -s -X POST "$BASE_URL/sessions/$SESSION_ID/facts" \
  -H "Content-Type: application/json" \
  -d '{
    "predicate": {
      "type": "compound",
      "value": "score",
      "args": [
        {"type": "atom", "value": "bob"},
        {"type": "number", "value": 87}
      ]
    }
  }' | jq -r '.status'

# Step 9: Test aggregation count(_, parent(X, Y), N)
echo "9. Testing: count(_, parent(X, Y), N)"
RESULT=$(curl -s -X POST "$BASE_URL/sessions/$SESSION_ID/query" \
  -H "Content-Type: application/json" \
  -d '{
    "goals": [{
      "type": "compound",
      "value": "count",
      "args": [
        {"type": "variable", "value": "_"},
        {
          "type": "compound",
          "value": "parent",
          "args": [
            {"type": "variable", "value": "X"},
            {"type": "variable", "value": "Y"}
          ]
        },
        {"type": "variable", "value": "N"}
      ]
    }]
  }')
echo "$RESULT" | jq -r '.solutions[0].bindings.N.value' | xargs -I {} echo "   Result: N = {}"

# Clean up
echo -e "\nCleaning up..."
curl -s -X DELETE "$BASE_URL/sessions/$SESSION_ID" > /dev/null
echo "✅ All tutorial tests completed!"