import {
  Button,
  Checkbox,
  Text,
  Center,
  Card,
  Table,
  useMantineTheme,
} from '@mantine/core';

import { useState } from 'react';

import {  Link } from 'react-router-dom';
import AppBase from './components/AppBase';

function Admin_maintenancesettings() {
  const theme = useMantineTheme();
  let MaintenanceSettings = [
    { name: 'CrawlPermissions', setting: 'true' },
    { name: 'PurgeOrphanedFile', setting: 'true' },
    { name: 'PurgeUser', setting: 'false' },
    { name: 'CrawlReplicas', setting: 'true' },
    { name: 'QuickCheck', setting: 'true' },
    { name: 'FullCheck', setting: 'false' },
    { name: 'DBCheck', setting: 'true' },
    { name: 'DefaultSettings', setting: 'false' },
  ];
  const [checked0, setChecked] = useState(() => {
    if (MaintenanceSettings[0].setting == 'true') {
      return true;
    }

    return false;
  });
  const [checked1, setChecked1] = useState(() => {
    if (MaintenanceSettings[1].setting == 'true') {
      return true;
    }

    return false;
  });
  const [checked2, setChecked2] = useState(() => {
    if (MaintenanceSettings[2].setting == 'true') {
      return true;
    }

    return false;
  });
  const [checked3, setChecked3] = useState(() => {
    if (MaintenanceSettings[3].setting == 'true') {
      return true;
    }

    return false;
  });
  const [checked4, setChecked4] = useState(() => {
    if (MaintenanceSettings[4].setting == 'true') {
      return true;
    }

    return false;
  });
  const [checked5, setChecked5] = useState(() => {
    if (MaintenanceSettings[5].setting == 'true') {
      return true;
    }

    return false;
  });
  const [checked6, setChecked6] = useState(() => {
    if (MaintenanceSettings[6].setting == 'true') {
      return true;
    }

    return false;
  });

  const [checked7, setChecked7] = useState(() => {
    if (MaintenanceSettings[7].setting == 'true') {
      return true;
    }

    return false;
  });

  return (
    <AppBase
      userType="admin"
      name="Alex Simmons"
      username="@alex"
      image="https://images.unsplash.com/photo-1496302662116-35cc4f36df92?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=2070&q=80"
    >
      <Center>
        <Card
          style={{
            marginLeft: '0%',
            height: '65vh',
            border: '1px solid ',
            marginTop: '3%',
            width: '60%',
            background:
              theme.colorScheme === 'dark'
                ? theme.colors.dark[8]
                : theme.white[0],
          }}
          shadow="sm"
          p="xl"
        >
          <Table striped verticalSpacing="md">
            <caption
              style={{
                fontWeight: '600',
                fontSize: '22px',
                color: 'black',
                textAlign: 'left',
                marginLeft: '1%',
              }}
            >
              {' '}
              <span style={{}}>Maintenance Settings</span>
            </caption>

            <thead></thead>
            <tbody style={{}}>
              <tr>
                <td>
                  {' '}
                  <Text style={{}}>
                    {' '}
                    Crawl the list of files to remove permissions from expired
                    users{' '}
                  </Text>
                </td>
                <td>
                  {' '}
                  <Checkbox
                    size="md"
                    checked={checked0}
                    onChange={(event) =>
                      setChecked(event.currentTarget.checked)
                    }
                  />{' '}
                </td>
              </tr>

              <tr>
                <td>
                  {' '}
                  <Text style={{}}> Purging orphaned files and shards </Text>
                </td>
                <td>
                  {' '}
                  <Checkbox
                    size="md"
                    id="1"
                    checked={checked1}
                    onChange={(event) =>
                      setChecked1(event.currentTarget.checked)
                    }
                  />
                </td>
              </tr>

              <tr>
                <td>
                  {' '}
                  <Text style={{}}> Purge a user and their files </Text>{' '}
                </td>
                <td>
                  {' '}
                  <Checkbox
                    size="md"
                    id="1"
                    checked={checked2}
                    onChange={(event) =>
                      setChecked2(event.currentTarget.checked)
                    }
                  />{' '}
                </td>
              </tr>

              <tr>
                <td>
                  {' '}
                  <Text style={{}}>
                    {' '}
                    Crawl all of the files to make sure it has full replicas
                  </Text>{' '}
                </td>
                <td>
                  {' '}
                  <Checkbox
                    size="md"
                    id="1"
                    checked={checked3}
                    onChange={(event) =>
                      setChecked3(event.currentTarget.checked)
                    }
                  />
                </td>
              </tr>

              <tr>
                <td>
                  {' '}
                  <Text style={{}}>
                    {' '}
                    Quick File Check (Only checks current versions of files to
                    see if it’s fine and is not corrupted){' '}
                  </Text>
                </td>
                <td>
                  {' '}
                  <Checkbox
                    size="md"
                    id="1"
                    checked={checked4}
                    onChange={(event) =>
                      setChecked4(event.currentTarget.checked)
                    }
                  />{' '}
                </td>
              </tr>

              <tr>
                <td>
                  {' '}
                  <Text style={{}}>
                    {' '}
                    Full File Check (Checks all fragments to ensure that it’s
                    not corrupted){' '}
                  </Text>{' '}
                </td>
                <td>
                  {' '}
                  <Checkbox
                    size="md"
                    id="1"
                    checked={checked5}
                    onChange={(event) =>
                      setChecked5(event.currentTarget.checked)
                    }
                  />
                </td>
              </tr>

              <tr>
                <td>
                  {' '}
                  <Text style={{}}> DB integrity Check </Text>{' '}
                </td>
                <td>
                  {' '}
                  <Checkbox
                    size="md"
                    id="1"
                    checked={checked6}
                    onChange={(event) =>
                      setChecked6(event.currentTarget.checked)
                    }
                  />
                </td>
              </tr>

              <tr>
                <td>
                  {' '}
                  <Text style={{}}> Use Default settings </Text>{' '}
                </td>
                <td>
                  {' '}
                  <Checkbox
                    size="md"
                    id="1"
                    checked={checked7}
                    onChange={(event) => [
                      setChecked(event.currentTarget.checked),
                      setChecked1(event.currentTarget.checked),
                      setChecked2(event.currentTarget.checked),
                      setChecked3(event.currentTarget.checked),
                      setChecked4(false),
                      setChecked5(event.currentTarget.checked),
                      setChecked6(event.currentTarget.checked),
                      setChecked7(event.currentTarget.checked),
                    ]}
                  />
                </td>
              </tr>
            </tbody>

            <td>
              <Button
                style={{ float: 'right' }}
                variant="default"
                color="dark"
                size="md"
                component={Link}
                to="/Admin_runmaintenance"
              >
                Save Settings
              </Button>
            </td>
          </Table>
        </Card>
      </Center>
    </AppBase>
  );
}

export default Admin_maintenancesettings;
