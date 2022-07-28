import {
  Flex,
  Heading,
  Stack,
  useColorModeValue,
  Text,
  Box,
  FormControl,
  FormLabel,
  Input,
  InputGroup,
  Button,
  Link,
  InputLeftAddon,
  FormHelperText,
} from "@chakra-ui/react";
import React, { useRef, useState } from "react";
import { AuthService, AuthType } from "./service";

type AuthProps = {
  setAuth: React.Dispatch<React.SetStateAction<AuthType | null>>;
};

export function Auth(props: AuthProps) {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<Error | null>();

  const nameRef = useRef<HTMLInputElement>(null);
  const emailRef = useRef<HTMLInputElement>(null);
  const channelRef = useRef<HTMLInputElement>(null);

  const handleSubmit = (event: React.FormEvent<HTMLInputElement>): void => {
    event.preventDefault();

    setIsLoading(true);

    if (error) {
      setError(null);
    }

    AuthService.login({
      name: nameRef.current?.value || "",
      email: emailRef.current?.value,
      channel: channelRef.current?.value || "lobby",
    })
      .then(props.setAuth)
      .catch((error) => {
        setIsLoading(false);
        setError(error);
      });
  };

  return (
    <Flex
      minH={"calc(100vh - 64px)"}
      align={"center"}
      justify={"center"}
      bg={useColorModeValue("gray.50", "gray.800")}
    >
      <Stack spacing={8} mx={"auto"} maxW={"lg"} py={12} px={6}>
        <Stack align={"center"}>
          <Heading fontSize={"4xl"} textAlign={"center"}>
            Login
          </Heading>
          <Text fontSize={"lg"} color={"gray.600"}>
            to enjoy of our awesome chat ✌️
          </Text>
        </Stack>
        <Box
          rounded={"lg"}
          bg={useColorModeValue("white", "gray.700")}
          boxShadow={"lg"}
          p={8}
        >
          <Stack spacing={4} as="form" onSubmit={handleSubmit}>
            <FormControl id="name" w={"sm"} isRequired>
              <FormLabel>Name</FormLabel>
              <Input type="text" name="name" ref={nameRef} />
            </FormControl>
            <FormControl id="email" size={"lg"}>
              <FormLabel>Email address</FormLabel>
              <Input type="email" name="email" ref={emailRef} />
              <FormHelperText>
                Optional. We only use it for getting your avatar via{" "}
                <Link href="https://gravatar.com" target={"_blank"}>
                  Gravatar
                </Link>
              </FormHelperText>
            </FormControl>
            <FormControl id="channel" isRequired>
              <FormLabel>Channel</FormLabel>
              <InputGroup>
                <InputLeftAddon children="#" />
                <Input
                  type="channel"
                  placeholder="lobby"
                  defaultValue={"lobby"}
                  ref={channelRef}
                />
              </InputGroup>
              <FormHelperText>
                You could create own channel, #lobby is default
              </FormHelperText>
            </FormControl>
            <Stack spacing={10} pt={2}>
              <Button
                type="submit"
                isLoading={isLoading}
                loadingText="Joining.."
                size="lg"
                bg={"blue.400"}
                color={"white"}
                _hover={{
                  bg: "blue.500",
                }}
              >
                Join
              </Button>
              {error ? <Text color="red">Error: {error.message}</Text> : null}
            </Stack>
          </Stack>
        </Box>
      </Stack>
    </Flex>
  );
}
