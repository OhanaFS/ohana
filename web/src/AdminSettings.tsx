import {
  Checkbox,
  Button,
  Text,
  Table,
  useMantineTheme,
  Modal,
  Radio,
  createStyles,
  TextInput,
  Textarea,
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
  },
}));

export function AdminSettings() {
  const theme = useMantineTheme();
  const { classes } = titleStyle();
  //retrieve from database
  let oldConfigurationSettings = [
    { name: 'clusterAlerts', setting: false },
    { name: 'sActionAlerts', setting: true },
    { name: 'supiciousAlerts', setting: false },
    { name: 'serverAlerts', setting: true },
    { name: 'sFileAlerts', setting: true },
    { name: 'BackupLocation', setting: 'C:\\Users\\admin' },
    { name: 'redundancy', setting: 'Low' },
  ];

  let newConfigurationSettings = [
    { name: 'clusterAlerts', setting: oldConfigurationSettings[0].setting },
    { name: 'sActionAlerts', setting: oldConfigurationSettings[1].setting },
    { name: 'supiciousAlerts', setting: oldConfigurationSettings[2].setting },
    { name: 'serverAlerts', setting: oldConfigurationSettings[3].setting },
    { name: 'sFileAlerts', setting: oldConfigurationSettings[4].setting },
    { name: 'BackupLocation', setting: oldConfigurationSettings[5].setting },
    { name: 'redundancy', setting: oldConfigurationSettings[6].setting },
  ];

  let [key, setKey] = useState('5q0L5mVB5mUlJjil');

  const [disable, setDisable] = useState(true);
  const [saveSettings, setsaveSettings] = useState(true);

  function generateRandomString() {
    var result = '';
    var characters =
      'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    var charactersLength = characters.length;
    for (var i = 0; i < 16; i++) {
      result += characters.charAt(Math.floor(Math.random() * charactersLength));
    }

    return result;
  }

  const [backupMod, setBackupMod] = useState(false);
  var [backupLocation, setBackupLocation] = useState(
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
  var [redundancyLevel, setredundancyLevel] = useState(
    oldConfigurationSettings[6].setting.toString()
  );
  const [redundancyTemp, setRedundancyTemp] = useState(
    oldConfigurationSettings[6].setting.toString()
  );

  function changeRedundancy() {
    redundancyLevel = redundancyTemp;
    newConfigurationSettings[6].setting = redundancyLevel;
    setredundancyLevel(redundancyTemp);
    validateRedundancy();
    console.log(settings);
    setVisible(false);
  }

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

  function changeClusterAlert() {
    setChecked((clusterAlerts) => !clusterAlerts);

    //detect that the first state is true but the end result is false
    if (clusterAlerts == true) {
      newConfigurationSettings[0].setting = false;
      validateClusterAlert();
    } else {
      //detect that the first state is false but the end result is true
      newConfigurationSettings[0].setting = true;
      validateClusterAlert();
    }
  }

  function changeActionAlerts() {
    setChecked1((sActionAlerts) => !sActionAlerts);
    //detect that the first state is true but the end result is false
    if (sActionAlerts == true) {
      newConfigurationSettings[1].setting = false;
      validateActionAlerts();
    } else {
      //detect that the first state is false but the end result is true
      newConfigurationSettings[1].setting = true;
      validateActionAlerts();
    }
  }

  function changeSupiciousAlerts() {
    setChecked2((supiciousAlerts) => !supiciousAlerts);
    //detect that the first state is true but the end result is false
    if (supiciousAlerts == true) {
      newConfigurationSettings[2].setting = false;
      validateSupiciousAlerts();
    } else {
      //detect that the first state is false but the end result is true
      newConfigurationSettings[2].setting = true;
      validateSupiciousAlerts();
    }
  }

  function changeServerAlerts() {
    setChecked3((serverAlerts) => !serverAlerts);
    //detect that the first state is true but the end result is false
    if (serverAlerts == true) {
      newConfigurationSettings[3].setting = false;
      validateServerAlerts();
    } else {
      //detect that the first state is false but the end result is true
      newConfigurationSettings[3].setting = true;
      validateServerAlerts();
    }
  }

  function changesFileAlerts() {
    setChecked4((sFileAlerts) => !sFileAlerts);
    //detect that the first state is true but the end result is false
    if (sFileAlerts == true) {
      newConfigurationSettings[4].setting = false;
      validateFileAlerts();
    } else {
      //detect that the first state is false but the end result is true
      newConfigurationSettings[4].setting = true;
      validateFileAlerts();
    }
  }
  /* there is 5 checkbox and 2 button, so each setting will bind to each item and if there is changes, 
     the useState will be false which enable the save button 
  */
  var [settingsA, setSettingsA] = useState(true);
  var [settingsB, setSettingsB] = useState(true);
  var [settingsC, setSettingsC] = useState(true);
  var [settingsD, setSettingsD] = useState(true);
  var [settingsE, setSettingsE] = useState(true);
  var [settingsF, setSettingsF] = useState(true);
  var [settingsG, setSettingsG] = useState(true);
  var settings = [
    settingsA,
    settingsB,
    settingsC,
    settingsD,
    settingsE,
    settingsF,
    settingsG,
  ];

  /* all the validation methods is to check if there is any changes, if there is changes 
     then each setting will be change to false



  */
  function validateClusterAlert() {
    if (
      oldConfigurationSettings[0].setting !==
      newConfigurationSettings[0].setting
    ) {
      if (settingsA == true) {
        settingsA = false;
        settings[0] = false;
        setSettingsA(false);
      }
      if (settings.indexOf(false) > -1) {
        setDisable(false);
      } else {
        setDisable(true);
      }
    } else {
      if (settingsA == false) {
        settingsA = true;
        settings[0] = true;
        setSettingsA(true);
      }

      if (settings.indexOf(false) > -1) {
        setDisable(false);
      } else {
        setDisable(true);
      }
    }
  }

  function validateActionAlerts() {
    if (
      oldConfigurationSettings[1].setting !==
      newConfigurationSettings[1].setting
    ) {
      if (settingsB == true) {
        settingsB = false;
        settings[1] = false;
        setSettingsB(false);
      }
      if (settings.indexOf(false) > -1) {
        setDisable(false);
      } else {
        setDisable(true);
      }
    } else {
      if (settingsB == false) {
        settingsB = true;
        settings[1] = true;
        setSettingsB(true);
      }
      if (settings.indexOf(false) > -1) {
        setDisable(false);
      } else {
        setDisable(true);
      }
    }
  }

  function validateSupiciousAlerts() {
    if (
      oldConfigurationSettings[2].setting !==
      newConfigurationSettings[2].setting
    ) {
      if (settingsC == true) {
        settingsC = false;
        settings[2] = false;
        setSettingsC(false);
      }
      if (settings.indexOf(false) > -1) {
        setDisable(false);
      } else {
        setDisable(true);
      }
    } else {
      if (settingsC == false) {
        settingsC = true;
        settings[2] = true;
        setSettingsC(true);
      }
      if (settings.indexOf(false) > -1) {
        setDisable(false);
      } else {
        setDisable(true);
      }
      console.log(settings);
    }
  }

  function validateServerAlerts() {
    if (
      oldConfigurationSettings[3].setting !==
      newConfigurationSettings[3].setting
    ) {
      if (settingsD == true) {
        settingsD = false;
        settings[3] = false;
        setSettingsD(false);
      }
      if (settings.indexOf(false) > -1) {
        setDisable(false);
      } else {
        setDisable(true);
      }
      console.log(settings);
    } else {
      if (settingsD == false) {
        settingsD = true;
        settings[3] = true;
        setSettingsD(true);
      }
      if (settings.indexOf(false) > -1) {
        setDisable(false);
      } else {
        setDisable(true);
      }
    }
  }

  function validateFileAlerts() {
    if (
      oldConfigurationSettings[4].setting !==
      newConfigurationSettings[4].setting
    ) {
      if (settingsE == true) {
        settingsE = false;
        settings[4] = false;
        setSettingsE(false);
      }
      if (settings.indexOf(false) > -1) {
        setDisable(false);
      } else {
        setDisable(true);
      }
    } else {
      if (settingsE == false) {
        settingsE = true;
        settings[4] = true;
        setSettingsE(true);
      }
      if (settings.indexOf(false) > -1) {
        setDisable(false);
      } else {
        setDisable(true);
      }
    }
  }
  function validateKey() {
    if (
      oldConfigurationSettings[5].setting !==
      newConfigurationSettings[5].setting
    ) {
      if (settingsF == true) {
        settingsF = false;
        settings[5] = false;
        setSettingsF(false);
      }
      if (settings.indexOf(false) > -1) {
        setDisable(false);
      } else {
        setDisable(true);
      }
      console.log(settings);
    } else {
      if (settingsF == false) {
        settingsF = true;
        settings[5] = true;
        setSettingsF(true);
      }
      if (settings.indexOf(false) > -1) {
        setDisable(false);
      } else {
        setDisable(true);
      }
    }
  }
  function validateRedundancy() {
    if (
      oldConfigurationSettings[6].setting !==
      newConfigurationSettings[6].setting
    ) {
      if (settingsG == true) {
        settingsG = false;
        settings[6] = false;
        setSettingsG(false);
      }
      if (settings.indexOf(false) > -1) {
        setDisable(false);
      } else {
        setDisable(true);
      }
      console.log(settings);
    } else {
      if (settingsG == false) {
        settingsG = true;
        settings[6] = true;
        setSettingsG(true);
      }
      if (settings.indexOf(false) > -1) {
        setDisable(false);
      } else {
        setDisable(true);
      }
    }
  }

  let [tempKey, setTempKey] = useState('');
  function generateKeys() {
    setTempKey(generateRandomString());
    setsaveSettings(false);
  }

  function saveKey() {
    tempKey == '' ? '' : setKey(tempKey);
    backup();
    setsaveSettings(true);
    backupLocation = backupTemp;
    newConfigurationSettings[5].setting = backupLocation;
    validateKey();
    setTempKey('');
    setBackupMod(false);
  }
  function downloadKey() {
    const fileData = JSON.stringify('key:' + key);
    const blob = new Blob([fileData], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.download = 'logs.txt';
    link.href = url;
    link.click();
  }

  // function to save all the pending changes
  function saveChanges() {
    setDisable(true);
  }

  return (
    <AppBase userType="admin">
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
                <Radio.Group
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
                    value="Low"
                    label="Low"
                    checked={redundancyTemp === 'Low'}
                  />
                </Radio.Group>
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
            onClose={() => [
              setBackupMod(false),
              setTempKey(''),
              setsaveSettings(true),
            ]}
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
                      style={{
                        marginRight: '100px',
                        height: '20px',
                        width: '200px',
                        fontSize: '10px',
                      }}
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
                  rightSection={
                    <Button
                      variant="default"
                      color="dark"
                      size="md"
                      style={{
                        marginRight: '100px',
                        height: '20px',
                        width: '200px',
                        fontSize: '10px',
                      }}
                      onClick={() => generateKeys()}
                    >
                      Generate New key
                    </Button>
                  }
                />

                <Textarea
                  label="New Location:"
                  radius="md"
                  size="lg"
                  required
                  onChange={(event) => [
                    setBackupTemp(event.target.value),
                    setsaveSettings(false),
                  ]}
                />
                <div
                  style={{
                    display: 'flex',
                    flexDirection: 'column',
                    marginBottom: '20px',
                    marginTop: '20px',
                  }}
                >
                  <Button
                    variant="default"
                    color="dark"
                    size="md"
                    style={{ marginTop: '20px', alignSelf: 'flex-end' }}
                    onClick={() => saveKey()}
                    disabled={saveSettings}
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
                    checked={clusterAlerts}
                    onChange={() => changeClusterAlert()}
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
                    onChange={changeActionAlerts}
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
                    onChange={changeSupiciousAlerts}
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
                    onChange={changeServerAlerts}
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
                    onChange={changesFileAlerts}
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
