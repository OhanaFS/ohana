import React from 'react';
import {
  UnstyledButton,
  Group,
  Avatar,
  Text,
  Box,
  useMantineTheme,
} from '@mantine/core';

interface FullName {
  image: string;
  name: string;
  username: string;
}

export function User(props: FullName) {
  const theme = useMantineTheme();

  return (
    <Box
      sx={{
        paddingTop: theme.spacing.sm,
        borderTop: `1px solid ${
          theme.colorScheme === 'dark'
            ? theme.colors.dark[4]
            : theme.colors.gray[2]
        }`,
      }}
    >
      <UnstyledButton
        sx={{
          display: 'block',
          width: '100%',
          padding: theme.spacing.xs,
          borderRadius: theme.radius.sm,
          color:
            theme.colorScheme === 'dark' ? theme.colors.dark[0] : theme.black,

          '&:hover': {
            backgroundColor:
              theme.colorScheme === 'dark'
                ? theme.colors.dark[6]
                : theme.colors.gray[0],
          },
        }}
      >
        <Group>
          <Avatar src={props.image} radius="xl" />
          <Box sx={{ flex: 1 }}>
            <Text size="sm" weight={500}>
              {props.name}
            </Text>
            <Text color="dimmed" size="xs">
              {props.username}
            </Text>
          </Box>

          {/* {theme.dir === 'ltr' ? <ChevronRight size={18} /> : <ChevronLeft size={18} />} */}
        </Group>
      </UnstyledButton>
    </Box>
  );
}
