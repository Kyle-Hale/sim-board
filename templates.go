package main

const leaderboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Lap Time Leaderboard</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: 'Helvetica Neue', Arial, sans-serif;
            background: #0A0E27;
            min-height: 100vh;
            padding: 24px;
            color: #FFFFFF;
            position: relative;
            overflow: hidden;
        }
        
        body::before {
            content: '';
            position: fixed;
            top: -50%;
            left: -50%;
            width: 200%;
            height: 200%;
            background: linear-gradient(180deg, #0A0E27 0%, #1A1F3A 50%, #0A0E27 100%);
            animation: backgroundShift 20s ease-in-out infinite;
            z-index: 0;
        }
        
        @keyframes backgroundShift {
            0%, 100% {
                transform: translate(0, 0) scale(1);
                opacity: 1;
            }
            25% {
                transform: translate(-2%, -2%) scale(1.05);
                opacity: 0.95;
            }
            50% {
                transform: translate(2%, 2%) scale(1.1);
                opacity: 0.9;
            }
            75% {
                transform: translate(-1%, 1%) scale(1.05);
                opacity: 0.95;
            }
        }
        
        body > * {
            position: relative;
            z-index: 1;
        }
        
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: #1A1F3A;
            border-radius: 0;
            border: 2px solid #ED1C24;
            box-shadow: 0 0 10px rgba(237, 28, 36, 0.15);
            overflow: hidden;
        }
        
        .header {
            background: linear-gradient(135deg, #0A0E27 0%, #1A1F3A 50%, #0A0E27 100%);
            border-bottom: 3px solid #ED1C24;
            color: #FFFFFF;
            padding: 32px 40px;
            text-align: center;
            position: relative;
        }
        
        .header::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background: linear-gradient(90deg, transparent 0%, rgba(237, 28, 36, 0.05) 50%, transparent 100%);
            pointer-events: none;
        }
        
        .header h1 {
            font-size: 2.5em;
            margin-bottom: 16px;
            font-weight: 900;
            color: #FFFFFF;
            text-transform: uppercase;
            letter-spacing: 3px;
            position: relative;
            z-index: 1;
            text-shadow: 0 2px 10px rgba(0, 0, 0, 0.5);
        }
        
        .header .track-name {
            font-size: 3em;
            font-weight: 900;
            letter-spacing: 2px;
            margin: 16px 0 0 0;
            color: #ED1C24;
            position: relative;
            z-index: 1;
            text-transform: uppercase;
            text-shadow: 0 0 8px rgba(237, 28, 36, 0.3);
            display: inline-block;
        }
        
        .header .track-car {
            font-size: 0.5em;
            font-weight: 700;
            color: rgba(255, 255, 255, 0.6);
            margin-right: 1px;
            vertical-align: baseline;
            text-transform: uppercase;
            letter-spacing: 2px;
            display: inline-block;
        }
        
        .header .track-at {
            font-size: 0.35em;
            font-weight: 400;
            color: rgba(255, 255, 255, 0.4);
            margin: 0 1px;
            vertical-align: baseline;
            text-transform: lowercase;
            display: inline-block;
        }
        
        }
        
        .leaderboard {
            padding: 32px;
            background: #1A1F3A;
        }
        
        .leaderboard-header {
            display: grid;
            grid-template-columns: 80px 1fr 200px;
            gap: 20px;
            padding: 20px 24px;
            background: #0A0E27;
            border: 2px solid #ED1C24;
            border-radius: 0;
            margin-bottom: 20px;
            font-weight: 900;
            color: #ED1C24;
            text-transform: uppercase;
            font-size: 0.875em;
            letter-spacing: 3px;
            box-shadow: 0 0 10px rgba(237, 28, 36, 0.15);
        }
        
        .top-three {
            margin-bottom: 16px;
        }
        
        .scrollable-entries {
            max-height: calc(100vh - 400px);
            overflow-y: auto;
            overflow-x: hidden;
        }
        
        .scrollable-entries::-webkit-scrollbar {
            width: 8px;
        }
        
        .scrollable-entries::-webkit-scrollbar-track {
            background: #1a1a1a;
            border-radius: 4px;
        }
        
        .scrollable-entries::-webkit-scrollbar-thumb {
            background: #ED1C24;
            border-radius: 4px;
        }
        
        .scrollable-entries::-webkit-scrollbar-thumb:hover {
            background: #FF2D3A;
        }
        
        .entry {
            display: grid;
            grid-template-columns: 80px 1fr 200px;
            gap: 20px;
            padding: 24px 28px;
            margin-bottom: 8px;
            background: #0A0E27;
            border: 1px solid rgba(237, 28, 36, 0.2);
            border-radius: 0;
            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.4);
            transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
            border-left: 5px solid transparent;
            cursor: pointer;
        }
        
        .entry:hover {
            transform: translateX(4px);
            box-shadow: 0 4px 20px rgba(237, 28, 36, 0.3);
            background: #1A1F3A;
            border-color: rgba(237, 28, 36, 0.4);
        }
        
        .scrollable-entries .entry {
            padding: 12px 20px;
            margin-bottom: 6px;
        }
        
        .entry.position-1 {
            border-left-color: #ED1C24;
            background: linear-gradient(90deg, rgba(237, 28, 36, 0.15) 0%, #0A0E27 100%);
            box-shadow: 0 0 15px rgba(237, 28, 36, 0.2);
            border: 1px solid #ED1C24;
        }
        
        .entry.position-1:hover {
            background: linear-gradient(90deg, rgba(237, 28, 36, 0.25) 0%, #1A1F3A 100%);
            box-shadow: 0 0 20px rgba(237, 28, 36, 0.25);
        }
        
        .entry.position-2 {
            border-left-color: #C0C0C0;
            background: linear-gradient(90deg, rgba(192, 192, 192, 0.1) 0%, #0A0E27 100%);
            border: 1px solid rgba(192, 192, 192, 0.3);
        }
        
        .entry.position-3 {
            border-left-color: #CD7F32;
            background: linear-gradient(90deg, rgba(205, 127, 50, 0.1) 0%, #0A0E27 100%);
            border: 1px solid rgba(205, 127, 50, 0.3);
        }
        
        .position {
            font-size: 2em;
            font-weight: 900;
            color: #ED1C24;
            display: flex;
            align-items: center;
            justify-content: center;
            text-shadow: 0 0 5px rgba(237, 28, 36, 0.3);
        }
        
        .scrollable-entries .position {
            font-size: 1.5em;
        }
        
        .position-1 .position {
            color: #ED1C24;
            text-shadow: 0 0 8px rgba(237, 28, 36, 0.4);
        }
        
        .position-2 .position {
            color: #C0C0C0;
            text-shadow: 0 0 5px rgba(192, 192, 192, 0.3);
        }
        
        .position-3 .position {
            color: #CD7F32;
            text-shadow: 0 0 5px rgba(205, 127, 50, 0.3);
        }
        
        .username {
            font-size: 1.5em;
            font-weight: 900;
            color: #FFFFFF;
            display: flex;
            align-items: center;
            justify-content: flex-start;
            letter-spacing: 1px;
            text-transform: uppercase;
            position: relative;
        }
        
        .scrollable-entries .username {
            font-size: 1.5em;
        }
        
        .time {
            font-size: 1.5em;
            font-weight: 900;
            color: #ED1C24;
            font-family: 'Courier New', monospace;
            display: flex;
            align-items: center;
            justify-content: flex-end;
            letter-spacing: 2px;
            text-shadow: 0 0 5px rgba(237, 28, 36, 0.3);
        }
        
        .scrollable-entries .time {
            font-size: 1.2em;
        }
        
        .empty-state {
            text-align: center;
            padding: 60px 20px;
            color: #666;
        }
        
        .empty-state h2 {
            font-size: 2em;
            margin-bottom: 10px;
            color: #ED1C24;
            font-weight: 900;
            text-transform: uppercase;
            letter-spacing: 2px;
        }
        
        .admin-link {
            position: fixed;
            bottom: 24px;
            right: 24px;
            background: #0A0E27;
            border: 2px solid #ED1C24;
            padding: 12px 24px;
            border-radius: 0;
            text-decoration: none;
            color: #ED1C24;
            font-weight: 900;
            font-size: 0.875em;
            text-transform: uppercase;
            letter-spacing: 2px;
            box-shadow: 0 0 15px rgba(237, 28, 36, 0.25);
            transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
            z-index: 10;
        }
        
        .admin-link:hover {
            background: #ED1C24;
            color: #FFFFFF;
            border-color: #ED1C24;
            box-shadow: 0 0 20px rgba(237, 28, 36, 0.4);
            transform: translateY(-2px);
        }
        
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>{{.Title}}</h1>
            <div class="track-name">{{.Track.Name}}</div>
        </div>
        <div class="leaderboard">
            <div class="leaderboard-header">
                <div>Rank</div>
                <div>Driver</div>
                <div style="text-align: right;">Time</div>
            </div>
            {{if .Entries}}
                <div class="top-three">
                    {{range $index, $entry := .Entries}}
                        {{if le $entry.Position 3}}
                        <div class="entry position-{{$entry.Position}}">
                            <div class="position">#{{$entry.Position}}</div>
                            <div class="username">
                                <span>{{$entry.Username}}</span>
                                {{if $entry.Car}}<span style="font-size: 0.7em; color: rgba(255, 255, 255, 0.5); font-weight: 400; text-transform: none; position: absolute; left: 50%; transform: translateX(-50%);">{{$entry.Car}}</span>{{end}}
                            </div>
                            <div class="time" style="display: flex; flex-direction: column; align-items: flex-end; gap: 2px;">
                                <span>{{$entry.Time}}</span>
                                {{if $.ShowAssistsLeaderboard}}
                                <span style="font-size: 0.6em; color: rgba(255, 255, 255, 0.5); font-weight: 500; text-transform: uppercase; letter-spacing: 0.6px; line-height: 1.1; white-space: nowrap; margin-top: 2px;">
                                    <span style="color: rgba(255, 255, 255, 0.4);">ABS:</span> {{if $entry.ABS}}<span style="color: rgba(237, 28, 36, 0.8);">ON</span>{{else}}<span style="color: rgba(76, 175, 80, 0.9);">OFF</span>{{end}}<span style="margin: 0 2px; color: rgba(255, 255, 255, 0.2);">|</span>
                                    <span style="color: rgba(255, 255, 255, 0.4);">TRANS:</span> {{if $entry.AutoTransmission}}<span style="color: rgba(237, 28, 36, 0.8);">AUTO</span>{{else}}<span style="color: rgba(76, 175, 80, 0.9);">MANUAL</span>{{end}}<span style="margin: 0 2px; color: rgba(255, 255, 255, 0.2);">|</span>
                                    <span style="color: rgba(255, 255, 255, 0.4);">TC:</span> {{if $entry.TractionControl}}<span style="color: rgba(237, 28, 36, 0.8);">ON</span>{{else}}<span style="color: rgba(76, 175, 80, 0.9);">OFF</span>{{end}}
                                </span>
                                {{end}}
                            </div>
                        </div>
                        {{end}}
                    {{end}}
                </div>
                {{$hasMore := gt (len .Entries) 3}}
                {{if $hasMore}}
                <div class="scrollable-entries" id="scrollable-leaderboard">
                    {{range $index, $entry := .Entries}}
                        {{if gt $entry.Position 3}}
                        <div class="entry position-{{$entry.Position}}">
                            <div class="position">#{{$entry.Position}}</div>
                            <div class="username">
                                <span>{{$entry.Username}}</span>
                                {{if $entry.Car}}<span style="font-size: 0.7em; color: rgba(255, 255, 255, 0.5); font-weight: 400; text-transform: none; position: absolute; left: 50%; transform: translateX(-50%);">{{$entry.Car}}</span>{{end}}
                            </div>
                            <div class="time" style="display: flex; flex-direction: column; align-items: flex-end; gap: 2px;">
                                <span>{{$entry.Time}}</span>
                                {{if $.ShowAssistsLeaderboard}}
                                <span style="font-size: 0.6em; color: rgba(255, 255, 255, 0.5); font-weight: 500; text-transform: uppercase; letter-spacing: 0.6px; line-height: 1.1; white-space: nowrap; margin-top: 2px;">
                                    <span style="color: rgba(255, 255, 255, 0.4);">ABS:</span> {{if $entry.ABS}}<span style="color: rgba(237, 28, 36, 0.8);">ON</span>{{else}}<span style="color: rgba(76, 175, 80, 0.9);">OFF</span>{{end}}<span style="margin: 0 2px; color: rgba(255, 255, 255, 0.2);">|</span>
                                    <span style="color: rgba(255, 255, 255, 0.4);">TRANS:</span> {{if $entry.AutoTransmission}}<span style="color: rgba(237, 28, 36, 0.8);">AUTO</span>{{else}}<span style="color: rgba(76, 175, 80, 0.9);">MANUAL</span>{{end}}<span style="margin: 0 2px; color: rgba(255, 255, 255, 0.2);">|</span>
                                    <span style="color: rgba(255, 255, 255, 0.4);">TC:</span> {{if $entry.TractionControl}}<span style="color: rgba(237, 28, 36, 0.8);">ON</span>{{else}}<span style="color: rgba(76, 175, 80, 0.9);">OFF</span>{{end}}
                                </span>
                                {{end}}
                            </div>
                        </div>
                        {{end}}
                    {{end}}
                </div>
                {{end}}
            {{else}}
                <div class="empty-state">
                    <h2>No lap times yet</h2>
                    <p>Be the first to set a time!</p>
                </div>
            {{end}}
        </div>
    </div>
    {{if .ShowAdminButton}}
    <a href="/admin" class="admin-link">⚙️ Admin</a>
    {{end}}
    <script>
        // Auto-refresh for track rotation (always enabled)
        (function() {
            const rotationInterval = {{.RotationInterval}};
            if (rotationInterval && rotationInterval >= 10) {
                // Calculate time until next rotation
                const now = Math.floor(Date.now() / 1000);
                const currentInterval = Math.floor(now / rotationInterval);
                const nextIntervalStart = (currentInterval + 1) * rotationInterval;
                const secondsUntilRefresh = nextIntervalStart - now;
                
                // Refresh the page at the next interval boundary
                setTimeout(function() {
                    window.location.reload();
                }, secondsUntilRefresh * 1000);
            }
        })();
        
        // Auto-scroll for entries after top 3
        (function() {
            const scrollable = document.getElementById('scrollable-leaderboard');
            if (!scrollable) return;
            
            const entries = scrollable.querySelectorAll('.entry');
            if (entries.length === 0) return;
            
            let currentIndex = 0;
            const scrollSpeed = 50; // pixels per second
            const pauseTime = 2000; // pause at top/bottom (ms)
            let isPaused = false;
            let scrollDirection = 1;
            
            function autoScroll() {
                if (isPaused) return;
                
                const maxScroll = scrollable.scrollHeight - scrollable.clientHeight;
                const currentScroll = scrollable.scrollTop;
                
                // Check if we've reached the bottom
                if (currentScroll >= maxScroll - 1) {
                    isPaused = true;
                    setTimeout(function() {
                        scrollable.scrollTo({
                            top: 0,
                            behavior: 'smooth'
                        });
                        setTimeout(function() {
                            isPaused = false;
                        }, pauseTime);
                    }, pauseTime);
                    return;
                }
                
                // Check if we're at the top
                if (currentScroll <= 1 && scrollDirection === -1) {
                    isPaused = true;
                    setTimeout(function() {
                        scrollDirection = 1;
                        isPaused = false;
                    }, pauseTime);
                    return;
                }
                
                // Scroll
                scrollable.scrollTop += (scrollSpeed / 60) * scrollDirection;
            }
            
            // Start scrolling after a brief pause
            setTimeout(function() {
                setInterval(autoScroll, 1000 / 60); // 60fps
            }, 1000);
        })();
    </script>
</body>
</html>`

const adminHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Admin Console - Lap Time Tracker</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: 'Helvetica Neue', Arial, sans-serif;
            background: #0A0E27;
            background: linear-gradient(180deg, #0A0E27 0%, #1A1F3A 100%);
            min-height: 100vh;
            padding: 24px;
            color: #FFFFFF;
        }
        
        .container {
            max-width: 1400px;
            margin: 0 auto;
            background: #1A1F3A;
            border-radius: 0;
            border: 2px solid #ED1C24;
            box-shadow: 0 0 20px rgba(237, 28, 36, 0.2);
            overflow: hidden;
        }
        
        .header {
            background: linear-gradient(135deg, #0A0E27 0%, #1A1F3A 50%, #0A0E27 100%);
            border-bottom: 3px solid #ED1C24;
            color: #FFFFFF;
            padding: 32px 40px;
            display: flex;
            justify-content: space-between;
            align-items: center;
            position: relative;
        }
        
        .header::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background: linear-gradient(90deg, transparent 0%, rgba(237, 28, 36, 0.05) 50%, transparent 100%);
            pointer-events: none;
        }
        
        .header h1 {
            font-size: 2.5em;
            font-weight: 900;
            color: #FFFFFF;
            position: relative;
            z-index: 1;
            letter-spacing: 3px;
            text-transform: uppercase;
            text-shadow: 0 2px 10px rgba(0, 0, 0, 0.5);
        }
        
        .back-link {
            color: #ED1C24;
            text-decoration: none;
            padding: 12px 24px;
            background: transparent;
            border: 2px solid #ED1C24;
            border-radius: 0;
            transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
            position: relative;
            z-index: 1;
            font-weight: 900;
            font-size: 0.875em;
            text-transform: uppercase;
            letter-spacing: 2px;
        }
        
        .back-link:hover {
            background: #ED1C24;
            color: #FFFFFF;
            border-color: #ED1C24;
            box-shadow: 0 0 20px rgba(237, 28, 36, 0.4);
            transform: translateY(-1px);
        }
        
        .content {
            padding: 32px;
            display: grid;
            grid-template-columns: 400px 1fr;
            gap: 24px;
            background: #1A1F3A;
        }
        
        .collapsible-header {
            display: flex;
            align-items: center;
            justify-content: space-between;
            cursor: pointer;
            user-select: none;
            padding: 12px 0;
            border-bottom: 2px solid rgba(237, 28, 36, 0.3);
            margin-bottom: 16px;
        }
        
        .collapsible-header:hover {
            border-bottom-color: rgba(237, 28, 36, 0.5);
        }
        
        .collapsible-header h2 {
            margin: 0;
            font-size: 1.5em;
        }
        
        .collapsible-toggle {
            color: #ED1C24;
            font-size: 1.2em;
            transition: transform 0.3s ease;
        }
        
        .collapsible-content {
            overflow: hidden;
            transition: max-height 0.3s ease;
        }
        
        .collapsible-content.collapsed {
            max-height: 0;
        }
        
        .collapsible-content.expanded {
            max-height: 5000px;
        }
        
        .section {
            background: #0A0E27;
            padding: 24px;
            border-radius: 0;
            border: 1px solid rgba(237, 28, 36, 0.3);
            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.4);
            transition: box-shadow 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        }
        
        .section:hover {
            box-shadow: 0 4px 16px rgba(237, 28, 36, 0.2);
            border-color: rgba(237, 28, 36, 0.5);
        }
        
        .section h2 {
            color: #ED1C24;
            margin-bottom: 24px;
            font-size: 1.75em;
            font-weight: 900;
            letter-spacing: 2px;
            text-transform: uppercase;
        }
        
        .form-group {
            margin-bottom: 20px;
        }
        
        .form-group label {
            display: block;
            margin-bottom: 8px;
            position: relative;
        }
        
        .tooltip-icon {
            display: inline-block;
            margin-left: 8px;
            color: rgba(237, 28, 36, 0.6);
            cursor: help;
            font-size: 0.9em;
            position: relative;
        }
        
        .tooltip-icon:hover {
            color: #ED1C24;
        }
        
        .tooltip-text {
            visibility: hidden;
            opacity: 0;
            background-color: #0A0E27;
            color: #FFFFFF;
            text-align: left;
            padding: 12px 16px;
            border: 2px solid #ED1C24;
            border-radius: 0;
            position: absolute;
            z-index: 10000;
            bottom: 125%;
            right: 0;
            white-space: normal;
            width: 250px;
            font-size: 0.875em;
            font-weight: 500;
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.6);
            transition: opacity 0.3s ease, visibility 0.3s ease;
            pointer-events: none;
        }
        
        .tooltip-icon:hover .tooltip-text {
            visibility: visible;
            opacity: 1;
        }
        
        .tooltip-text::after {
            content: "";
            position: absolute;
            top: 100%;
            right: 20px;
            border-width: 8px;
            border-style: solid;
            border-color: #ED1C24 transparent transparent transparent;
        }
            font-weight: 700;
            color: #ED1C24;
            text-transform: uppercase;
            letter-spacing: 1px;
            font-size: 0.875em;
        }
        
        * {
            -webkit-tap-highlight-color: transparent;
        }
        
        input,
        select,
        input[type="text"],
        input[type="text"]:focus,
        input[type="text"]:hover,
        select:focus,
        select:hover {
            -webkit-appearance: none !important;
            -moz-appearance: none !important;
            appearance: none !important;
            background-image: none !important;
        }
        
        .form-group input,
        .form-group select,
        #lap-time,
        #lap-track,
        #lap-user {
            width: 100%;
            padding: 14px 16px;
            border: 2px solid rgba(237, 28, 36, 0.3);
            background: #1A1F3A !important;
            color: #FFFFFF !important;
            border-radius: 0;
            font-size: 0.9375em;
            font-weight: 500;
            transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
            font-family: inherit;
            box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.3);
            -webkit-appearance: none !important;
            -moz-appearance: none !important;
            appearance: none !important;
            box-sizing: border-box;
            margin: 0;
        }
        
        .form-group input[type="text"],
        .form-group input[type="text"]:focus,
        .form-group input[type="text"]:hover,
        #lap-time,
        #lap-time:focus,
        #lap-time:hover {
            -webkit-appearance: none !important;
            -moz-appearance: none !important;
            appearance: none !important;
            background: #1A1F3A !important;
            background-color: #1A1F3A !important;
        }
        
        .form-group input::placeholder {
            color: rgba(255, 255, 255, 0.3);
            font-weight: 400;
            opacity: 1;
        }
        
        .form-group input::-webkit-input-placeholder {
            color: rgba(255, 255, 255, 0.3);
            font-weight: 400;
        }
        
        .form-group input::-moz-placeholder {
            color: rgba(255, 255, 255, 0.3);
            font-weight: 400;
            opacity: 1;
        }
        
        .form-group input:-ms-input-placeholder {
            color: rgba(255, 255, 255, 0.3);
            font-weight: 400;
        }
        
        .form-group input:focus,
        .form-group select:focus {
            outline: none;
            border-color: #ED1C24;
            background: #1A1F3A;
            box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.3), 0 0 0 3px rgba(237, 28, 36, 0.2), 0 0 20px rgba(237, 28, 36, 0.15);
        }
        
        .form-group select,
        #lap-track,
        #lap-user {
            background: #1A1F3A !important;
            background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 12 12'%3E%3Cpath fill='%23ED1C24' d='M6 9L1 4h10z'/%3E%3C/svg%3E") !important;
            background-repeat: no-repeat !important;
            background-position: right 16px center !important;
            background-size: 12px !important;
            padding-right: 40px;
            cursor: pointer;
            -webkit-appearance: none !important;
            -moz-appearance: none !important;
            appearance: none !important;
        }
        
        .form-group select::-ms-expand,
        #lap-track::-ms-expand,
        #lap-user::-ms-expand {
            display: none !important;
        }
        
        select::-ms-expand {
            display: none !important;
        }
        
        .form-group select:hover {
            border-color: rgba(237, 28, 36, 0.6);
            background: #1A1F3A;
        }
        
        .form-group select:focus {
            background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 12 12'%3E%3Cpath fill='%23FFFFFF' d='M6 9L1 4h10z'/%3E%3C/svg%3E");
        }
        
        .form-group select option {
            background: #1A1F3A;
            color: #FFFFFF;
            padding: 12px;
            border: none;
            -webkit-appearance: none;
            -moz-appearance: none;
            appearance: none;
        }
        
        .form-group input:hover {
            border-color: rgba(237, 28, 36, 0.6);
            background: #1A1F3A;
        }
        
        .form-group input[type="text"]:focus {
            background: #1A1F3A;
        }
        
        input[type="checkbox"] {
            width: 20px;
            height: 20px;
            appearance: none;
            -webkit-appearance: none;
            -moz-appearance: none;
            background: #0A0E27;
            border: 2px solid #ED1C24;
            border-radius: 0;
            cursor: pointer;
            position: relative;
            transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        }
        
        input[type="checkbox"]:hover {
            border-color: #ED1C24;
            box-shadow: 0 0 10px rgba(237, 28, 36, 0.3);
        }
        
        input[type="checkbox"]:checked {
            background: #ED1C24;
            border-color: #ED1C24;
        }
        
        input[type="checkbox"]:checked::after {
            content: '✓';
            position: absolute;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
            color: #FFFFFF;
            font-size: 14px;
            font-weight: 900;
        }
        
        .btn {
            padding: 12px 24px;
            background: #0A0E27;
            color: #ED1C24;
            border: 2px solid #ED1C24;
            border-radius: 0;
            font-size: 0.875em;
            font-weight: 900;
            cursor: pointer;
            transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
            text-transform: uppercase;
            letter-spacing: 2px;
        }
        
        .btn:hover {
            background: #ED1C24;
            color: #FFFFFF;
            border-color: #ED1C24;
            box-shadow: 0 0 20px rgba(237, 28, 36, 0.4);
            transform: translateY(-2px);
        }
        
        .btn:active {
            transform: translateY(0);
            box-shadow: 0 0 10px rgba(237, 28, 36, 0.3);
        }
        
        .btn-secondary {
            background: #0A0E27;
            border-color: rgba(255, 255, 255, 0.3);
            color: rgba(255, 255, 255, 0.8);
        }
        
        .btn-secondary:hover {
            background: rgba(255, 255, 255, 0.1);
            color: #FFFFFF;
            border-color: rgba(255, 255, 255, 0.5);
        }
        
        .btn-danger {
            background: transparent;
            border: none;
            color: #dc3545;
            font-size: 1em;
            padding: 4px 6px;
            cursor: pointer;
            opacity: 0.5;
            transition: opacity 0.2s, transform 0.2s;
            border-radius: 4px;
            line-height: 1;
            min-width: auto;
            width: auto;
        }
        
        .btn-danger:hover {
            opacity: 1;
            background: rgba(220, 53, 69, 0.1);
            transform: scale(1.1);
        }
        
        .btn-danger:active {
            transform: scale(0.95);
        }
        
        .lap-time-item .btn-danger {
            font-size: 0.85em;
            padding: 2px 2px;
            flex-shrink: 0;
        }
        
        .track-item .actions {
            display: flex;
            gap: 10px;
            align-items: center;
        }
        
        .lap-time-item {
            display: grid;
            grid-template-columns: 1fr auto;
            gap: 15px;
            align-items: center;
        }
        
        .lap-time-item .time-container {
            display: flex;
            align-items: center;
            gap: 10px;
            justify-content: flex-end;
        }
        
        .tracks-list {
            margin-top: 20px;
        }
        
        .track-item {
            background: #0A0E27;
            padding: 16px 20px;
            margin-bottom: 8px;
            border-radius: 0;
            border: 1px solid rgba(237, 28, 36, 0.2);
            display: flex;
            justify-content: space-between;
            align-items: center;
            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.4);
            transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
            border-left: 5px solid transparent;
        }
        
        .track-item:hover {
            background: #1A1F3A;
            box-shadow: 0 4px 12px rgba(237, 28, 36, 0.2);
            transform: translateX(4px);
            border-color: rgba(237, 28, 36, 0.4);
        }
        
        .track-item.active {
            border-left-color: #ED1C24;
            background: rgba(237, 28, 36, 0.1);
            box-shadow: 0 0 15px rgba(237, 28, 36, 0.2);
            border: 1px solid #ED1C24;
        }
        
        .track-item .track-name {
            font-weight: 700;
            color: #FFFFFF;
            font-size: 0.9375em;
        }
        
        .track-item .track-status {
            color: #ED1C24;
            font-weight: 700;
            font-size: 0.8125em;
            text-transform: uppercase;
            letter-spacing: 1px;
        }
        
        .track-edit-form {
            display: flex;
            align-items: center;
            gap: 8px;
            margin-top: 8px;
        }
        
        .track-edit-form input {
            background: #1e1e1e;
            border: 1px solid rgba(255, 255, 255, 0.1);
            color: #e0e0e0;
            border-radius: 4px;
        }
        
        .lap-times-list {
            margin-top: 20px;
            max-height: 400px;
            overflow-y: auto;
        }
        
        .lap-time-item {
            background: #0A0E27;
            padding: 16px 20px;
            margin-bottom: 8px;
            border-radius: 0;
            border: 1px solid rgba(237, 28, 36, 0.2);
            display: grid;
            grid-template-columns: 1fr auto;
            gap: 15px;
            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.4);
            transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        }
        
        .lap-time-item:hover {
            background: #1A1F3A;
            box-shadow: 0 4px 12px rgba(237, 28, 36, 0.2);
            border-color: rgba(237, 28, 36, 0.4);
        }
        
        .lap-time-item .username {
            font-weight: 700;
            color: #FFFFFF;
            font-size: 0.9375em;
            text-transform: uppercase;
            letter-spacing: 1px;
        }
        
        .lap-time-item .time {
            font-family: 'Courier New', monospace;
            font-weight: 900;
            color: #ED1C24;
            text-align: right;
            font-size: 0.9375em;
            letter-spacing: 1px;
        }
        
        .message {
            padding: 15px;
            border-radius: 8px;
            margin-bottom: 20px;
            display: none;
        }
        
        .message.success {
            background: rgba(34, 197, 94, 0.1);
            color: #22c55e;
            border: 1px solid #22c55e;
            display: block;
        }
        
        .message.error {
            background: rgba(237, 28, 36, 0.1);
            color: #ED1C24;
            border: 1px solid #ED1C24;
            display: block;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>⚙️ Admin Console</h1>
            <a href="/" class="back-link">← Back to Leaderboard</a>
        </div>
        <div class="content">
            <div style="display: flex; flex-direction: column; gap: 24px;">
                <div class="section">
                    <div class="collapsible-header" onclick="toggleCollapsible('track-management')">
                        <h2>Track Management</h2>
                        <span class="collapsible-toggle" id="track-management-toggle">▼</span>
                    </div>
                    <div class="collapsible-content collapsed" id="track-management">
                        <div id="message-track" class="message"></div>
                        <div class="form-group">
                            <label for="track-name">Track Name</label>
                            <input type="text" id="track-name" placeholder="Enter track name">
                        </div>
                        <button class="btn" onclick="addTrack()">Add Track</button>
                        
                        <div class="tracks-list">
                            <h3 style="margin-top: 30px; margin-bottom: 15px; color: #ED1C24; font-weight: 900; text-transform: uppercase; letter-spacing: 2px; font-size: 1em;">Tracks</h3>
                            <div id="tracks-list"></div>
                        </div>
                        
                    </div>
                </div>
                
                <div class="section">
                    <div class="collapsible-header" onclick="toggleCollapsible('user-management')">
                        <h2>User Management</h2>
                        <span class="collapsible-toggle" id="user-management-toggle">▼</span>
                    </div>
                    <div class="collapsible-content collapsed" id="user-management">
                        <div id="message-user" class="message"></div>
                        <div class="form-group">
                            <label for="user-name">Username</label>
                            <input type="text" id="user-name" placeholder="Enter username">
                        </div>
                        <button class="btn" onclick="addUser()">Add User</button>
                        
                        <div class="users-list">
                            <h3 style="margin-top: 30px; margin-bottom: 15px; color: #ED1C24; font-weight: 900; text-transform: uppercase; letter-spacing: 2px; font-size: 1em;">Users</h3>
                            <div id="users-list"></div>
                        </div>
                    </div>
                </div>
                
                <div class="section">
                    <div class="collapsible-header" onclick="toggleCollapsible('leaderboard-management')">
                        <h2>Leaderboard Management</h2>
                        <span class="collapsible-toggle" id="leaderboard-management-toggle">▼</span>
                    </div>
                    <div class="collapsible-content collapsed" id="leaderboard-management">
                        <div id="message-settings" class="message"></div>
                        <div style="display: flex; flex-direction: column; gap: 16px;">
                            <label style="display: flex; align-items: center; cursor: pointer; margin: 0;">
                                <input type="checkbox" id="show-admin-button" style="width: 20px; height: 20px; margin-right: 8px; cursor: pointer; flex-shrink: 0; box-sizing: border-box;">
                                <span style="font-weight: 500; color: #ED1C24; text-transform: uppercase; letter-spacing: 1px; font-size: 0.875em;">Show Admin Button on Leaderboard</span>
                            </label>
                            <label style="display: flex; align-items: center; cursor: pointer; margin: 0;">
                                <input type="checkbox" id="show-assists-leaderboard" style="width: 20px; height: 20px; margin-right: 8px; cursor: pointer; flex-shrink: 0; box-sizing: border-box;">
                                <span style="font-weight: 500; color: #ED1C24; text-transform: uppercase; letter-spacing: 1px; font-size: 0.875em;">Show Assists on Leaderboard</span>
                            </label>
                            <div class="form-group" style="margin-top: 8px;">
                                <label for="leaderboard-title" style="font-weight: 500; color: #ED1C24; text-transform: uppercase; letter-spacing: 1px; font-size: 0.875em; margin-bottom: 8px; display: block;">Leaderboard Title</label>
                                <input type="text" id="leaderboard-title" placeholder="Sim Racing Leaderboard" style="width: 100%; padding: 10px 14px; border: 2px solid rgba(237, 28, 36, 0.3); background: #1A1F3A; color: #FFFFFF; border-radius: 0; font-size: 0.9375em; font-weight: 500; font-family: inherit; box-sizing: border-box;">
                            </div>
                            
                            <div class="form-group" style="margin-top: 16px;">
                                <div class="collapsible-header" onclick="toggleCollapsible('track-rotation')" style="cursor: pointer; padding: 8px 0; margin-bottom: 12px;">
                                    <span style="font-weight: 500; color: #ED1C24; text-transform: uppercase; letter-spacing: 1px; font-size: 0.875em;">Track Rotation</span>
                                    <span class="collapsible-toggle" id="track-rotation-toggle" style="float: right; color: #ED1C24;">▼</span>
                                </div>
                                <div class="collapsible-content collapsed" id="track-rotation" style="display: flex; flex-direction: column; gap: 12px;">
                                    <div class="form-group">
                                        <label for="rotation-interval" style="font-weight: 500; color: #ED1C24; text-transform: uppercase; letter-spacing: 1px; font-size: 0.875em; margin-bottom: 8px; display: block;">Rotation Interval (seconds)</label>
                                        <input type="number" id="rotation-interval" min="10" value="60" placeholder="60" style="width: 100%; padding: 10px 14px; border: 2px solid rgba(237, 28, 36, 0.3); background: #1A1F3A; color: #FFFFFF; border-radius: 0; font-size: 0.9375em; font-weight: 500; font-family: inherit; box-sizing: border-box;">
                                    </div>
                                    <div class="form-group">
                                        <div style="margin-bottom: 8px;">
                                            <label style="font-weight: 500; color: #ED1C24; text-transform: uppercase; letter-spacing: 1px; font-size: 0.875em; margin-bottom: 8px; display: block;">
                                                Tracks on Leaderboard
                                                <span class="tooltip-icon">
                                                    ℹ️
                                                    <span class="tooltip-text">All selected tracks will rotate. If only one track is selected, it will remain static.</span>
                                                </span>
                                            </label>
                                            <div style="margin-bottom: 8px;">
                                                <button type="button" class="btn" onclick="selectAllRotationTracks()" style="padding: 4px 12px; font-size: 0.75em; margin-right: 8px;">Select All</button>
                                                <button type="button" class="btn" onclick="deselectAllRotationTracks()" style="padding: 4px 12px; font-size: 0.75em;">Deselect All</button>
                                            </div>
                                        </div>
                                        <div id="rotation-tracks-list" style="max-height: 200px; overflow-y: auto; border: 1px solid rgba(237, 28, 36, 0.3); padding: 12px; background: #0A0E27;">
                                            <!-- Track checkboxes will be populated here -->
                                        </div>
                                    </div>
                                </div>
                            </div>
                            
                            <button class="btn" onclick="saveSettings()" style="padding: 8px 16px; font-size: 0.75em;">Save Settings</button>
                        </div>
                    </div>
                </div>
            </div>
            
            <div class="section" style="grid-column: 2;">
                <h2>Lap Time Management</h2>
                <div id="message-lap" class="message"></div>
                <div class="form-group">
                    <label for="lap-track">Select Track</label>
                    <select id="lap-track">
                        <option value="">Loading...</option>
                    </select>
                </div>
                <div class="form-group">
                    <label for="lap-user">Select User</label>
                    <select id="lap-user">
                        <option value="">Loading...</option>
                    </select>
                </div>
                <div class="form-group">
                    <label for="lap-time">Lap Time (e.g., 1:23.456 or 83.4)</label>
                    <input type="text" id="lap-time" placeholder="1:23.456 or 83.4">
                </div>
                <div class="form-group">
                    <label for="lap-car">Car (optional)</label>
                    <input type="text" id="lap-car" placeholder="e.g., Ferrari, Porsche">
                </div>
                <div class="form-group" style="margin-top: 16px;">
                    <div class="collapsible-header" onclick="toggleCollapsible('sim-assists')" style="cursor: pointer; padding: 8px 0;">
                        <span style="font-weight: 500; color: #ED1C24; text-transform: uppercase; letter-spacing: 1px; font-size: 0.875em;">Sim Assists</span>
                        <span class="collapsible-toggle" id="sim-assists-toggle" style="float: right; color: #ED1C24;">▼</span>
                    </div>
                    <div class="collapsible-content collapsed" id="sim-assists" style="display: flex; flex-direction: column; gap: 12px; margin-top: 12px;">
                        <label style="display: flex; align-items: center; cursor: pointer; margin: 0;">
                            <input type="checkbox" id="lap-abs" style="width: 20px; height: 20px; margin-right: 8px; cursor: pointer; flex-shrink: 0; box-sizing: border-box;">
                            <span style="font-weight: 500; color: #FFFFFF; font-size: 0.875em;">ABS</span>
                        </label>
                        <label style="display: flex; align-items: center; cursor: pointer; margin: 0;">
                            <input type="checkbox" id="lap-mt" style="width: 20px; height: 20px; margin-right: 8px; cursor: pointer; flex-shrink: 0; box-sizing: border-box;">
                            <span style="font-weight: 500; color: #FFFFFF; font-size: 0.875em;">Auto Transmission</span>
                        </label>
                        <label style="display: flex; align-items: center; cursor: pointer; margin: 0;">
                            <input type="checkbox" id="lap-tc" style="width: 20px; height: 20px; margin-right: 8px; cursor: pointer; flex-shrink: 0; box-sizing: border-box;">
                            <span style="font-weight: 500; color: #FFFFFF; font-size: 0.875em;">Traction Control</span>
                        </label>
                    </div>
                </div>
                <button class="btn" onclick="addOrUpdateLapTime()" id="lap-time-button">Add/Update Lap Time</button>
                
                <div class="lap-times-list">
                    <h3 style="margin-top: 30px; margin-bottom: 15px; color: #ED1C24; font-weight: 900; text-transform: uppercase; letter-spacing: 2px; font-size: 1em;">Lap Times</h3>
                    <div id="lap-times-list"></div>
                </div>
            </div>
        </div>
    </div>
    
    <script>
        let tracks = [];
        let users = [];
        let activeTrackId = null;
        let existingLapTimes = {};
        
        async function loadUsers() {
            try {
                const response = await fetch('/api/users');
                users = await response.json();
                
                const select = document.getElementById('lap-user');
                select.innerHTML = '<option value="">Select a user...</option>' + 
                    users.map(u => '<option value="' + u.id + '">' + u.username + '</option>').join('');
                
                const list = document.getElementById('users-list');
                if (users.length === 0) {
                    list.innerHTML = '<p style="color: #999; text-align: center; padding: 20px;">No users yet</p>';
                } else {
                    list.innerHTML = users.map(user => 
                        '<div class="lap-time-item" style="margin-bottom: 8px;">' +
                            '<div class="username" style="flex: 1;">' + user.username + '</div>' +
                            '<div class="time-container">' +
                                '<button class="btn btn-secondary" onclick="editUser(' + user.id + ', \'' + user.username.replace(/'/g, "\\'") + '\')" title="Edit user" style="padding: 4px 8px; font-size: 0.75em; margin-right: 4px;">✏️</button>' +
                                '<button class="btn btn-danger" onclick="deleteUser(' + user.id + ', \'' + user.username.replace(/'/g, "\\'") + '\')" title="Delete user">🗑️</button>' +
                            '</div>' +
                        '</div>'
                    ).join('');
                }
            } catch (error) {
                console.error('Error loading users:', error);
            }
        }
        
        async function addUser() {
            const username = document.getElementById('user-name').value.trim();
            if (!username) {
                showMessage('message-user', 'Please enter a username', 'error');
                return;
            }
            
            try {
                const response = await fetch('/api/users', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ username })
                });
                
                if (response.ok) {
                    document.getElementById('user-name').value = '';
                    showMessage('message-user', 'User added successfully!', 'success');
                    loadUsers();
                } else {
                    const error = await response.text();
                    showMessage('message-user', 'Error adding user: ' + error, 'error');
                }
            } catch (error) {
                showMessage('message-user', 'Error adding user', 'error');
            }
        }
        
        async function editUser(userId, currentUsername) {
            const newUsername = prompt('Enter new username:', currentUsername);
            if (!newUsername || newUsername.trim() === '') return;
            
            try {
                const response = await fetch('/api/users/' + userId, {
                    method: 'PUT',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ username: newUsername.trim() })
                });
                
                if (response.ok) {
                    showMessage('message-user', 'User updated successfully!', 'success');
                    loadUsers();
                } else {
                    showMessage('message-user', 'Error updating user', 'error');
                }
            } catch (error) {
                showMessage('message-user', 'Error updating user', 'error');
            }
        }
        
        async function deleteUser(userId, username) {
            if (!confirm('Delete user "' + username + '"? This will also delete all their lap times.')) {
                return;
            }
            
            try {
                const response = await fetch('/api/users/' + userId, {
                    method: 'DELETE'
                });
                
                if (response.ok) {
                    showMessage('message-user', 'User deleted successfully!', 'success');
                    loadUsers();
                    const trackId = document.getElementById('lap-track').value;
                    if (trackId) {
                        loadLapTimes(parseInt(trackId));
                    }
                } else {
                    showMessage('message-user', 'Error deleting user', 'error');
                }
            } catch (error) {
                showMessage('message-user', 'Error deleting user', 'error');
            }
        }
        
        async function loadTracks() {
            try {
                const response = await fetch('/api/tracks');
                tracks = await response.json();
                activeTrackId = tracks.find(t => t.is_active)?.id || null;
                
                // Update track select
                const select = document.getElementById('lap-track');
                select.innerHTML = tracks.map(t => {
                    return '<option value="' + t.id + '">' + t.name + '</option>';
                }).join('');
                
                // Update tracks list
                const list = document.getElementById('tracks-list');
                list.innerHTML = tracks.map(track => {
                    const isEditing = false;
                    return '<div class="track-item ' + (track.is_active ? 'active' : '') + '" id="track-item-' + track.id + '">' +
                        '<div>' +
                            '<div class="track-name" id="track-name-' + track.id + '">' + track.name + '</div>' +
                            '<div class="track-edit-form" id="track-edit-' + track.id + '" style="display: none;">' +
                                '<input type="text" id="edit-name-' + track.id + '" placeholder="Track Name" value="' + track.name.replace(/'/g, "&#39;") + '" style="width: 200px; margin-right: 8px; padding: 4px; font-size: 0.85em;">' +
                                '<button class="btn btn-secondary" onclick="saveTrack(' + track.id + ')" style="padding: 4px 12px; font-size: 0.75em; margin-right: 4px;">Save</button>' +
                                '<button class="btn btn-secondary" onclick="cancelEdit(' + track.id + ')" style="padding: 4px 12px; font-size: 0.75em;">Cancel</button>' +
                            '</div>' +
                        '</div>' +
                        '<div class="actions">' +
                            '<button class="btn btn-secondary" onclick="editTrack(' + track.id + ')" title="Edit track" style="padding: 4px 8px; font-size: 0.75em; margin-right: 4px;">✏️</button>' +
                            '<button class="btn btn-danger" onclick="deleteTrack(' + track.id + ', \'' + track.name.replace(/'/g, "\\'") + '\')" title="Delete track">🗑️</button>' +
                        '</div>' +
                    '</div>';
                }).join('');
                
                if (activeTrackId) {
                    loadLapTimes(activeTrackId);
                }
            } catch (error) {
                console.error('Error loading tracks:', error);
            }
        }
        
        async function checkExistingLapTime(trackId, userId) {
            if (!trackId || !userId) return null;
            try {
                const response = await fetch('/api/laptimes?track_id=' + trackId);
                const lapTimes = await response.json();
                const existing = lapTimes.find(lt => lt.user_id == userId);
                return existing || null;
            } catch (error) {
                return null;
            }
        }
        
        async function addTrack() {
            const name = document.getElementById('track-name').value.trim();
            if (!name) {
                showMessage('message-track', 'Please enter a track name', 'error');
                return;
            }
            
            try {
                const response = await fetch('/api/tracks', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ name })
                });
                
                if (response.ok) {
                    document.getElementById('track-name').value = '';
                    showMessage('message-track', 'Track added successfully!', 'success');
                    loadTracks();
                } else {
                    showMessage('message-track', 'Error adding track', 'error');
                }
            } catch (error) {
                showMessage('message-track', 'Error adding track', 'error');
            }
        }
        
        async function setActiveTrack(trackId) {
            try {
                const response = await fetch('/api/tracks/active', {
                    method: 'PUT',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ track_id: trackId })
                });
                
                if (response.ok) {
                    showMessage('message-track', 'Active track updated!', 'success');
                    loadTracks();
                } else {
                    showMessage('message-track', 'Error updating active track', 'error');
                }
            } catch (error) {
                showMessage('message-track', 'Error updating active track', 'error');
            }
        }
        
        async function loadLapTimes(trackId) {
            if (!trackId) return;
            
            try {
                const response = await fetch('/api/laptimes?track_id=' + trackId);
                const lapTimes = await response.json();
                
                existingLapTimes = {};
                lapTimes.forEach(lt => {
                    existingLapTimes[lt.user_id] = lt;
                });
                
                updateLapTimeButton();
                
                const list = document.getElementById('lap-times-list');
                if (lapTimes.length === 0) {
                    list.innerHTML = '<p style="color: #999; text-align: center; padding: 20px;">No lap times for this track</p>';
                } else {
                    // Always show assists in admin console
                    list.innerHTML = lapTimes.map(lt => {
                        const absOn = lt.abs;
                        const transAuto = lt.auto_transmission;
                        const tcOn = lt.traction_control;
                        const carDisplay = lt.car ? '<span style="font-size: 0.75em; color: rgba(255, 255, 255, 0.5); margin-left: 8px;">(' + lt.car + ')</span>' : '';
                        return '<div class="lap-time-item">' +
                            '<div class="username">' + lt.username + carDisplay + '</div>' +
                            '<div class="time-container">' +
                                '<div class="time" style="display: flex; flex-direction: column; align-items: flex-end; gap: 6px;">' +
                                    '<span>' + lt.time + '</span>' +
                                    '<span style="font-size: 0.7em; color: rgba(255, 255, 255, 0.5); font-weight: 500; text-transform: uppercase; letter-spacing: 0.8px; line-height: 1.2; white-space: nowrap;">' +
                                        '<span style="color: rgba(255, 255, 255, 0.4);">ABS:</span> ' + 
                                        (absOn ? '<span style="color: rgba(237, 28, 36, 0.8);">ON</span>' : '<span style="color: rgba(76, 175, 80, 0.9);">OFF</span>') + 
                                        ' <span style="margin: 0 6px; color: rgba(255, 255, 255, 0.2);">|</span> ' +
                                        '<span style="color: rgba(255, 255, 255, 0.4);">TRANS:</span> ' + 
                                        (transAuto ? '<span style="color: rgba(237, 28, 36, 0.8);">AUTO</span>' : '<span style="color: rgba(76, 175, 80, 0.9);">MANUAL</span>') + 
                                        ' <span style="margin: 0 6px; color: rgba(255, 255, 255, 0.2);">|</span> ' +
                                        '<span style="color: rgba(255, 255, 255, 0.4);">TC:</span> ' + 
                                        (tcOn ? '<span style="color: rgba(237, 28, 36, 0.8);">ON</span>' : '<span style="color: rgba(76, 175, 80, 0.9);">OFF</span>') +
                                    '</span>' +
                                '</div>' +
                                '<button class="btn btn-danger" onclick="deleteLapTime(' + lt.id + ', ' + trackId + ')" title="Delete lap time">🗑️</button>' +
                            '</div>' +
                        '</div>';
                    }).join('');
                }
            } catch (error) {
                console.error('Error loading lap times:', error);
            }
        }
        
        function updateLapTimeButton() {
            const trackId = document.getElementById('lap-track').value;
            const userId = document.getElementById('lap-user').value;
            
            // Load existing lap time assists and car if editing
            if (trackId && userId && existingLapTimes[userId]) {
                const existing = existingLapTimes[userId];
                document.getElementById('lap-abs').checked = existing.abs || false;
                document.getElementById('lap-mt').checked = existing.auto_transmission || false;
                document.getElementById('lap-tc').checked = existing.traction_control || false;
                document.getElementById('lap-car').value = existing.car || '';
                document.getElementById('lap-time').value = existing.time || '';
            } else {
                document.getElementById('lap-abs').checked = false;
                document.getElementById('lap-mt').checked = false;
                document.getElementById('lap-tc').checked = false;
                document.getElementById('lap-car').value = '';
                document.getElementById('lap-time').value = '';
            }
            const button = document.getElementById('lap-time-button');
            
            if (!trackId || !userId) {
                button.textContent = 'Add/Update Lap Time';
                return;
            }
            
            const existing = existingLapTimes[userId];
            if (existing) {
                button.textContent = 'Update Lap Time (' + existing.time + ')';
            } else {
                button.textContent = 'Add Lap Time';
            }
        }
        
        function normalizeTime(timeStr) {
            timeStr = timeStr.trim().replace(/\s+/g, '');
            
            if (!timeStr) return null;
            
            let minutes = 0;
            let seconds = 0;
            let milliseconds = 0;
            
            if (timeStr.includes(':')) {
                const parts = timeStr.split(':');
                minutes = parseInt(parts[0]) || 0;
                const secParts = parts[1].split('.');
                seconds = parseInt(secParts[0]) || 0;
                if (secParts[1]) {
                    const msStr = secParts[1].padEnd(3, '0').substring(0, 3);
                    milliseconds = parseInt(msStr) || 0;
                }
            } else if (timeStr.includes('.')) {
                const parts = timeStr.split('.');
                seconds = parseInt(parts[0]) || 0;
                if (parts[1]) {
                    const msStr = parts[1].padEnd(3, '0').substring(0, 3);
                    milliseconds = parseInt(msStr) || 0;
                }
            } else {
                seconds = parseInt(timeStr) || 0;
            }
            
            if (isNaN(minutes) || isNaN(seconds) || isNaN(milliseconds)) {
                return null;
            }
            
            while (seconds >= 60) {
                minutes++;
                seconds -= 60;
            }
            
            const mm = String(minutes).padStart(2, '0');
            const ss = String(seconds).padStart(2, '0');
            const mmm = String(milliseconds).padStart(3, '0');
            
            return mm + ':' + ss + '.' + mmm;
        }
        
        async function addOrUpdateLapTime() {
            const trackId = document.getElementById('lap-track').value;
            const userId = document.getElementById('lap-user').value;
            const timeInput = document.getElementById('lap-time').value.trim();
            
            if (!trackId || !userId || !timeInput) {
                showMessage('message-lap', 'Please fill in all fields', 'error');
                return;
            }
            
            const time = normalizeTime(timeInput);
            if (!time) {
                showMessage('message-lap', 'Invalid time format. Examples: 1:23.456, 83.4, 1:30', 'error');
                return;
            }
            
            try {
                const abs = document.getElementById('lap-abs').checked;
                const mt = document.getElementById('lap-mt').checked;
                const tc = document.getElementById('lap-tc').checked;
                const car = document.getElementById('lap-car').value.trim();
                const response = await fetch('/api/laptimes/add', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ 
                        track_id: parseInt(trackId), 
                        user_id: parseInt(userId), 
                        time,
                        car: car,
                        abs: abs,
                        auto_transmission: mt,
                        traction_control: tc
                    })
                });
                
                if (response.ok) {
                    document.getElementById('lap-time').value = '';
                    document.getElementById('lap-car').value = '';
                    document.getElementById('lap-abs').checked = false;
                    document.getElementById('lap-mt').checked = false;
                    document.getElementById('lap-tc').checked = false;
                    const existing = existingLapTimes[userId];
                    if (existing) {
                        showMessage('message-lap', 'Lap time updated successfully!', 'success');
                    } else {
                        showMessage('message-lap', 'Lap time added successfully!', 'success');
                    }
                    loadLapTimes(parseInt(trackId));
                    updateLapTimeButton();
                } else {
                    const error = await response.text();
                    showMessage('message-lap', 'Error: ' + error, 'error');
                }
            } catch (error) {
                showMessage('message-lap', 'Error adding/updating lap time', 'error');
            }
        }
        
        function showMessage(id, message, type) {
            const msgEl = document.getElementById(id);
            msgEl.textContent = message;
            msgEl.className = 'message ' + type;
            setTimeout(() => {
                msgEl.className = 'message';
            }, 3000);
        }
        
        async function loadSettings() {
            try {
                const response = await fetch('/api/settings/admin-button');
                const data = await response.json();
                document.getElementById('show-admin-button').checked = data.value === 'true';
                
                const assistsResponse = await fetch('/api/settings/show-assists-leaderboard');
                const assistsData = await assistsResponse.json();
                document.getElementById('show-assists-leaderboard').checked = assistsData.value === 'true';
                
                const titleResponse = await fetch('/api/settings/leaderboard-title');
                const titleData = await titleResponse.json();
                document.getElementById('leaderboard-title').value = titleData.value || 'Sim Racing Leaderboard';
                
                // Load rotation settings
                const rotationResponse = await fetch('/api/settings/track-rotation');
                const rotationData = await rotationResponse.json();
                const rotationIntervalEl = document.getElementById('rotation-interval');
                if (rotationIntervalEl) rotationIntervalEl.value = rotationData.interval || '60';
                
                // Load tracks for rotation selection (this will restore selected tracks)
                await loadTracksForRotation();
            } catch (error) {
                console.error('Error loading settings:', error);
            }
        }
        
        async function loadTracksForRotation() {
            try {
                const response = await fetch('/api/tracks');
                const tracks = await response.json();
                const container = document.getElementById('rotation-tracks-list');
                container.innerHTML = '';
                
                // Get selected tracks from rotation settings
                const rotationResponse = await fetch('/api/settings/track-rotation');
                let selectedTracks = [];
                if (rotationResponse.ok) {
                    const rotationData = await rotationResponse.json();
                    if (rotationData.tracks && rotationData.tracks !== '' && rotationData.tracks !== 'all') {
                        selectedTracks = rotationData.tracks.split(',').map(t => t.trim()).filter(t => t !== '');
                    } else if (rotationData.tracks === 'all') {
                        // If "all", select all tracks
                        selectedTracks = tracks.map(t => t.id.toString());
                    }
                }
                
                // If no selected tracks, default to active track or first track
                if (selectedTracks.length === 0) {
                    const activeTrack = tracks.find(t => t.is_active);
                    if (activeTrack) {
                        selectedTracks = [activeTrack.id.toString()];
                    } else if (tracks.length > 0) {
                        selectedTracks = [tracks[0].id.toString()];
                    }
                }
                
                tracks.forEach(track => {
                    const label = document.createElement('label');
                    label.style.cssText = 'display: flex; align-items: center; cursor: pointer; margin: 0 0 8px 0;';
                    const checkbox = document.createElement('input');
                    checkbox.type = 'checkbox';
                    checkbox.value = track.id;
                    checkbox.checked = selectedTracks.includes(track.id.toString());
                    checkbox.style.cssText = 'width: 20px; height: 20px; margin-right: 8px; cursor: pointer; flex-shrink: 0; box-sizing: border-box;';
                    const span = document.createElement('span');
                    span.style.cssText = 'font-weight: 500; color: #FFFFFF; font-size: 0.875em;';
                    span.textContent = track.name;
                    label.appendChild(checkbox);
                    label.appendChild(span);
                    container.appendChild(label);
                });
            } catch (error) {
                console.error('Error loading tracks for rotation:', error);
            }
        }
        
        function selectAllRotationTracks() {
            const checkboxes = document.querySelectorAll('#rotation-tracks-list input[type="checkbox"]');
            checkboxes.forEach(cb => cb.checked = true);
        }
        
        function deselectAllRotationTracks() {
            const checkboxes = document.querySelectorAll('#rotation-tracks-list input[type="checkbox"]');
            checkboxes.forEach(cb => cb.checked = false);
        }
        
        async function saveSettings() {
            const showButton = document.getElementById('show-admin-button').checked;
            const showAssists = document.getElementById('show-assists-leaderboard').checked;
            const title = document.getElementById('leaderboard-title').value.trim() || 'Sim Racing Leaderboard';
            
            // Get rotation settings
            const rotationInterval = document.getElementById('rotation-interval').value || '60';
            const selected = Array.from(document.querySelectorAll('#rotation-tracks-list input[type="checkbox"]:checked'))
                .map(cb => cb.value);
            const rotationTracks = selected.length > 0 ? selected.join(',') : '';
            
            try {
                const response = await fetch('/api/settings/admin-button', {
                    method: 'PUT',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ value: showButton ? 'true' : 'false' })
                });
                if (!response.ok) {
                    const errorText = await response.text();
                    showMessage('message-settings', 'Error saving admin button setting: ' + errorText, 'error');
                    return;
                }
                
                const assistsResponse = await fetch('/api/settings/show-assists-leaderboard', {
                    method: 'PUT',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ value: showAssists ? 'true' : 'false' })
                });
                if (!assistsResponse.ok) {
                    const errorText = await assistsResponse.text();
                    showMessage('message-settings', 'Error saving assists setting: ' + errorText, 'error');
                    return;
                }
                
                const titleResponse = await fetch('/api/settings/leaderboard-title', {
                    method: 'PUT',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ value: title })
                });
                if (!titleResponse.ok) {
                    const errorText = await titleResponse.text();
                    showMessage('message-settings', 'Error saving title setting: ' + errorText, 'error');
                    return;
                }
                
                const rotationResponse = await fetch('/api/settings/track-rotation', {
                    method: 'PUT',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        enabled: 'true',
                        tracks: rotationTracks,
                        interval: rotationInterval
                    })
                });
                if (!rotationResponse.ok) {
                    const errorText = await rotationResponse.text();
                    showMessage('message-settings', 'Error saving rotation settings: ' + errorText, 'error');
                    return;
                }
                
                showMessage('message-settings', 'Settings saved successfully!', 'success');
                // Reload lap times to update display
                const trackId = document.getElementById('lap-track').value;
                if (trackId) {
                    loadLapTimes(parseInt(trackId));
                }
            } catch (error) {
                showMessage('message-settings', 'Error saving settings: ' + error.message, 'error');
            }
        }
        
        function toggleCollapsible(id) {
            const content = document.getElementById(id);
            const toggle = document.getElementById(id + '-toggle');
            if (content.classList.contains('collapsed')) {
                content.classList.remove('collapsed');
                content.classList.add('expanded');
                toggle.textContent = '▲';
            } else {
                content.classList.remove('expanded');
                content.classList.add('collapsed');
                toggle.textContent = '▼';
            }
        }
        
        // Global function to update rotation display
        function updateRotationDisplay() {
            const rotationAll = document.getElementById('rotation-all');
            const rotationCustom = document.getElementById('rotation-custom');
            const rotationTracksList = document.getElementById('rotation-tracks-list');
            
            if (rotationAll && rotationCustom && rotationTracksList) {
                if (rotationAll.checked) {
                    rotationTracksList.style.display = 'none';
                } else if (rotationCustom.checked) {
                    rotationTracksList.style.display = 'block';
                }
            }
        }
        
        // Load tracks, users, and settings on page load
        loadTracks();
        loadUsers();
        loadSettings();
        
        // Update lap times when track selection changes
        document.getElementById('lap-track').addEventListener('change', function() {
            loadLapTimes(parseInt(this.value));
        });
        
        // Update button text when user selection changes
        document.getElementById('lap-user').addEventListener('change', function() {
            updateLapTimeButton();
        });
        
        async function deleteTrack(trackId, trackName) {
            if (!confirm('Are you sure you want to delete track "' + trackName + '"? This will also delete all lap times for this track.')) {
                return;
            }
            
            try {
                const response = await fetch('/api/tracks/' + trackId, {
                    method: 'DELETE'
                });
                
                if (response.ok) {
                    showMessage('message-track', 'Track deleted successfully!', 'success');
                    loadTracks();
                } else {
                    const error = await response.text();
                    showMessage('message-track', 'Error deleting track: ' + error, 'error');
                }
            } catch (error) {
                showMessage('message-track', 'Error deleting track', 'error');
            }
        }
        
        async function deleteLapTime(lapTimeId, trackId) {
            if (!confirm('Are you sure you want to delete this lap time?')) {
                return;
            }
            
            try {
                const response = await fetch('/api/laptimes/' + lapTimeId, {
                    method: 'DELETE'
                });
                
                if (response.ok) {
                    showMessage('message-lap', 'Lap time deleted successfully!', 'success');
                    loadLapTimes(trackId);
                } else {
                    const error = await response.text();
                    showMessage('message-lap', 'Error deleting lap time: ' + error, 'error');
                }
            } catch (error) {
                showMessage('message-lap', 'Error deleting lap time', 'error');
            }
        }
        
        function editTrack(trackId) {
            const nameDiv = document.getElementById('track-name-' + trackId);
            const editForm = document.getElementById('track-edit-' + trackId);
            if (nameDiv && editForm) {
                nameDiv.style.display = 'none';
                editForm.style.display = 'block';
            }
        }
        
        function cancelEdit(trackId) {
            const nameDiv = document.getElementById('track-name-' + trackId);
            const editForm = document.getElementById('track-edit-' + trackId);
            if (nameDiv && editForm) {
                nameDiv.style.display = 'block';
                editForm.style.display = 'none';
            }
        }
        
        async function saveTrack(trackId) {
            const name = document.getElementById('edit-name-' + trackId).value.trim();
            
            if (!name) {
                showMessage('message-track', 'Please enter a track name', 'error');
                return;
            }
            
            try {
                const response = await fetch('/api/tracks/' + trackId, {
                    method: 'PUT',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ id: trackId, name })
                });
                
                if (response.ok) {
                    showMessage('message-track', 'Track updated successfully!', 'success');
                    loadTracks();
                } else {
                    const error = await response.text();
                    showMessage('message-track', 'Error updating track: ' + error, 'error');
                }
            } catch (error) {
                showMessage('message-track', 'Error updating track', 'error');
            }
        }
    </script>
</body>
</html>`
