import { renderLogin } from './login.js';
import { connectWebSocket, setOnMessage } from './ws.js';
import { renderRooms, setRoomsState } from './rooms.js';
import { renderChat, setChatState } from './chat.js';

let username = '';
let currentRoom = '';
let rooms = ['General', 'Random', 'Tech'];
let messages = {};

function handleLogin(name, password) {
    username = name;
    password = password;

    // Send POST request to /login
    fetch('/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username: name, password: password })
    })
    .then(async res => {
        if (res.ok) {
            connectWebSocket(onMessage);
            setRoomsState({ username, rooms, currentRoom, messages, onRoomChange });
            renderRooms();
        } else {
            const err = await res.text();
            alert('Login failed: ' + err);
        }
    })
    .catch(() => {
        alert('Network error during login.');
    });
}

function onRoomChange(room) {
    currentRoom = room;
    if (!messages[currentRoom]) messages[currentRoom] = [];
    setChatState({ username, currentRoom, messages, onSendMessage });
    renderChat();
}

function onSendMessage(text) {
    if (window.ws && window.ws.readyState === WebSocket.OPEN) {
        const msg = { user: username, room: currentRoom, text };
        window.ws.send(JSON.stringify(msg));
    }
}

function onMessage(msg) {
    if (!messages[msg.room]) messages[msg.room] = [];
    messages[msg.room].push(msg);
    if (msg.room === currentRoom) {
        setChatState({ username, currentRoom, messages, onSendMessage });
        renderChat();
    }
}

window.onload = () => {
    renderLogin(handleLogin);
};