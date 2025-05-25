import { renderLogin } from './login.js';
import { connectWebSocket, setOnMessage } from './ws.js';
import { renderRooms, setRoomsState } from './rooms.js';
import { renderChat, setChatState } from './chat.js';

let username = '';
let currentRoom = '';
let rooms = ['General', 'Random', 'Tech'];
let messages = {};
let changeEventType = 'change_room'
let sendMessageEventType = "send_message"
let receiveMessageEventType = "receive_message"

/**
 * Event is used to wrap all messages Send and Recieved
 * on the Websocket
 * The type is used as a RPC
 * */
class Event {
    // Each Event needs a Type
    // The payload is not required
    constructor(type, payload) {
        this.type = type;
        this.payload = payload;
    }
}

/**
 * SendMessageEvent is used to send messages to other clients
 * */
class SendMessageEvent {
    constructor(message, from, room) {
        this.message = message;
        this.from = from;
        this.room = room
    }
}
/**
 * ReceiveMessageEvent is messages comming from clients
 * */
class ReceiveMessageEvent {
    constructor(message, from, sent, room) {
        this.message = message;
        this.from = from;
        this.sent = sent;
        this.room = room
    }
}

/**
 * ChangeChatRoomEvent is used to switch chatroom
 * */
class ChangeChatRoomEvent {
    constructor(room) {
        this.room = room;
    }
}

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
    let changeEvent = new ChangeChatRoomEvent(currentRoom);
    sendEvent(changeEventType, changeEvent);
}

function onSendMessage(text) {
    if (window.ws && window.ws.readyState === WebSocket.OPEN) {
        // const msg = { user: username, room: currentRoom, text };
        // window.ws.send(JSON.stringify(msg));
        let outgoingEvent = new SendMessageEvent(text, username, currentRoom);
        sendEvent(sendMessageEventType, outgoingEvent)
    }
}

function onMessage(eventData) {
    console.log(eventData)
    const eventMsg = Object.assign(new Event, eventData);
    switch(eventMsg.type){
        case receiveMessageEventType:
            var msg = eventMsg.payload
            if (!messages[msg.room]) messages[msg.room] = [];
            var date = new Date(msg.sent);
            msg.sent = date
            messages[msg.room].push(msg);
            if (msg.room === currentRoom) {
                setChatState({ username, currentRoom, messages, onSendMessage });
                renderChat();
            }
            break;
        default:
            alert("unsupported message type");
            break;
    }
}

function sendEvent(eventName, payload) {
    // Create a event Object with a event named send_message
    const msg = new Event(eventName, payload);
    console.log(JSON.stringify(msg))
    // Format as JSON and send
    window.ws.send(JSON.stringify(msg));
}

window.onload = () => {
    renderLogin(handleLogin);
};