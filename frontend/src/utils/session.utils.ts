import { Session, User } from "@/blizzflow/backend/domain/model";

export const SessionUtils = {
  saveSession(session: Session, user: Partial<User>): void {
    localStorage.setItem("session", JSON.stringify({ session, user }));
  },

  getSession(): { session: Session; user: User } | null {
    const sessionStr = localStorage.getItem("session");
    console.clear();
    console.log(sessionStr);
    return sessionStr ? JSON.parse(sessionStr) : null;
  },

  clearSession(): void {
    localStorage.removeItem("session");
  },
};
