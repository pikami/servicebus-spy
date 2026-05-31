import type { Message } from "./api-client";

export type FavoriteMessage = {
  uiId: string;
  queueOrTopic: string;
  subject: string;
  body: string;
  contentType: string;
};

export type SendMessageFormState = {
  queueOrTopic: string;
  subject: string;
  body: string;
  contentType: string;
};

const STORAGE_KEY = "servicebus-spy-favorites";

export function loadFavorites(): FavoriteMessage[] {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) return [];
    const parsed = JSON.parse(raw);
    return Array.isArray(parsed) ? parsed : [];
  } catch {
    return [];
  }
}

export function saveFavorites(favorites: FavoriteMessage[]): void {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(favorites));
}

export function messageToFavorite(message: Message): FavoriteMessage {
  const body =
    message.applicationDataType === "string"
      ? String(message.applicationData ?? "")
      : JSON.stringify(message.applicationData ?? null);

  return {
    uiId: crypto.randomUUID(),
    queueOrTopic: message.destination,
    subject: message.subject ?? "",
    body,
    contentType: message.applicationDataType,
  };
}

export function favoriteToFormState(
  favorite: FavoriteMessage,
): SendMessageFormState {
  return {
    queueOrTopic: favorite.queueOrTopic,
    subject: favorite.subject,
    body: favorite.body,
    contentType: favorite.contentType,
  };
}

export function formStateToFavorite(
  form: SendMessageFormState,
): FavoriteMessage {
  return {
    uiId: crypto.randomUUID(),
    queueOrTopic: form.queueOrTopic.trim(),
    subject: form.subject,
    body: form.body,
    contentType: form.contentType,
  };
}
