package main

const loginTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Login - Jacob's Flight Tracker</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background: #f5f5f5;
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
            margin: 0;
        }
        .login-container {
            background: white;
            padding: 40px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            width: 100%;
            max-width: 400px;
        }
        h1 {
            margin: 0 0 30px 0;
            color: #333;
            text-align: center;
            font-size: 32px;
        }
        .form-group {
            margin-bottom: 20px;
        }
        label {
            display: block;
            margin-bottom: 8px;
            color: #555;
            font-weight: 500;
        }
        input {
            width: 100%;
            padding: 12px;
            border: 2px solid #ddd;
            border-radius: 4px;
            font-size: 16px;
            box-sizing: border-box;
            transition: border-color 0.3s;
        }
        input:focus {
            outline: none;
            border-color: #007bff;
        }
        button {
            width: 100%;
            padding: 14px;
            background: #007bff;
            color: white;
            border: none;
            border-radius: 4px;
            font-size: 16px;
            font-weight: 600;
            cursor: pointer;
            transition: background 0.3s;
        }
        button:hover {
            background: #0056b3;
        }
        .error {
            background: #fee;
            color: #c33;
            padding: 12px;
            border-radius: 4px;
            margin-bottom: 20px;
            border: 1px solid #fcc;
        }
        .info {
            margin-top: 20px;
            padding: 15px;
            background: #f8f9fa;
            border-radius: 4px;
            font-size: 14px;
            color: #666;
            border: 1px solid #dee2e6;
        }
        .info strong {
            color: #333;
        }
        .info code {
            background: #e9ecef;
            padding: 2px 6px;
            border-radius: 3px;
            font-family: monospace;
        }
    </style>
</head>
<body>
    <div class="login-container">
        <h1>Jacob's Flight Tracker</h1>
        {{if .Error}}
        <div class="error">{{.Error}}</div>
        {{end}}
        <form method="POST" action="/login">
            <div class="form-group">
                <label for="username">Username</label>
                <input type="text" id="username" name="username" required autofocus>
            </div>
            <div class="form-group">
                <label for="password">Password</label>
                <input type="password" id="password" name="password" required>
            </div>
            <button type="submit">Login</button>
        </form>
    </div>
