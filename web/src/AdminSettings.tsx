import {
  Checkbox,
  Button,
  Text,
  Center,
  Table,
  Card,
  useMantineTheme,
} from '@mantine/core';
import { useState } from 'react';
import AppBase from './components/AppBase';

export function AdminSettings() {
  const theme = useMantineTheme();
  //retrieve from database
  let ConfigurationSettings = [
    { name: 'clusterAlerts', setting: 'true' },
    { name: 'sActionAlerts', setting: 'true' },
    { name: 'supiciousAlerts', setting: 'false' },
    { name: 'serverAlerts', setting: 'true' },
    { name: 'sFileAlerts', setting: 'true' },
    { name: 'BackupLocation', setting: 'C:\\Users\\admin' },
    { name: 'redundancy', setting: 'Low' },
  ];

  const [clusterAlerts, setChecked] = useState(() => {
    if (ConfigurationSettings[0].setting === 'true') {
      return true;
    }
    return false;
  });
  const [sActionAlerts, setChecked1] = useState(() => {
    if (ConfigurationSettings[1].setting === 'true') {
      return true;
    }

    return false;
  });
  const [supiciousAlerts, setChecked2] = useState(() => {
    if (ConfigurationSettings[2].setting === 'true') {
      return true;
    }

    return false;
  });
  const [serverAlerts, setChecked3] = useState(() => {
    if (ConfigurationSettings[3].setting === 'true') {
      return true;
    }

    return false;
  });
  const [sFileAlerts, setChecked4] = useState(() => {
    if (ConfigurationSettings[4].setting === 'true') {
      return true;
    }

    return false;
  });

  let currentLocation = ConfigurationSettings[5].setting;
  let redundancy = ConfigurationSettings[6].setting;

  const [disable, setDisable] = useState(true);

  function checkUser(){

    setDisable((prevValue) => prevValue);
  }
  
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
            alignItems: 'flex-start',
          }}
            
            > 
        <div className="settings">
            <Table striped verticalSpacing="md"     >
              <caption
                style={{ fontWeight: '600', fontSize: '22px', color: 'black' }}
              >
  
                <span style={{ textAlign: 'center' }}>
                  Notification Settings
                </span>
              </caption>

              <thead></thead>
              <tbody style={{}}>
                <tr>
                  <td
                  
                    style={{
                      textAlign: 'left',
                      fontWeight: '400',
                      fontSize: '18px',
                      color: 'black',
                    }}
                  >

                    <Text style={{}}> Allow Cluster health alerts </Text>
                  </td>
                  <td style={{ position:'relative'}}>
                      <Checkbox
                        size="md"
                        style={{     
                          position: 'absolute',
                          top:'18px',
                          right: '16px'}}
                        checked={clusterAlerts}
                        onChange={(event) => [
                          setChecked(event.currentTarget.checked),
                          setDisable(event.currentTarget.checked),
                        ]}
                      />{' '}
                   {' '}
                  </td>
                </tr>

                <tr>
                  <td
                   
                    style={{
                      textAlign: 'left',
                      fontWeight: '400',
                      fontSize: '18px',
                      color: 'black',
                    }}
                  >
                    {' '}
                    <Text style={{}}> Allow server offline alerts </Text>
                  </td>
                  <td style={{ position:'relative'}}>
             
                      <Checkbox
                        size="md"
                        id="1"
                        style={{     
                          position: 'absolute',
                          top:'18px',
                          right: '16px'}}
                        checked={sActionAlerts}
                        onChange={(event) => [
                          setChecked1(event.currentTarget.checked),
                          setDisable(event.currentTarget.checked),
                        ]}
                      />
                   
                  </td>
                </tr>

                <tr>
                  <td
                  
                    style={{
                      textAlign: 'left',
                      fontWeight: '400',
                      fontSize: '18px',
                      color: 'black',
                    }}
                  >
                    {' '}
                    <Text style={{}}> Allow supicious action alerts </Text>{' '}
                  </td>
                  <td style={{ position:'relative'}}>
               
                      <Checkbox
                        size="md"
                        id="1"
                        style={{     
                          position: 'absolute',
                          top:'18px',
                          right: '16px'}}
                        checked={supiciousAlerts}
                        onChange={(event) => [
                          setChecked2(event.currentTarget.checked),
                          setDisable(event.currentTarget.checked),
                        ]}
                      />{' '}
                
                  </td>
                </tr>

                <tr>
                  <td
                   
                    style={{
                      textAlign: 'left',
                      fontWeight: '400',
                      fontSize: '18px',
                      color: 'black',
                    }}
                  >
                    {' '}
                    <span style={{}}> Allow server full alert </span>{' '}
                  </td>
                  <td style={{ position:'relative'}} >
                  
                      <Checkbox
                        size="md"
                        id="1"
                        style={{     
                          position: 'absolute',
                          top:'18px',
                          right: '16px'}}
                        checked={serverAlerts}
                        onChange={(event) => [
                          setChecked3(event.currentTarget.checked),
                          setDisable(event.currentTarget.checked),
                        ]}
                      />{' '}
                  
                  </td>
                </tr>

                <tr>
                  <td
                  
                    style={{
                      textAlign: 'left',
                      fontWeight: '400',
                      fontSize: '18px',
                      color: 'black',
                    }}
                  >
                    {' '}
                    <span style={{}}> Allow supicious file alerts </span>
                  </td>
                
                  <td  style={{ position:'relative'}}>
                   
                      {' '}
                      <Checkbox
                        size="md"
                        id="1"
                        style={{     
                          position: 'absolute',
                          top:'18px',
                          right: '16px'}}
                        checked={sFileAlerts}
                        onChange={(event) => [
                          setChecked4(event.currentTarget.checked),
                          setDisable(event.currentTarget.checked),
                        ]}
                      />
                  
                  </td>
                 
                </tr>

                <tr>
                  <td
                    style={{
                      textAlign: 'left',
                      fontWeight: '400',
                      fontSize: '18px',
                      color: 'black',
                    }}
                  >
                    {' '}
                    Backup encryption key{' '}
                    <Text weight={700}>
                      Current Location:{currentLocation}{' '}
                    </Text>{' '}
                  </td>
                  <td>
                    {' '}
                    <Button
                      style={{ float: 'right' }}
                      variant="default"
                      color="dark"
                      size="md"
                    >
                      {' '}
                      Backup
                    </Button>
                  </td>
                </tr>

                <tr>
                  <td
                    style={{
                      textAlign: 'left',
                      fontWeight: '400',
                      fontSize: '18px',
                      color: 'black',
                    }}
                  >
                    {' '}
                    Change the redundancy level of the files{' '}
                    <Text weight={700}>
                      Current redundancy level:{redundancy}{' '}
                    </Text>{' '}
                  </td>
                  <td>
                    {' '}
                    <Button
                      style={{ float: 'right' }}
                      variant="default"
                      color="dark"
                      size="md"
                    >
                      {' '}
                      Change
                    </Button>
                  </td>
                </tr>
              </tbody>
              <tfoot>
             
                <td colSpan={2} style={{ position:'relative'}}>
                  <Button
                    disabled={disable}
                    style={{     
                      position: 'absolute',
                      top:'18px',
                      right: '16px'}}
                    variant="default"
                    color="dark"
                    size="md"
                  >
                    Save Pending Changes
                  </Button>
                </td>
              </tfoot>
            </Table>
      </div>
          </div>
      </AppBase>
    </>
  );
}


