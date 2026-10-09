import WebSocket from "ws";

const port = 1234;
const ws = new WebSocket(`ws://localhost:${port}`);

ws.on('open', () => {
	console.log('[Client] connected.');
	ws.send('Hi, this is clienTest!');
});

ws.on('message', (data) => {
	console.log(`Recieved a message from Server: ${data}`);
});