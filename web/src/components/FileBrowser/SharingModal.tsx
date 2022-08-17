import React from 'react';
import {
  ActionIcon,
  Anchor,
  Autocomplete,
  Avatar,
  Button,
  CopyButton,
  createStyles,
  Group,
  Menu,
  Modal,
  Stack,
  Switch,
  Text,
  Tooltip,
  UnstyledButton,
} from '@mantine/core';
import { IconCheck, IconCopy, IconDots, IconX } from '@tabler/icons';
import { useQueryUser } from '../../api/auth';
import { EntryType, useQueryFileMetadata } from '../../api/file';
import {
  getSharingLinkURL,
  useMutateCreateSharingLink,
  useMutateDeleteSharingLink,
  useQueryFileSharingLinks,
} from '../../api/sharing';

export type SharingModalProps = {
  fileId: string;
  opened: boolean;
  onClose: () => any;
};

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

const SharingModal = (props: SharingModalProps) => {
  const { fileId, opened, onClose } = props;
  const { classes } = useStyles();

  const qUser = useQueryUser();
  const qFile = useQueryFileMetadata(fileId);
  const qSharedLinks = useQueryFileSharingLinks(fileId);
  const sharedLinks = qSharedLinks.data ?? [];

  const [isPublicShared, setIsPublicShared] = React.useState(
    sharedLinks.length > 0
  );

  const mCreateSharedLink = useMutateCreateSharingLink();
  const mDeleteSharedLink = useMutateDeleteSharingLink();

  const removeAllSharedLinks = async () => {
    for (const sharedLink of sharedLinks) {
      const link = sharedLink.shortened_link;
      await mDeleteSharedLink.mutateAsync({ link, fileId });
    }
  };

  React.useEffect(() => {
    const hasPublicLinks = sharedLinks.length > 0;

    if (isPublicShared && !hasPublicLinks) mCreateSharedLink.mutate({ fileId });
    else if (!isPublicShared && hasPublicLinks) removeAllSharedLinks();
    else if (isPublicShared !== hasPublicLinks)
      setIsPublicShared(hasPublicLinks);
  }, [isPublicShared, sharedLinks]);

  return (
    <Modal
      centered
      opened={opened}
      onClose={onClose}
      size="lg"
      title={`Share "${qFile.data?.file_name}"`}
      overflow="outside"
      styles={(theme) => ({ title: { fontSize: theme.fontSizes.xl } })}
    >
      <Stack>
        <Autocomplete
          transition="pop-top-left"
          placeholder="Add people and groups"
          data={['Nobody']}
        />

        <Text weight="bold" pt="sm">
          People with access
        </Text>
        <UnstyledButton>
          <Group className={classes.userButton}>
            <Avatar size={40} color="blue" radius="xl">
              {(qUser.data?.name || '')
                .split(' ')
                .map((part) => part.substring(0, 1).toUpperCase())
                .join('')}
            </Avatar>
            <div>
              <Text>{qUser.data?.name} (You)</Text>
              <Text size="xs" color="dimmed">
                {qUser.data?.email}
              </Text>
            </div>
          </Group>
        </UnstyledButton>

        <Text weight="bold" pt="sm">
          General access
        </Text>

        <Switch
          disabled={qSharedLinks.isLoading || mCreateSharedLink.isLoading}
          label={
            isPublicShared
              ? 'Anyone with the link can view'
              : 'Enable public link'
          }
          onChange={(e) => setIsPublicShared(e.currentTarget.checked)}
          checked={isPublicShared}
        />

        {(qSharedLinks.data ?? []).map((link) => (
          <Group position="apart" key={link.shortened_link}>
            <Anchor
              target="_blank"
              href={getSharingLinkURL(link.shortened_link, 'preview')}
            >
              {link.shortened_link}
            </Anchor>
            <Group>
              <CopyButton
                value={getSharingLinkURL(link.shortened_link, 'preview')}
                timeout={2000}
              >
                {({ copied, copy }) => (
                  <Tooltip
                    label={copied ? 'Copied' : 'Copy'}
                    withArrow
                    position="left"
                  >
                    <ActionIcon color={copied ? 'teal' : 'gray'} onClick={copy}>
                      {copied ? (
                        <IconCheck size={16} />
                      ) : (
                        <IconCopy size={16} />
                      )}
                    </ActionIcon>
                  </Tooltip>
                )}
              </CopyButton>
              <Menu shadow="md" width={200}>
                <Menu.Target>
                  <ActionIcon color="gray">
                    <IconDots size={16} />
                  </ActionIcon>
                </Menu.Target>
                <Menu.Dropdown>
                  <Menu.Item
                    icon={<IconX size={14} />}
                    color="red"
                    onClick={() =>
                      mDeleteSharedLink.mutate({
                        fileId,
                        link: link.shortened_link,
                      })
                    }
                    disabled={mDeleteSharedLink.isLoading}
                  >
                    Unshare
                  </Menu.Item>
                </Menu.Dropdown>
              </Menu>
            </Group>
          </Group>
        ))}

        <Group position="right">
          <Button onClick={onClose}>Done</Button>
        </Group>
      </Stack>
    </Modal>
  );
};

export default SharingModal;
