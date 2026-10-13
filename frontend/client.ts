const port = 1234;
const ws = new WebSocket(`ws://localhost:${port}`);

type ChatMessage = {
	username: string;
	message: string;
};

export const messages: ChatMessage[] = [];

ws.onopen = () => {
	console.log('[Client] connected.');
};

ws.onmessage = (event: MessageEvent) => {
	try {
		const data = JSON.parse(event.data) as ChatMessage;

		messages.push(data);
	} catch {
		console.error('[Client] Invalid message recieved.');
	}
};

ws.onerror = (error) => {
	console.error('[Client] WebSocket error:', error);
};

ws.onclose = () => {
	console.log('[Client] disconnected.');
};

export function sendMessage(data: ChatMessage) {
	if (ws.readyState !== WebSocket.OPEN) {
		console.error('[Client] WebSocket is not conected.');
		return;
	}
	ws.send(JSON.stringify(data));
}
