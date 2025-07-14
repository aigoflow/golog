package main

const loginTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Prolog Engine - Login</title>
    <style>
        :root {
            --bg-primary: #1a1a1a;
            --bg-secondary: #2a2a2a;
            --text-primary: #00ff00;
            --text-secondary: #00ffff;
            --text-error: #ff0000;
            --border-color: #00ff00;
            --button-bg: #00ff00;
            --button-text: #1a1a1a;
            --button-hover: #00ffff;
            --input-bg: #1a1a1a;
        }
        body {
            font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
            background: var(--bg-primary);
            color: var(--text-primary);
            margin: 0;
            padding: 0;
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
        }
        .login-container {
            background: var(--bg-secondary);
            border: 2px solid var(--border-color);
            border-radius: 8px;
            padding: 2rem;
            max-width: 400px;
            width: 100%;
        }
        .logo {
            text-align: center;
            font-size: 1.5rem;
            margin-bottom: 1rem;
            color: var(--text-primary);
        }
        .form-group {
            margin-bottom: 1rem;
        }
        label {
            display: block;
            margin-bottom: 0.5rem;
            color: var(--text-primary);
        }
        input[type="password"] {
            width: 100%;
            padding: 0.5rem;
            background: var(--input-bg);
            border: 1px solid var(--border-color);
            color: var(--text-primary);
            font-family: inherit;
            box-sizing: border-box;
        }
        input[type="password"]:focus {
            outline: none;
            border-color: var(--text-secondary);
            box-shadow: 0 0 5px var(--text-secondary);
        }
        button {
            width: 100%;
            padding: 0.75rem;
            background: var(--button-bg);
            color: var(--button-text);
            border: none;
            cursor: pointer;
            font-family: inherit;
            font-weight: bold;
        }
        button:hover {
            background: var(--button-hover);
        }
        .error {
            color: var(--text-error);
            margin-top: 0.5rem;
            text-align: center;
        }
    </style>
</head>
<body>
    <div class="login-container">
        <div class="logo">🧠 Prolog Engine</div>
        <form method="POST" action="/ui/login">
            <div class="form-group">
                <label for="password">Password:</label>
                <input type="password" id="password" name="password" required>
            </div>
            <button type="submit">Login</button>
            {{if .error}}
            <div class="error">{{.error}}</div>
            {{end}}
        </form>
    </div>
