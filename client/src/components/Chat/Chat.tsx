import { Flex, Stack, useColorModeValue } from "@chakra-ui/react";
import { useState, useEffect, useCallback } from "react";

import useWebSocket, { ReadyState } from "react-use-websocket";
import { AuthService } from "../Auth/service";

export const Chat = () => {
  const [socketUrl, setSocketUrl] = useState("ws://localhost:8080/channel");
  const [messageHistory, setMessageHistory] = useState<Array<MessageEvent>>([]);

  const { sendMessage, lastMessage, readyState } = useWebSocket(socketUrl, {
    queryParams: {
      token: AuthService.token || "",
    },
    shouldReconnect: (closeEvent) => {
      // TODO: impl proper logic
      return true;
    },
  });

  useEffect(() => {
    if (lastMessage !== null) {
      setMessageHistory((prev) => prev.concat(lastMessage));
    }
  }, [lastMessage, setMessageHistory]);

  const handleClickSendMessage = useCallback(() => sendMessage("Hello"), []);

  const connectionStatus = {
    [ReadyState.CONNECTING]: "Connecting",
    [ReadyState.OPEN]: "Open",
    [ReadyState.CLOSING]: "Closing",
    [ReadyState.CLOSED]: "Closed",
    [ReadyState.UNINSTANTIATED]: "Uninstantiated",
  }[readyState];

  return (
    <Flex
      minH={"calc(100vh - 64px)"}
      align={"center"}
      justify={"center"}
      bg={useColorModeValue("gray.50", "gray.800")}
    >
      <Stack spacing={8} mx={"auto"} maxW={"lg"} py={12} px={6}>
        <Stack align={"center"}>
          <button
            onClick={handleClickSendMessage}
            disabled={readyState !== ReadyState.OPEN}
          >
            Click Me to send 'Hello'
          </button>
          <span>The WebSocket is currently {connectionStatus}</span>
          {lastMessage ? <span>Last message: {lastMessage.data}</span> : null}
          <ul>
            {messageHistory.map((message, idx) => (
              <span key={idx}>{message ? message.data : null}</span>
            ))}
          </ul>
        </Stack>
      </Stack>
    </Flex>
  );
};
