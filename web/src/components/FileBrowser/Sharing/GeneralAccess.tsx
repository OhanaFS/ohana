import React from 'react';
import {
  Switch,
  Group,
  Anchor,
  CopyButton,
  ActionIcon,
  Menu,
  Tooltip,
  Text,
} from '@mantine/core';
import { IconCheck, IconCopy, IconDots, IconPlus, IconX } from '@tabler/icons';
import {
  EntryType,
  FileMetadata,
  useQueryFileMetadata,
} from '../../../api/file';
import {
  getSharingLinkURL,
  useMutateCreateSharingLink,
  useMutateDeleteSharingLink,
  useQueryFileSharingLinks,
} from '../../../api/sharing';

export type GeneralAccessProps = {
  fileId: string;
};

const GeneralAccess = ({ fileId }: GeneralAccessProps) => {
  const qFile = useQueryFileMetadata(fileId);
  const file = qFile.data as FileMetadata | undefined;

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
    <>
      <Text weight="bold" pt="sm">
        General access
      </Text>

      <Group position="apart">
        <Switch
          disabled={
            file?.entry_type === EntryType.Folder ||
            qSharedLinks.isLoading ||
            mCreateSharedLink.isLoading
          }
          label={
            file?.entry_type === EntryType.Folder
              ? 'Folders cannot be shared publicly'
              : isPublicShared
              ? 'Anyone with the link can view'
              : 'Enable public link'
          }
          onChange={(e) => setIsPublicShared(e.currentTarget.checked)}
          checked={isPublicShared}
        />
        {isPublicShared && (
          <Tooltip label="Create a new link" withArrow position="left">
            <ActionIcon onClick={() => mCreateSharedLink.mutate({ fileId })}>
              <IconPlus size={16} />
            </ActionIcon>
          </Tooltip>
        )}
      </Group>

      {sharedLinks.map((link) => (
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
                    {copied ? <IconCheck size={16} /> : <IconCopy size={16} />}
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
    </>
  );
};

export default GeneralAccess;
