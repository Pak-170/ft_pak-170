import { get, writable } from 'svelte/store';

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

export type AppNotification = {
  id: number;
  type: 'error' | 'message';
  message: string;
  chatID?: number;
};

export const users = writable<AuthUser[]>([]);
export const chats = writable<Chat[]>([]);
export const messages = writable<ChatMessage[]>([]);
export const currentUser = writable<AuthUser | null>(null);

export const selectedChatID = writable<number | null>(null);
export const unreadMessages = writable<Record<number, number>>({});
export const notifications = writable<AppNotification[]>([]);
export const chatPageActive = writable(false);

let ws: WebSocket | null = null;
let connectionPromise: Promise<void> | null = null;
let notificationID = 0;

export function addErrorNotification(message: string) {
  notifications.update((current) => [
    ...current,
    {
      id: ++notificationID,
      type: 'error',
      message,
    },
  ]);
}

export function dismissNotification(id: number) {
  notifications.update((current) =>
    current.filter((notification) => notification.id !== id),
  );
}

export function clearMessageNotifications() {
  notifications.update((current) =>
    current.filter((notification) => notification.type !== 'message'),
  );
}

export function markChatAsRead(chatID: number) {
  unreadMessages.update((current) => {
    const updated = { ...current };
    delete updated[chatID];
    return updated;
  });

  notifications.update((current) =>
    current.filter(
      (notification) =>
        notification.type !== 'message' || notification.chatID !== chatID,
    ),
  );
}

function getChatName(chatID: number) {
  const chat = get(chats).find((item) => item.ID === chatID);

  if (!chat) {
    return `Chat ${chatID}`;
  }

  if (chat.Type === 'individual') {
    const user = chat.Participants?.find(
      (participant) => participant.ID !== get(currentUser)?.ID,
    );

    return user?.alias || user?.username || 'Private Chat';
  }

  return chat.Name || `Chat ${chatID}`;
}

function handleIncomingMessage(message: ChatMessage) {
  messages.update((current) => [...current, message]);

  const user = get(currentUser);

  if (!user || message.Username === user.username) {
    return;
  }

  if (get(selectedChatID) === message.ChatID) {
    return;
  }

  const currentCount = get(unreadMessages)[message.ChatID] ?? 0;
  const newCount = currentCount + 1;

  unreadMessages.update((current) => ({
    ...current,
    [message.ChatID]: newCount,
  }));

  if (get(chatPageActive)) {
    return;
  }

  const chatName = getChatName(message.ChatID);
  const notificationMessage =
    `You have ${newCount} new messages from ${chatName}`;

  notifications.update((current) => {
    const existing = current.find(
      (notification) =>
        notification.type === 'message' &&
        notification.chatID === message.ChatID,
    );

    if (existing) {
      return current.map((notification) =>
        notification.id === existing.id
          ? { ...notification, message: notificationMessage }
          : notification,
      );
    }

    return [
      ...current,
      {
        id: ++notificationID,
        type: 'message',
        chatID: message.ChatID,
        message: notificationMessage,
      },
    ];
  });
}

function handleIncomingData(data: any) {
  if (Array.isArray(data)) {
    data.forEach(handleIncomingData);
    return;
  }

  if (!data) return;

  if (typeof data.ChatID === 'number') {
    handleIncomingMessage(data as ChatMessage);
    return;
  }

  if (typeof data.ID === 'number' && Array.isArray(data.Participants)) {
    chats.update((current) => {
      const index = current.findIndex((chat) => chat.ID === data.ID);

      if (index === -1) {
        return [...current, data];
      }

      const updated = [...current];
      updated[index] = data;
      return updated;
    });
  }
}

export async function connectWebSocket() {
  if (ws?.readyState === WebSocket.OPEN) {
    return;
  }

  if (connectionPromise) {
    return connectionPromise;
  }

  connectionPromise = (async () => {
    let connectionEstablished = false;

    try {
      const response = await fetch('/api/me');

      if (!response.ok) {
        throw new Error('Could not authenticate user.');
      }

      const user: AuthUser = await response.json();
      currentUser.set(user);

      const socket = new WebSocket('/api/chat');
      ws = socket;

      await new Promise<void>((resolve, reject) => {
        socket.onopen = () => {
          connectionEstablished = true;
          resolve();
        };

        socket.onmessage = (event) => {
          try {
            handleIncomingData(JSON.parse(event.data));
          } catch (error) {
            console.error('Invalid WebSocket message:', error);
          }
        };

        socket.onerror = () => {
          if (!connectionEstablished) {
            reject(new Error('WebSocket connection failed.'));
          }
        };

        socket.onclose = () => {
          if (ws === socket) {
            ws = null;
          }

          if (connectionEstablished) {
            addErrorNotification('Chat connection was lost.');
          }
        };
      });
    } catch (error) {
      if (!connectionEstablished) {
        addErrorNotification('Could not connect to the chat server.');
      }

      throw error;
    }
  })();

  try {
    await connectionPromise;
  } finally {
    connectionPromise = null;
  }
}

export async function fetchUsers() {
  try {
    const response = await fetch('/api/users');

    if (!response.ok) {
      throw new Error('Could not fetch users.');
    }

    const data: AuthUser[] = await response.json();
    users.set(data);
  } catch (error) {
    addErrorNotification('Could not load users.');
    throw error;
  }
}

export async function fetchChats() {
  try {
    const response = await fetch('/api/chats');

    if (!response.ok) {
      throw new Error('Could not fetch chats.');
    }

    const data: Chat[] = await response.json();
    chats.set(data);
  } catch (error) {
    addErrorNotification('Could not load chats.');
    throw error;
  }
}

export function disconnectWebSocket() {
  const socket = ws;

  ws = null;
  connectionPromise = null;

  if (socket) {
    socket.close();
  }

  currentUser.set(null);
  users.set([]);
  chats.set([]);
  messages.set([]);
  selectedChatID.set(null);
  unreadMessages.set({});
  notifications.set([]);
  chatPageActive.set(false);
}

export async function logout() {
  try {
    const response = await fetch('/api/logout', {
      method: 'POST',
    });

    if (!response.ok) {
      throw new Error('Could not log out.');
    }

    disconnectWebSocket();
  } catch (error) {
    addErrorNotification('Could not log out.');
    throw error;
  }
}

export async function createChat(data: CreateChatData): Promise<Chat> {
  try {
    const response = await fetch('/api/chats', {
      method: 'POST',
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
  } catch (error) {
    addErrorNotification('Could not create chat.');
    throw error;
  }
}

export function sendMessage(
  chatID: number,
  username: string,
  message: string,
): boolean {
  if (!ws || ws.readyState !== WebSocket.OPEN) {
    addErrorNotification('Could not send message. Chat is disconnected.');
    return false;
  }

  const newMessage: ChatMessage = {
    infoType: 'chat',
    ChatID: chatID,
    Username: username,
    Msg: message,
    SendTime: new Date().toISOString(),
  };

  try {
    ws.send(JSON.stringify(newMessage));
    messages.update((current) => [...current, newMessage]);
    return true;
  } catch (error) {
    addErrorNotification('Could not send message.');
    return false;
  }
}