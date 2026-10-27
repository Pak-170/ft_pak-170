import { writable } from 'svelte/store';

type ChatMessage = {
  ChatID: number;
  Username: string;
  Msg: string;
  SendTime: string;
};

type AuthUser = {
  ID: number;
  username: string;
  alias: string | null;
};

export const messages = writable<ChatMessage[]>([]);
export const currentUser = writable<AuthUser | null>(null);

let ws: WebSocket | null = null;

export async function connectWebSocket() {
  if (
    ws &&
    (ws.readyState === WebSocket.OPEN || ws.readyState === WebSocket.CONNECTING)
  ) {
    return;
  }

  try {
    const response = await fetch('/api/me', {
      credentials: 'include',
    });
    if (!response.ok) {
      throw new Error('User is not authenticated');
    }

    const user = (await response.json()) as AuthUser;

    currentUser.set(user);

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';

    ws = new WebSocket(
      `${protocol}//${window.location.host}/api/chat?user-id=${user.ID}`,
    );

    ws.onopen = () => {
      console.log('[Client] connected.');
    };

    ws.onmessage = (event: MessageEvent) => {
      try {
        const data = JSON.parse(event.data) as ChatMessage | ChatMessage[];

        messages.update((currentMessages) => {
          if (Array.isArray(data)) {
            const existingMessages = new Set(
              currentMessages.map(
                (msg) =>
                  `${msg.ChatID}-${msg.SendTime}-${msg.Username}-${msg.Msg}`,
              ),
            );
            const newMessages = data.filter(
              (msg) =>
                !existingMessages.has(
                  `${msg.ChatID}-${msg.SendTime}-${msg.Username}-${msg.Msg}`,
                ),
            );
            return [...currentMessages, ...newMessages];
          }
          return [...currentMessages, data];
        });
        console.log('[Client] Messages received:', data);
      } catch {
        console.error('[Client] Invalid message received.');
      }
    };

    ws.onerror = (error) => {
      console.error('[Client] WebSocket error:', error);
    };

    ws.onclose = () => {
      console.log('[Client] disconnected.');
      ws = null;
    };
  } catch (error) {
    console.error('[Client] Could not connect: ', error);
    currentUser.set(null);
  }
}

export function sendMessage(chatID: number, username: string, message: string) {
  if (!ws || ws.readyState !== WebSocket.OPEN) {
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

  messages.update((currentMessages) => [...currentMessages, data]);
}
