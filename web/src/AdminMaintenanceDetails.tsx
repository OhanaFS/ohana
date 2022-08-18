import { Accordion, Text } from '@mantine/core';
import { useParams } from 'react-router-dom';
import { useQueryGetMaintenanceRecordsID } from './api/maintenance';
import AppBase from './components/AppBase';

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
          <Accordion>
            {qMaintenanceDetails.data?.orphaned_files_check ? (
              <Accordion.Item value="Orphaned Shards Check">
                <Accordion.Control>Orphaned Shards Check</Accordion.Control>
                <Accordion.Panel></Accordion.Panel>
              </Accordion.Item>
            ) : null}
          </Accordion>
        </div>
      </div>
    </AppBase>
  );
};

export default AdminMaintenanceDetails;
