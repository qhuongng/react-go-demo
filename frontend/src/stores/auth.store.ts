import { create } from "zustand";

interface AuthState {
    id: number;
    accessToken: string;

    setAuth: (id: number, token: string) => void;
    removeAuth: () => void;
    setAccessToken: (token: string) => void;
}

const useAuthStore = create<AuthState>((set) => ({
    id: 0,
    accessToken: "",

    setAuth: (id, token) => set({ id: id, accessToken: token }),
    removeAuth: () => set({ id: 0, accessToken: "" }),
    setAccessToken: (token) => set({ accessToken: token }),
}));

export default useAuthStore;
