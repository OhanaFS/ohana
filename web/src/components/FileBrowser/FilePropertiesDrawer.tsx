import React from 'react';
import {
  AspectRatio,
  Drawer,
  Accordion,
  ScrollArea,
  Loader,
  Button,
  Stack,
  Group,
} from '@mantine/core';
import { IconShare, IconStar } from '@tabler/icons';
import { EntryType, useQueryFileMetadata } from '../../api/file';
import FilePreview from './FilePreview';
import FileVersions from './Properties/FileVersions';
import PasswordForm from './Properties/PasswordForm';
import PropertiesTable from './Properties/PropertiesTable';
import SharingProperties from './Properties/Sharing';
import SharingModal from './SharingModal';

export type FilePropertiesDrawerProps = {
  fileId: string;
  onClose: () => void;
};

const FilePropertiesDrawer = (props: FilePropertiesDrawerProps) => {
  const { fileId, onClose } = props;
  const [isSharingModalOpen, setIsSharingModalOpen] = React.useState(false);
  const qFile = useQueryFileMetadata(fileId);

  return (
    <>
      <Drawer
        opened={fileId !== ''}
        onClose={() => {
          setIsSharingModalOpen(false);
          onClose();
        }}
        title={qFile.data?.file_name}
        padding="lg"
        position="right"
        size="xl"
        styles={{
          title: {
            overflow: 'hidden',
            textOverflow: 'ellipsis',
            whiteSpace: 'nowrap',
          },
        }}
      >
        <ScrollArea className="h-full">
          <Stack spacing="md">
            {qFile.data?.entry_type === EntryType.File ? (
              <AspectRatio ratio={16 / 9}>
                <FilePreview fileId={fileId} />
              </AspectRatio>
            ) : null}

            <Group grow>
              <Button
                leftIcon={<IconStar size={16} fill="currentColor" />}
                variant="light"
              >
                Add to favorites
              </Button>
              <Button
                leftIcon={<IconShare size={16} fill="currentColor" />}
                variant="light"
                onClick={() => setIsSharingModalOpen(true)}
              >
                Share
              </Button>
            </Group>

            <Accordion defaultValue="properties">
              <Accordion.Item value="properties">
                <Accordion.Control>Properties</Accordion.Control>
                <Accordion.Panel>
                  {qFile.data ? (
                    <PropertiesTable metadata={qFile.data} />
                  ) : (
                    <Loader />
                  )}
                </Accordion.Panel>
              </Accordion.Item>
              {qFile.data?.entry_type === EntryType.File ? (
                <>
                  <Accordion.Item value="password">
                    <Accordion.Control>Password</Accordion.Control>
                    <Accordion.Panel>
                      <PasswordForm fileId={fileId} />
                    </Accordion.Panel>
                  </Accordion.Item>
                  <Accordion.Item value="versioning">
                    <Accordion.Control>File Versions</Accordion.Control>
                    <Accordion.Panel>
                      <FileVersions fileId={fileId} />
                    </Accordion.Panel>
                  </Accordion.Item>
                </>
              ) : null}
            </Accordion>
          </Stack>
        </ScrollArea>
      </Drawer>
      <SharingModal
        fileId={fileId}
        opened={isSharingModalOpen}
        onClose={() => setIsSharingModalOpen(false)}
      />
    </>
  );
};

export default FilePropertiesDrawer;
