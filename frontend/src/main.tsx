import React, { FC, useState } from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import ProtectedRoute from "./components/auth/protected-route";
import { CallbackPage, Home, Login, NotFound, PurchasePage } from "./pages";
import { AuthProvider } from "./providers/auth-provider";
import SignUpPage from "./pages/sign-up";
import "./globals.css";
import { Toaster } from "sonner";

const App: FC = () => {
  const [isAuthenticated] = useState<boolean>(false);

  return (
    <Router>
      <AuthProvider>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/sign-in" element={<Login />} />
          <Route path="/sign-up" element={<SignUpPage />} />
          <Route path="/purchase" element={<PurchasePage />} />
          <Route path="/callback" element={<CallbackPage />} />
          <Route
            path="/protected"
            element={
              <ProtectedRoute isAuthenticated={isAuthenticated}>
                <div>Protected Content</div>
              </ProtectedRoute>
            }
          />
          <Route path="*" element={<NotFound />} />
        </Routes>
        <Toaster richColors />
      </AuthProvider>
    </Router>
  );
};

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