</body>
</html>
`

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Shuttle Flight Tracker</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 1000px;
            margin: 30px auto;
            padding: 20px;
            background: #f5f5f5;
        }
        .header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 30px;
        }
        h1 {
            font-size: 36px;
            margin: 0;
            color: #333;
        }
        .logout-btn {
            padding: 10px 20px;
            background: #dc3545;
            color: white;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            text-decoration: none;
            font-size: 14px;
        }
        .logout-btn:hover {
            background: #c82333;
        }
        .add-flight {
            background: white;
            padding: 20px;
            border-radius: 8px;
            margin-bottom: 30px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .form-row {
            display: flex;
            gap: 15px;
            align-items: center;
            flex-wrap: wrap;
        }
        .checkbox-group {
            display: flex;
            gap: 20px;
            padding: 12px 15px;
            background: #f8f9fa;
            border-radius: 4px;
            border: 2px solid #ddd;
            height: 48px;
            box-sizing: border-box;
        }
        .checkbox-group label {
            display: flex;
            align-items: center;
            gap: 6px;
            cursor: pointer;
            font-size: 15px;
            user-select: none;
        }
        .checkbox-group input[type="checkbox"] {
            width: 18px;
            height: 18px;
            cursor: pointer;
        }
        input {
            padding: 12px;
            font-size: 16px;
            border: 2px solid #ddd;
            border-radius: 4px;
        }
        input[type="text"] {
            width: 200px;
        }
        input[type="number"] {
            width: 80px;
        }
        button {
            padding: 12px 24px;
            font-size: 16px;
            background: #007bff;
            color: white;
            border: none;
            border-radius: 4px;
            cursor: pointer;
        }
        button:hover {
            background: #0056b3;
        }
        .remove-btn {
            background: #6c757d;
            padding: 8px 16px;
            font-size: 14px;
        }
        .remove-btn:hover {
            background: #5a6268;
        }
        .error {
            color: #dc3545;
            padding: 10px;
            margin: 10px 0;
        }
        .flights-grid {
            display: grid;
            gap: 15px;
        }
        .flight-row {
            position: relative;
        }
        .note-container {
            position: absolute;
            right: 100%;
            margin-right: 15px;
            top: 50%;
            transform: translateY(-50%);
            width: 180px;
        }
        .note-bubble {
            position: relative;
            background: #fff9e6;
            border: 2px solid #ffd700;
            border-radius: 12px;
            padding: 12px;
            font-size: 13px;
            color: #333;
            min-height: 40px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            word-wrap: break-word;
        }
        .note-bubble:after {
            content: '';
            position: absolute;
            right: -10px;
            top: 50%;
            transform: translateY(-50%);
            width: 0;
            height: 0;
            border-left: 10px solid #ffd700;
            border-top: 8px solid transparent;
            border-bottom: 8px solid transparent;
        }
        .note-bubble:before {
            content: '';
            position: absolute;
            right: -7px;
            top: 50%;
            transform: translateY(-50%);
            width: 0;
            height: 0;
            border-left: 8px solid #fff9e6;
            border-top: 6px solid transparent;
            border-bottom: 6px solid transparent;
            z-index: 1;
        }
        .note-bubble.empty {
            background: #f8f9fa;
            border: 2px dashed #dee2e6;
            color: #999;
            font-style: italic;
        }
        .note-bubble.empty:after {
            border-left-color: #dee2e6;
        }
        .note-bubble.empty:before {
            border-left-color: #f8f9fa;
        }
        .note-edit-form {
            position: absolute;
            right: 100%;
            margin-right: 15px;
            top: 50%;
            transform: translateY(-50%);
            width: 180px;
            z-index: 10;
            display: none;
        }
        .flight-row:hover .note-edit-form {
            display: block;
        }
        .flight-row:hover .note-container {
            display: none;
        }
        .note-input {
            width: 100%;
            padding: 12px;
            font-size: 13px;
            border: 2px solid #ffd700;
            border-radius: 12px;
            font-family: Arial, sans-serif;
            background: #fff9e6;
            min-height: 60px;
            box-shadow: 0 4px 8px rgba(0,0,0,0.15);
            resize: vertical;
        }
        .note-input:focus {
            outline: none;
            border-color: #ffb700;
        }
        .flight-card {
            background: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            display: grid;
            grid-template-columns: 150px 1fr 250px 100px;
            align-items: center;
            gap: 20px;
        }
        .flight-number {
            font-size: 32px;
            font-weight: bold;
            color: #333;
        }
        .flight-details {
            font-size: 14px;
            color: #666;
        }
        .flight-details p {
            margin: 5px 0;
        }
        .arrival-time {
            text-align: right;
        }
        .expected-time {
            font-size: 42px;
            font-weight: bold;
            color: #28a745;
            line-height: 1;
        }
        .expected-time.delayed {
            color: #ffc107;
        }
        .scheduled-time {
            font-size: 14px;
            color: #999;
            margin-top: 5px;
        }
        .badge {
            display: inline-block;
            padding: 4px 8px;
            border-radius: 4px;
            font-size: 12px;
            font-weight: bold;
            text-transform: uppercase;
            margin-right: 5px;
        }
        .badge.pickup {
            background: #d1ecf1;
            color: #0c5460;
        }
        .badge.dropoff {
            background: #fff3cd;
            color: #856404;
        }
        .badge.both {
            background: #d4edda;
            color: #155724;
        }
        .badge.active {
            background: #d4edda;
            color: #155724;
        }
        .badge.landed {
            background: #cce5ff;
            color: #004085;
        }
        .badge.scheduled {
            background: #e2e3e5;
            color: #383d41;
        }
        .crew-count {
            font-weight: bold;
            color: #007bff;
        }
        .empty-state {
            text-align: center;
            padding: 60px;
            color: #999;
            font-size: 18px;
        }
        .auto-refresh {
            color: #28a745;
            font-size: 12px;
            margin-left: 10px;
        }
    </style>
    <script>
        {{if not .IsDemo}}
        // Auto-refresh every 5 minutes (disabled for demo users)
        setTimeout(function() {
            location.reload();
        }, 300000);
        {{end}}

        function saveNote(textarea) {
            const form = textarea.closest('form');
            const formData = new FormData(form);
            fetch('/update-note', {
                method: 'POST',
                body: formData
            }).then(() => {
                location.reload();
            });
        }
    </script>
</head>
<body>
    <div class="header">
        <h1>Today's Flights
            {{if .IsDemo}}
            <span style="color: #ffc107; font-size: 16px;">(Demo Mode - Sample Data)</span>
            {{else}}
            <span class="auto-refresh">● Auto-refresh: 5 min</span>
            {{end}}
        </h1>
        <a href="/logout" class="logout-btn">Logout</a>
    </div>

    <div class="add-flight">
        <form method="POST" action="/add">
            <div class="form-row">
                <input type="text" name="flight_number" placeholder="Flight # (e.g., AA100)" required>
                <div class="checkbox-group">
                    <label>
                        <input type="checkbox" name="is_pickup" checked>
                        <span>Pickup</span>
                    </label>
                    <label>
                        <input type="checkbox" name="is_dropoff">
                        <span>Dropoff</span>
                    </label>
                </div>
                <input type="number" name="crew_count" placeholder="# Crew" min="1" value="1" required>
                <button type="submit">Add Flight</button>
            </div>
        </form>
        {{if .Error}}
        <div class="error">{{.Error}}</div>
        {{end}}
    </div>

    {{if .Flights}}
    <div class="flights-grid">
        {{range .Flights}}
        <div class="flight-row">
            <div class="note-container">
                <div class="{{if .Note}}note-bubble{{else}}note-bubble empty{{end}}">
                    {{if .Note}}{{.Note}}{{else}}Click to add note{{end}}
                </div>
            </div>
            <form method="POST" action="/update-note" class="note-edit-form">
                <input type="hidden" name="id" value="{{.ID}}">
                <textarea name="note" class="note-input" placeholder="Add a note..." onblur="saveNote(this)">{{.Note}}</textarea>
            </form>
            <div class="flight-card">
                <div class="flight-number">{{.FlightNumber}}</div>
                <div class="flight-details">
                    <p><strong>{{.Airline}}</strong></p>
                    <p>
                        <span class="badge {{.Status}}">{{.Status}}</span>
                    </p>
                    <p><span class="crew-count">{{.CrewCount}} crew</span></p>
                    {{if .IsDelayed}}
                    <p style="color: #ffc107; font-weight: bold;">+{{.Delay}} min delay</p>
                    {{end}}
                </div>
                <div class="arrival-time">
                    <div class="expected-time {{if .IsDelayed}}delayed{{end}}">{{.ExpectedArrival}}</div>
                    <div class="scheduled-time">scheduled: {{.ScheduledArrival}}</div>
                    <div style="margin-top: 10px;">
                        <span class="badge {{.Type}}">{{.Type}}</span>
                    </div>
                </div>
                <div>
                    <form method="POST" action="/remove">
                        <input type="hidden" name="id" value="{{.ID}}">
                        <button type="submit" class="remove-btn">Remove</button>
                    </form>
                </div>
            </div>
        </div>
        {{end}}
    </div>
    {{else}}
    <div class="empty-state">
        No flights added yet. Add a flight number above to get started.
    </div>
    {{end}}
</body>
</html>
`
