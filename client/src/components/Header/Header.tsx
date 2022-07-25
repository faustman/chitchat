import {
  Box,
  Flex,
  Heading,
  Stack,
  useColorModeValue,
  Text,
  Icon,
  Menu,
  MenuButton,
  Button,
  Avatar,
  Center,
  MenuDivider,
  MenuList,
  MenuItem,
} from "@chakra-ui/react";
import { ColorModeSwitcher } from "./ColorModeSwitcher";
import { HiChatAlt2 } from "react-icons/hi";
import { AuthType, UserType } from "../Auth/service";
import { MouseEventHandler } from "react";

type HeaderProps = {
  auth: AuthType | null;
  logout: MouseEventHandler;
};

type UserMenuProps = {
  user: UserType;
  logout: MouseEventHandler;
};

const UserMenu = (props: UserMenuProps) => (
  <Menu>
    <MenuButton
      as={Button}
      rounded={"full"}
      variant={"link"}
      cursor={"pointer"}
      minW={0}
    >
      <Avatar size={"sm"} src={props.user.avatar} />
    </MenuButton>
    <MenuList alignItems={"center"}>
      <br />
      <Center>
        <Avatar size={"2xl"} src={props.user.avatar} />
      </Center>
      <br />
      <Center>
        <p>{props.user.name}</p>
      </Center>
      <br />
      <MenuDivider />
      <MenuItem as="button" onClick={props.logout}>
        Logout
      </MenuItem>
    </MenuList>
  </Menu>
);

export const Header = (props: HeaderProps) => (
  <>
    <Box bg={useColorModeValue("gray.100", "gray.900")} px={4}>
      <Flex h={16} alignItems={"center"} justifyContent={"space-between"}>
        <Box>
          <Heading size="lg">
            <Icon as={HiChatAlt2} w={8} h={8} />
            CHITCHAT
          </Heading>
        </Box>
        {props.auth ? (
          <Box>
            <Text fontSize={"lg"} color={"gray.600"}>
              #{props.auth?.channel}
            </Text>
          </Box>
        ) : null}
        <Flex alignItems={"center"}>
          {props.auth ? (
            <UserMenu user={props.auth.user} logout={props.logout} />
          ) : null}

          <Stack direction={"row"} spacing={7}>
            <ColorModeSwitcher justifySelf="flex-end" />
          </Stack>
        </Flex>
      </Flex>
    </Box>
  </>
);
