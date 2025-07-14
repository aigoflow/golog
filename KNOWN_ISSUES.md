# Known Issues and Limitations

## Rule Resolution Limitations

### Grandparent Rule Issue
The Prolog engine currently has limitations with complex rule resolution, particularly with transitive rules like the grandparent relationship.

**Symptom:**
```prolog
?- parent(tom, bob).
?- parent(bob, alice).
?- grandparent(X, Z) :- parent(X, Y), parent(Y, Z).
?- grandparent(tom, X)
% Returns: No successful solutions (should return X = alice)
```

**Workaround:**
For now, you can work around this by using direct queries:
```prolog
?- parent(tom, Y), parent(Y, Z)
% This will correctly return Y = bob, Z = alice
```

**Status:** This is a known limitation in the engine's rule resolution mechanism that requires deeper investigation into the unification and backtracking implementation.

## UI Issues (Fixed)

### ✅ Query Parsing (Fixed)
Previously, queries with commas inside parentheses were incorrectly parsed. This has been fixed.

### ✅ Rule Body Parsing (Fixed)  
Rule bodies were being split incorrectly at commas. This has been fixed by using proper parenthesis-aware parsing.

### ✅ Help Command (Fixed)
The `help` command was not implemented as a builtin predicate. This has been added.

## Testing

Run `./test_tutorial.sh` to verify the current functionality and see which features are working correctly.