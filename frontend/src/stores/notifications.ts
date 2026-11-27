import { writable } from 'svelte/store';

export type AppNotification = {
  id: number;
  message: string;
};

export const unreadByChat = writable<Record<number, number>>({});
export const activeChatID = writable<number | null>(null);
export const errorNotifications = writable<AppNotification[]>([]);

let notificationID = 0;

export function setActiveChat(chatID: number | null) {
  activeChatID.set(chatID);

  if (chatID !== null) {
    unreadByChat.update((current) => {
      const updated = { ...current };
      delete updated[chatID];
      return updated;
    });
  }
}

export function addUnreadMessage(chatID: number) {
  unreadByChat.update((current) => ({
    ...current,
    [chatID]: (current[chatID] ?? 0) + 1,
  }));
}

export function addErrorNotification(message: string) {
  const id = ++notificationID;

  errorNotifications.update((current) => [...current, { id, message }]);
}

export function removeErrorNotification(id: number) {
  errorNotifications.update((current) =>
    current.filter((notification) => notification.id !== id),
  );
}

export function clearNotifications() {
  unreadByChat.set({});
  activeChatID.set(null);
  errorNotifications.set([]);
}
