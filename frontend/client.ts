const ws = new WebSocket(`/api/chat?user-id=1`);

type ChatMessage = {
  ChatID: number;
  Username: string;
  Msg: string;
  SendTime: string;
};

export const messages: ChatMessage[] = [
  {
    ChatID: 1,
    Username: 'yuri',
    Msg: 'Olá a todos!',
    SendTime: '2026-10-16T10:00:00Z',
  },
  {
    ChatID: 1,
    Username: 'yuri',
    Msg: 'Olá a todos!',
    SendTime: '2026-10-16T10:00:00Z',
  },
  {
    ChatID: 1,
    Username: 'yuri',
    Msg: 'Olá a todos!',
    SendTime: '2026-10-16T10:00:00Z',
  },
  {
    ChatID: 1,
    Username: 'yuri',
    Msg: 'Olá a todos!',
    SendTime: '2026-10-16T10:00:00Z',
  },
  {
    ChatID: 1,
    Username: 'yuri',
    Msg: 'Olá a todos!',
    SendTime: '2026-10-16T10:00:00Z',
  },
  {
    ChatID: 1,
    Username: 'yuri',
    Msg: 'Olá a todos!',
    SendTime: '2026-10-16T10:00:00Z',
  },
  {
    ChatID: 1,
    Username: 'yuri',
    Msg: 'Olá a todos!',
    SendTime: '2026-10-16T10:00:00Z',
  },
  {
    ChatID: 1,
    Username: 'pepe',
    Msg: 'Boas! Como está a correr o projeto?',
    SendTime: '2026-10-16T10:01:00Z',
  },
  {
    ChatID: 1,
    Username: 'yuri',
    Msg: 'Já consegui ligar o frontend ao WebSocket.',
    SendTime: '2026-10-16T10:02:00Z',
  },
  {
    ChatID: 2,
    Username: 'maria',
    Msg: 'Alguém pode rever o meu código?',
    SendTime: '2026-10-16T10:05:00Z',
  },
  {
    ChatID: 2,
    Username: 'yuri',
    Msg: 'Claro! Envia-me o link do repositório.',
    SendTime: '2026-10-16T10:06:00Z',
  },
  {
    ChatID: 3,
    Username: 'carlos',
    Msg: 'A reunião começa às 15:00.',
    SendTime: '2026-10-16T10:10:00Z',
  },
  {
    ChatID: 1,
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
    Username: username,
    Msg: message,
    SendTime: new Date().toISOString(),
  };

  ws.send(JSON.stringify(data));
}
