import { Window } from "@wailsio/runtime";

export function SetWIndowFullScreen() {
  Window.SetTitle("Blizzflow | Dashboard");
  Window.Fullscreen();
  Window.SetAlwaysOnTop(true);
  Window.SetPosition(0, 0);
}

export function SetWIndowUnFullScreen() {
  Window.UnFullscreen();
  Window.SetAlwaysOnTop(false);
  Window.Center();
}

export function setWindowSignUpWindow(user: boolean) {
  Window.SetTitle(user ? "Blizzflow | Sign In" : "Blizzflow | Sign Up");
  SetWIndowUnFullScreen();
  Window.SetSize(user ? 400 : 800, 600);
  Window.SetResizable(false);
}

export function setWindowPurchaseWindow() {
  Window.SetTitle("Blizzflow | Purchase License");
  SetWIndowUnFullScreen();
  Window.SetSize(640, 480);
  Window.SetResizable(false);
}