</body>
</html>`

const uiTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        
        /* Dark Theme (Default) */
        :root {
            --bg-primary: #1a1a1a;
            --bg-secondary: #2a2a2a;
            --bg-tertiary: #333;
            --text-primary: #00ff00;
            --text-secondary: #00ffff;
            --text-error: #ff0000;
            --text-warning: #ffff00;
            --text-muted: #aaa;
            --border-color: #00ff00;
            --scrollbar-thumb: #00ff00;
            --scrollbar-thumb-hover: #00ffff;
            --scrollbar-track: #0a0a0a;
            --button-bg: #00ff00;
            --button-text: #1a1a1a;
            --button-hover: #00ffff;
            --input-bg: #1a1a1a;
            --modal-overlay: rgba(0, 0, 0, 0.8);
        }
        
        /* Light Theme */
        body.light-theme {
            --bg-primary: #ffffff;
            --bg-secondary: #f5f5f5;
            --bg-tertiary: #e0e0e0;
            --text-primary: #2d7a2d;
            --text-secondary: #1976d2;
            --text-error: #d32f2f;
            --text-warning: #f57c00;
            --text-muted: #666;
            --border-color: #2d7a2d;
            --scrollbar-thumb: #2d7a2d;
            --scrollbar-thumb-hover: #1976d2;
            --scrollbar-track: #e0e0e0;
            --button-bg: #2d7a2d;
            --button-text: #ffffff;
            --button-hover: #1976d2;
            --input-bg: #ffffff;
            --modal-overlay: rgba(0, 0, 0, 0.5);
        }
        
        body { 
            font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace; 
            background: var(--bg-primary); 
            color: var(--text-primary);
            height: 100vh; 
            overflow: hidden;
            transition: background-color 0.3s ease, color 0.3s ease;
        }
        .container { display: flex; height: 100vh; }
        .sidebar { width: 250px; background: var(--bg-secondary); border-right: 2px solid var(--border-color); padding: 1rem; overflow-y: auto; }
        .main-content { flex: 1; display: flex; flex-direction: column; }
        .header { background: var(--bg-secondary); border-bottom: 2px solid var(--border-color); padding: 1rem; display: flex; justify-content: space-between; align-items: center; }
        .terminal-container { 
            flex: 1; 
            background: var(--bg-primary); 
            position: relative; 
            overflow: hidden;
            display: flex;
            flex-direction: column;
        }
        .help-sidebar { 
            width: 300px; 
            background: var(--bg-secondary); 
            border-left: 2px solid var(--border-color); 
            padding: 1rem; 
            overflow-y: auto; 
            transition: width 0.3s ease;
            scrollbar-width: thin;
            scrollbar-color: var(--scrollbar-thumb) var(--bg-secondary);
        }
        .help-sidebar::-webkit-scrollbar {
            width: 8px;
        }
        .help-sidebar::-webkit-scrollbar-track {
            background: var(--bg-secondary);
        }
        .help-sidebar::-webkit-scrollbar-thumb {
            background: var(--scrollbar-thumb);
            border-radius: 4px;
        }
        .help-sidebar::-webkit-scrollbar-thumb:hover {
            background: var(--scrollbar-thumb-hover);
        }
        .help-sidebar.collapsed { width: 0; padding: 0; overflow: hidden; }
        .terminal { 
            width: 100%; 
            height: 100%; 
            background: var(--bg-primary); 
            color: var(--text-primary); 
            font-family: inherit; 
            font-size: 14px; 
            padding: 1rem; 
            border: none; 
            outline: none; 
            resize: none; 
            overflow-y: scroll; 
            overflow-x: hidden;
            scrollbar-width: thin;
            scrollbar-color: var(--scrollbar-thumb) var(--bg-primary);
        }
        .terminal::-webkit-scrollbar {
            width: 12px;
        }
        .terminal::-webkit-scrollbar-track {
            background: var(--scrollbar-track);
            border: 1px solid var(--border-color);
        }
        .terminal::-webkit-scrollbar-thumb {
            background: var(--scrollbar-thumb);
            border-radius: 6px;
            border: 1px solid var(--scrollbar-track);
        }
        .terminal::-webkit-scrollbar-thumb:hover {
            background: var(--scrollbar-thumb-hover);
        }
        .session-list { margin-bottom: 1rem; }
        .session-item { 
            background: var(--bg-primary); 
            border: 1px solid var(--bg-tertiary); 
            margin-bottom: 0.5rem; 
            padding: 0.5rem; 
            cursor: pointer; 
            transition: all 0.2s; 
        }
        .session-item:hover { 
            border-color: var(--border-color); 
            background: var(--bg-tertiary); 
        }
        .session-item.active { 
            border-color: var(--border-color); 
            background: var(--bg-tertiary); 
        }
        .session-name { 
            font-weight: bold; 
            color: var(--text-secondary); 
        }
        .session-desc { 
            font-size: 0.8rem; 
            color: var(--text-muted); 
            margin-top: 0.2rem; 
        }
        .session-buttons {
            display: flex;
            flex-direction: column;
            gap: 0.5rem;
        }
        .btn { 
            background: var(--button-bg); 
            color: var(--button-text); 
            border: none; 
            padding: 0.5rem 1rem; 
            cursor: pointer; 
            font-family: inherit; 
            margin: 0.2rem;
            transition: background-color 0.3s ease;
            width: 100%;
            box-sizing: border-box;
        }
        .btn:hover { 
            background: var(--button-hover); 
        }
        .btn-small { padding: 0.2rem 0.5rem; font-size: 0.8rem; }
        .input-group { margin-bottom: 0.5rem; }
        .input-group input { 
            width: 100%; 
            background: var(--input-bg); 
            border: 1px solid var(--bg-tertiary); 
            color: var(--text-primary); 
            padding: 0.5rem; 
            font-family: inherit;
            transition: border-color 0.3s ease;
        }
        .input-group input:focus { 
            outline: none; 
            border-color: var(--border-color); 
        }
        .header-buttons { display: flex; gap: 0.5rem; align-items: center; }
        .toggle-help, .theme-toggle, .clear-terminal { 
            background: var(--bg-tertiary); 
            color: var(--text-primary); 
            border: 1px solid var(--border-color); 
            padding: 0.5rem 1rem; 
            cursor: pointer; 
            font-family: inherit;
            transition: all 0.3s ease;
            height: 40px;
            display: flex;
            align-items: center;
            justify-content: center;
            border-radius: 4px;
            white-space: nowrap;
        }
        .toggle-help:hover, .theme-toggle:hover, .clear-terminal:hover { 
            background: var(--button-bg); 
            color: var(--button-text); 
        }
        .theme-toggle {
            font-size: 1.2rem;
            width: 40px;
            padding: 0.5rem;
        }
        .help-section { margin-bottom: 1.5rem; }
        .help-title { 
            color: var(--text-secondary); 
            font-weight: bold; 
            margin-bottom: 0.5rem; 
            border-bottom: 1px solid var(--bg-tertiary); 
            padding-bottom: 0.2rem; 
        }
        .example { 
            background: var(--bg-primary); 
            border: 1px solid var(--bg-tertiary); 
            padding: 0.5rem; 
            margin: 0.5rem 0; 
            cursor: pointer; 
            transition: all 0.2s; 
        }
        .example:hover { 
            border-color: var(--border-color); 
        }
        .example-title { color: var(--text-secondary); font-weight: bold; margin-bottom: 0.2rem; }
        .example-code { color: var(--text-warning); font-size: 0.9rem; margin-bottom: 0.2rem; }
        .example-desc { color: var(--text-muted); font-size: 0.8rem; }
        .terminal-output { white-space: pre-wrap; word-wrap: break-word; }
        .prompt { color: var(--text-secondary); }
        .error { color: var(--text-error); }
        .success { color: var(--text-primary); }
        .warning { color: var(--text-warning); }
        .modal { 
            display: none; 
            position: fixed; 
            z-index: 1000; 
            left: 0; 
            top: 0; 
            width: 100%; 
            height: 100%; 
            background-color: var(--modal-overlay); 
        }
        .modal-content { 
            background-color: var(--bg-secondary); 
            margin: 15% auto; 
            padding: 20px; 
            border: 2px solid var(--border-color); 
            width: 400px; 
            max-width: 90%; 
        }
        .close { 
            color: var(--text-muted); 
            float: right; 
            font-size: 28px; 
            font-weight: bold; 
            cursor: pointer; 
        }
        .close:hover { 
            color: var(--text-primary); 
        }
        .help-tabs { display: flex; margin-bottom: 1rem; }
        .tab-button { background: var(--bg-tertiary); color: var(--text-primary); border: none; padding: 0.5rem 1rem; cursor: pointer; font-family: inherit; margin-right: 0.2rem; }
        .tab-button.active { background: var(--button-bg); color: var(--button-text); }
        .tab-button:hover { background: var(--bg-secondary); }
        .tab-button.active:hover { background: var(--button-hover); color: var(--button-text); }
        .tab-content { }
        .tutorial-progress { background: var(--bg-primary); padding: 0.5rem; margin-bottom: 1rem; border: 1px solid var(--bg-tertiary); display: flex; justify-content: space-between; align-items: center; }
        .tutorial-step { margin-bottom: 1.5rem; }
        .step-title { color: var(--text-secondary); font-weight: bold; margin-bottom: 0.5rem; }
        .step-desc { color: var(--text-muted); margin-bottom: 0.5rem; font-size: 0.9rem; }
        .step-command { background: var(--bg-primary); border: 2px solid var(--bg-tertiary); padding: 1rem; cursor: pointer; transition: all 0.2s; }
        .step-command:hover { border-color: var(--border-color); }
        .step-command.completed { border-color: var(--border-color); opacity: 0.7; }
        .cmd-text { color: var(--text-warning); font-weight: bold; margin-bottom: 0.5rem; }
        .cmd-expected { color: var(--text-primary); font-size: 0.8rem; }
        .tutorial-complete { text-align: center; padding: 2rem; background: var(--bg-primary); border: 2px solid var(--border-color); }
    </style>
</head>
<body>
    <div class="container">
        <div class="sidebar">
            <h3>Sessions</h3>
            <div class="session-list" id="sessionList"></div>
            <div class="session-buttons">
                <button class="btn" onclick="showCreateSessionModal()">New Session</button>
                <button class="btn btn-small" onclick="deleteCurrentSession()">Delete Current</button>
                <button class="btn btn-small" onclick="clearCache()">Clear Cache</button>
            </div>
            <div style="margin-top: 2rem;">
                <h4>Current Session</h4>
                <div id="currentSessionInfo">No session selected</div>
            </div>
        </div>
        
        <div class="main-content">
            <div class="header">
                <h2>🧠 Prolog Engine REPL</h2>
                <div class="header-buttons">
                    <button class="theme-toggle" onclick="toggleTheme()" title="Toggle Dark/Light Mode">
                        <span id="themeIcon">🌙</span>
                    </button>
                    <button class="clear-terminal" onclick="clearTerminal()" title="Clear Terminal">Clear</button>
                    <button class="toggle-help" onclick="toggleHelp()">Toggle Help</button>
                </div>
            </div>
            <div class="terminal-container">
                <div id="terminal" class="terminal" contenteditable="true"></div>
            </div>
        </div>
        
        <div class="help-sidebar" id="helpSidebar">
            <div class="help-tabs">
                <button class="tab-button active" onclick="showTab('help')">Help</button>
                <button class="tab-button" onclick="showTab('tutorial')">Tutorial</button>
            </div>
            
            <div id="helpTab" class="tab-content">
                <h3>Help & Examples</h3>
                <div class="help-section">
                <div class="help-title">Basic Facts</div>
                <div class="example" onclick="insertExample('fact(atom)')">
                    <div class="example-title">Simple Fact</div>
                    <div class="example-code">fact(atom).</div>
                    <div class="example-desc">Declares a simple fact</div>
                </div>
                <div class="example" onclick="insertExample('parent(tom, bob)')">
                    <div class="example-title">Relationship</div>
                    <div class="example-code">parent(tom, bob).</div>
                    <div class="example-desc">Tom is parent of Bob</div>
                </div>
            </div>
            
            <div class="help-section">
                <div class="help-title">Rules</div>
                <div class="example" onclick="insertExample('grandparent(X, Z) :- parent(X, Y), parent(Y, Z)')">
                    <div class="example-title">Grandparent Rule</div>
                    <div class="example-code">grandparent(X, Z) :- parent(X, Y), parent(Y, Z).</div>
                    <div class="example-desc">Defines grandparent relationship</div>
                </div>
            </div>
            
            <div class="help-section">
                <div class="help-title">Queries</div>
                <div class="example" onclick="insertExample('?- parent(tom, X)')">
                    <div class="example-title">Find Children</div>
                    <div class="example-code">?- parent(tom, X).</div>
                    <div class="example-desc">Find all children of Tom</div>
                </div>
            </div>
            
            <div class="help-section">
                <div class="help-title">Built-in Predicates</div>
                <div class="example" onclick="insertExample('?- =(X, test)')">
                    <div class="example-title">Unification</div>
                    <div class="example-code">?- =(X, test).</div>
                    <div class="example-desc">Unify X with 'test'</div>
                </div>
                <div class="example" onclick="insertExample('?- now(X)')">
                    <div class="example-title">Current Time</div>
                    <div class="example-code">?- now(X).</div>
                    <div class="example-desc">Get current timestamp</div>
                </div>
            </div>
            
            <div class="help-section">
                <div class="help-title">Aggregation</div>
                <div class="example" onclick="insertExample('?- count(_, parent(X, Y), N)')">
                    <div class="example-title">Count</div>
                    <div class="example-code">?- count(_, parent(X, Y), N).</div>
                    <div class="example-desc">Count all parent relationships</div>
                </div>
            </div>
            </div>
            
            <div id="tutorialTab" class="tab-content" style="display: none;">
                <h3>Interactive Tutorial</h3>
                <p><strong>Click on commands below to insert them into the terminal. All steps remain visible to show progress.</strong></p>
                
                <div class="tutorial-step" id="step1">
                    <div class="step-title">1. Add Basic Facts</div>
                    <div class="step-desc">Let's start by adding some family relationship facts.</div>
                    <div class="step-command" onclick="insertExample('parent(tom, bob).')">
                        <div class="cmd-text">parent(tom, bob).</div>
                        <div class="cmd-expected">Expected: "Fact added."</div>
                    </div>
                </div>
                
                <div class="tutorial-step" id="step2">
                    <div class="step-title">2. Add More Facts</div>
                    <div class="step-desc">Add another parent relationship.</div>
                    <div class="step-command" onclick="insertExample('parent(bob, alice).')">
                        <div class="cmd-text">parent(bob, alice).</div>
                        <div class="cmd-expected">Expected: "Fact added."</div>
                    </div>
                </div>
                
                <div class="tutorial-step" id="step3">
                    <div class="step-title">3. Query Facts</div>
                    <div class="step-desc">Find children of Tom.</div>
                    <div class="step-command" onclick="insertExample('parent(tom, X)')">
                        <div class="cmd-text">parent(tom, X)</div>
                        <div class="cmd-expected">Expected: X = bob</div>
                    </div>
                </div>
                
                <div class="tutorial-step" id="step4">
                    <div class="step-title">4. Add a Rule</div>
                    <div class="step-desc">Define grandparent relationship.</div>
                    <div class="step-command" onclick="insertExample('grandparent(X, Z) :- parent(X, Y), parent(Y, Z).')">
                        <div class="cmd-text">grandparent(X, Z) :- parent(X, Y), parent(Y, Z).</div>
                        <div class="cmd-expected">Expected: "Rule added."</div>
                    </div>
                </div>
                
                <div class="tutorial-step" id="step5">
                    <div class="step-title">5. Query Rule</div>
                    <div class="step-desc">Find grandchildren of Tom.</div>
                    <div class="step-command" onclick="insertExample('grandparent(tom, X)')">
                        <div class="cmd-text">grandparent(tom, X)</div>
                        <div class="cmd-expected">Expected: X = alice</div>
                    </div>
                </div>
                
                <div class="tutorial-step" id="step6">
                    <div class="step-title">6. Unification</div>
                    <div class="step-desc">Test the unification builtin =.</div>
                    <div class="step-command" onclick="insertExample('=(X, alice)')">
                        <div class="cmd-text">=(X, alice)</div>
                        <div class="cmd-expected">Expected: X = alice</div>
                    </div>
                </div>
                
                <div class="tutorial-step" id="step7">
                    <div class="step-title">7. Type Checking</div>
                    <div class="step-desc">Check if alice is an atom.</div>
                    <div class="step-command" onclick="insertExample('atom(alice)')">
                        <div class="cmd-text">atom(alice)</div>
                        <div class="cmd-expected">Expected: Yes (1)</div>
                    </div>
                </div>
                
                <div class="tutorial-step" id="step8">
                    <div class="step-title">8. Variable Check</div>
                    <div class="step-desc">Check if X is a variable.</div>
                    <div class="step-command" onclick="insertExample('var(X)')">
                        <div class="cmd-text">var(X)</div>
                        <div class="cmd-expected">Expected: Yes (1)</div>
                    </div>
                </div>
                
                <div class="tutorial-step" id="step9">
                    <div class="step-title">9. Number Type</div>
                    <div class="step-desc">Add numeric facts and check number type.</div>
                    <div class="step-command" onclick="insertExample('age(alice, 25).')">
                        <div class="cmd-text">age(alice, 25).</div>
                        <div class="cmd-expected">Expected: "Fact added."</div>
                    </div>
                </div>
                
                <div class="tutorial-step" id="step10">
                    <div class="step-title">10. Number Check</div>
                    <div class="step-desc">Check if 25 is a number.</div>
                    <div class="step-command" onclick="insertExample('number(25)')">
                        <div class="cmd-text">number(25)</div>
                        <div class="cmd-expected">Expected: Yes (1)</div>
                    </div>
                </div>
                
                <div class="tutorial-step" id="step11">
                    <div class="step-title">11. Add Scores</div>
                    <div class="step-desc">Add score facts for aggregation tests.</div>
                    <div class="step-command" onclick="insertExample('score(alice, 95).')">
                        <div class="cmd-text">score(alice, 95).</div>
                        <div class="cmd-expected">Expected: "Fact added."</div>
                    </div>
                </div>
                
                <div class="tutorial-step" id="step12">
                    <div class="step-title">12. More Scores</div>
                    <div class="step-desc">Add another score.</div>
                    <div class="step-command" onclick="insertExample('score(bob, 87).')">
                        <div class="cmd-text">score(bob, 87).</div>
                        <div class="cmd-expected">Expected: "Fact added."</div>
                    </div>
                </div>
                
                <div class="tutorial-step" id="step13">
                    <div class="step-title">13. Count Aggregation</div>
                    <div class="step-desc">Count all parent relationships.</div>
                    <div class="step-command" onclick="insertExample('count(_, parent(X, Y), N)')">
                        <div class="cmd-text">count(_, parent(X, Y), N)</div>
                        <div class="cmd-expected">Expected: N = 2</div>
                    </div>
                </div>
                
                <div class="tutorial-step" id="step14">
                    <div class="step-title">14. Sum Aggregation</div>
                    <div class="step-desc">Sum all scores (Y values).</div>
                    <div class="step-command" onclick="insertExample('sum(Y, score(X, Y), Total)')">
                        <div class="cmd-text">sum(Y, score(X, Y), Total)</div>
                        <div class="cmd-expected">Expected: Total = 182</div>
                    </div>
                </div>
                
                <div class="tutorial-step" id="step15">
                    <div class="step-title">15. Max Aggregation</div>
                    <div class="step-desc">Find the maximum score.</div>
                    <div class="step-command" onclick="insertExample('max(Y, score(X, Y), Max)')">
                        <div class="cmd-text">max(Y, score(X, Y), Max)</div>
                        <div class="cmd-expected">Expected: Max = 95</div>
                    </div>
                </div>
                
                <div class="tutorial-step" id="step16">
                    <div class="step-title">16. Min Aggregation</div>
                    <div class="step-desc">Find the minimum score.</div>
                    <div class="step-command" onclick="insertExample('min(Y, score(X, Y), Min)')">
                        <div class="cmd-text">min(Y, score(X, Y), Min)</div>
                        <div class="cmd-expected">Expected: Min = 87</div>
                    </div>
                </div>
                
                <div class="tutorial-step" id="step17">
                    <div class="step-title">17. Date/Time</div>
                    <div class="step-desc">Get current timestamp.</div>
                    <div class="step-command" onclick="insertExample('now(T)')">
                        <div class="cmd-text">now(T)</div>
                        <div class="cmd-expected">Expected: T = current timestamp</div>
                    </div>
                </div>
                
                <div class="tutorial-step" id="step18">
                    <div class="step-title">18. Complex Query</div>
                    <div class="step-desc">Find all parent-child pairs.</div>
                    <div class="step-command" onclick="insertExample('parent(X, Y)')">
                        <div class="cmd-text">parent(X, Y)</div>
                        <div class="cmd-expected">Expected: Multiple solutions</div>
                    </div>
                </div>
                
                <div class="tutorial-step" id="step19">
                    <div class="step-title">19. Help Command</div>
                    <div class="step-desc">Show available commands.</div>
                    <div class="step-command" onclick="insertExample('help')">
                        <div class="cmd-text">help</div>
                        <div class="cmd-expected">Expected: Help text displayed</div>
                    </div>
                </div>
                
                <div class="tutorial-step" id="step20">
                    <div class="step-title">20. Clear Cache</div>
                    <div class="step-desc">Clear the memoization cache.</div>
                    <div class="step-command" onclick="insertExample('clear_cache')">
                        <div class="cmd-text">clear_cache</div>
                        <div class="cmd-expected">Expected: "Cache cleared"</div>
                    </div>
                </div>
            </div>
        </div>
    </div>
    
    <div id="createSessionModal" class="modal">
        <div class="modal-content">
            <span class="close" onclick="hideCreateSessionModal()">&times;</span>
            <h3>Create New Session</h3>
            <div class="input-group">
                <label>Name:</label>
                <input type="text" id="sessionName" placeholder="Enter session name">
            </div>
            <div class="input-group">
                <label>Description:</label>
                <input type="text" id="sessionDesc" placeholder="Enter session description">
            </div>
            <button class="btn" onclick="createSession()">Create</button>
            <button class="btn" onclick="hideCreateSessionModal()">Cancel</button>
        </div>
    </div>

    <script src="/ui/js"></script>
</body>
</html>`

