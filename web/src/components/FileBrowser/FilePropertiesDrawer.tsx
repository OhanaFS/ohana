import {
  AspectRatio,
  Drawer,
  Table,
  Accordion,
  Text,
  PasswordInput,
  Button,
  Group,
  ScrollArea,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { showNotification } from '@mantine/notifications';
import { IconX } from '@tabler/icons';
import {
  EntryType,
  MetadataKeyMap,
  useMutateUpdateFileMetadata,
  useQueryFileMetadata,
} from '../../api/file';
import FilePreview from './FilePreview';
import FileVersions from './Properties/FileVersions';
import PasswordForm from './Properties/PasswordForm';

export type FilePropertiesDrawerProps = {
  fileId: string;
  onClose: () => void;
};

const FilePropertiesDrawer = (props: FilePropertiesDrawerProps) => {
  const { fileId, onClose } = props;
  const qFile = useQueryFileMetadata(fileId);
  const qFileMeta = useQueryFileMetadata(fileId);
  const mFileMeta = useMutateUpdateFileMetadata();

  const form = useForm({
    initialValues: {
      password: '',
      password_c: '',
    },
  });

  return (
    <Drawer
      opened={fileId !== ''}
      onClose={onClose}
      title={qFile.data?.file_name}
      padding="lg"
      position="right"
      size="xl"
    >
      <ScrollArea className="h-full">
        <AspectRatio ratio={16 / 9}>
          <FilePreview fileId={fileId} />
        </AspectRatio>
        <Accordion defaultValue="properties">
          <Accordion.Item value="properties">
            <Accordion.Control>Properties</Accordion.Control>
            <Accordion.Panel>
              <Table>
                <thead>
                  <tr>
                    <th>Property</th>
                    <th>Value</th>
                  </tr>
                </thead>
                <tbody>
                  {Object.keys(qFile.data || {})
                    .map((key) => key as keyof typeof MetadataKeyMap)
                    .filter((key) =>
                      (
                        [
                          'file_name',
                          'size',
                          'created_time',
                          'modified_time',
                          'version_no',
                        ] as Array<keyof typeof MetadataKeyMap>
                      ).includes(key)
                    )
                    .map((key) => (
                      <tr key={key}>
                        <td>{MetadataKeyMap[key]}</td>
                        <td>{(qFile.data as any)[key]}</td>
                      </tr>
                    ))}
                </tbody>
              </Table>
            </Accordion.Panel>
          </Accordion.Item>
          {qFile.data?.entry_type === EntryType.File ? (
            <>
              <Accordion.Item value="password">
                <Accordion.Control>Password</Accordion.Control>
                <Accordion.Panel>
                  <PasswordForm fileID={fileId} />
                </Accordion.Panel>
              </Accordion.Item>
              <Accordion.Item value="versioning">
                <Accordion.Control>File Versions</Accordion.Control>
                <Accordion.Panel>
                  <FileVersions fileId={fileId} />
                </Accordion.Panel>
              </Accordion.Item>
            </>
          ) : null}{' '}
        </Accordion>
      </ScrollArea>
    </Drawer>
  );
};

export default FilePropertiesDrawer;
