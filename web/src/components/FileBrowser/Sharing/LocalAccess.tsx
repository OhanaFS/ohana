import React from 'react';
import {
  Autocomplete,
  Text,
  Group,
  Avatar,
  SelectItemProps,
  Button,
  Accordion,
  Switch,
} from '@mantine/core';
import { useQueryUser, useQueryUsers, WhoamiResponse } from '../../../api/auth';
import {
  EntryType,
  FileMetadata,
  FilePermission,
  Permission,
  useMutateAddFilePermissions,
  useMutateDeleteFilePermissions,
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
  const qFile = useQueryFileMetadata(perm.file_id);
  const is_folder =
    (qFile.data as FileMetadata | undefined)?.entry_type === EntryType.Folder;
  const isYou = perm.User.user_id === qUser.data?.user_id;
  const [switches, setSwitches] = React.useState<Permission>({
    can_read: perm.can_read,
    can_write: perm.can_write,
    can_execute: perm.can_execute,
    can_share: perm.can_share,
    can_audit: perm.can_audit,
  });
  const [isSaved, setIsSaved] = React.useState(true);

  const mUpdate = useMutateUpdateFilePermissions();
  const mDelete = useMutateDeleteFilePermissions();

  React.useEffect(() => {
    if (!isSaved) {
      mUpdate
        .mutateAsync({
          ...switches,
          file_id: perm.file_id,
          permission_id: perm.permission_id,
          is_folder,
        })
        .then(() => setIsSaved(true));
    }
  }, [switches, isSaved]);

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
              disabled={isYou || mUpdate.isLoading}
              label="Can view"
              checked={switches.can_read}
              onChange={(e) => {
                setSwitches({
                  ...switches,
                  can_read: e.currentTarget.checked,
                });
                setIsSaved(false);
              }}
            />
            <Switch
              disabled={isYou || mUpdate.isLoading}
              label="Can edit"
              checked={switches.can_write}
              onChange={(e) => {
                setSwitches({
                  ...switches,
                  can_write: e.currentTarget.checked,
                });
                setIsSaved(false);
              }}
            />
            <Switch
              disabled={isYou || mUpdate.isLoading}
              label="Can share"
              checked={switches.can_share}
              onChange={(e) => {
                setSwitches({
                  ...switches,
                  can_share: e.currentTarget.checked,
                });
                setIsSaved(false);
              }}
            />
          </Group>
          <Button
            compact
            color="red"
            variant="subtle"
            disabled={isYou}
            onClick={() =>
              mDelete.mutate({
                file_id: perm.file_id,
                permission_id: perm.permission_id,
                is_folder,
              })
            }
          >
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
  const is_folder =
    (qFile.data as FileMetadata | undefined)?.entry_type === EntryType.Folder;

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
      can_audit: false,

      is_folder,
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
