// Dynamically construct the WebSocket URL
const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'; // Use 'wss' for secure connections, 'ws' for non-secure
const host = window.location.host; // Gets the host (domain or localhost:port)
const socket = new WebSocket(`${protocol}://${host}/chat`); // WebSocket endpoint

socket.onopen = () => {
    console.log('WebSocket connected!');
};

socket.onmessage = (event) => {
    const message = JSON.parse(event.data);
    displayMessage(message);
};

function sendMessage(receiverID, messageContent) {
    const message = {
        senderID: currentUserID,  // Replace this with the actual sender's user ID
        receiverID: receiverID,
        message: messageContent
    };
    socket.send(JSON.stringify(message));
}

function displayMessage(message) {
    // Render the message in the UI
    const chatBox = document.getElementById('chat-box');
    const messageDiv = document.createElement('div');
    messageDiv.textContent = `${message.senderID}: ${message.message}`;
    chatBox.appendChild(messageDiv);
}
