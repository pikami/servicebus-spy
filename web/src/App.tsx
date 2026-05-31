import { useEffect, useState } from "react";
import { apiClient, type Message } from "./lib/api-client";
import { SiteHeader } from "@/components/site-header";
import { MessageTable } from "@/components/message-table/message-table";
import { Toaster } from "@/components/ui/sonner";

function App() {
  const [messages, setMessages] = useState<Message[]>([]);
  useEffect(() => {
    apiClient.getMessages().then((messages) => {
      setMessages(messages);
    });
  }, []);

  return (
    <>
      <div
        data-slot="layout"
        className="bg-background relative z-10 flex min-h-svh flex-col"
      >
        <SiteHeader />
        <div className="flex w-full items-center justify-center p-6 md:p-10">
          <div className="w-full max-w-7xl">
            <MessageTable messages={messages} />
          </div>
        </div>
      </div>
      <Toaster />
    </>
  );
}

export default App;
