import {
  Checkbox,
  Button,
  Text,
  Table,
  useMantineTheme,
  Modal,
  Radio,
  RadioGroup,
  createStyles,
  TextInput,
  Textarea
} from '@mantine/core';
import { useState } from 'react';
import AppBase from './components/AppBase';

const titleStyle = createStyles(() => ({
  title: {
    fontWeight: 600,
    fontSize: '24px',

  },
  input: {
    fontWeight: 600,
    fontSize: '24px',

  }
}));
export function AdminSettings() {
  const theme = useMantineTheme();
  const { classes } = titleStyle();
  //retrieve from database
  let oldConfigurationSettings = [
    { name: 'clusterAlerts', setting: true },
    { name: 'sActionAlerts', setting: true },
    { name: 'supiciousAlerts', setting: false },
    { name: 'serverAlerts', setting: true },
    { name: 'sFileAlerts', setting: true },
    { name: 'BackupLocation', setting: 'C:\\Users\\admin' },
    { name: 'redundancy', setting: 'Low' },
  ];

  let newConfigurationSettings = [
    { name: 'clusterAlerts', setting: true },
    { name: 'sActionAlerts', setting: true },
    { name: 'supiciousAlerts', setting: false },
    { name: 'serverAlerts', setting: true },
    { name: 'sFileAlerts', setting: true },
    { name: 'BackupLocation', setting: 'C:\\Users\\admin' },
    { name: 'redundancy', setting: 'Low' },
  ];

  const [disable, setDisable] = useState(true);
  const [saveSettings,setsaveSettings] = useState(true);

  function generateRandomString() {
    var result = "";
    var characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    var charactersLength = characters.length;
    for (var i = 0; i < 16; i++) {
      result += characters.charAt(Math.floor(Math.random() *
        charactersLength));
    }

    return result;
  }

  function add(item: string) {
    if (oldConfigurationSettings[0].setting === true) {
      newConfigurationSettings[0].setting =
        !newConfigurationSettings[0].setting;
    } else {
      newConfigurationSettings[0].setting = true;
    }

    console.log('using add method');
    console.log('oldConfigurationSettings ', oldConfigurationSettings);
    console.log('newConfigurationSettings ', newConfigurationSettings);
    validate();
  }

  function remove(item: string) {
    if (newConfigurationSettings[0].setting === true) {
      newConfigurationSettings[0].setting = false;
    } else {
      newConfigurationSettings[0].setting = true;
    }
    console.log('using remove method');
    console.log('oldConfigurationSettings ', oldConfigurationSettings);
    console.log('newConfigurationSettings ', newConfigurationSettings);

    validate();
  }

  function validate() {
    if (
      oldConfigurationSettings[0].setting !==
      newConfigurationSettings[0].setting
    ) {
      setDisable(false);
    } else {
      setDisable(true);
    }
  }

  const [backupMod, setBackupMod] = useState(false);
  const [backupLocation, setBackupLocation] = useState(
    oldConfigurationSettings[5].setting.toString()
  );
  const [backupTemp, setBackupTemp] = useState(
    oldConfigurationSettings[5].setting.toString()
  );

  function backup() {
    setBackupLocation(backupTemp);
    setBackupMod(false);
  }

  const [redundancyMod, setVisible] = useState(false);
  const [redundancyLevel, setredundancyLevel] = useState(
    oldConfigurationSettings[6].setting.toString()
  );
  const [redundancyTemp, setRedundancyTemp] = useState(
    oldConfigurationSettings[6].setting.toString()
  );

  function changeRedundancy() {
    setredundancyLevel(redundancyTemp);
    setVisible(false);
  }

  function saveChanges() { }

  const [clusterAlerts, setChecked] = useState(() => {
    if (oldConfigurationSettings[0].setting === true) {
      return true;
    }
    return false;
  });
  const [sActionAlerts, setChecked1] = useState(() => {
    if (oldConfigurationSettings[1].setting === true) {
      return true;
    }

    return false;
  });
  const [supiciousAlerts, setChecked2] = useState(() => {
    if (oldConfigurationSettings[2].setting === true) {
      return true;
    }

    return false;
  });
  const [serverAlerts, setChecked3] = useState(() => {
    if (oldConfigurationSettings[3].setting === true) {
      return true;
    }

    return false;
  });
  const [sFileAlerts, setChecked4] = useState(() => {
    if (oldConfigurationSettings[4].setting === true) {
      return true;
    }

    return false;
  });

  
  let [key,setKey] =  useState("5q0L5mVB5mUlJjil");
  let [tempKey, setTempKey] = useState("");
  function generateKeys() {
    setTempKey(generateRandomString());
    setsaveSettings(false);
  }

  function saveKey() {
    tempKey==""? "":setKey(tempKey);
    
    backup();
    
    setTempKey("");
    setBackupMod(false);
  }
  function downloadKey() {

    const fileData = JSON.stringify("key:" + key);
    const blob = new Blob([fileData], { type: "text/plain" });
    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.download = "logs.txt";
    link.href = url;
    link.click();

    /* after download, delete away all the logs?
    setlogs(current =>
      current.filter(logs => {
        return null;
      }),
    );
     */
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
          <Modal
            classNames={{
              title: classes.title,
            }}
            centered
            size={400}
            opened={redundancyMod}
            title="Redundancy Level"
            onClose={() => setVisible(false)}
          >
            <div
              style={{
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'left',
              }}
            >
              <div
                style={{
                  display: 'flex',
                  flexDirection: 'row',
                }}
              >
                <RadioGroup
                  orientation="vertical"
                  label={
                    <span style={{ fontSize: '16px' }}>
                      Choose the Redundancy Level
                    </span>
                  }
                  spacing="xl"
                  required
                  onChange={(value) => setRedundancyTemp(value)}
                >
                  <Radio
                    value="High"
                    label="High"
                    checked={redundancyTemp === 'High'}
                  />
                  <Radio
                    value="Medium"
                    label="Medium"
                    checked={redundancyTemp === 'Medium'}
                  />

                  <Radio
                    value="low"
                    label="Low"
                    checked={redundancyTemp === 'Low'}
                  />
                </RadioGroup>
              </div>
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
                  }}
                  onClick={() => changeRedundancy()}
                >
                  Submit
                </Button>
              </div>
            </div>
          </Modal>
          <Modal
            centered
            size={600}
            opened={backupMod}
            title={
              <span style={{ fontSize: '22px', fontWeight: 550 }}>
                {' '}
                Backup Key
              </span>
            }
            onClose={() => [setBackupMod(false),setTempKey("")]}
          >
            <div
              style={{
                display: 'flex',
                flexDirection: 'column',
                height: '100%',
              }}
            >
              <div
                style={{
                  display: 'flex',
                  flexDirection: 'column',
                  justifyContent: 'center',
                  backgroundColor: 'white',
                }}
              >
                <caption
                  style={{
                    textAlign: 'center',
                    fontWeight: 600,
                    fontSize: '24px',
                    color: 'black',
                    marginBottom: '20px',
                    alignSelf: 'center',
                  }}
                ></caption>
                <TextInput
                  classNames={{
                    input: classes.input,
                  }}
                  label="Current Key"
                  radius="xs"
                  size="md"
                  required
                  value={key}
                  disabled={true}
           
                  rightSection={

                    <Button
                      variant="default"
                      color="dark"
                      size="md"
                      style={{ marginRight: '100px', height: '20px', width: '200px', fontSize: '10px' }}
                      onClick={() => downloadKey()}
                    >
                      Download key
                    </Button>
                  }

                />
                <TextInput
                  classNames={{
                    input: classes.input,
                  }}
                  label="New Key"
                  radius="xs"
                  size="md"
                  required
                  value={tempKey}
                  disabled={true}
                  onChange={(event) => setTempKey(event.currentTarget.value)}
                  rightSection={<Button
                    variant="default"
                    color="dark"
                    size="md"
                    style={{ marginRight: '100px', height: '20px', width: '200px', fontSize: '10px' }}

                    onClick={() => generateKeys()}
                  >
                    Generate New key
                  </Button>}




                />

              <Textarea
                  label="New Location:"
                  radius="md"
                  size="lg"
                  required
                  
                  onChange={(event) => [setBackupTemp(event.target.value),setsaveSettings(false)]}
                />
                <div style={{
                  display: 'flex',
                  flexDirection: 'column',
                  marginBottom: '20px',
                  marginTop: '20px'
                }}>
                  <Button
                    variant="default"
                    color="dark"
                    size="md"
                    style={{ marginTop: '20px',alignSelf:'flex-end' }}
                    onClick={() => saveKey()}
                    disabled= {saveSettings}
                  >
                    Save Settings
                  </Button>

                </div>
              </div>
            </div>
          </Modal>
          <Table striped>
            <caption
              style={{
                fontWeight: '600',
                fontSize: '22px',
                color: 'black',
                margin: '15px',
              }}
            >
              <span
                style={{
                  textAlign: 'center',
                }}
              >
                Notification Settings
              </span>
            </caption>

            <thead></thead>
            <tbody>
              <tr style={{}}>
                <div
                  style={{
                    height: '50px',
                    display: 'flex',
                    flexDirection: 'row',
                    justifyContent: 'space-between',
                    marginLeft: '10px',
                  }}
                >
                  <Text id="settingText"> Allow Cluster health alerts </Text>
                  <Checkbox
                    size="md"
                    style={{
                      marginRight: '50px',
                    }}
                    id="clusterAlerts"
                    checked={clusterAlerts}
                    onChange={(event) => [
                      event.currentTarget.checked
                        ? add('clusterAlerts')
                        : remove('clusterAlerts'),
                      setChecked(event.currentTarget.checked),
                    ]}
                  />
                </div>
              </tr>
              <tr>
                <div
                  style={{
                    display: 'flex',
                    height: '50px',
                    flexDirection: 'row',
                    justifyContent: 'space-between',
                    marginLeft: '10px',
                  }}
                >
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
                    ]}
                  />
                </div>
              </tr>
              <tr>
                <div
                  style={{
                    display: 'flex',
                    height: '50px',
                    flexDirection: 'row',
                    justifyContent: 'space-between',
                    marginLeft: '10px',
                  }}
                >
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
                    ]}
                  />{' '}
                </div>
              </tr>
              <tr>
                <div
                  style={{
                    display: 'flex',
                    height: '50px',
                    flexDirection: 'row',
                    justifyContent: 'space-between',
                    marginLeft: '10px',
                  }}
                >
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
                    ]}
                  />{' '}
                </div>
              </tr>
              <tr>
                <div
                  style={{
                    display: 'flex',
                    height: '50px',
                    flexDirection: 'row',
                    justifyContent: 'space-between',
                    marginLeft: '10px',
                  }}
                >
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
                    ]}
                  />
                </div>
              </tr>
              <tr>
                <div
                  style={{
                    display: 'flex',
                    height: '50px',
                    flexDirection: 'row',
                    justifyContent: 'space-between',
                    marginLeft: '10px',
                  }}
                >
                  <span>
                    <Text id="settingText"> Backup encryption key</Text>{' '}
                    <Text weight={700} style={{ marginLeft: '10px' }}>
                      Location: {backupLocation}{' '}
                    </Text>{' '}
                  </span>{' '}
                  <Button
                    style={{ marginRight: '10px', marginTop: '20px' }}
                    variant="default"
                    color="dark"
                    size="md"
                    onClick={() => setBackupMod(true)}
                  >
                    Backup
                  </Button>
                </div>
              </tr>
              <tr>
                <div
                  style={{
                    display: 'flex',
                    height: '50px',
                    flexDirection: 'row',
                    justifyContent: 'space-between',
                    marginLeft: '10px',
                  }}
                >
                  {' '}
                  <span style={{ marginTop: '10px' }}>
                    <Text id="settingText">
                      {' '}
                      Change the redundancy level of the files
                    </Text>{' '}
                    <Text weight={700} style={{ marginLeft: '10px' }}>
                      Redundancy Level: {redundancyLevel}{' '}
                    </Text>{' '}
                  </span>{' '}
                  <Button
                    style={{ marginRight: '10px', marginTop: '20px' }}
                    variant="default"
                    color="dark"
                    size="md"
                    onClick={() => setVisible(true)}
                  >
                    Change
                  </Button>
                </div>
              </tr>
            </tbody>
            <tfoot>
              <div
                style={{
                  display: 'flex',
                  flexDirection: 'column',
                }}
              >
                <Button
                  disabled={disable}
                  style={{
                    alignSelf: 'flex-end',
                    marginTop: '30px',
                    marginRight: '5px',
                  }}
                  onClick={() => saveChanges()}
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
