import {
  Button,
  Textarea,
  Checkbox,
  Text,
  Divider,
  NumberInput,
  Group,
  TextInput,
  PasswordInput,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { showNotification } from '@mantine/notifications';
import { useState } from 'react';
import { useMutatePostFileKey, useMutateUpdateStitch } from './api/maintenance';
import AppBase from './components/AppBase';

export function AdminConfiguration() {
  const mRotateKey = useMutatePostFileKey();
  const mUpdateStitch = useMutateUpdateStitch();
  //function will be rotate key
  function rotateKey() {}

  const keyRotationForm = useForm({
    initialValues: {
      file_id: '',
      password: '',
    },
  });

  const form = useForm({
    initialValues: {
      data_shards: 2,
      parity_shards: 1,
      key_threshold: 2,
    },
  });

  return (
    <AppBase userType="admin">
      <div
        style={{
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'flex-start',
          height: '100%',
        }}
      >
        <div
          style={{
            display: 'flex',
            border: '1px solid #ccc',
            flexDirection: 'column',
            justifyContent: 'center',
            alignItems: 'center',
            width: '90%',
            backgroundColor: 'white',
            borderRadius: '10px',
            padding: '20px',
            maxWidth: '500px',
          }}
        >
          <caption
            style={{
              textAlign: 'center',
              fontWeight: 600,
              fontSize: '24px',
              color: 'black',
              marginBottom: '20px',
              alignSelf: 'center',
            }}
          >
            Settings
          </caption>
          <div className="flex flex-col w-full">
            <Divider
              my="xs"
              label="Rotate Key"
              variant="dotted"
              labelPosition="center"
            />
            <form
              className="mt-3"
              onSubmit={keyRotationForm.onSubmit((values) => {
                values.file_id === ''
                  ? showNotification({
                      message: 'File ID required',
                    })
                  : mRotateKey
                      .mutateAsync(values)
                      .then((e) =>
                        showNotification({
                          message: 'Success',
                        })
                      )
                      .catch((e) =>
                        showNotification({
                          message: String(e.message),
                        })
                      );
              })}
            >
              <TextInput
                label="File ID *"
                placeholder="Please enter the file ID"
                {...keyRotationForm.getInputProps('file_id')}
              />
              <PasswordInput
                className="mt-2"
                placeholder="Password"
                label="Password"
                {...keyRotationForm.getInputProps('password')}
              />
              <Group position="right" mt="lg">
                <Button type="submit">Submit</Button>
              </Group>
            </form>
          </div>

          <div className="flex flex-col w-full mt-5">
            <Divider
              my="xs"
              label="Set Redundancy Level"
              variant="dotted"
              labelPosition="center"
            />
            <form
              className="mt-3"
              onSubmit={form.onSubmit((values) => {
                console.log(values);
                mUpdateStitch
                  .mutateAsync(values)
                  .then((e) =>
                    showNotification({
                      message: 'Success',
                    })
                  )
                  .catch((e) =>
                    showNotification({
                      message: String(e.message),
                    })
                  );
              })}
            >
              <NumberInput
                label="Number of Data Shards"
                description="From 1 to 10"
                max={10}
                min={1}
                {...form.getInputProps('data_shards')}
              />
              <NumberInput
                className="mt-2"
                label="Number of Parity Shards"
                description="From 1 to 10"
                max={10}
                min={1}
                {...form.getInputProps('parity_shards')}
              />
              <NumberInput
                className="mt-2"
                label="Key Threshold Value"
                description="From 1 to 10"
                max={10}
                min={1}
                {...form.getInputProps('key_threshold')}
              />
              <Group position="right" mt="lg">
                <Button type="submit">Submit</Button>
              </Group>
            </form>
          </div>
        </div>
      </div>
    </AppBase>
  );
}
