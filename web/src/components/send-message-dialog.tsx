import { useEffect, useState } from "react";
import { toast } from "sonner";
import { apiClient } from "@/lib/api-client";
import {
  formStateToFavorite,
  type SendMessageFormState,
} from "@/lib/favorites";
import { useFavorites } from "@/lib/favorites-context";
import { cn } from "@/lib/utils";
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
import { Field, FieldGroup, FieldLabel } from "./ui/field";
import { Input } from "./ui/input";

export type { SendMessageFormState };

interface SendMessageDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  initialValues?: SendMessageFormState;
}

const defaultFormState: SendMessageFormState = {
  subject: "",
  body: "",
  queueOrTopic: "",
  contentType: "application/json",
};

export function SendMessageDialog({
  open,
  onOpenChange,
  initialValues,
}: SendMessageDialogProps) {
  const [form, setForm] = useState(defaultFormState);
  const [isSending, setIsSending] = useState(false);
  const { addFavorite } = useFavorites();

  useEffect(() => {
    if (open) {
      setForm(initialValues ?? defaultFormState);
    }
  }, [open, initialValues]);

  const resetForm = () => setForm(defaultFormState);

  const handleOpenChange = (nextOpen: boolean) => {
    if (!nextOpen) {
      resetForm();
    }
    onOpenChange(nextOpen);
  };

  const handleSend = async () => {
    if (!form.queueOrTopic.trim()) {
      toast.error("Topic/queue name is required");
      return;
    }

    setIsSending(true);
    try {
      await apiClient.sendMessage({
        subject: form.subject,
        body: form.body,
        queueOrTopic: form.queueOrTopic.trim(),
        contentType: form.contentType,
      });
      toast.success("Message sent");
      handleOpenChange(false);
    } catch (error) {
      console.error(error);
      toast.error("Failed to send message");
    } finally {
      setIsSending(false);
    }
  };

  const handleSaveAsFavorite = () => {
    if (!form.queueOrTopic.trim()) {
      toast.error("Topic/queue name is required");
      return;
    }

    addFavorite(formStateToFavorite(form));
    toast.success("Saved to favorites");
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>Send Message</DialogTitle>
          <DialogDescription>
            Send a message to a Service Bus topic or queue.
          </DialogDescription>
        </DialogHeader>

        <FieldGroup>
          <Field>
            <FieldLabel htmlFor="queueOrTopic">Topic/Queue</FieldLabel>
            <Input
              id="queueOrTopic"
              value={form.queueOrTopic}
              onChange={(e) =>
                setForm((current) => ({
                  ...current,
                  queueOrTopic: e.target.value,
                }))
              }
              placeholder="my-topic-or-queue"
            />
          </Field>

          <Field>
            <FieldLabel htmlFor="subject">Subject</FieldLabel>
            <Input
              id="subject"
              value={form.subject}
              onChange={(e) =>
                setForm((current) => ({ ...current, subject: e.target.value }))
              }
              placeholder="Message subject"
            />
          </Field>

          <Field>
            <FieldLabel htmlFor="contentType">Content Type</FieldLabel>
            <Input
              id="contentType"
              value={form.contentType}
              onChange={(e) =>
                setForm((current) => ({
                  ...current,
                  contentType: e.target.value,
                }))
              }
            />
          </Field>

          <Field>
            <FieldLabel htmlFor="body">Body</FieldLabel>
            <textarea
              id="body"
              value={form.body}
              onChange={(e) =>
                setForm((current) => ({ ...current, body: e.target.value }))
              }
              placeholder='{"key": "value"}'
              rows={8}
              className={cn(
                "placeholder:text-muted-foreground border-input w-full min-w-0 rounded-md border bg-transparent px-3 py-2 text-base shadow-xs transition-[color,box-shadow] outline-none md:text-sm",
                "focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px]",
              )}
            />
          </Field>
        </FieldGroup>

        <DialogFooter className="gap-2 sm:justify-between">
          <Button
            type="button"
            variant="secondary"
            onClick={handleSaveAsFavorite}
            disabled={isSending}
          >
            Save as Favorite
          </Button>
          <div className="flex gap-2">
            <DialogClose asChild>
              <Button variant="outline" disabled={isSending}>
                Cancel
              </Button>
            </DialogClose>
            <Button onClick={handleSend} disabled={isSending}>
              {isSending ? "Sending..." : "Send"}
            </Button>
          </div>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
