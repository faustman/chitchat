import * as React from "react"
import { ChakraProvider, theme } from "@chakra-ui/react";

import { Header } from "./components/Header/Header";
import { Auth } from "./components/Auth/Auth";
import { Progress } from "./components/Progress/Progress";
import { AuthService, AuthType } from "./components/Auth/service";
import { Chat } from "./components/Chat/Chat";

export const App = () => {
  const [auth, setAuth] = React.useState<AuthType | null>(null);
  const [progress, setProgress] = React.useState<string | null>("Initialze..");
  const [initMessages, setInitMessages] = React.useState([]);

  // Check Auth first
  React.useEffect(() => {
    AuthService.auth().then(setAuth);
  }, []);

  // Load MSGs
  React.useEffect(() => {
    if (auth) {
      setProgress("Loading messages..");

      fetch("/messages?token=" + AuthService.token)
        .then((r) => r.json())
        .then((data) => setInitMessages(data.messages || []))
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
        <Chat initMsg={initMessages} />
      )}
    </ChakraProvider>
  );
};
