let state = {};

export function setChatState(newState) {
    state = { ...state, ...newState };
}

export function renderChat() {
    document.getElementById('chatArea').innerHTML = `
        <h3>Room: ${state.currentRoom}</h3>
        <div class="messages" id="messages">
            ${(state.messages[state.currentRoom] || []).map(msg => `
                <div class="message"><b>${msg.user}:</b> ${msg.text}</div>
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