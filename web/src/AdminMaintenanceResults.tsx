import {
  Button,
  Table,
  ScrollArea,
} from '@mantine/core';
import { useState } from 'react';
import { Link } from 'react-router-dom';
import AppBase from './components/AppBase';

export interface LogsDetails {
  groupList: Array<any>;
  tableHeader: Array<[string]>;
  tableBody: Array<[string]>;
  caption: string;
  pointerEvents: boolean;
}
export function AdminMaintenanceResults(props: LogsDetails) {
  
  let [logs, setValue] = useState(props.groupList);
  const ths = props.tableHeader.map((items) => (
    <th style={{ fontWeight: '600' }}>{items}</th>
  ));

  const rows = logs.map((items) => (
    <tr>
      <td>{items['Maintenance date']}</td>
      <td>{items['Total Files']}</td>
      <td>{items['Start Time']}</td>
      <td>{items['End Time']}</td>
      <td>{items['Maintenance Type']}</td>
    </tr>
  ));


  return (
    <>
      <AppBase
        userType="admin"
        name="Alex Simmons"
        username="@alex"
        image="https://images.unsplash.com/photo-1496302662116-35cc4f36df92?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=2070&q=80"
      >
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
            height: '85%',
            width: '90%',
            marginTop: '10px',
            marginLeft: '20px',
          }}
        >
          <Table
            id="maintenanceLogsTable"
            captionSide="top"
            striped
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
              {props.caption}
            </caption>
            <thead>{ths}</thead>
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
