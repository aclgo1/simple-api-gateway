import { create } from "zustand";
import { persist } from "zustand/middleware";
import { jwtDecode } from "jwt-decode";

export const useAuthStore = create(
  persist(
    (set) => ({
      accessToken: null,
      refreshToken: null,
      userId: null,
      isInitializing: true, // 1. Flag para controlar o carregamento inicial do storage

      setTokens: (accessToken, refreshToken) => {
        let userId = null;
        if (accessToken) {
          try {
            const decoded = jwtDecode(accessToken);
            userId = decoded.id;
          } catch (error) {
            console.error("error parser JWT:", error);
          }
        }
        set({ accessToken, refreshToken, userId, isInitializing: false });
      },

      logout: () =>
        set({
          accessToken: null,
          refreshToken: null,
          userId: null,
          isInitializing: false,
        }),
    }),
    {
      name: "auth-storage",
      // 2. Executa quando o Zustand termina de ler os dados do localStorage
      onRehydrateStorage: () => (state) => {
        if (state) {
          state.isInitializing = false;
        }
      },
    },
  ),
);
