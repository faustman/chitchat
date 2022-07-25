import {
  Flex,
  Heading,
  Stack,
  useColorModeValue,
  Text,
  Spinner,
} from "@chakra-ui/react";

type ProgressProps = {
  text?: string | null;
};

export const Progress = (props: ProgressProps) => (
  <Flex
    minH={"calc(100vh - 64px)"}
    align={"center"}
    justify={"center"}
    bg={useColorModeValue("gray.50", "gray.800")}
  >
    <Stack spacing={8} mx={"auto"} maxW={"lg"} py={12} px={6}>
      <Stack align={"center"}>
        <Spinner
          thickness="4px"
          speed="0.65s"
          emptyColor="gray.200"
          color="blue.500"
          size="xl"
        />
      </Stack>
      {props.text ? (
        <Text fontSize={"md"} color={"gray.600"}>
          {props.text}
        </Text>
      ) : null}
    </Stack>
  </Flex>
);
