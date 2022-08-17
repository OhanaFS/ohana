import React from 'react';
import {
  Autocomplete,
  Text,
  UnstyledButton,
  Group,
  Avatar,
  createStyles,
  SelectItemProps,
  Stack,
  Menu,
  Button,
  Accordion,
  Switch,
} from '@mantine/core';
import { useQueryUser, useQueryUsers, WhoamiResponse } from '../../../api/auth';
import {
  FilePermission,
  Permission,
  useMutateAddFilePermissions,
  useMutateUpdateFilePermissions,
  useQueryFileMetadata,
  useQueryFilePermissions,
} from '../../../api/file';

type UserButtonProps = {
  name: string;
  email: string;
  isYou?: boolean;
};

const UserLayout = React.forwardRef<HTMLDivElement, UserButtonProps>(
  ({ name, email, isYou, ...rest }: UserButtonProps, ref) => {
    return (
      <div ref={ref} {...rest}>
        <Group noWrap>
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

const UserButton = ({ perm }: { perm: FilePermission }) => {
  const qUser = useQueryUser();
  const isYou = perm.User.user_id === qUser.data?.user_id;
  const [switches, setSwitches] = React.useState<Permission>({
    can_read: false,
    can_write: false,
    can_execute: false,
    can_share: false,
  });

  const mUpdate = useMutateUpdateFilePermissions();
  React.useEffect(() => {}, [mUpdate.data]);

  return (
    <Accordion.Item value={perm.permission_id}>
      <Accordion.Control>
        <UserLayout
          name={perm.User.name}
          email={perm.User.email}
          isYou={isYou}
        />
      </Accordion.Control>
      <Accordion.Panel>
        <Group position="apart">
          <Group>
            <Switch
              disabled={isYou}
              label="Can view"
              checked={switches.can_read}
              onChange={(e) =>
                setSwitches({
                  ...switches,
                  can_read: e.currentTarget.checked,
                })
              }
            />
            <Switch
              disabled={isYou}
              label="Can edit"
              checked={switches.can_write}
              onChange={(e) =>
                setSwitches({
                  ...switches,
                  can_write: e.currentTarget.checked,
                })
              }
            />
            <Switch
              disabled={isYou}
              label="Can share"
              checked={switches.can_share}
              onChange={(e) =>
                setSwitches({
                  ...switches,
                  can_share: e.currentTarget.checked,
                })
              }
            />
          </Group>
          <Button compact color="red" variant="subtle" disabled={isYou}>
            Revoke access
          </Button>
        </Group>
      </Accordion.Panel>
    </Accordion.Item>
  );
};

const UserButtons = ({ permissions }: { permissions: FilePermission[] }) => {
  return (
    <Accordion variant="filled">
      {permissions.map((perm) => (
        <UserButton key={perm.permission_id} perm={perm} />
      ))}
    </Accordion>
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
      <UserButtons permissions={currentPerms} />
    </>
  );
};

export default LocalAccess;
