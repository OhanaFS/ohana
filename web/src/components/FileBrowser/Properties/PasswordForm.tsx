import React from 'react';
import { Button, Group, PasswordInput, Text } from '@mantine/core';
import { useForm } from '@mantine/form';

import {
  useMutateUpdateFileMetadata,
  useQueryFileMetadata,
} from '../../../api/file';
import { showNotification } from '@mantine/notifications';
import { IconX } from '@tabler/icons';

type PasswordFormProps = {
  fileId: string;
};

const PasswordForm = (props: PasswordFormProps) => {
  const qFileMeta = useQueryFileMetadata(props.fileId);
  const mFileMeta = useMutateUpdateFileMetadata();

  const form = useForm({
    initialValues: {
      password: '',
      password_c: '',
    },
  });

  return (
    <>
      <Text>
        {qFileMeta.data?.password_protected
          ? 'Password Protected: Enter your old password and new password to change it'
          : 'No Password: Set a password below'}
      </Text>
      <form
        onSubmit={form.onSubmit((values) => {
          if (qFileMeta.data?.password_protected) {
            mFileMeta.mutate({
              file_id: props.fileId,
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
              file_id: props.fileId,
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
    </>
  );
};

export default PasswordForm;
