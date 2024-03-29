import { Button, Checkbox, Group, Table, Text } from '@mantine/core';
import { useForm } from '@mantine/form';
import AppBase from './components/AppBase';
import { useMutateStartMainRecordsID } from './api/maintenance';
import { useNavigate } from 'react-router-dom';
import { showNotification } from '@mantine/notifications';

export function AdminRunMaintenance() {
  const form = useForm({
    initialValues: {
      full_shards_check: false,
      quick_shards_check: false,
      missing_shards_check: false,
      orphaned_shards_check: false,
      orphaned_files_check: false,
      permission_check: false,
      delete_fragments: false,
    },
  });

  const inputFields = [
    { id: 'full_shards_check', desc: 'Check all the shards for integrity' },
    {
      id: 'quick_shards_check',
      desc: 'Check integrity for current version of all shards',
    },
    { id: 'orphaned_shards_check', desc: 'Check if there are garbage shards' },
    { id: 'orphaned_files_check', desc: 'Check if there are garbage files' },
  ];

  const mStartMaintenance = useMutateStartMainRecordsID();
  const navigate = useNavigate();

  return (
    <>
      <AppBase userType="admin">
        <div className="flex justify-center">
          <div className="w-auto p-8 bg-white flex flex-col items-center rounded-lg">
            <Text weight={500} className="text-3xl" mb="lg">
              Run Maintenance
            </Text>
            <form
              onSubmit={form.onSubmit((values) => {
                Object.values(values).some((val) => val)
                  ? mStartMaintenance
                      .mutateAsync(values)
                      .then((e) => navigate(`/maintenance/${e.id}`))
                  : showNotification({
                      title: 'Nothing Checked',
                      message: 'Please check at least one field!',
                    });
              })}
            >
              <Table horizontalSpacing="xl" verticalSpacing="xs" fontSize="xl">
                <tbody>
                  {inputFields.map((field, i) => (
                    <tr key={i}>
                      <td>
                        <Text>{field.desc}</Text>
                      </td>
                      <td>
                        <Checkbox
                          mt="md"
                          {...form.getInputProps(field.id, {
                            type: 'checkbox',
                          })}
                        />
                      </td>
                    </tr>
                  ))}
                </tbody>
              </Table>

              <Group position="right" mt="lg">
                <Button type="submit">Perform Maintenance</Button>
              </Group>
            </form>
          </div>
        </div>
      </AppBase>
    </>
  );
}
