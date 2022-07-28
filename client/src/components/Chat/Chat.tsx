import {
  Box,
  createIcon,
  Flex,
  HStack,
  IconButton,
  Input,
  useColorModeValue,
  VStack,
  Text,
  Avatar,
  useToast,
  UseToastOptions,
  Heading,
  StackDivider,
  Icon,
  AvatarBadge,
  Badge,
} from "@chakra-ui/react";
import { useState, useEffect, useCallback, useRef } from "react";
import { HiUsers } from "react-icons/hi";

import useWebSocket, { ReadyState } from "react-use-websocket";
import { AuthService, AuthType, UserType } from "../Auth/service";

export const SendIcon = createIcon({
  displayName: "SendIcon",
  viewBox: "0 0 20 20",
  // path can also be an array of elements, if you have multiple paths, lines, shapes, etc.
  path: (
    <path d="m 17.448976,10.895396 a 1,1 0 0 0 0,-1.7879996 l -13.9999998,-7 a 1,1 0 0 0 -1.409,1.169 l 1.429,5 a 1,1 0 0 0 0.962,0.725 h 4.571 a 1,1 0 1 1 0,1.9999996 h -4.571 a 1,1 0 0 0 -0.962,0.725 l -1.428,5 a 1,1 0 0 0 1.408,1.17 l 13.9999998,-7 z" />
  ),
});

export type ChannelMessage = {
  type: string;
  from_user: UserType;
  text: string;
  sent_at: number;
};

const connectionStatus = {
  [ReadyState.CONNECTING]: "Connecting",
  [ReadyState.OPEN]: "Connected",
  [ReadyState.CLOSING]: "Closing",
  [ReadyState.CLOSED]: "Closed",
  [ReadyState.UNINSTANTIATED]: "Uninstantiated",
};

type ChatProps = {
  initMsg?: Array<ChannelMessage>;
  initUsers?: Array<UserType>;
  auth: AuthType;
};

