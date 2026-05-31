export type Message = {
  destination: string;
  role: "sender" | "receiver";
  messageId: string;
  subject: string;

  uiId: string;
  fullDestination: string;
  subscriber: string | null;

  applicationProperties: Record<string, any>;
  applicationData: any;
  applicationDataType: "string" | "json" | "bytes" | "unknown";
  messageProperties: Record<string, any>;
};

export type ApiMessage = {
  linkInfo: {
    role: "sender" | "receiver";
    sourceAddress: string;
    targetAddress: string;
  };
  messageAnnotations: Record<string, any>;
  messageProperties: Record<string, any>;
  applicationData: any;
  applicationDataType: "string" | "json" | "bytes" | "unknown";
  applicationProperties: Record<string, any>;
  raw: string;
};

export type ServiceBusMessage = {
  queueOrTopic: string;
  body: string;
  contentType: string;
  subject: string;
};

class APIClient {
  private baseUrl: string;

  constructor(baseUrl: string = "/api") {
    this.baseUrl = baseUrl;
  }

  async getMessages(): Promise<Message[]> {
    const response = await fetch(`${this.baseUrl}/messages`);
    const apiMessages = await response.json();
    return apiMessages.map(mapApiMessageToMessage);
  }

  async sendMessage(message: ServiceBusMessage): Promise<void> {
    const response = await fetch(`${this.baseUrl}/messages/send`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(message),
    });

    if (!response.ok) {
      console.error(await response.text());
      throw new Error(`Failed to send message: ${response.statusText}`);
    }
  }
}

export const apiClient = new APIClient();

const mapApiMessageToMessage = (apiMessage: ApiMessage): Message => {
  const shortDestination = apiMessage.linkInfo.targetAddress
    .replace(/^\//, "")
    .replace(/\/Subscriptions\/.+$/, "")
    .replace(
      /-[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}-Receiver$/,
      "",
    );

  const subscriber = apiMessage.linkInfo.targetAddress.includes("Subscriptions")
    ? apiMessage.linkInfo.targetAddress
        .split("/Subscriptions/")[1]
        .replace(
          /-[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}-Receiver$/,
          "",
        )
    : null;

  return {
    destination: shortDestination,
    role: apiMessage.linkInfo.role,
    messageId:
      apiMessage.messageProperties["message-id"] ??
      `unknown-${crypto.randomUUID()}`,
    subject: apiMessage.messageProperties["subject"],

    uiId: crypto.randomUUID(),
    fullDestination: apiMessage.linkInfo.targetAddress,
    subscriber: subscriber ?? null,

    applicationProperties: apiMessage.applicationProperties?.["_fields"] ?? {},
    applicationData: apiMessage.applicationData,
    applicationDataType: apiMessage.applicationDataType,
    messageProperties: apiMessage.messageProperties,
  };
};
