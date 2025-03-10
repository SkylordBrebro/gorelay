
function updateStatus() {
    fetch('/api/status')
        .then(response => response.json())
        .then(data => {
            Object.keys(data).forEach(alias => {
                const client = data[alias];
                const clientBox = document.getElementById('client-' + alias);
                if (clientBox) {
                    updateClientBox(clientBox, client);
                }
            });
        });
}

function updateLogs() {
    fetch('/api/logs')
        .then(response => response.json())
        .then(data => {
            Object.keys(data).forEach(alias => {
                const logs = data[alias];
                const logBox = document.getElementById('logs-' + alias);
                if (logBox) {
                    updateLogBox(logBox, logs);
                }
            });
        });
}

function updateClientBox(box, client) {
    const status = box.querySelector('.client-status');
    if (client.connected) {
        status.textContent = 'Online';
        status.className = 'client-status status-online';
    } else {
        status.textContent = 'Offline';
        status.className = 'client-status status-offline';
    }

    const details = box.querySelector('.client-details');
    if (client.stats) {
        details.innerHTML = '';
        Object.keys(client.stats).forEach(key => {
            const value = client.stats[key];
            const item = document.createElement('div');
            item.className = 'detail-item';
            item.textContent = key + ': ' + value;
            details.appendChild(item);
        });
    }
}

function updateLogBox(box, logs) {
    if (!logs || !logs.length) return;
    
    const newLogs = logs.slice(-100); // Keep last 100 logs
    box.innerHTML = '';
    
    newLogs.forEach(log => {
        const entry = document.createElement('div');
        entry.className = 'log-entry log-' + log.level;
        entry.textContent = '[' + new Date(log.timestamp).toLocaleTimeString() + '] ' + log.message;
        box.appendChild(entry);
    });
    
    box.scrollTop = box.scrollHeight;
}

function reconnectClient(alias) {
    fetch('/reconnect/' + alias, { method: 'POST' })
        .then(response => {
            if (response.ok) {
                console.log('Reconnect request sent for ' + alias);
            } else {
                console.error('Failed to send reconnect request for ' + alias);
            }
        });
}

// Update status and logs every second
setInterval(updateStatus, 1000);
setInterval(updateLogs, 1000);
