import {
  createContext,
  useCallback,
  useContext,
  useMemo,
  useState,
  type ReactNode,
} from "react";
import {
  loadFavorites,
  saveFavorites,
  type FavoriteMessage,
} from "./favorites";

type FavoritesContextValue = {
  favorites: FavoriteMessage[];
  addFavorite: (favorite: FavoriteMessage) => void;
  removeFavorite: (uiId: string) => void;
};

const FavoritesContext = createContext<FavoritesContextValue | null>(null);

export function FavoritesProvider({ children }: { children: ReactNode }) {
  const [favorites, setFavorites] = useState<FavoriteMessage[]>(loadFavorites);

  const addFavorite = useCallback((favorite: FavoriteMessage) => {
    setFavorites((current) => {
      const next = [...current, favorite];
      saveFavorites(next);
      return next;
    });
  }, []);

  const removeFavorite = useCallback((uiId: string) => {
    setFavorites((current) => {
      const next = current.filter((favorite) => favorite.uiId !== uiId);
      saveFavorites(next);
      return next;
    });
  }, []);

  const value = useMemo(
    () => ({ favorites, addFavorite, removeFavorite }),
    [favorites, addFavorite, removeFavorite],
  );

  return (
    <FavoritesContext.Provider value={value}>
      {children}
    </FavoritesContext.Provider>
  );
}

export function useFavorites() {
  const context = useContext(FavoritesContext);
  if (!context) {
    throw new Error("useFavorites must be used within a FavoritesProvider");
  }
  return context;
}
