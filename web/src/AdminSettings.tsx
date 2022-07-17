import {
  Checkbox,
  Button,
  Text,
  Table,
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

  function checkUser() {

    setDisable((prevValue) => prevValue);
  }

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
        <div className="settings">
          <Table striped      >
            <caption
              style={{ fontWeight: '600', fontSize: '22px', color: 'black', margin: '15px' }}
            >
              <span style={{ textAlign: 'center' }}>
                Notification Settings
              </span>
            </caption>
            
            <thead></thead>
            <tbody>
             
              <tr style={{}}>
                <div  style={{
                  height: '50px',
                  display: 'flex',
                  flexDirection: 'row',
                  justifyContent: 'space-between',
                  marginLeft:'10px',
                }}>
                  <Text id="settingText" > Allow Cluster health alerts </Text>
                  <Checkbox
                    size="md"
                    style={{
                      marginRight: '50px',
                    }}
                    checked={clusterAlerts}
                    onChange={(event) => [
                      setChecked(event.currentTarget.checked),
                      setDisable(event.currentTarget.checked),
                    ]}
                  />
                </div>
              </tr>
              <tr>
                <div style={{
                  display: 'flex',
                  height: '50px',
                  flexDirection: 'row',
                  justifyContent: 'space-between',
                  marginLeft:'10px',
                }}>

                  <Text id="settingText"> Allow server offline alerts </Text>
                  <Checkbox
                    size="md"
                    id="1"
                    style={{
                      marginRight: '50px',
                    }}
                    checked={sActionAlerts}
                    onChange={(event) => [
                      setChecked1(event.currentTarget.checked),
                      setDisable(event.currentTarget.checked),
                    ]}
                  />
                </div>
              </tr>
              <tr>
                <div style={{
                  display: 'flex',
                  height: '50px',
                  flexDirection: 'row',
                  justifyContent: 'space-between',
                  marginLeft:'10px',
                }}>
                  <Text id="settingText"> Allow supicious action alerts </Text>{' '}
                  <Checkbox
                    size="md"
                    id="1"
                    style={{
                      marginRight: '50px',
                    }}
                    checked={supiciousAlerts}
                    onChange={(event) => [
                      setChecked2(event.currentTarget.checked),
                      setDisable(event.currentTarget.checked),
                    ]}
                  />{' '}
                </div>
              </tr>
              <tr>
                <div style={{
                  display: 'flex',
                  height: '50px',
                  flexDirection: 'row',
                  justifyContent: 'space-between',
                  marginLeft:'10px',
                }}>

                  <Text id="settingText"> Allow server full alert </Text>
                  <Checkbox
                    size="md"
                    id="1"
                    style={{
                      marginRight: '50px',
                    }}
                    checked={serverAlerts}
                    onChange={(event) => [
                      setChecked3(event.currentTarget.checked),
                      setDisable(event.currentTarget.checked),
                    ]}
                  />{' '}

                </div>
              </tr>
              <tr>
                <div style={{
                  display: 'flex',
                  height: '50px',
                  flexDirection: 'row',
                  justifyContent: 'space-between',
                  marginLeft:'10px',
                }}>
                  <Text id="settingText"> Allow supicious file alerts </Text>{' '}
                  <Checkbox
                    size="md"
                    id="1"
                    style={{
                      marginRight: '50px',
                    }}
                    checked={sFileAlerts}
                    onChange={(event) => [
                      setChecked4(event.currentTarget.checked),
                      setDisable(event.currentTarget.checked),
                    ]}
                  />
                </div>
              </tr>
              <tr>
                <div style={{
                  display: 'flex',
                  height: '50px',
                  flexDirection: 'row',
                  justifyContent: 'space-between',
                  marginLeft:'10px',
                }}>
                    <span>
                  <Text id="settingText">   Backup encryption key</Text>{' '}
                  <Text weight={700} style={{marginLeft:'10px'}}>
                    Current Location:{currentLocation}{' '}
                  </Text>{' '}
                  </span>
                  {' '}
                  <Button
                    style={{  marginRight: '10px',marginTop:'20px'}}
                    variant="default"
                    color="dark"
                    size="md"
                  >
                    {' '}
                    Backup
                  </Button>
                </div>
              </tr>
              <tr>
                <div style={{
                  display: 'flex',
                  height: '50px',
                  flexDirection: 'row',
                  justifyContent: 'space-between',
                  marginLeft:'10px',
                }}>
                  {' '}
                  <span style={{marginTop:'10px'}}>
                  <Text id="settingText">     Change the redundancy level of the files</Text>{' '}
                  <Text weight={700} style={{marginLeft:'10px'}}>
                    Current redundancy level:{redundancy}{' '}
                  </Text>{' '}
                  </span>
                  {' '}
                  <Button
                    style={{   marginRight: '10px',marginTop:'20px' }}
                    variant="default"
                    color="dark"
                    size="md"
                  >
                    {' '}
                    Change
                  </Button>
                </div>
              </tr>
          
            </tbody>
            <tfoot>
              <div style={{
                display: 'flex',
                flexDirection: 'column',
              }}>
                <Button
                  disabled={disable}
                  style={{
                      alignSelf:'flex-end',
                      marginTop:'30px',
                      marginRight:'5px'
                      
                  }}
                  variant="default"
                  color="dark"
                  size="md"
                >
                  Save Pending Changes
                </Button>
           
              </div>
            </tfoot>
          </Table>
        </div>
      </div>
    </AppBase>

  );
}


