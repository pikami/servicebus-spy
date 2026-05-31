import { ModeToggle } from "./mode-toggle";

export function SiteHeader() {
  return (
    <header className="bg-background sticky top-0 z-50 w-full">
      <div className="container-wrapper flex items-center justify-between px-6 py-2">
        <h1 className="text-2xl font-bold">ServiceBus Spy</h1>
        <ModeToggle />
      </div>
    </header>
  );
}
