const params = new URLSearchParams(window.location.hash.split('?')[1] ?? '');
const userId = params.get('user-id') ?? '1';
const ws = new WebSocket(`/api/chat?user-id=${userId}`);

type ChatMessage = {
  ChatID: number;
  UserID: number;
  Username: string;
  Msg: string;
  SendTime: string;
};

export const messages: ChatMessage[] = [
  {
    ChatID: 1,
	UserID: 1,
    Username: 'yuri',
    Msg: 'Olá a todos!',
    SendTime: '2026-10-16T10:00:00Z',
  },
  {
    ChatID: 2,
	UserID: 1,
    Username: 'yuri',
    Msg: 'Olá a todos!',
    SendTime: '2026-10-16T10:00:00Z',
  },
  {
    ChatID: 3,
	UserID: 1,
    Username: 'yuri',
    Msg: 'Olá a todos!',
    SendTime: '2026-10-16T10:00:00Z',
  },
];

ws.onopen = () => {
  console.log('[Client] connected.');
};

ws.onmessage = (event: MessageEvent) => {
  try {
    const data = JSON.parse(event.data) as ChatMessage | ChatMessage[];

    if (Array.isArray(data)) {
      messages.push(...data);
    } else {
      messages.push(data);
    }
  } catch {
    console.error('[Client] Invalid message received.');
  }
};

ws.onerror = (error) => {
  console.error('[Client] WebSocket error:', error);
};

ws.onclose = () => {
  console.log('[Client] disconnected.');
};

export function sendMessage(chatID: number, username: string, message: string) {
  if (ws.readyState !== WebSocket.OPEN) {
    console.error('[Client] WebSocket is not connected.');
    return;
  }

  const data: ChatMessage = {
    ChatID: chatID,
	UserID: Number(userId),
    Username: username,
    Msg: message,
    SendTime: new Date().toISOString(),
  };

  ws.send(JSON.stringify(data));
}
