import { Accordion, Progress, Text } from '@mantine/core';
import { useParams } from 'react-router-dom';
import { useQueryGetMaintenanceRecordsID } from './api/maintenance';
import AppBase from './components/AppBase';
import { formatDateTime } from './shared/util';

const AdminMaintenanceDetails = () => {
  const params = useParams();
  const qMaintenanceDetails = useQueryGetMaintenanceRecordsID(
    Number(params.id)
  );
  console.log(qMaintenanceDetails.data);
  return (
    <AppBase userType="admin">
      <div className="flex justify-center">
        <div className="w-3/4 flex flex-col bg-white rounded-md p-5">
          <Text className="self-center text-2xl mb-5">Maintenance Summary</Text>
          <Progress
            value={qMaintenanceDetails.data?.progress}
            label={String(qMaintenanceDetails.data?.progress)}
            size="xl"
            mb="md"
          />

          <Accordion>
            {qMaintenanceDetails.data?.orphaned_files_check ? (
              <Accordion.Item value="Orphaned Files Check">
                <Accordion.Control>Orphaned Files Check</Accordion.Control>
                <Accordion.Panel>
                  <Text>
                    Start Time:
                    {' ' +
                      formatDateTime(
                        qMaintenanceDetails.data?.orphaned_files_results[0]
                          .start_time
                      )}
                  </Text>
                  <Text>
                    End Time:
                    {' ' +
                      formatDateTime(
                        qMaintenanceDetails.data?.orphaned_files_results[0]
                          .end_time
                      )}
                  </Text>
                  <Text>
                    Status:{' '}
                    {qMaintenanceDetails.data?.orphaned_files_results[0]
                      .in_progress
                      ? 'In Progress'
                      : 'Done'}
                  </Text>
                  {qMaintenanceDetails.data?.orphaned_files_results[0].msg ? (
                    <Text>
                      Message:{' '}
                      {qMaintenanceDetails.data?.orphaned_files_results[0].msg}
                    </Text>
                  ) : null}
                </Accordion.Panel>
              </Accordion.Item>
            ) : null}
            {qMaintenanceDetails.data?.quick_shards_health_check ? (
              <Accordion.Item value="Check integrity for current version of shards">
                <Accordion.Control>
                  Check integrity for current version of shards
                </Accordion.Control>
                <Accordion.Panel>
                  <Accordion>
                    {qMaintenanceDetails.data.quick_shards_health_progress.map(
                      (item, index) => (
                        <Accordion.Item key={index} value={item.server_name}>
                          <Accordion.Control>
                            {item.server_name}
                          </Accordion.Control>
                          <Accordion.Panel>
                            <Text>
                              Start Time:
                              {' ' + formatDateTime(item.start_time)}
                            </Text>
                            <Text>
                              End Time:
                              {' ' + formatDateTime(item.end_time)}
                            </Text>
                            <Text>
                              Status:{' '}
                              {item.in_progress ? 'In Progress' : 'Done'}
                            </Text>
                            {item.msg ? <Text>Message: {item.msg}</Text> : null}
                          </Accordion.Panel>
                        </Accordion.Item>
                      )
                    )}
                  </Accordion>
                </Accordion.Panel>
              </Accordion.Item>
            ) : null}
            {qMaintenanceDetails.data?.all_files_shards_health_check ? (
              <Accordion.Item value="Full Shards Check">
                <Accordion.Control>Full Shards Check</Accordion.Control>
                <Accordion.Panel>
                  <Accordion>
                    {qMaintenanceDetails.data.all_files_shards_health_progress.map(
                      (item, index) => (
                        <Accordion.Item key={index} value={item.server_name}>
                          <Accordion.Control>
                            {item.server_name}
                          </Accordion.Control>
                          <Accordion.Panel>
                            <Text>
                              Start Time:
                              {' ' + formatDateTime(item.start_time)}
                            </Text>
                            <Text>
                              End Time:
                              {' ' + formatDateTime(item.end_time)}
                            </Text>
                            <Text>
                              Status:{' '}
                              {item.in_progress ? 'In Progress' : 'Done'}
                            </Text>
                            {item.msg ? <Text>Message: {item.msg}</Text> : null}
                          </Accordion.Panel>
                        </Accordion.Item>
                      )
                    )}
                  </Accordion>
                </Accordion.Panel>
              </Accordion.Item>
            ) : null}
            {qMaintenanceDetails.data?.orphaned_shards_check ? (
              <Accordion.Item value="Orphaned Shards Check">
                <Accordion.Control>Orphaned Shards Check</Accordion.Control>
                <Accordion.Panel>
                  <Accordion>
                    {qMaintenanceDetails.data.orphaned_shards_progress.map(
                      (item, index) => (
                        <Accordion.Item key={index} value={item.server_name}>
                          <Accordion.Control>
                            {item.server_name}
                          </Accordion.Control>
                          <Accordion.Panel>
                            <Text>
                              Start Time:
                              {' ' + formatDateTime(item.start_time)}
                            </Text>
                            <Text>
                              End Time:
                              {' ' + formatDateTime(item.end_time)}
                            </Text>
                            <Text>
                              Status:{' '}
                              {item.in_progress ? 'In Progress' : 'Done'}
                            </Text>
                            {item.msg ? <Text>Message: {item.msg}</Text> : null}
                          </Accordion.Panel>
                        </Accordion.Item>
                      )
                    )}
                  </Accordion>
                </Accordion.Panel>
              </Accordion.Item>
            ) : null}
          </Accordion>
        </div>
      </div>
    </AppBase>
  );
};

export default AdminMaintenanceDetails;
