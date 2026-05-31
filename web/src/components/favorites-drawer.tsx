import { Star, Trash2 } from "lucide-react";
import type { FavoriteMessage } from "@/lib/favorites";
import { useFavorites } from "@/lib/favorites-context";
import { Button } from "./ui/button";
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerDescription,
  DrawerHeader,
  DrawerTitle,
} from "./ui/drawer";

interface FavoritesDrawerProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSendFavorite: (favorite: FavoriteMessage) => void;
}

export function FavoritesDrawer({
  open,
  onOpenChange,
  onSendFavorite,
}: FavoritesDrawerProps) {
  const { favorites, removeFavorite } = useFavorites();

  const handleSend = (favorite: FavoriteMessage) => {
    onOpenChange(false);
    onSendFavorite(favorite);
  };

  return (
    <Drawer open={open} onOpenChange={onOpenChange} direction="right">
      <DrawerContent className="data-[vaul-drawer-direction=right]:sm:max-w-md">
        <DrawerHeader>
          <DrawerTitle className="flex items-center gap-2">
            <Star className="size-4" />
            Favorites
          </DrawerTitle>
          <DrawerDescription>
            Saved messages you can quickly send to a queue or topic.
          </DrawerDescription>
        </DrawerHeader>

        <div className="flex-1 overflow-y-auto px-4 pb-4">
          {favorites.length === 0 ? (
            <p className="text-muted-foreground py-8 text-center text-sm">
              No favorites yet. Save a message from the table or send dialog.
            </p>
          ) : (
            <ul className="flex flex-col gap-3">
              {favorites.map((favorite) => (
                <li
                  key={favorite.uiId}
                  className="border-border flex flex-col gap-2 rounded-lg border p-3"
                >
                  <div className="min-w-0">
                    <p className="truncate font-medium">
                      {favorite.subject || "No subject"}
                    </p>
                    <p className="text-muted-foreground truncate text-sm">
                      {favorite.queueOrTopic}
                    </p>
                    <p className="text-muted-foreground mt-1 line-clamp-2 text-xs">
                      {favorite.body}
                    </p>
                  </div>
                  <div className="flex gap-2">
                    <Button
                      size="sm"
                      className="flex-1"
                      onClick={() => handleSend(favorite)}
                    >
                      Send
                    </Button>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => removeFavorite(favorite.uiId)}
                      aria-label="Remove favorite"
                    >
                      <Trash2 className="size-4" />
                    </Button>
                  </div>
                </li>
              ))}
            </ul>
          )}
        </div>

        <div className="border-border border-t p-4">
          <DrawerClose asChild>
            <Button variant="outline" className="w-full">
              Close
            </Button>
          </DrawerClose>
        </div>
      </DrawerContent>
    </Drawer>
  );
}
