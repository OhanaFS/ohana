import React from 'react';
import {
  Autocomplete,
  Text,
  UnstyledButton,
  Group,
  Avatar,
  createStyles,
  SelectItemProps,
} from '@mantine/core';
import { useQueryUser, useQueryUsers, WhoamiResponse } from '../../../api/auth';
import {
  useMutateAddFilePermissions,
  useQueryFileMetadata,
  useQueryFilePermissions,
} from '../../../api/file';

const useStyles = createStyles((theme) => ({
  userButton: {
    padding: theme.spacing.sm,
    transition: 'all 0.1s',
    borderRadius: theme.radius.xl,
    '&:hover': {
      backgroundColor: theme.colors.blue[0],
    },
  },
}));

type UserButtonProps = {
  name: string;
  email: string;
  isYou?: boolean;
  withStyles?: boolean;
};

const UserLayout = React.forwardRef<HTMLDivElement, UserButtonProps>(
  ({ name, email, isYou, withStyles, ...rest }: UserButtonProps, ref) => {
    const { classes } = useStyles();
    return (
      <div ref={ref} {...rest}>
        <Group className={withStyles ? classes.userButton : ''} noWrap>
          <Avatar size={40} color="blue" radius="xl">
            {name
              .split(' ')
              .map((part) => part.substring(0, 1).toUpperCase())
              .join('')}
          </Avatar>
          <div>
            <Text>
              {name}
              {isYou ? ' (You)' : ''}
            </Text>
            <Text size="xs" color="dimmed">
              {email}
            </Text>
          </div>
        </Group>
      </div>
    );
  }
);

const UserButton = (props: UserButtonProps) => {
  return (
    <UnstyledButton>
      <UserLayout {...props} withStyles />
    </UnstyledButton>
  );
};

interface ItemProps extends SelectItemProps {
  name: string;
  email: string;
  user_id: string;
}

const AutocompleteItem = React.forwardRef<HTMLDivElement, ItemProps>(
  ({ name, email, user_id, ...rest }: ItemProps, ref) => (
    <UserLayout ref={ref} name={name} email={email} {...rest} />
  )
);

const LocalAccess = ({ fileId }: { fileId: string }) => {
  const qUser = useQueryUser();
  const qFile = useQueryFileMetadata(fileId);
  const qUsers = useQueryUsers();
  const qPerms = useQueryFilePermissions(fileId);

  const mCreatePerms = useMutateAddFilePermissions();

  const autocompleteData = (qUsers.data ?? []).map((user) => ({
    ...user,
    value: user.name,
  }));

  const handleShare = (item: WhoamiResponse) => {
    mCreatePerms.mutate({
      file_id: fileId,
      users: [item.user_id],
      groups: [],

      can_read: true,
      can_write: false,
      can_execute: false,
      can_share: false,
    });
    console.log('got new stuff!', { item });
  };

  const currentPerms = (qPerms.data ?? []).filter(
    (perm) => perm.version_no === qFile.data?.version_no
  );

  return (
    <>
      <Autocomplete
        transition="pop"
        transitionDuration={100}
        transitionTimingFunction="ease"
        placeholder="Add people and groups"
        data={autocompleteData}
        itemComponent={AutocompleteItem}
        filter={(value, item) => {
          const searchTerm = value.toLowerCase().trim();
          const fullItem = item as unknown as WhoamiResponse;
          return (
            fullItem.name.toLowerCase().includes(searchTerm) ||
            fullItem.email.toLowerCase().includes(searchTerm) ||
            fullItem.user_id.toLowerCase().includes(searchTerm)
          );
        }}
        onItemSubmit={(item) => handleShare(item as unknown as WhoamiResponse)}
        disabled={mCreatePerms.isLoading}
      />

      <Text weight="bold" pt="sm">
        People with access
      </Text>
      {currentPerms.map((perm) => (
        <UserButton
          key={perm.permission_id}
          name={perm.User.name}
          email={perm.User.email}
          isYou={perm.user_id === qUser.data?.user_id}
        />
      ))}
    </>
  );
};

export default LocalAccess;
