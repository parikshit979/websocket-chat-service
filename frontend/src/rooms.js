let state = {};

export function setRoomsState(newState) {
    state = { ...state, ...newState };
}

export function renderRooms() {
    document.getElementById('app').innerHTML = `
        <h2>Welcome, ${state.username}</h2>
        <div class="chat-room-list">
            <strong>Rooms:</strong>
            ${state.rooms.map(room => `
                <button ${room === state.currentRoom ? 'disabled' : ''} data-room="${room}">${room}</button>
            `).join('')}
        </div>
        <div id="chatArea"></div>
    `;
    document.querySelectorAll('.chat-room-list button').forEach(btn => {
        btn.onclick = () => {
            state.onRoomChange(btn.dataset.room);
        };
    });
    if (state.currentRoom) {
        state.onRoomChange(state.currentRoom);
    }
}