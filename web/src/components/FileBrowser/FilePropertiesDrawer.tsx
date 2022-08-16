import {
  AspectRatio,
  Drawer,
  Table,
  Accordion,
  ScrollArea,
  Loader,
} from '@mantine/core';
import {
  EntryType,
  MetadataKeyMap,
  useQueryFileMetadata,
} from '../../api/file';
import FilePreview from './FilePreview';
import FileVersions from './Properties/FileVersions';
import PasswordForm from './Properties/PasswordForm';
import PropertiesTable from './Properties/PropertiesTable';
import SharingProperties from './Properties/Sharing';

export type FilePropertiesDrawerProps = {
  fileId: string;
  onClose: () => void;
};

const FilePropertiesDrawer = (props: FilePropertiesDrawerProps) => {
  const { fileId, onClose } = props;
  const qFile = useQueryFileMetadata(fileId);

  return (
    <Drawer
      opened={fileId !== ''}
      onClose={onClose}
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
        {qFile.data?.entry_type === EntryType.File ? (
          <AspectRatio ratio={16 / 9}>
            <FilePreview fileId={fileId} />
          </AspectRatio>
        ) : null}
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
              <Accordion.Item value="sharing">
                <Accordion.Control>Sharing</Accordion.Control>
                <Accordion.Panel>
                  <SharingProperties fileId={fileId} />
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
      </ScrollArea>
    </Drawer>
  );
};

export default FilePropertiesDrawer;
