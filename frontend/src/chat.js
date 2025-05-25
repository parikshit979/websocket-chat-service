let state = {};

export function setChatState(newState) {
    state = { ...state, ...newState };
}

function formatDate(date) {
    const d = date.getDate();
    const daySuffix = (n) => {
        if (n > 3 && n < 21) return 'th';
        switch (n % 10) {
            case 1: return 'st';
            case 2: return 'nd';
            case 3: return 'rd';
            default: return 'th';
        }
    };
    const month = date.toLocaleString('default', { month: 'long' });
    const year = date.getFullYear();
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    const seconds = String(date.getSeconds()).padStart(2, '0');
    return `${d}${daySuffix(d)} ${month} ${year} ${hours}:${minutes}:${seconds}`;
}

export function renderChat() {
    document.getElementById('chatArea').innerHTML = `
        <h3>Room: ${state.currentRoom}</h3>
        <div class="messages" id="messages">
            ${(state.messages[state.currentRoom] || []).map(msg => `
                <div class="message"><b>${formatDate(msg.sent)}</b> <b>${msg.from}:</b> ${msg.message}</div>
            `).join('')}
        </div>
        <input id="msgInput" placeholder="Type a message..." autocomplete="off" />
        <button id="sendBtn">Send</button>
    `;
    const chatArea = document.getElementById('chatArea');
    chatArea.scrollTop = chatArea.scrollHeight;
    document.getElementById('sendBtn').onclick = sendMessage;
    document.getElementById('msgInput').onkeydown = e => {
        if (e.key === 'Enter') sendMessage();
    };
}

function sendMessage() {
    const input = document.getElementById('msgInput');
    const text = input.value.trim();
    if (text) {
        state.onSendMessage(text);
        input.value = '';
    }
}