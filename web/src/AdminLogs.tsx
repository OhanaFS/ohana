import {
  Card,
  Grid,
  useMantineTheme,
  Button,
  Center,
  ScrollArea,
  Table,
} from '@mantine/core';
import { BrowserRouter as Router, Link } from 'react-router-dom';
import { useState } from 'react';

export interface LogsDetails {
  groupList: Array<any>;
  addObjectLabel: string;
  deleteObjectLabel: string;
  tableHeader: Array<[string]>;
  tableBody: Array<[string]>;
  caption: string;
  pointerEvents: boolean;
  consoleWidth: number;
  consoleHeight: number;
}

export function AdminLogs(props: LogsDetails) {
  const theme = useMantineTheme();
  let [logs, setValue] = useState(props.groupList);

  const ths = props.tableHeader.map((items) => (
    <th
      style={{
        width: '20%',
        textAlign: 'left',
        fontWeight: '700',
        fontSize: '16px',
        color: 'black',
      }}
    >
      {items}
    </th>
  ));

  const rows = logs.map((items) => (
    <tr>
      <td
        width="20%"
        style={{
          textAlign: 'left',
          fontWeight: '400',
          fontSize: '16px',
          color: 'black',
        }}
      >
        {items['Maintenance date']}
      </td>
      <td
        width="20%"
        style={{
          textAlign: 'left',
          fontWeight: '400',
          fontSize: '16px',
          color: 'black',
        }}
      >
        {items['Total Files']}
      </td>
      <td
        width="20%"
        style={{
          textAlign: 'left',
          fontWeight: '400',
          fontSize: '16px',
          color: 'black',
        }}
      >
        {items['Start Time']}
      </td>
      <td
        width="20%"
        style={{
          textAlign: 'left',
          fontWeight: '400',
          fontSize: '16px',
          color: 'black',
        }}
      >
        {items['End Time']}
      </td>
      <td
        width="20%"
        style={{
          textAlign: 'left',
          fontWeight: '400',
          fontSize: '16px',
          color: 'black',
        }}
      >
        {items['Maintenance Type']}
      </td>
    </tr>
  ));

  return (
    <>
      <Center style={{ marginLeft: '15%' }}>
        <Grid style={{ width: props.consoleWidth + 'vh' }}>
          <Grid.Col span={12}>
            <Card
              style={{
                marginLeft: '0%',
                height: props.consoleHeight + 'vh',
                border: '1px solid ',
                marginTop: '4%',
                background:
                  theme.colorScheme === 'dark'
                    ? theme.colors.dark[8]
                    : theme.white,
              }}
              shadow="sm"
              p="xl"
            >
              <Card.Section
                style={{ textAlign: 'left', marginLeft: '0%' }}
              ></Card.Section>

              <ScrollArea style={{ height: '90%', marginTop: '1%' }}>
                <Table
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
                    Maintenance Records
                  </caption>
                  <thead>{ths}</thead>
                  <tbody>{rows}</tbody>
                </Table>
              </ScrollArea>

              <div style={{ display: 'flex' }}>
                <Button
                  variant="default"
                  color="dark"
                  size="md"
                  style={{ marginLeft: 'auto', marginTop: '3%' }}
                  component={Link}
                  to="/runmaintenance"
                >
                  Perform Maintenance
                </Button>
              </div>
            </Card>
          </Grid.Col>
        </Grid>
      </Center>
    </>
  );
}
