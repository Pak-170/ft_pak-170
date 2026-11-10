import { writable } from 'svelte/store';

type ChatMessage = {
  infoType?: 'chat';
  ChatID: number;
  Username: string;
  Msg: string;
  SendTime: string;
};

type AuthUser = {
  infoType?: 'chat';
  ID: number;
  username: string;
  alias: string | null;
};

type Chat = {
  infoType?: 'chat';
  ID: number;
  Name: string;
  Type: 'individual' | 'group';
  Participants: AuthUser[];
};

type CreateChatData = {
  infoType: 'chat';
  name?: string;
  type: 'individual' | 'group';
  userIDs: number[];
};

export const users = writable<AuthUser[]>([]);
export const chats = writable<Chat[]>([]);
export const messages = writable<ChatMessage[]>([]);
export const currentUser = writable<AuthUser | null>(null);

let ws: WebSocket | null = null;

function handleIncomingData(data: any) {
  if (Array.isArray(data)) {
    data.forEach(handleIncomingData);
    return;
  }

  if (!data) return;

  if (typeof data.ChatID === 'number') {
    messages.update((current) => [...current, data]);
  } else if (typeof data.ID === 'number' && Array.isArray(data.Participants)) {
    chats.update((current) => {
      const index = current.findIndex((chat) => chat.ID === data.ID);

      if (index === -1) return [...current, data];

      const updated = [...current];
      updated[index] = data;
      return updated;
    });
  }
}

export async function connectWebSocket() {
  if (ws?.readyState === WebSocket.OPEN) return;

  const response = await fetch('/api/me', {
    credentials: 'include',
  });

  if (!response.ok) {
    throw new Error('Could not authenticate user.');
  }

  const user: AuthUser = await response.json();
  currentUser.set(user);

  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  ws = new WebSocket(`${protocol}//${window.location.host}/api/chat`);

  ws.onmessage = (event) => {
    try {
      handleIncomingData(JSON.parse(event.data));
    } catch (error) {
      console.error('Invalid WebSocket message:', error);
    }
  };

  await new Promise<void>((resolve, reject) => {
    ws!.onopen = () => resolve();
    ws!.onerror = () => reject(new Error('WebSocket connection failed.'));
  });
}

export async function fetchUsers() {
  const response = await fetch('/api/users', {
    credentials: 'include',
  });

  if (!response.ok) {
    throw new Error('Could not fetch users.');
  }

  const data: AuthUser[] = await response.json();
  users.set(data);
}

export async function fetchChats() {
  const response = await fetch('/api/chats', {
    credentials: 'include',
  });

  if (!response.ok) {
    throw new Error('Could not fetch chats.');
  }

  const data: Chat[] = await response.json();
  chats.set(data);
}

export function disconnectWebSocket() {
  if (ws) {
    ws.close();
    ws = null;
  }
  users.set([]);
  chats.set([]);
  messages.set([]);
}

export async function createChat(data: CreateChatData): Promise<Chat> {
  const response = await fetch('/api/chats', {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  });

  if (!response.ok) {
    throw new Error('Could not create chat.');
  }

  const chat: Chat = await response.json();

  chats.update((current) => [...current, chat]);

  return chat;
}

export function sendMessage(chatID: number, username: string, message: string) {
  if (!ws || ws.readyState !== WebSocket.OPEN) return;

  const newMessage: ChatMessage = {
    infoType: 'chat',
    ChatID: chatID,
    Username: username,
    Msg: message,
    SendTime: new Date().toISOString(),
  };

  ws.send(JSON.stringify(newMessage));
  messages.update((current) => [...current, newMessage]);
}
