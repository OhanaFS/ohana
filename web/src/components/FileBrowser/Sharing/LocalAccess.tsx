import { forwardRef } from 'react';
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

const UserLayout = forwardRef<HTMLDivElement, UserButtonProps>(
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

const AutocompleteItem = forwardRef<HTMLDivElement, ItemProps>(
  ({ name, email, user_id, ...rest }: ItemProps, ref) => (
    <UserLayout ref={ref} name={name} email={email} {...rest} />
  )
);

const LocalAccess = ({ fileId }: { fileId: string }) => {
  const qUser = useQueryUser();
  const qUsers = useQueryUsers();

  const autocompleteData = (qUsers.data ?? []).map((user) => ({
    ...user,
    value: user.name,
  }));

  const handleShare = (item: WhoamiResponse) => {
    console.log('got new stuff!', { item });
  };

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
      />

      <Text weight="bold" pt="sm">
        People with access
      </Text>
      <UserButton
        name={qUser.data?.name || ''}
        email={qUser.data?.email || ''}
        isYou
      />
    </>
  );
};

export default LocalAccess;
