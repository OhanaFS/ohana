import {
  AspectRatio,
  Drawer,
  Table,
  Accordion,
  Text,
  PasswordInput,
  Button,
  Group,
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
          <Accordion.Item value="password">
            <Accordion.Control>Password</Accordion.Control>
            <Accordion.Panel>
              <Text>
                {qFileMeta.data?.password_protected
                  ? 'Password Protected: Enter your old password and new password to change it'
                  : 'No Password: Set a password below'}
              </Text>
              <form
                onSubmit={form.onSubmit((values) => {
                  if (qFileMeta.data?.password_protected) {
                    mFileMeta.mutate({
                      file_id: fileId,
                      old_password: values.password,
                      new_password: values.password_c,
                      password_modification: true,
                      password_hint: 'test',
                    });
                  } else {
                    if (values.password !== values.password_c) {
                      console.log('error');
                      showNotification({
                        title: 'Password Mismatch',
                        message: "The two password fields don't match",
                        icon: <IconX />,
                        color: 'red',
                      });
                      return;
                    }
                    mFileMeta.mutate({
                      file_id: fileId,
                      new_password: values.password,
                      password_modification: true,
                      password_protected: true,
                      password_hint: 'test',
                    });
                  }
                })}
              >
                <PasswordInput
                  placeholder={
                    qFileMeta.data?.password_protected
                      ? 'Enter Current Password'
                      : 'Create New Password'
                  }
                  {...form.getInputProps('password')}
                />
                <PasswordInput
                  placeholder={
                    qFileMeta.data?.password_protected
                      ? 'Enter New Password'
                      : 'Confirm Password'
                  }
                  {...form.getInputProps('password_c')}
                />
                <Group position="right" mt="md">
                  <Button className="" type="submit">
                    Submit
                  </Button>
                </Group>
              </form>
            </Accordion.Panel>
          </Accordion.Item>
        ) : null}
      </Accordion>
    </Drawer>
  );
};

export default FilePropertiesDrawer;
