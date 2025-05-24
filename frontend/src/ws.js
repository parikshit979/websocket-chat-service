let onMessageCallback = null;

export function connectWebSocket(onMessage) {
    // Replace with your actual WSS server URL
    window.ws = new WebSocket("ws://" + document.location.host + "/ws");
    onMessageCallback = onMessage;
    window.ws.onopen = () => console.log('WebSocket connected');
    window.ws.onmessage = (event) => {
        const msg = JSON.parse(event.data);
        if (onMessageCallback) onMessageCallback(msg);
    };
    window.ws.onclose = () => alert('WebSocket disconnected');
}

export function setOnMessage(cb) {
    onMessageCallback = cb;
}