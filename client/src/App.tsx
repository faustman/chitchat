import * as React from "react"
import { ChakraProvider, theme } from "@chakra-ui/react";

import { Header } from "./components/Header/Header";
import { Auth } from "./components/Auth/Auth";
import { Progress } from "./components/Progress/Progress";
import { AuthService, AuthType, UserType } from "./components/Auth/service";
import { Chat, ChannelMessage } from "./components/Chat/Chat";
import { ChatService } from "./components/Chat/service";

export const App = () => {
  const [auth, setAuth] = React.useState<AuthType | null>(null);
  const [progress, setProgress] = React.useState<string | null>("Initialze..");
  const [initMessages, setInitMessages] = React.useState<Array<ChannelMessage>>(
    []
  );
  const [initUsers, setInitUsers] = React.useState<Array<UserType>>([]);

  const checkAuth = () => {
    AuthService.auth().then((auth) => {
      setAuth(auth);
      if (!auth) {
        // Show login page
        setProgress(null);
      }
    });
  };

  // Check Auth first
  React.useEffect(() => {
    // Listen for login/logout from other windows
    window.addEventListener(
      "storage",
      (event) => {
        checkAuth();
      },
      { once: true }
    );

    checkAuth();
  }, []);

  // Load MSGs and Users
  React.useEffect(() => {
    if (auth) {
      setProgress("Loading messages..");

      // small feedback of process..
      setTimeout(() => {
        setProgress((prev) => (prev ? "Loading users.." : null));
      }, 500);

      Promise.all([
        ChatService.fetchMessages().then(setInitMessages),
        ChatService.fetchUsers().then(setInitUsers),
      ])
        .then(() => setProgress(null))
        .catch(console.error);
    }
  }, [auth]);

  const logout = () => {
    AuthService.reset();
    setAuth(null);
  };

  return (
    <ChakraProvider theme={theme}>
      <Header auth={auth} logout={logout} />
      {progress ? (
        <Progress text={progress} />
      ) : !auth ? (
        <Auth setAuth={setAuth} />
      ) : (
        <Chat initUsers={initUsers} initMsg={initMessages} auth={auth} />
      )}
    </ChakraProvider>
  );
}
