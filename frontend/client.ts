import { writable, get } from 'svelte/store';

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

function isObject(value: unknown): value is Record<string, unknown> {
  return value !== null && typeof value === 'object' && !Array.isArray(value);
}

function isChatMessage(value: unknown): value is ChatMessage {
  if (!isObject(value)) return false;
  return (
    typeof value.ChatID === 'number' &&
    typeof value.Username === 'string' &&
    typeof value.Msg === 'string' &&
    typeof value.SendTime === 'string'
  );
}

function isChat(value: unknown): value is Chat {
  if (!isObject(value)) return false;

  return (
    typeof value.ID === 'number' &&
    typeof value.Name === 'string' &&
    (value.Type === 'individual' || value.Type === 'group') &&
    Array.isArray(value.Participants)
  );
}

function addMessages(newMessages: ChatMessage[]) {
  messages.update((currentMessages) => {
    const existingMessages = new Set(
      currentMessages.map(
        (msg) => `${msg.ChatID}-${msg.SendTime}-${msg.Username}-${msg.Msg}`,
      ),
    );

    const uniqueMessages = newMessages.filter((msg) => {
      const key = `${msg.ChatID}-${msg.SendTime}-${msg.Username}-${msg.Msg}`;
      if (existingMessages.has(key)) return false;
      existingMessages.add(key);
      return true;
    });

    return [...currentMessages, ...uniqueMessages];
  });
}

function updateChat(chat: Chat) {
  chats.update((currentChats) => [
    ...currentChats.filter((item) => item.ID !== chat.ID),
    chat,
  ]);
}

function handleIncomingData(data: unknown) {
  if (data === null || data === undefined) {
    console.warn('[Client] Ignoring null WebSocket data.');
    return;
  }

  if (Array.isArray(data)) {
    for (const item of data) {
      handleIncomingData(item);
    }
    return;
  }

  if (isChatMessage(data)) {
    addMessages([data]);
    console.log('[Client] Message received:', data);
    return;
  }

  if (isChat(data)) {
    updateChat(data);
    console.log('[Client] Chat received:', data);
    return;
  }

  console.warn('[Client] Unknown WebSocket data:', data);
}

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

    ws = new WebSocket(`${protocol}//${window.location.host}/api/chat`);

    ws.onopen = () => {
      console.log('[Client] Connected.');
    };

    ws.onmessage = (event: MessageEvent) => {
      try {
        const data: unknown = JSON.parse(event.data);
        handleIncomingData(data);
      } catch (error) {
        console.error('[Client] Invalid WebSocket JSON:', error);
      }
    };

    ws.onerror = (error) => {
      console.error('[Client] WebSocket error:', error);
    };

    ws.onclose = () => {
      console.log('[Client] Disconnected.');
      ws = null;
    };
  } catch (error) {
    console.error('[Client] Could not connect:', error);
    currentUser.set(null);
  }
}

export async function fetchUsers(): Promise<AuthUser[]> {
  const response = await fetch('/api/users', {
    credentials: 'include',
  });

  if (!response.ok) {
    throw new Error('Could not load users');
  }

  const result: unknown = await response.json();
  const data = Array.isArray(result) ? (result as AuthUser[]) : [];

  users.set(data);
  return data;
}

export async function fetchChats(): Promise<Chat[]> {
  const response = await fetch('/api/chats', {
    credentials: 'include',
  });

  if (!response.ok) {
    throw new Error('Could not load chats');
  }

  const result: unknown = await response.json();
  const data = Array.isArray(result) ? (result as Chat[]) : [];

  chats.set(data);
  return data;
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
    throw new Error('Could not create chat');
  }

  const updatedChats = await fetchChats();

  const currentUserID = get(currentUser)?.ID;

  const expectedParticipants = [
    ...new Set([
      ...(currentUserID !== undefined ? [currentUserID] : []),
      ...data.userIDs,
    ]),
  ].sort((a, b) => a - b);

  const matchingChats = updatedChats
    .filter((chat) => {
      const participantIDs = chat.Participants.map(
        (participant) => participant.ID,
      ).sort((a, b) => a - b);

      return (
        chat.Type === data.type &&
        chat.Name === (data.name ?? '') &&
        participantIDs.length === expectedParticipants.length &&
        participantIDs.every((id, index) => id === expectedParticipants[index])
      );
    })
    .sort((a, b) => b.ID - a.ID);

  const createdChat = matchingChats[0];

  if (!createdChat) {
    throw new Error('Chat was created, but could not be found.');
  }

  updateChat(createdChat);

  return createdChat;
}

export function sendMessage(chatID: number, username: string, message: string) {
  if (!ws || ws.readyState !== WebSocket.OPEN) {
    console.error('[Client] WebSocket is not connected.');
    return;
  }

  const data: ChatMessage = {
    infoType: 'chat',
    ChatID: chatID,
    Username: username,
    Msg: message,
    SendTime: new Date().toISOString(),
  };

  ws.send(JSON.stringify(data));

  addMessages([data]);
}
