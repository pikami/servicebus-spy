import { apiClient, type Message } from "@/lib/api-client";
import { Button } from "./ui/button";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "./ui/dialog";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "./ui/accordion";
import { CodeBlock } from "./ui/code-block";

interface MessageDetailsDialogProps {
  message: Message | null;
  onClose: () => void;
}

const resendMessage = async (message: Message) => {
  try {
    await apiClient.sendMessage({
      queueOrTopic: message.fullDestination,
      body: JSON.stringify(message.applicationData),
      contentType: message.applicationDataType,
      subject: message.subject,
    });
  } catch (error) {
    console.error(error);
  }
};

export function MessageDetailsDialog({
  message,
  onClose,
}: MessageDetailsDialogProps) {
  return (
    <Dialog open={message !== null} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="max-h-[80vh] sm:max-w-7xl overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Message Details</DialogTitle>
          <DialogDescription>
            Full destination: {message?.fullDestination}
          </DialogDescription>
        </DialogHeader>

        <Accordion
          type="multiple"
          className="w-full"
          defaultValue={["application-data"]}
        >
          <AccordionItem value="application-data">
            <AccordionTrigger>Application Data</AccordionTrigger>
            <AccordionContent>
              <CodeBlock>
                {JSON.stringify(message?.applicationData, null, 2)}
              </CodeBlock>
            </AccordionContent>
          </AccordionItem>

          <AccordionItem value="application-properties">
            <AccordionTrigger>Application Properties</AccordionTrigger>
            <AccordionContent>
              <CodeBlock>
                {JSON.stringify(message?.applicationProperties, null, 2)}
              </CodeBlock>
            </AccordionContent>
          </AccordionItem>

          <AccordionItem value="message-properties">
            <AccordionTrigger>Message Properties</AccordionTrigger>
            <AccordionContent>
              <CodeBlock>
                {JSON.stringify(message?.messageProperties, null, 2)}
              </CodeBlock>
            </AccordionContent>
          </AccordionItem>
        </Accordion>

        <DialogFooter>
          <DialogClose asChild>
            <Button variant="outline">Close</Button>
          </DialogClose>
          <Button variant="outline" onClick={() => resendMessage(message!)}>
            Resend Message
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
