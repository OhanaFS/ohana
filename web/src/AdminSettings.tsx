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

// style the specific label
const titleStyle = createStyles(() => ({
  // style the title for redundancy level modal
  title: {
    fontWeight: 600,
    fontSize: '24px',
  },

  //style the text inside current Key (backup key modal)
  input: {
    fontWeight: 600,
    fontSize: '24px',
  },
}));

/*
 *  so for the useState right,i notice that if u use the method to change the variable data, it will not 
 *  immediately change the data but wait until the next time u access the variable then it will show u the updated data.
 *  so no choice i just change the data straight from variable instead of using the method to change the data.
 *  then if u dont use the method to change the state but change straigth from the variable, next time u access the variable, it will use
 *  back the previous state.
 *  so i change the state of variable and also the variable data.
*/


export function AdminSettings() {
  // variable that will be binded to
  const { classes } = titleStyle();

  const data = [false, true, false, true, false, 'C:\\Users\\admin', 'Low'];

  // for some reason, you need to use state change to change the data inside the oldConfiguration
  var [oclusterAlerts, setClusterAlerts] = useState(data[0]);
  var [oActionAlerts, setoActionAlerts] = useState(data[1]);
  var [oSupiciousAlerts, setoSupiciousAlerts] = useState(data[2]);
  var [oServerAlerts, setoServerAlerts] = useState(data[3]);
  var [oFileAlerts, setoFileAlerts] = useState(data[4]);
  var [oBackupLocation, setoBackupLocation] = useState(data[5].toString());
  var [oredundancy, setoredundancy] = useState(data[6].toString());

  //retrieve from database
  var oldConfigurationSettings = [
    { name: 'clusterAlerts', setting: oclusterAlerts },
    { name: 'sActionAlerts', setting: oActionAlerts },
    { name: 'supiciousAlerts', setting: oSupiciousAlerts },
    { name: 'serverAlerts', setting: oServerAlerts },
    { name: 'sFileAlerts', setting: oFileAlerts },
    { name: 'BackupLocation', setting: oBackupLocation },
    { name: 'redundancy', setting: oredundancy },
  ];

  // bind the default newConfigurationSettings to oldConfigurationSettings
  var newConfigurationSettings = [
    { name: 'clusterAlerts', setting: oldConfigurationSettings[0].setting },
    { name: 'sActionAlerts', setting: oldConfigurationSettings[1].setting },
    { name: 'supiciousAlerts', setting: oldConfigurationSettings[2].setting },
    { name: 'serverAlerts', setting: oldConfigurationSettings[3].setting },
    { name: 'sFileAlerts', setting: oldConfigurationSettings[4].setting },
    {
      name: 'BackupLocation',
      setting: oldConfigurationSettings[5].setting.toString(),
    },
    {
      name: 'redundancy',
      setting: oldConfigurationSettings[6].setting.toString(),
    },
  ];

  // key variable
  var [key, setKey] = useState('5q0L5mVB5mUlJjil');

  // Variable that will decide whether the save pending changes button is disabled
  const [disable, setDisable] = useState(true);

  //  Variable that will decide whether  save settings button is disabled
  const [saveSettings, setsaveSettings] = useState(true);

  // function that will create a 16 letter key
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

  // Variable that will decide whether the backupModal visibility is true or false
  const [backupMod, setBackupMod] = useState(false);

  // Variable that is bind to the backup key location
  var [backupLocation, setBackupLocation] = useState(
    oldConfigurationSettings[5].setting.toString()
  );

  // temp location for backuplocation as only when the user save the settings, then it save the changes
  const [backupTemp, setBackupTemp] = useState(
    oldConfigurationSettings[5].setting.toString()
  );

  /* after the user save the location, it will set the temp backup variable to the actual backup variable and 
     sets the backupModal visibility to false*/
  function backup() {
    setBackupLocation(backupTemp);
    setBackupMod(false);
  }

  // Variable that will decide whether the redundancyMod visibility is true or false
  const [redundancyMod, setVisible] = useState(false);

  // Variable that is bind to the redundancy Level
  var [redundancyLevel, setredundancyLevel] = useState(
    oldConfigurationSettings[6].setting.toString()
  );

  // temp location for redundancy Level as only when the user save the redundancy Level, then it save the changes
  var [redundancyTemp, setRedundancyTemp] = useState(
    oldConfigurationSettings[6].setting.toString()
  );

  /* after the user save the redundancy Level, it will set the temp redundancy Level variable to the actual redundancy Level variable and 
     sets the redundancyMod visibility to false*/
  function changeRedundancy() {
    redundancyLevel = redundancyTemp;
    setredundancyLevel(redundancyTemp);
    newConfigurationSettings[6].setting = redundancyLevel;
    validateRedundancy();
    setVisible(false);
  }

  // validate and save the key (works the same as how the checkbox is validated)
  function saveKey() {
    tempKey == '' ? '' : setKey(tempKey),key=tempKey;
    backup();
    setsaveSettings(true);
    backupLocation = backupTemp;
    setBackupLocation(backupTemp);
    newConfigurationSettings[5].setting=backupTemp;
    validateKey();
    //reset temp key to default
    setTempKey('');
    setBackupMod(false);
  }

  // each of these variable is binded, so if the data retrieve from database is true, the checkbox will be ticked
  var [clusterAlerts, setChecked] = useState(() => {
    if (oldConfigurationSettings[0].setting === true) {
      return true;
    }

    return false;
  });

  var [sActionAlerts, setChecked1] = useState(() => {
    if (oldConfigurationSettings[1].setting === true) {
      return true;
    }

    return false;
  });
  var [supiciousAlerts, setChecked2] = useState(() => {
    if (oldConfigurationSettings[2].setting === true) {
      return true;
    }

    return false;
  });
  var [serverAlerts, setChecked3] = useState(() => {
    if (oldConfigurationSettings[3].setting === true) {
      return true;
    }

    return false;
  });
  var [sFileAlerts, setChecked4] = useState(() => {
    if (oldConfigurationSettings[4].setting === true) {
      return true;
    }

    return false;
  });

  // this function will check if any changes is made to the first checkbox
  function changeClusterAlert() {
    // set whether it is ticked or not, so if the checkbox is ticked, it will untick the checkbox
    clusterAlerts = !clusterAlerts;
    setChecked((clusterAlerts) => !clusterAlerts);

    if (clusterAlerts == true) {
      newConfigurationSettings[0].setting = true;
      validateClusterAlert();
    } else {
      newConfigurationSettings[0].setting = false;
      validateClusterAlert();
    }
  }

  // this function will check if any changes is made to the second checkbox
  function changeActionAlerts() {
    // set whether it is ticked or not, so if the checkbox is ticked, it will untick the checkbox
    sActionAlerts = !sActionAlerts;
    setChecked1((sActionAlerts) => !sActionAlerts);
    //detect that the first state is true but the end result is false
    if (sActionAlerts == true) {
      newConfigurationSettings[1].setting = true;
      validateActionAlerts();
    } else {
      //detect that the first state is false but the end result is true
      newConfigurationSettings[1].setting = false;
      validateActionAlerts();
    }
  }

  // this function will check if any changes is made to the third checkbox
  function changeSupiciousAlerts() {
    // set whether it is ticked or not, so if the checkbox is ticked, it will untick the checkbox
    supiciousAlerts = !supiciousAlerts;
    setChecked2((supiciousAlerts) => !supiciousAlerts);
    //detect that the first state is true but the end result is false
    if (supiciousAlerts == true) {
      newConfigurationSettings[2].setting = true;
      validateSupiciousAlerts();
    } else {
      //detect that the first state is false but the end result is true
      newConfigurationSettings[2].setting = false;
      validateSupiciousAlerts();
    }
  }

  // this function will check if any changes is made to the fourth checkbox
  function changeServerAlerts() {
    // set whether it is ticked or not, so if the checkbox is ticked, it will untick the checkbox
    serverAlerts = !serverAlerts;
    setChecked3((serverAlerts) => !serverAlerts);
    //detect that the first state is true but the end result is false
    if (serverAlerts == true) {
      newConfigurationSettings[3].setting = true;
      validateServerAlerts();
    } else {
      //detect that the first state is false but the end result is true
      newConfigurationSettings[3].setting = false;
      validateServerAlerts();
    }
  }

  // this function will check if any changes is made to the last checkbox
  function changesFileAlerts() {
    // set whether it is ticked or not, so if the checkbox is ticked, it will untick the checkbox
    sFileAlerts = !sFileAlerts;
    setChecked4((sFileAlerts) => !sFileAlerts);
    //detect that the first state is true but the end result is false
    if (sFileAlerts == true) {
      newConfigurationSettings[4].setting = true;
      validateFileAlerts();
    } else {
      //detect that the first state is false but the end result is true
      newConfigurationSettings[4].setting = false;
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

  // variable that is bind to New key textfield
  let [tempKey, setTempKey] = useState('');
  // generate the key and show on the textfield, and enable the save button
  function generateKeys() {
    setTempKey(generateRandomString());
    setsaveSettings(false);
  }


  //download the key
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
    if (settingsA == false) {
      newConfigurationSettings[0].setting =
        !newConfigurationSettings[0].setting;

      oldConfigurationSettings[0].setting = newConfigurationSettings[0].setting;
      settingsA = true;
      setSettingsA(true);
      settings[0] = true;
      settings = [true, true, true, true, true, true, true];

      if (oldConfigurationSettings[0].setting === true) {
        setClusterAlerts(true);
        oclusterAlerts = true;
        setChecked(true);
      } else {
        setClusterAlerts(false);
        oclusterAlerts = false;
        setChecked(false);
      }
    }
    if (settingsB == false) {
      newConfigurationSettings[1].setting =
        !newConfigurationSettings[1].setting;

      oldConfigurationSettings[1].setting = newConfigurationSettings[1].setting;
      settingsB = true;
      setSettingsB(true);
      settings[1] = true;
      settings = [true, true, true, true, true, true, true];

      if (oldConfigurationSettings[1].setting === true) {
        setoActionAlerts(true);
        oActionAlerts = true;
        setChecked1(true);
      } else {
        setoActionAlerts(false);
        oActionAlerts = false;
        setChecked1(false);
      }
    }
    if (settingsC == false) {
      newConfigurationSettings[2].setting =
        !newConfigurationSettings[2].setting;

      oldConfigurationSettings[2].setting = newConfigurationSettings[2].setting;
      settingsC = true;
      setSettingsC(true);
      settings[2] = true;
      settings = [true, true, true, true, true, true, true];

      if (oldConfigurationSettings[2].setting === true) {
        setoSupiciousAlerts(true);
        oSupiciousAlerts = true;
        setChecked2(true);
      } else {
        setoSupiciousAlerts(false);
        oSupiciousAlerts = false;
        setChecked2(false);
      }
    }
    if (settingsD == false) {
      newConfigurationSettings[3].setting =
        !newConfigurationSettings[3].setting;

      oldConfigurationSettings[3].setting = newConfigurationSettings[3].setting;
      settingsD = true;
      setSettingsD(true);
      settings[3] = true;
      settings = [true, true, true, true, true, true, true];

      if (oldConfigurationSettings[3].setting === true) {
        setoServerAlerts(true);
        oServerAlerts = true;
        setChecked3(true);
      } else {
        setoServerAlerts(false);
        oServerAlerts = false;
        setChecked3(false);
      }
    }

    if (settingsE == false) {
      newConfigurationSettings[4].setting =
        !newConfigurationSettings[4].setting;

      oldConfigurationSettings[4].setting = newConfigurationSettings[4].setting;
      settingsE = true;
      setSettingsE(true);
      settings[4] = true;
      settings = [true, true, true, true, true, true, true];

      if (oldConfigurationSettings[4].setting === true) {
        setoFileAlerts(true);
        oFileAlerts = true;
        setChecked4(true);
      } else {
        setoFileAlerts(false);
        oFileAlerts = false;
        setChecked4(false);
      }
    }

    //backup
    if (settingsF == false) {
    
      oBackupLocation=backupLocation;
      setoBackupLocation(backupLocation);
      oldConfigurationSettings[5].setting=backupLocation;
      newConfigurationSettings[5].setting=backupLocation;
      settingsF = true;
      setSettingsF(true);
      settings[5] = true;
      settings = [true, true, true, true, true, true, true];

    }

    // redunacncy
    if (settingsG == false) {
      oredundancy=redundancyLevel;
      setoredundancy(redundancyLevel);
      oldConfigurationSettings[6].setting=redundancyLevel;
      newConfigurationSettings[6].setting=redundancyLevel;
      settingsG = true;
      setSettingsG(true);
      settings[5] = true;
      settings = [true, true, true, true, true, true, true];

    }
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
                  value={redundancyTemp}
                  onChange={(value) => (setRedundancyTemp(value),redundancyTemp="Low")}
                >
                  <Radio
                    value="High"
                    label="High"
               
                  />
                  <Radio
                    value="Medium"
                    label="Medium"
                
                  />

                  <Radio
                    value="Low"
                    label="Low"
                
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
