"use client";

import {
  createContext,
  useCallback,
  useEffect,
  useMemo,
  useState,
} from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { SessionService } from "@/blizzflow/backend/domain/services/session";
import {
  Login,
  Logout,
  RecoverPassword,
  Register,
  SetSecurityQuestions,
} from "@/blizzflow/backend/domain/services/auth/authservice";
import { LicenseService } from "@/blizzflow/backend/domain/services/license";
import { LicenseHandler } from "@/blizzflow/backend/domain/handlers/license";
import { UserService } from "@/blizzflow/backend/domain/services/user";
import { SessionUtils } from "@/utils/session.utils";
import { User, Session } from "@/blizzflow/backend/domain/model";
import {
  SetWIndowFullScreen,
  setWindowPurchaseWindow,
  setWindowSignUpWindow,
} from "@/lib/utils";
interface AuthState {
  isAuthenticated: boolean;
  user: Partial<User> | null;
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
  setLicenseStatus: (status: boolean) => void;
  checkSession: () => Promise<boolean>;
}

export const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const navigate = useNavigate();
  const [license, setLicense] = useState<boolean>(false);
  const { pathname } = useLocation();
  const [authState, setAuthState] = useState<AuthState>({
    isAuthenticated: false,
    user: null,
    session: null,
  });

  const checkSession = useCallback(async () => {
    const savedSession = SessionUtils.getSession();
    if (!savedSession) return false;

    const isValid = await SessionService.ValidateSession(
      savedSession.session.ID
    );
    if (isValid) {
      setAuthState({
        isAuthenticated: true,
        user: savedSession.user,
        session: savedSession.session,
      });
      return true;
    }

    SessionUtils.clearSession();
    return false;
  }, []);

  // Check license status
  useEffect(() => {
    const validateLicense = async () => {
      try {
        const licenseData = await LicenseHandler.ReadLicense();
        const isValid = await LicenseService.ValidateLicense(licenseData);
        setLicense(isValid);

        if (!isValid) {
          setWindowPurchaseWindow();
          navigate("/purchase", { viewTransition: true });
        }
      } catch (error) {
        console.error("License validation failed:", error);
        setLicense(false);
        setWindowPurchaseWindow();
        navigate("/purchase", { viewTransition: true });
      }
    };

    if (!license) validateLicense();
  }, [license, navigate]);

  // Authentication check
  useEffect(() => {
    const validateAuth = async () => {
      try {
        const isSessionValid = await checkSession();

        // Redirect to dashboard if everything is valid
        if (isSessionValid && license) {
          SetWIndowFullScreen();
          navigate("/", { viewTransition: true });
          return;
        }

        // Existing authentication logic
        if (
          !isSessionValid &&
          pathname !== "/sign-up" &&
          pathname !== "/sign-in"
        ) {
          const user = await UserService.GetUserByUsername(
            localStorage.getItem("username") || ""
          );
          setWindowSignUpWindow(user ? true : false);
          navigate(user ? "/sign-in" : "/sign-up", { viewTransition: true });
        }
      } catch (error) {
        console.error("Authentication validation failed:", error);
        navigate("/sign-up", { viewTransition: true });
      }
    };

    if (license) validateAuth();
  }, [license, pathname, navigate, checkSession]);

  const login = useCallback(
    async (username: string, password: string) => {
      try {
        const session = await Login(username, password);
        if (!session) throw new Error("Invalid credentials");

        const user = { ID: session.UserID, Username: username };
        SessionUtils.saveSession(session, user);

        setAuthState({
          isAuthenticated: true,
          user,
          session,
        });
      } catch (error) {
        console.error("Login failed:", error);
        throw error;
      }
    },
    [setAuthState]
  );

  const register = useCallback(
    async (username: string, password: string) => {
      try {
        await Register(username, password);
        localStorage.setItem("username", username); // Save username for redirect
        await login(username, password);
      } catch (error) {
        console.error("Registration failed:", error);
        throw error;
      }
    },
    [login]
  );

  const logout = useCallback(async () => {
    try {
      if (authState.session) {
        await new Promise((resolve) =>
          setTimeout(() => {
            resolve("");
          }, 2000)
        );
        await Logout(authState.session.ID);
        navigate("/callback", { viewTransition: true });
      }
    } finally {
      SessionUtils.clearSession();
      setAuthState({
        isAuthenticated: false,
        user: null,
        session: null,
      });
    }
  }, [authState.session, setAuthState]);

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

  const contextValue = useMemo(
    () => ({
      ...authState,
      login,
      register,
      logout,
      setSecurityQuestions,
      recoverPassword,
      checkSession,
      setLicenseStatus: (status: boolean) => {
        setLicense(status);
        if (!status) navigate("/purchase", { viewTransition: true });
      },
    }),
    [authState, login, register, logout, navigate]
  );

  return (
    <AuthContext.Provider value={contextValue}>{children}</AuthContext.Provider>
  );
}
