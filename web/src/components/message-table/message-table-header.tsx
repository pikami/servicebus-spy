import type { Message } from "@/lib/api-client";
import type { Table } from "@tanstack/react-table";
import { useEffect, useMemo } from "react";
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
} from "@/components/ui/dropdown-menu";
import { DropdownMenuTrigger } from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";
import { Checkbox } from "../ui/checkbox";
import { Field } from "../ui/field";
import { Label } from "../ui/label";
import { Input } from "../ui/input";

interface MessageIdFilterValue {
  distinct: boolean;
  uiIds: string[];
}

function MessageIdFilter({
  table,
  data,
}: {
  table: Table<Message>;
  data: Message[];
}) {
  const messageIdFilterValue = table
    .getColumn("messageId")
    ?.getFilterValue() as MessageIdFilterValue;

  const distinctUiIdsByMessageId = useMemo(() => {
    const seen = new Set<string>();
    return data
      .filter(
        (message) =>
          !seen.has(message.messageId) && seen.add(message.messageId),
      )
      .map((message) => message.uiId);
  }, [data]);
  useEffect(() => {
    table.getColumn("messageId")?.setFilterValue((x: MessageIdFilterValue) => ({
      ...x,
      uiIds: distinctUiIdsByMessageId,
    }));
  }, [distinctUiIdsByMessageId]);

  return (
    <Field orientation="horizontal">
      <Checkbox
        id="distinct-checkbox"
        name="distinct-checkbox"
        checked={messageIdFilterValue?.distinct}
        onCheckedChange={(value) => {
          table
            .getColumn("messageId")
            ?.setFilterValue((x: MessageIdFilterValue) => ({
              ...x,
              distinct: !!value,
            }));
        }}
      />
      <Label htmlFor="distinct-checkbox">Distinct Messages</Label>
    </Field>
  );
}

function DestinationFilter({
  table,
  data,
}: {
  table: Table<Message>;
  data: Message[];
}) {
  const destinations = useMemo(() => {
    return [...new Set(data.map((message) => message.destination))];
  }, [data]);

  const destinationFilterValue =
    (table.getColumn("destination")?.getFilterValue() as string[]) ?? [];

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" className="ml-auto">
          Destination
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        {destinations.map((destination) => (
          <DropdownMenuCheckboxItem
            key={destination}
            checked={destinationFilterValue.includes(destination)}
            onCheckedChange={(value) => {
              if (value) {
                table
                  .getColumn("destination")
                  ?.setFilterValue([...destinationFilterValue, destination]);
              } else {
                table
                  .getColumn("destination")
                  ?.setFilterValue(
                    destinationFilterValue.filter((d) => d !== destination),
                  );
              }
            }}
          >
            {destination}
          </DropdownMenuCheckboxItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

function RoleFilter({ table }: { table: Table<Message> }) {
  const roleFilterValue =
    (table.getColumn("role")?.getFilterValue() as string[]) ?? [];

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" className="ml-auto">
          Role
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        {["sender", "receiver"].map((role) => (
          <DropdownMenuCheckboxItem
            key={role}
            checked={roleFilterValue?.includes(role)}
            onCheckedChange={(value) => {
              if (value) {
                table
                  .getColumn("role")
                  ?.setFilterValue([...roleFilterValue, role]);
              } else {
                table
                  .getColumn("role")
                  ?.setFilterValue(roleFilterValue.filter((d) => d !== role));
              }
            }}
          >
            {role}
          </DropdownMenuCheckboxItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

function SubjectFilter({
  table,
  data,
}: {
  table: Table<Message>;
  data: Message[];
}) {
  const subjects = useMemo(() => {
    return [...new Set(data.map((message) => message.subject))];
  }, [data]);

  const subjectFilterValue =
    (table.getColumn("subject")?.getFilterValue() as string[]) ?? [];

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" className="ml-auto">
          Subject
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        {subjects.map((subject) => (
          <DropdownMenuCheckboxItem
            key={subject}
            checked={subjectFilterValue.includes(subject)}
            onCheckedChange={(value) => {
              if (value) {
                table
                  .getColumn("subject")
                  ?.setFilterValue([...subjectFilterValue, subject]);
              } else {
                table
                  .getColumn("subject")
                  ?.setFilterValue(
                    subjectFilterValue.filter((d) => d !== subject),
                  );
              }
            }}
          >
            {subject}
          </DropdownMenuCheckboxItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

function ApplicationDataFilter({ table }: { table: Table<Message> }) {
  const applicationDataFilterValue =
    (table.getColumn("applicationDataPreview")?.getFilterValue() as string) ??
    "";

  return (
    <Field orientation="horizontal">
      <Label htmlFor="applicationDataFilter">Application Data</Label>
      <Input
        type="text"
        id="applicationDataFilter"
        value={applicationDataFilterValue}
        onChange={(e) =>
          table
            .getColumn("applicationDataPreview")
            ?.setFilterValue(e.target.value)
        }
      />
    </Field>
  );
}

export function TableHeader({
  table,
  data,
}: {
  table: Table<Message>;
  data: Message[];
}) {
  return (
    <div className="flex items-center gap-2">
      <MessageIdFilter table={table} data={data} />
      <DestinationFilter table={table} data={data} />
      <RoleFilter table={table} />
      <SubjectFilter table={table} data={data} />
      <ApplicationDataFilter table={table} />
    </div>
  );
}
