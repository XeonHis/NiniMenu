import { create } from "zustand"

interface AuthState {
  token: string | null
  isLoggedIn: boolean
  login: (token: string) => void
  logout: () => void
}

export const useAuthStore = create<AuthState>((set) => ({
  token: localStorage.getItem("token"),
  isLoggedIn: !!localStorage.getItem("token"),
  login: (token: string) => {
    localStorage.setItem("token", token)
    set({ token, isLoggedIn: true })
  },
  logout: () => {
    localStorage.removeItem("token")
    set({ token: null, isLoggedIn: false })
  },
}))
