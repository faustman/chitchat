import * as React from "react"
import { ChakraProvider, theme, Text } from "@chakra-ui/react";

import { Header } from "./components/Header/Header";
import { Auth } from "./components/Auth/Auth";
import { Progress } from "./components/Progress/Progress";
import { AuthService, AuthType } from "./components/Auth/service";

export const App = () => {
  const [auth, setAuth] = React.useState<AuthType | null>(null);
  const [progress, setProgress] = React.useState<string | null>("Iinitialze..");

  // Check Auth first
  React.useEffect(() => {
    AuthService.auth()
      .then(setAuth)
      .then(() => setProgress(null));
  }, []);

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
        <Text>Logged in!!</Text>
      )}
    </ChakraProvider>
  );
};
