import { useState } from "react";
import { ModeToggle } from "./mode-toggle";
import { SendMessageDialog } from "./send-message-dialog";
import { Button } from "./ui/button";

export function SiteHeader() {
  const [sendDialogOpen, setSendDialogOpen] = useState(false);

  return (
    <>
      <header className="bg-background sticky top-0 z-50 w-full">
        <div className="container-wrapper flex items-center justify-between px-6 py-2">
          <h1 className="text-2xl font-bold">ServiceBus Spy</h1>
          <div className="flex items-center gap-2">
            <Button variant="outline" onClick={() => setSendDialogOpen(true)}>
              Send Message
            </Button>
            <ModeToggle />
          </div>
        </div>
      </header>
      <SendMessageDialog
        open={sendDialogOpen}
        onOpenChange={setSendDialogOpen}
      />
    </>
  );
}
