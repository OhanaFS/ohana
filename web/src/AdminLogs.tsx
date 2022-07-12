import {
  useMantineTheme,
  Button,
  ScrollArea,
  Table,
} from '@mantine/core';
import { useScrollIntoView } from '@mantine/hooks';
import { BrowserRouter as Router, Link } from 'react-router-dom';
import { useState } from 'react';
import './assets/styles.css';

export interface LogsDetails {
  groupList: Array<any>;
  tableHeader: Array<[string]>;
  tableBody: Array<[string]>;
  caption: string;
  pointerEvents: boolean;
  consoleWidth: number;
  consoleHeight: number;
}

export function AdminLogs(props: LogsDetails) {
  const theme = useMantineTheme();
  const { scrollableRef } = useScrollIntoView();

  let [logs, setValue] = useState(props.groupList);

  const ths = props.tableHeader.map((items) => (
    <th 
     
    >
      {items}
    </th>
  ));

  const rows = logs.map((items) => (
    <tr>
      <td      
       
      >
        {items['Maintenance date']}
      </td>
      <td
       
      >
        {items['Total Files']}
      </td>
      <td
      
       
      >
        {items['Start Time']}
      </td>
      <td
      
      >
        {items['End Time']}
      </td>
      <td
       
      >
        {items['Maintenance Type']}
      </td>
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
        <ScrollArea style={{ height: '80%', width: '100%', marginTop: '1%' }}>
          <Table id='maintenanceLogsTable' 
            captionSide="top"
              striped
              highlightOnHover
              verticalSpacing="sm">
                
            <caption
              style={{
                textAlign: 'center',
                fontWeight: '600',
                fontSize: '24px',
                color: 'black',
                marginTop: '2%',
              }}
            >
              {props.caption}{' '}
            </caption>
            <thead>{ths}</thead>

            <tbody>{rows}</tbody>
          </Table>
        </ScrollArea>

        <div style={{ position: 'relative' }}>
      

          <td
            style={{
              position: 'absolute',
              top: '5px',
              right: '16px',
            }}
          >
            {' '}
            <Button
                  variant="default"
                  color="dark"
                  size="md"
                  style={{
                    position: 'absolute',
                    top: '18px',
                    right: '10px'
                  }}
                  component={Link}
                  to="/runmaintenance"
                >
                  Perform Maintenance
                </Button>
          </td>
        </div>
      </div>
    </div>

  );
}

