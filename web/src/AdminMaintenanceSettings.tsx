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

export function AdminMaintenanceSettings() {
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
      

      <div
          style={{
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'flex-start',
           
          }}

        >
          <div className="maintenanceSettings">
          
            <Table striped verticalSpacing="xs" id='maintenanceSettings' >
              <caption  >
                {' '} <div style={{marginTop:'10px'}}>    Maintenance Settings
              
            
                </div>
              </caption>

              <thead></thead>
              <tbody style={{}}>
                <tr>
                  <td>

                      Crawl the list of files to remove permissions from expired
                      users
            
                  </td>
                  <td>
        
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
               
                  Purging orphaned files and shards 
                  </td>
                  <td>
             
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
                   Purge a user and their files 
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
                  
                      {' '}
                      Crawl all of the files to make sure it has full replicas
                  
                  </td>
                  <td>
                    {' '}
                    <Checkbox
                      size="md"
                      id="2"
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
                
                      {' '}
                      Quick File Check (Only checks current versions of files to
                      see if it’s fine and is not corrupted){' '}
                 
                  </td>
                  <td>
                    {' '}
                    <Checkbox
                      size="md"
                      id="3"
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
               
                      {' '}
                      Full File Check (Checks all fragments to ensure that it’s
                      not corrupted){' '}
                  
                  </td>
                  <td>
                    {' '}
                    <Checkbox
                      size="md"
                      id="4"
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
                    DB integrity Check 
                  </td>
                  <td>
                    {' '}
                    <Checkbox
                      size="md"
                      id="5"
                      checked={checked6}
                      onChange={(event) =>
                        setChecked6(event.currentTarget.checked)
                      }
                    />
                  </td>
                </tr>

                <tr>
                  <td> </td>
                  <td> </td>
                </tr>
              </tbody>
              <div style={{ position: 'relative' }}>
              <td  >
                <Button
                    style={{
                      position: 'absolute',
                      top: '0px',
                      right: '0px'
                    }}
                  variant="default"
                  color="dark"
                  size="md"
                  component={Link}
                  to="/runmaintenance"
                >
                  Save Settings
                </Button>
              </td>
              </div>
            </Table>
       </div>
       </div>
    </AppBase>
  );
}
