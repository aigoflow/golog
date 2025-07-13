#!/bin/bash

# Prolog Server Demo - Interactive demonstration
# Shows off the key features with readable output

BASE_URL="http://localhost:8080"

# Colors
BLUE='\033[0;34m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}🔮 Prolog Server Demo${NC}"
echo "====================="

# Helper to make pretty API calls
call_api() {
    local method="$1"
    local endpoint="$2"
    local data="$3"
    local description="$4"
    
    echo -e "\n${YELLOW}📡 $description${NC}"
    echo "   $method $endpoint"
    if [ -n "$data" ]; then
        echo "   Data: $(echo "$data" | jq -c .)"
    fi
    
    response=$(curl -s -X "$method" "$BASE_URL$endpoint" \
        -H "Content-Type: application/json" \
        -d "$data" | jq .)
    
    echo -e "${GREEN}   Response:${NC}"
    echo "$response" | sed 's/^/   /'
}

echo -e "\n${BLUE}🏗️  Building Knowledge Base${NC}"
echo "============================"

# Add family facts
call_api "POST" "/facts" '{
  "predicate": {
    "type": "compound",
    "value": "parent",
    "args": [
      {"type": "atom", "value": "john"},
      {"type": "atom", "value": "mary"}
    ]
  }
}' "Adding: parent(john, mary)"

call_api "POST" "/facts" '{
  "predicate": {
    "type": "compound", 
    "value": "parent",
    "args": [
      {"type": "atom", "value": "mary"},
      {"type": "atom", "value": "ann"}
    ]
  }
}' "Adding: parent(mary, ann)"

call_api "POST" "/facts" '{
  "predicate": {
    "type": "compound",
    "value": "parent", 
    "args": [
      {"type": "atom", "value": "bob"},
      {"type": "atom", "value": "john"}
    ]
  }
}' "Adding: parent(bob, john)"

call_api "POST" "/facts" '{
  "predicate": {
    "type": "compound",
    "value": "likes",
    "args": [
      {"type": "atom", "value": "mary"},
      {"type": "atom", "value": "pizza"}
    ]
  }
}' "Adding: likes(mary, pizza)"

# Add grandparent rule
call_api "POST" "/rules" '{
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
}' "Adding rule: grandparent(X,Z) :- parent(X,Y), parent(Y,Z)"

echo -e "\n${BLUE}🔍 Querying Knowledge Base${NC}"
echo "=========================="

# Simple fact query
call_api "POST" "/query" '{
  "goals": [{
    "type": "compound",
    "value": "likes",
    "args": [
      {"type": "atom", "value": "mary"},
      {"type": "atom", "value": "pizza"}
    ]
  }]
}' "Query: likes(mary, pizza) - exact match"

# Variable query  
call_api "POST" "/query" '{
  "goals": [{
    "type": "compound",
    "value": "parent",
    "args": [
      {"type": "variable", "value": "X"},
      {"type": "atom", "value": "mary"}
    ]
  }]
}' "Query: parent(X, mary) - find parents of mary"

# Multiple variables
call_api "POST" "/query" '{
  "goals": [{
    "type": "compound", 
    "value": "parent",
    "args": [
      {"type": "variable", "value": "Parent"},
      {"type": "variable", "value": "Child"}
    ]
  }]
}' "Query: parent(Parent, Child) - find all parent relationships"

# Rule-based query
call_api "POST" "/query" '{
  "goals": [{
    "type": "compound",
    "value": "grandparent",
    "args": [
      {"type": "atom", "value": "bob"},
      {"type": "atom", "value": "ann"}
    ]
  }]
}' "Query: grandparent(bob, ann) - using rule inference"

# Query with variable in rule
call_api "POST" "/query" '{
  "goals": [{
    "type": "compound",
    "value": "grandparent", 
    "args": [
      {"type": "variable", "value": "GP"},
      {"type": "atom", "value": "ann"}
    ]
  }]
}' "Query: grandparent(GP, ann) - find grandparents of ann"

echo -e "\n${BLUE}🔧 Testing Built-ins${NC}"
echo "==================="

# Unification
call_api "POST" "/query" '{
  "goals": [{
    "type": "compound",
    "value": "=",
    "args": [
      {"type": "variable", "value": "X"},
      {"type": "atom", "value": "hello"}
    ]
  }]
}' "Built-in: =(X, hello) - unification"

# Type checking
call_api "POST" "/query" '{
  "goals": [{
    "type": "compound",
    "value": "atom",
    "args": [
      {"type": "atom", "value": "mary"}
    ]
  }]
}' "Built-in: atom(mary) - type checking"

call_api "POST" "/query" '{
  "goals": [{
    "type": "compound", 
    "value": "var",
    "args": [
      {"type": "variable", "value": "X"}
    ]
  }]
}' "Built-in: var(X) - variable checking"

# Complex query with multiple goals
call_api "POST" "/query" '{
  "goals": [
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
      "value": "likes", 
      "args": [
        {"type": "variable", "value": "Y"},
        {"type": "atom", "value": "pizza"}
      ]
    }
  ]
}' "Complex: parent(X,Y), likes(Y,pizza) - multi-goal query"

echo -e "\n${GREEN}🎉 Demo Complete!${NC}"
echo "The server demonstrates:"
echo "  ✅ Fact storage and retrieval"
echo "  ✅ Rule definition and inference"  
echo "  ✅ Variable unification"
echo "  ✅ Backtracking search"
echo "  ✅ Built-in predicates"
echo "  ✅ Complex multi-goal queries"