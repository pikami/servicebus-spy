import { DataTable } from "@/components/data-table";
import { MessageDetailsDialog } from "@/components/message-details-dialog";
import type { Message } from "@/lib/api-client";
import type { ColumnDef, ColumnFiltersState, Row } from "@tanstack/react-table";
import { useState } from "react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
} from "@/components/ui/dropdown-menu";
import { DropdownMenuTrigger } from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";
import { MoreHorizontal } from "lucide-react";
import { TableHeader } from "./message-table-header";
import { resendMessage } from "@/lib/utils";

interface MessageTableProps {
  messages: Message[];
}

export function MessageTable({ messages }: MessageTableProps) {
  const [messageDetails, setMessageDetails] = useState<Message | null>(null);

  const columns: ColumnDef<Message>[] = [
    {
      id: "messageId",
      header: "Message ID",
      accessorKey: "messageId",
      filterFn: (row: Row<Message>, _, filterValue: any) => {
        return (
          !filterValue.distinct || filterValue.uiIds.includes(row.original.uiId)
        );
      },
    },
    {
      header: "Destination",
      accessorKey: "destination",
      filterFn: "arrIncludesSome",
    },
    {
      header: "Subscriber",
      accessorKey: "subscriber",
      filterFn: "arrIncludesSome",
    },
    {
      id: "role",
      header: "Role",
      accessorKey: "role",
      filterFn: "arrIncludesSome",
    },
    {
      header: "Subject",
      accessorKey: "subject",
      filterFn: "arrIncludesSome",
    },
    {
      id: "applicationDataPreview",
      header: "Application Data",
      accessorFn: (row) => JSON.stringify(row.applicationData),
      cell: ({ row }) => {
        return (
          <pre title={JSON.stringify(row.original.applicationData)}>
            {JSON.stringify(row.original.applicationData).substring(0, 10)}...
          </pre>
        );
      },
      filterFn: "includesString",
    },
    {
      id: "actions",
      cell: ({ row }) => {
        return (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="h-8 w-8 p-0">
                <span className="sr-only">Open menu</span>
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuLabel>Actions</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem onClick={() => setMessageDetails(row.original)}>
                View message details
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => resendMessage(row.original)}>
                Resend message
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        );
      },
    },
  ];

  const initialFilterValues: ColumnFiltersState = [
    {
      id: "messageId",
      value: {
        distinct: true,
        uiIds: [],
      },
    },
  ];

  return (
    <>
      <DataTable
        columns={columns}
        data={messages}
        initialFilterValues={initialFilterValues}
        headerComponent={TableHeader}
      />
      <MessageDetailsDialog
        message={messageDetails}
        onClose={() => setMessageDetails(null)}
      />
    </>
  );
}
