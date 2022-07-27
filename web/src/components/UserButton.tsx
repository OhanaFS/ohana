import React, { forwardRef } from 'react';
import { UnstyledButton, Group, Avatar, Text } from '@mantine/core';
import { IconChevronRight, IconLogout } from '@tabler/icons';

interface UserButtonProps extends React.ComponentPropsWithoutRef<'button'> {
  image?: string;
  name: string;
  email: string;
  icon?: React.ReactNode;
}

export const UserButton = forwardRef<HTMLButtonElement, UserButtonProps>(
  ({ image, name, email, icon, ...others }: UserButtonProps, ref) => (
    <UnstyledButton
      ref={ref}
      sx={(theme) => ({
        display: 'block',
        width: '100%',
        padding: theme.spacing.md,
        color:
          theme.colorScheme === 'dark' ? theme.colors.dark[0] : theme.black,

        '&:hover': {
          backgroundColor:
            theme.colorScheme === 'dark'
              ? theme.colors.dark[8]
              : theme.colors.gray[0],
        },
      })}
      {...others}
    >
      <Group>
        <Avatar src={image} radius="xl" />

        <div style={{ flex: 1 }}>
          <Text size="sm" weight={500}>
            {name}
          </Text>

          <Text color="dimmed" size="xs">
            {email}
          </Text>
        </div>
        <a href="/auth/logout">
          {icon || <IconLogout size={16} color={'#bf4042'} />}
        </a>
      </Group>
    </UnstyledButton>
  )
);
