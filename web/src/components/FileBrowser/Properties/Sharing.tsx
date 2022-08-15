import {
  ActionIcon,
  Anchor,
  CopyButton,
  Group,
  Menu,
  Stack,
  Text,
  Tooltip,
} from '@mantine/core';
import { IconCheck, IconCopy, IconDots, IconPlus, IconX } from '@tabler/icons';
import {
  getSharingLinkURL,
  useMutateCreateSharingLink,
  useMutateDeleteSharingLink,
  useQueryFileSharingLinks,
} from '../../../api/sharing';

export type SharingPropertiesProps = {
  fileId: string;
};

const SharingProperties = (props: SharingPropertiesProps) => {
  const { fileId } = props;
  const qSharedLinks = useQueryFileSharingLinks(fileId);
  const mCreateSharedLink = useMutateCreateSharingLink();
  const mDeleteSharedLink = useMutateDeleteSharingLink();

  return (
    <Stack>
      <Group position="apart">
        <Text>Share link</Text>
        <Tooltip label="Create sharing link" withArrow position="left">
          <ActionIcon
            onClick={() => mCreateSharedLink.mutate({ fileId })}
            disabled={mCreateSharedLink.isLoading}
          >
            <IconPlus />
          </ActionIcon>
        </Tooltip>
      </Group>
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
    </Stack>
  );
};

export default SharingProperties;
