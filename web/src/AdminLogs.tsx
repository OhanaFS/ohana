import { Button, ScrollArea, Table } from '@mantine/core';
import { BrowserRouter as Router, Link } from 'react-router-dom';
import { useState } from 'react';
import './assets/styles.css';

export interface LogsDetails {

   // all the data needed from database 
  groupList: Array<any>;
  tableHeader: Array<[string]>;
  tableBody: Array<[string]>;
  caption: string;
  //whether the group can be press
  pointerEvents: boolean;
}

export function AdminLogs(props: LogsDetails) {

   // variable that show all the logs inside the props.groupList
  let [logs, setValue] = useState(props.groupList);

  // display table header that is from props
  const ths = props.tableHeader.map((items) => (
    <th style={{ fontWeight: '600' }}>{items}</th>
  ));
  
  // display all the rows that is from props
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
  );
}