const jsContent = `
// Global state
let currentSession = null;
let sessions = [];
let terminalHistory = [];
let historyIndex = -1;
let apiKey = null;

// Initialize
document.addEventListener('DOMContentLoaded', function() {
    initTerminal();
    loadSessions();
    checkApiKey();
});

function checkApiKey() {
    fetch('/api/v1/sessions')
        .then(response => {
            if (response.status === 401) {
                apiKey = prompt('API Key required for backend access:');
            }
            return response;
        })
        .catch(error => {
            console.error('Error checking API key:', error);
        });
}

function getHeaders() {
    const headers = { 'Content-Type': 'application/json' };
    if (apiKey) {
        headers['Authorization'] = 'Bearer ' + apiKey;
    }
    return headers;
}

function initTerminal() {
    const terminal = document.getElementById('terminal');
    terminal.innerHTML = 
        '<span class="success">🧠 Prolog Engine REPL v2.0</span><br>' +
        '<span class="warning">Welcome! Select a session to start or create a new one.</span><br>' +
        '<span class="warning">Type "help" for available commands.</span><br><br>' +
        '<span class="prompt">?- </span>';
    
    terminal.addEventListener('keydown', handleTerminalInput);
    
    // Ensure terminal stays focused
    terminal.addEventListener('click', function() {
        terminal.focus();
    });
    
    // Focus terminal when clicking anywhere in the terminal container
    document.querySelector('.terminal-container').addEventListener('click', function() {
        terminal.focus();
    });
    
    terminal.focus();
}

function handleTerminalInput(event) {
    if (event.key === 'Enter') {
        event.preventDefault();
        const terminal = document.getElementById('terminal');
        const content = terminal.textContent || terminal.innerText;
        const lines = content.split('\n');
        const lastLine = lines[lines.length - 1];
        
        if (lastLine.startsWith('?- ') || lastLine.includes('?- ')) {
            const promptIndex = lastLine.lastIndexOf('?- ');
            const input = lastLine.substring(promptIndex + 3).trim();
            if (input) {
                // Add to history
                terminalHistory.push(input);
                historyIndex = terminalHistory.length;
                console.log('Added to history:', input, 'History:', terminalHistory);
                processInput(input);
            }
        } else {
            appendToTerminal('<br><span class="prompt">?- </span>');
        }
    } else if (event.key === 'ArrowUp') {
        console.log('ArrowUp detected! History:', terminalHistory, 'Index:', historyIndex);
        event.preventDefault();
        event.stopPropagation();
        navigateHistory(-1);
        return false;
    } else if (event.key === 'ArrowDown') {
        console.log('ArrowDown detected! History:', terminalHistory, 'Index:', historyIndex);
        event.preventDefault();
        event.stopPropagation();
        navigateHistory(1);
        return false;
    }
}

function appendToTerminal(content) {
    const terminal = document.getElementById('terminal');
    terminal.innerHTML += content;
    
    // Ensure cursor is at the end
    const range = document.createRange();
    const sel = window.getSelection();
    range.selectNodeContents(terminal);
    range.collapse(false);
    sel.removeAllRanges();
    sel.addRange(range);
    
    // Scroll to bottom after a short delay to ensure content is rendered
    setTimeout(() => {
        terminal.scrollTop = terminal.scrollHeight;
    }, 10);
}

function processInput(input) {
    if (!input.trim()) return;
    
    appendToTerminal('<br>');
    
    if (input === 'help') {
        // Show help text directly
        showHelp();
        return;
    }
    
    if (input === 'clear') {
        clearTerminal();
        return;
    }
    
    if (input === 'sessions') {
        listSessions();
        return;
    }
    
    if (input === 'clear_cache') {
        clearCache();
        return;
    }
    
    if (!currentSession) {
        appendToTerminal('<span class="error">No session selected! Please select or create a session first.</span><br>');
        appendToTerminal('<span class="prompt">?- </span>');
        return;
    }
    
    executePrologInput(input);
}

function executePrologInput(input) {
    const trimmed = input.trim();
    
    // Check if it starts with ?- (query prefix)
    if (trimmed.startsWith('?-')) {
        // Remove ?- prefix and execute as query
        const queryStr = trimmed.substring(2).trim();
        // Remove trailing . if present
        const cleanQuery = queryStr.endsWith('.') ? queryStr.slice(0, -1) : queryStr;
        executeQuery(cleanQuery);
    } else if (trimmed.endsWith('.')) {
        // It's a fact or rule
        if (trimmed.includes(':-')) {
            addRule(trimmed);
        } else {
            addFact(trimmed);
        }
    } else {
        // No period, assume it's a query
        executeQuery(trimmed);
    }
}

function addFact(factStr) {
    const fact = parseFact(factStr.slice(0, -1));
    
    fetch('/api/v1/sessions/' + currentSession.id + '/facts', {
        method: 'POST',
        headers: getHeaders(),
        body: JSON.stringify({ predicate: fact })
    })
    .then(response => response.json())
    .then(data => {
        if (data.error) {
            appendToTerminal('<span class="error">Error: ' + data.error + '</span><br>');
        } else {
            appendToTerminal('<span class="success">Fact added.</span><br>');
        }
        appendToTerminal('<span class="prompt">?- </span>');
    })
    .catch(error => {
        console.error('Fetch error:', error);
        appendToTerminal('<span class="error">Error: ' + error.message + ' (Check console for details)</span><br>');
        appendToTerminal('<span class="prompt">?- </span>');
    });
}

function addRule(ruleStr) {
    const parts = ruleStr.slice(0, -1).split(':-');
    if (parts.length !== 2) {
        appendToTerminal('<span class="error">Invalid rule format</span><br>');
        appendToTerminal('<span class="prompt">?- </span>');
        return;
    }
    
    const head = parseTerm(parts[0].trim());
    const bodyTerms = parseGoals(parts[1].trim());
    
    const rule = {
        head: head,
        body: bodyTerms
    };
    
    fetch('/api/v1/sessions/' + currentSession.id + '/rules', {
        method: 'POST',
        headers: getHeaders(),
        body: JSON.stringify(rule)
    })
    .then(response => response.json())
    .then(data => {
        if (data.error) {
            appendToTerminal('<span class="error">Error: ' + data.error + '</span><br>');
        } else {
            appendToTerminal('<span class="success">Rule added.</span><br>');
        }
        appendToTerminal('<span class="prompt">?- </span>');
    })
    .catch(error => {
        console.error('Fetch error:', error);
        appendToTerminal('<span class="error">Error: ' + error.message + ' (Check console for details)</span><br>');
        appendToTerminal('<span class="prompt">?- </span>');
    });
}

function executeQuery(queryStr) {
    // Parse goals properly, respecting parentheses
    const goals = parseGoals(queryStr);
    
    const query = { goals: goals };
    
    fetch('/api/v1/sessions/' + currentSession.id + '/query', {
        method: 'POST',
        headers: getHeaders(),
        body: JSON.stringify(query)
    })
    .then(response => response.json())
    .then(data => {
        if (data.error) {
            appendToTerminal('<span class="error">Error: ' + data.error + '</span><br>');
        } else {
            displayQueryResults(data.solutions);
        }
        appendToTerminal('<span class="prompt">?- </span>');
    })
    .catch(error => {
        console.error('Fetch error:', error);
        appendToTerminal('<span class="error">Error: ' + error.message + ' (Check console for details)</span><br>');
        appendToTerminal('<span class="prompt">?- </span>');
    });
}

function displayQueryResults(solutions) {
    if (!solutions || solutions.length === 0) {
        appendToTerminal('<span class="warning">No solutions found.</span><br>');
        return;
    }
    
    let successCount = 0;
    solutions.forEach((solution, index) => {
        if (solution.success) {
            successCount++;
            if (solution.bindings && Object.keys(solution.bindings).length > 0) {
                appendToTerminal('<span class="success">Solution ' + successCount + ':</span><br>');
                for (const variable in solution.bindings) {
                    const binding = solution.bindings[variable];
                    appendToTerminal('  ' + variable + ' = ' + formatTerm(binding) + '<br>');
                }
            } else {
                appendToTerminal('<span class="success">Yes (' + successCount + ')</span><br>');
            }
        }
    });
    
    if (successCount === 0) {
        appendToTerminal('<span class="warning">No successful solutions.</span><br>');
    } else {
        appendToTerminal('<span class="success">Found ' + successCount + ' solution(s).</span><br>');
    }
}

function parseTerm(str) {
    str = str.trim();
    
    if (/^[A-Z_][a-zA-Z0-9_]*$/.test(str)) {
        return { type: 'variable', value: str };
    }
    
    if (/^-?\d+(\.\d+)?$/.test(str)) {
        return { type: 'number', value: parseFloat(str) };
    }
    
    // First try to match operators like =, <, >, etc.
    const operatorMatch = str.match(/^([=<>!@#$%^&*+\-\/\\]+)\((.*)\)$/);
    if (operatorMatch) {
        const functor = operatorMatch[1];
        const argsStr = operatorMatch[2];
        const args = parseArgs(argsStr);
        return { type: 'compound', value: functor, args: args };
    }
    
    // Then try regular functors
    const match = str.match(/^([a-z][a-zA-Z0-9_]*)\((.*)\)$/);
    if (match) {
        const functor = match[1];
        const argsStr = match[2];
        const args = parseArgs(argsStr);
        return { type: 'compound', value: functor, args: args };
    }
    
    return { type: 'atom', value: str };
}

function parseArgs(argsStr) {
    if (!argsStr.trim()) return [];
    
    const args = [];
    let current = '';
    let parenLevel = 0;
    
    for (let i = 0; i < argsStr.length; i++) {
        const char = argsStr[i];
        if (char === ',' && parenLevel === 0) {
            args.push(parseTerm(current.trim()));
            current = '';
        } else {
            if (char === '(') parenLevel++;
            if (char === ')') parenLevel--;
            current += char;
        }
    }
    
    if (current.trim()) {
        args.push(parseTerm(current.trim()));
    }
    
    return args;
}

function parseGoals(goalsStr) {
    if (!goalsStr.trim()) return [];
    
    const goals = [];
    let current = '';
    let parenLevel = 0;
    
    for (let i = 0; i < goalsStr.length; i++) {
        const char = goalsStr[i];
        if (char === ',' && parenLevel === 0) {
            goals.push(parseTerm(current.trim()));
            current = '';
        } else {
            if (char === '(') parenLevel++;
            if (char === ')') parenLevel--;
            current += char;
        }
    }
    
    if (current.trim()) {
        goals.push(parseTerm(current.trim()));
    }
    
    return goals;
}

function parseFact(str) {
    return parseTerm(str);
}

function formatTerm(term) {
    if (!term) return 'null';
    
    switch (term.type) {
        case 'atom':
            return term.value;
        case 'variable':
            return term.value;
        case 'number':
            return term.value.toString();
        case 'compound':
            const args = term.args ? term.args.map(formatTerm).join(', ') : '';
            return term.value + '(' + args + ')';
        default:
            return JSON.stringify(term);
    }
}

function showHelp() {
    appendToTerminal('<span class="success">Available commands:</span><br>' +
        '  help          - Show this help<br>' +
        '  clear         - Clear terminal<br>' +
        '  sessions      - List all sessions<br><br>' +
        '<span class="success">Prolog syntax:</span><br>' +
        '  fact(atom).                    - Add a fact<br>' +
        '  rule(X) :- condition(X).       - Add a rule<br>' +
        '  ?- query(X)                    - Execute a query<br><br>' +
        '<span class="success">Examples:</span><br>' +
        '  parent(tom, bob).              - Tom is parent of Bob<br>' +
        '  grandparent(X,Z) :- parent(X,Y), parent(Y,Z).<br>' +
        '  ?- parent(tom, X)              - Find children of Tom<br>' +
        '  ?- count(_, parent(X,Y), N)    - Count parent relationships<br><br>' +
        '<span class="success">Built-ins:</span><br>' +
        '  =(X, value), atom(X), var(X), number(X)<br>' +
        '  now(X), count(..), sum(..), max(..), min(..)<br><br>');
    appendToTerminal('<span class="prompt">?- </span>');
}

function clearTerminal() {
    const terminal = document.getElementById('terminal');
    terminal.innerHTML = '';
    
    if (currentSession) {
        terminal.innerHTML = '<div><span class="success">Session: ' + currentSession.name + '</span></div>';
    } else {
        terminal.innerHTML = '<div><span class="warning">No session selected. Create or select a session.</span></div>';
    }
    
    terminal.innerHTML += '<div><span class="prompt">?- </span></div>';
    
    // Position cursor after the prompt
    const range = document.createRange();
    const sel = window.getSelection();
    range.selectNodeContents(terminal);
    range.collapse(false);  // false = collapse to end
    sel.removeAllRanges();
    sel.addRange(range);
    
    terminal.focus();
}

function navigateHistory(direction) {
    console.log('navigateHistory called with direction:', direction);
    console.log('Current history:', terminalHistory);
    console.log('Current index:', historyIndex);
    
    if (terminalHistory.length === 0) {
        console.log('No history items!');
        return;
    }
    
    historyIndex += direction;
    if (historyIndex < 0) historyIndex = 0;
    if (historyIndex >= terminalHistory.length) historyIndex = terminalHistory.length;
    
    console.log('New index:', historyIndex);
    
    const terminal = document.getElementById('terminal');
    const content = terminal.innerHTML;
    const promptText = '<span class="prompt">?- </span>';
    const lastPromptIndex = content.lastIndexOf(promptText);
    
    if (lastPromptIndex !== -1) {
        // Get everything before and including the prompt
        const beforePrompt = content.substring(0, lastPromptIndex + promptText.length);
        // Get the history item or empty string if at the end
        const historyItem = historyIndex < terminalHistory.length ? terminalHistory[historyIndex] : '';
        // Update terminal content
        terminal.innerHTML = beforePrompt + historyItem;
        
        // Position cursor at the end
        const range = document.createRange();
        const sel = window.getSelection();
        range.selectNodeContents(terminal);
        range.collapse(false);
        sel.removeAllRanges();
        sel.addRange(range);
        
        // Ensure terminal is focused
        terminal.focus();
    }
}

function loadSessions() {
    fetch('/api/v1/sessions', { headers: getHeaders() })
        .then(response => response.json())
        .then(data => {
            if (data.sessions) {
                sessions = data.sessions;
                renderSessions();
                if (sessions.length > 0 && !currentSession) {
                    selectSession(sessions[0]);
                }
            }
        })
        .catch(error => {
            console.error('Error loading sessions:', error);
        });
}

function renderSessions() {
    const sessionList = document.getElementById('sessionList');
    sessionList.innerHTML = '';
    
    sessions.forEach(session => {
        const sessionDiv = document.createElement('div');
        sessionDiv.className = 'session-item';
        sessionDiv.onclick = () => selectSession(session);
        
        if (currentSession && currentSession.id === session.id) {
            sessionDiv.classList.add('active');
        }
        
        sessionDiv.innerHTML = 
            '<div class="session-name">' + session.name + '</div>' +
            '<div class="session-desc">' + (session.description || 'No description') + '</div>';
        
        sessionList.appendChild(sessionDiv);
    });
}

function selectSession(session) {
    currentSession = session;
    renderSessions();
    updateCurrentSessionInfo();
    appendToTerminal('<br><span class="success">Session switched to: ' + session.name + '</span><br>');
    appendToTerminal('<span class="prompt">?- </span>');
}

function updateCurrentSessionInfo() {
    const info = document.getElementById('currentSessionInfo');
    if (currentSession) {
        info.innerHTML = 
            '<strong>' + currentSession.name + '</strong><br>' +
            '<small>' + (currentSession.description || 'No description') + '</small><br>' +
            '<small>Created: ' + new Date(currentSession.created_at).toLocaleDateString() + '</small>';
    } else {
        info.innerHTML = 'No session selected';
    }
}

function showCreateSessionModal() {
    document.getElementById('createSessionModal').style.display = 'block';
    document.getElementById('sessionName').focus();
}

function hideCreateSessionModal() {
    document.getElementById('createSessionModal').style.display = 'none';
    document.getElementById('sessionName').value = '';
    document.getElementById('sessionDesc').value = '';
}

function createSession() {
    const name = document.getElementById('sessionName').value.trim();
    const description = document.getElementById('sessionDesc').value.trim();
    
    if (!name) {
        alert('Session name is required');
        return;
    }
    
    const sessionRequest = { name: name, description: description };
    
    fetch('/api/v1/sessions', {
        method: 'POST',
        headers: getHeaders(),
        body: JSON.stringify(sessionRequest)
    })
    .then(response => response.json())
    .then(data => {
        if (data.error) {
            alert('Error creating session: ' + data.error);
        } else {
            hideCreateSessionModal();
            // Automatically select the new session
            selectSession(data);
            loadSessions();
            appendToTerminal('<br><span class="success">Session "' + name + '" created and activated!</span><br>');
            appendToTerminal('<span class="prompt">?- </span>');
        }
    })
    .catch(error => {
        alert('Error creating session: ' + error.message);
    });
}

function deleteCurrentSession() {
    if (!currentSession) {
        alert('No session selected');
        return;
    }
    
    if (!confirm('Delete session "' + currentSession.name + '"? This cannot be undone.')) {
        return;
    }
    
    fetch('/api/v1/sessions/' + currentSession.id, {
        method: 'DELETE',
        headers: getHeaders()
    })
    .then(response => response.json())
    .then(data => {
        if (data.error) {
            alert('Error deleting session: ' + data.error);
        } else {
            appendToTerminal('<br><span class="warning">Session "' + currentSession.name + '" deleted.</span><br>');
            const deletedSessionId = currentSession.id;
            currentSession = null;
            
            // Remove the deleted session from the local array
            sessions = sessions.filter(s => s.id !== deletedSessionId);
            
            // Re-render sessions and update UI
            renderSessions();
            updateCurrentSessionInfo();
            
            // Clear terminal and show appropriate message
            const terminal = document.getElementById('terminal');
            if (sessions.length === 0) {
                terminal.innerHTML = '<div><span class="warning">No sessions available. Create a new session to start.</span></div>';
                terminal.innerHTML += '<div><span class="prompt">?- </span></div>';
            } else {
                // Auto-select the first remaining session
                selectSession(sessions[0]);
                appendToTerminal('<span class="prompt">?- </span>');
            }
        }
    })
    .catch(error => {
        alert('Error deleting session: ' + error.message);
    });
}

function clearCache() {
    fetch('/api/v1/cache/clear', {
        method: 'POST',
        headers: getHeaders()
    })
    .then(response => response.json())
    .then(data => {
        if (data.error) {
            appendToTerminal('<span class="error">Error clearing cache: ' + data.error + '</span><br>');
        } else {
            appendToTerminal('<span class="success">Cache cleared.</span><br>');
        }
        appendToTerminal('<span class="prompt">?- </span>');
    })
    .catch(error => {
        appendToTerminal('<span class="error">Error clearing cache: ' + error.message + '</span><br>');
        appendToTerminal('<span class="prompt">?- </span>');
    });
}

function listSessions() {
    appendToTerminal('<span class="success">Available sessions:</span><br>');
    sessions.forEach((session, index) => {
        const marker = currentSession && currentSession.id === session.id ? ' [CURRENT]' : '';
        appendToTerminal('  ' + (index + 1) + '. ' + session.name + ' - ' + (session.description || 'No description') + marker + '<br>');
    });
    appendToTerminal('<span class="prompt">?- </span>');
}

function toggleHelp() {
    const sidebar = document.getElementById('helpSidebar');
    sidebar.classList.toggle('collapsed');
}

function toggleTheme() {
    const body = document.body;
    const themeIcon = document.getElementById('themeIcon');
    
    if (body.classList.contains('light-theme')) {
        body.classList.remove('light-theme');
        themeIcon.textContent = '🌙';
        localStorage.setItem('theme', 'dark');
    } else {
        body.classList.add('light-theme');
        themeIcon.textContent = '☀️';
        localStorage.setItem('theme', 'light');
    }
}

// Load saved theme preference
document.addEventListener('DOMContentLoaded', function() {
    const savedTheme = localStorage.getItem('theme');
    if (savedTheme === 'light') {
        document.body.classList.add('light-theme');
        document.getElementById('themeIcon').textContent = '☀️';
    }
});

function insertExample(example) {
    const terminal = document.getElementById('terminal');
    const content = terminal.innerHTML;
    const promptText = '<span class="prompt">?- </span>';
    const lastPromptIndex = content.lastIndexOf(promptText);
    
    if (lastPromptIndex !== -1) {
        const beforePrompt = content.substring(0, lastPromptIndex + promptText.length);
        terminal.innerHTML = beforePrompt + example;
        
        const range = document.createRange();
        const sel = window.getSelection();
        range.selectNodeContents(terminal);
        range.collapse(false);
        sel.removeAllRanges();
        sel.addRange(range);
        
        terminal.focus();
    }
}

window.onclick = function(event) {
    const modal = document.getElementById('createSessionModal');
    if (event.target === modal) {
        hideCreateSessionModal();
    }
}

document.addEventListener('keydown', function(event) {
    if (event.target.id === 'sessionName' || event.target.id === 'sessionDesc') {
        if (event.key === 'Enter') {
            createSession();
        }
    }
});

// Tab functionality
function showTab(tabName) {
    // Hide all tab contents
    document.getElementById('helpTab').style.display = 'none';
    document.getElementById('tutorialTab').style.display = 'none';
    
    // Remove active class from all buttons
    document.querySelectorAll('.tab-button').forEach(btn => btn.classList.remove('active'));
    
    // Show selected tab and activate button
    if (tabName === 'help') {
        document.getElementById('helpTab').style.display = 'block';
        document.querySelector('.tab-button:first-child').classList.add('active');
    } else if (tabName === 'tutorial') {
        document.getElementById('tutorialTab').style.display = 'block';
        document.querySelector('.tab-button:last-child').classList.add('active');
    }
}
`