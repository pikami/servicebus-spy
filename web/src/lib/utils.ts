import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";
import { apiClient, type Message } from "./api-client";
import { toast } from "sonner";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export const resendMessage = async (message: Message) => {
  try {
    await apiClient.sendMessage({
      queueOrTopic: message.destination,
      body: JSON.stringify(message.applicationData),
      contentType: message.applicationDataType,
      subject: message.subject,
    });
    toast.success("Message resent");
  } catch (error) {
    console.error(error);
    toast.error("Failed to resend message");
  }
};
