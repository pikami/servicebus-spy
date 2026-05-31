import { useState } from "react";
import { Star } from "lucide-react";
import { favoriteToFormState, type FavoriteMessage, type SendMessageFormState } from "@/lib/favorites";
import { useFavorites } from "@/lib/favorites-context";
import { ModeToggle } from "./mode-toggle";
import { FavoritesDrawer } from "./favorites-drawer";
import { SendMessageDialog } from "./send-message-dialog";
import { Button } from "./ui/button";

export function SiteHeader() {
  const [sendDialogOpen, setSendDialogOpen] = useState(false);
  const [favoritesDrawerOpen, setFavoritesDrawerOpen] = useState(false);
  const [sendDialogInitialValues, setSendDialogInitialValues] = useState<
    SendMessageFormState | undefined
  >();
  const { favorites } = useFavorites();

  const openSendDialog = (initialValues?: SendMessageFormState) => {
    setSendDialogInitialValues(initialValues);
    setSendDialogOpen(true);
  };

  const handleSendDialogOpenChange = (open: boolean) => {
    setSendDialogOpen(open);
    if (!open) {
      setSendDialogInitialValues(undefined);
    }
  };

  const handleSendFavorite = (favorite: FavoriteMessage) => {
    openSendDialog(favoriteToFormState(favorite));
  };

  return (
    <>
      <header className="bg-background sticky top-0 z-50 w-full">
        <div className="container-wrapper flex items-center justify-between px-6 py-2">
          <h1 className="text-2xl font-bold">ServiceBus Spy</h1>
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              onClick={() => setFavoritesDrawerOpen(true)}
            >
              <Star className="size-4" />
              Favorites
              {favorites.length > 0 ? (
                <span className="bg-primary text-primary-foreground ml-1 rounded-full px-1.5 py-0.5 text-xs">
                  {favorites.length}
                </span>
              ) : null}
            </Button>
            <Button variant="outline" onClick={() => openSendDialog()}>
              Send Message
            </Button>
            <ModeToggle />
          </div>
        </div>
      </header>
      <FavoritesDrawer
        open={favoritesDrawerOpen}
        onOpenChange={setFavoritesDrawerOpen}
        onSendFavorite={handleSendFavorite}
      />
      <SendMessageDialog
        open={sendDialogOpen}
        onOpenChange={handleSendDialogOpenChange}
        initialValues={sendDialogInitialValues}
      />
    </>
  );
}
