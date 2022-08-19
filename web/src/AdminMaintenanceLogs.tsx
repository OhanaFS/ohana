import { ScrollArea, Button, Table, ActionIcon } from '@mantine/core';
import { IconInfoCircle } from '@tabler/icons';
import { Link, useNavigate } from 'react-router-dom';
import { useQueryGetMaintenanceRecords } from './api/maintenance';
import AppBase from './components/AppBase';
import { formatDateTime } from './shared/util';

export function AdminMaintenanceLogs() {
  const qGetMaintenanceRecords = useQueryGetMaintenanceRecords(0, '', '', '');
  const maintenanceRecords = qGetMaintenanceRecords.data ?? [];

  const navigate = useNavigate();

  // variable that show all the logs inside the props.groupList
  const logsHeader = ['Maintenance date', 'Time Taken', 'Issues', ''];

  // display table header that is from props
  const ths = logsHeader.map((items, i) => (
    <th key={i} style={{ fontWeight: '600' }}>
      {items}
    </th>
  ));

  // display all the rows that is from props
  const rows = maintenanceRecords.map((items, i) => (
    <tr key={i}>
      <td>{formatDateTime(items.start_time)}</td>
      <td>{items.total_time_taken}</td>
      <td>{items.status_msg ? items.status_msg : 'None'}</td>
      <td>
        <ActionIcon
          onClick={() => navigate(`/maintenance/${items.id}`)}
          variant="transparent"
        >
          <IconInfoCircle size={25} />
        </ActionIcon>
      </td>
    </tr>
  ));
  return (
    <>
      <AppBase userType="admin">
        <div
          style={{
            display: 'flex',
            justifyContent: 'center',
            height: '85vh',
          }}
        >
          <div className="maintenanceLogs">
            <ScrollArea
              style={{
                height: '90%',
                width: '100%',
                marginTop: '10px',
                padding: '10px',
              }}
            >
              <Table
                id="maintenanceLogsTable"
                captionSide="top"
                highlightOnHover
                verticalSpacing="sm"
              >
                <caption
                  style={{
                    textAlign: 'center',
                    fontWeight: '600',
                    fontSize: '24px',
                    color: 'black',
                  }}
                >
                  Maintenance Records
                </caption>
                <thead>
                  <tr>{ths}</tr>
                </thead>
                <tbody>{rows}</tbody>
              </Table>
            </ScrollArea>

            <div
              style={{
                display: 'flex',
                flexDirection: 'column',
              }}
            >
              <Button
                variant="default"
                color="dark"
                size="md"
                style={{
                  alignSelf: 'flex-end',
                  marginRight: '15px',
                  marginTop: '10px',
                }}
                component={Link}
                to="/runmaintenance"
              >
                Perform Maintenance
              </Button>
            </div>
          </div>
        </div>
      </AppBase>
    </>
  );
}

export default AdminMaintenanceLogs;
function useState(groupList: any): [any, any] {
  throw new Error('Function not implemented.');
}