export const Chat = (props: ChatProps) => {
  const [messageHistory, setMessageHistory] = useState<Array<ChannelMessage>>(
    props.initMsg || []
  );

  const initUsers = (props.initUsers || [])
    .concat(props.auth.user) // Current user always online
    .reduce((users, user) => ({ ...users, [user.id]: user }), {});
  const [usersOnline, setUsersOnline] = useState<{
    [key: string]: UserType;
  }>(initUsers);

  const messageTextRef = useRef<HTMLInputElement>(null);
  const bottomAnchorRef = useRef<HTMLInputElement>(null);

  const { sendMessage, lastJsonMessage, readyState } = useWebSocket(
    "ws://localhost:8080/channel",
    {
      queryParams: {
        token: AuthService.token || "",
      },
      shouldReconnect: (closeEvent) => {
        // TODO: impl proper logic
        console.log(closeEvent);

        return true;
      },
    }
  );

  // Listen to new ws messages and apply it to state
  useEffect(() => {
    if (lastJsonMessage !== null) {
      const msg = lastJsonMessage as ChannelMessage;

      if (msg.type === "message") {
        setMessageHistory((prev) => prev.concat(msg));
      }

      if (msg.type === "join") {
        const user = msg.from_user;
        setUsersOnline((prev) => ({ ...prev, [user.id]: user }));
      }

      if (msg.type === "leave") {
        setUsersOnline((prev) => {
          delete prev[msg.from_user.id];

          return prev;
        });
      }
    }
  }, [lastJsonMessage, setMessageHistory]);

  // After new message scroll to the buttom
  useEffect(() => {
    if (bottomAnchorRef.current) {
      bottomAnchorRef.current.scrollIntoView({ behavior: "smooth" });
    }
  }, [messageHistory]);

  const handleClickSendMessage = useCallback(
    (event: React.FormEvent<HTMLInputElement>) => {
      event.preventDefault();

      const text = messageTextRef.current?.value;

      if (!text) {
        return;
      }

      sendMessage(text);
      messageTextRef.current.value = "";
    },
    []
  );

  // WS Connection Status toast
  const toast = useToast();

  useEffect(() => {
    if (!toast.isActive(readyState)) {
      if (readyState === ReadyState.OPEN) {
        toast.closeAll();
      }

      const toastOptions = {
        id: readyState,
        title: `${connectionStatus[readyState]}`,
        position: "top",
        status: "success",
        isClosable: true,
        duration: 2000,
      } as UseToastOptions;

      if (readyState === ReadyState.CLOSED) {
        toastOptions.status = "error";
        toastOptions.duration = null;
        toastOptions.isClosable = false;
      }

      if (
        readyState === ReadyState.CONNECTING ||
        readyState === ReadyState.CLOSING
      ) {
        toastOptions.status = "loading";
      }

      toast(toastOptions);
    }
  }, [readyState, toast]);

  const msgBoxColor = useColorModeValue("gray.400", "gray.500");

  return (
    <Flex
      minH={"calc(100vh - 64px)"}
      position={"relative"}
      height={"calc(100vh - 64px)"}
      bg={useColorModeValue("gray.50", "gray.800")}
    >
      <Flex
        bg={useColorModeValue("gray.200", "gray.700")}
        rounded={"lg"}
        flexGrow={"1"}
        flexDirection={"column"}
        m={2}
      >
        <Flex
          flexGrow={"1"}
          flexDirection={"column"}
          overflowY={"scroll"}
          shrink={"1"}
          bg={useColorModeValue("gray.300", "gray.600")}
          m={2}
          rounded={"lg"}
        >
          <Flex flexShrink={"1"} flexGrow={"1"} />
          {messageHistory.map((message, idx) => (
            <Box key={idx} m={2}>
              <HStack align={"flex-start"}>
                <Box mt={"4px"}>
                  <Avatar
                    size="sm"
                    name={message.from_user.name}
                    src={message.from_user.avatar}
                  />
                </Box>
                <VStack align={"flex-start"}>
                  <Box w={"100%"}>
                    <Text fontSize="xs">{message.from_user.name}</Text>
                  </Box>
                  <Box bg={msgBoxColor} p={2} rounded={"lg"} m={0}>
                    <Text>{message.text}</Text>
                  </Box>
                </VStack>
                <Box mt={"20px"}>
                  <Text fontSize="xs">
                    {new Date(message.sent_at).toLocaleString()}
                  </Text>
                </Box>
              </HStack>
            </Box>
          ))}
          <Box ref={bottomAnchorRef}></Box>
        </Flex>
        <Flex
          bg={useColorModeValue("gray.200", "gray.700")}
          rounded={"lg"}
          height={"46px"}
          m={2}
          direction={"column"}
        >
          <HStack as={"form"} onSubmit={handleClickSendMessage}>
            <Input placeholder="Message.." ref={messageTextRef} />
            <IconButton
              type={"submit"}
              // colorScheme="teal"
              aria-label="Send"
              // size="lg"
              icon={<SendIcon />}
            />
          </HStack>
        </Flex>
      </Flex>
      <Box
        bg={useColorModeValue("gray.200", "gray.700")}
        w={"20%"}
        m={2}
        rounded={"lg"}
      >
        <VStack
          divider={
            <StackDivider
              borderColor={useColorModeValue("gray.50", "gray.800")}
            />
          }
          spacing={2}
          align="stretch"
        >
          <HStack justifyContent={"center"} textAlign={"center"} mt={"8px"}>
            <Icon as={HiUsers} w={6} h={6} />
            <Heading size={"md"}>Online</Heading>
          </HStack>
          <VStack align={"left"} ml={"8px"}>
            {Object.entries(usersOnline).map(([key, user]) => (
              <HStack key={key}>
                <Avatar size="sm" name={user.name} src={user.avatar}>
                  <AvatarBadge boxSize="1.25em" bg="green.500" />
                </Avatar>
                <Text>{user.name}</Text>
                {user.name === props.auth.user.name ? (
                  <Text as="i">(you)</Text>
                ) : null}
              </HStack>
            ))}
          </VStack>
        </VStack>
      </Box>
    </Flex>
  );
};
