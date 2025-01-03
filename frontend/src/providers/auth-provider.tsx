"use client";

import { createContext, useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { Window } from "@wailsio/runtime";
import { SessionService } from "@/blizzflow/backend/domain/services/session";
import {
  Login,
  Logout,
  RecoverPassword,
  Register,
  SetSecurityQuestions,
} from "@/blizzflow/backend/domain/services/auth/authservice";

interface User {
  ID: number;
  Username: string;
}

interface Session {
  ID: number;
  UserID: number;
  CreatedAt: string;
}

interface AuthState {
  isAuthenticated: boolean;
  user: User | null;
  session: Session | null;
}

interface AuthContextType extends AuthState {
  login: (username: string, password: string) => Promise<void>;
  register: (username: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  setSecurityQuestions: (
    username: string,
    questions: Record<string, string>
  ) => Promise<void>;
  recoverPassword: (
    username: string,
    answers: Record<string, string>,
    newPassword: string
  ) => Promise<void>;
}

export const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const navigate = useNavigate();
  const { pathname } = useLocation();
  const [authState, setAuthState] = useState<AuthState>({
    isAuthenticated: false,
    user: null,
    session: null,
  });

  useEffect(() => {
    const checkAuth = async () => {
      try {
        const sessionStr = localStorage.getItem("session");
        if (sessionStr) {
          const session = JSON.parse(sessionStr);
          const isValid = await SessionService.ValidateSession(session.ID);
          if (isValid) {
            setAuthState({
              isAuthenticated: true,
              user: session.user,
              session: session,
            });
          } else {
            localStorage.removeItem("session");
            if (pathname !== "/sign-up" && pathname !== "/sign-in") {
              Window.Center();
              window.resizeTo(400, 300);
              navigate("/sign-up", {
                viewTransition: true,
              });
            }
          }
        } else if (pathname !== "/sign-up" && pathname !== "/sign-in") {
          Window.Center();
          window.resizeTo(400, 300);
          navigate("/sign-up", {
            viewTransition: true,
          });
        }
      } catch (error) {
        console.error("Auth check failed:", error);
        localStorage.removeItem("session");
        if (pathname !== "/sign-up" && pathname !== "/sign-in") {
          Window.Center();
          window.resizeTo(400, 300);
          navigate("/sign-up", {
            viewTransition: true,
          });
        }
      }
    };
    checkAuth();
  }, [pathname, navigate]);

  const login = async (username: string, password: string) => {
    try {
      const session = await Login(username, password);
      if (session) {
        const user = { ID: session.UserID, Username: username };
        localStorage.setItem("session", JSON.stringify({ ...session, user }));
        setAuthState({
          isAuthenticated: true,
          user,
          session,
        });
      } else {
        throw new Error("Invalid credentials");
      }
    } catch (error) {
      console.error("Login failed:", error);
      throw error;
    }
  };

  const register = async (username: string, password: string) => {
    try {
      await Register(username, password);
      await login(username, password);
    } catch (error) {
      console.error("Registration failed:", error);
      throw error;
    }
  };

  const logout = async () => {
    try {
      const session = authState.session;
      if (session) {
        await Logout(session.ID);
      }
    } finally {
      localStorage.removeItem("session");
      setAuthState({
        isAuthenticated: false,
        user: null,
        session: null,
      });
    }
  };

  const setSecurityQuestions = async (
    username: string,
    questions: Record<string, string>
  ) => {
    await SetSecurityQuestions(username, questions);
  };

  const recoverPassword = async (
    username: string,
    answers: Record<string, string>,
    newPassword: string
  ) => {
    await RecoverPassword(username, answers, newPassword);
  };

  const value = {
    ...authState,
    login,
    register,
    logout,
    setSecurityQuestions,
    recoverPassword,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
