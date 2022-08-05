import {
  useMantineTheme,
  Checkbox,
  Button,
  Table,
  Text,
  ActionIcon,
  Modal,
} from '@mantine/core';
import { Link, useNavigate } from 'react-router-dom';
import { useState } from 'react';
import { Settings } from 'tabler-icons-react';
import AppBase from './components/AppBase';

export function AdminRunMaintenance() {
  const theme = useMantineTheme();

  /* get all the settings from database and amount of time needed for each settings


  




  */
  let MaintenanceSettings = [
    { name: 'CrawlPermissions', setting: 'true', time: 10 },
    { name: 'PurgeOrphanedFile', setting: 'true', time: 10 },
    { name: 'PurgeUser', setting: 'false', time: 10 },
    { name: 'CrawlReplicas', setting: 'true', time: 10 },
    { name: 'QuickCheck', setting: 'true', time: 10 },
    { name: 'FullCheck', setting: 'false', time: 10 },
    { name: 'DBCheck', setting: 'true', time: 10 },
  ];

  // Maintenance settings that will be run
  var [firstCheck, setFirstCheck] = useState(() => {
    if (MaintenanceSettings[0].setting === 'true') {
      return true;
    }

    return false;
  });

  var [secondCheck, setSecondCheck] = useState(() => {
    if (MaintenanceSettings[1].setting === 'true') {
      return true;
    }

    return false;
  });
  var [thirdCheck, setThirdCheck] = useState(() => {
    if (MaintenanceSettings[2].setting === 'true') {
      return true;
    }

    return false;
  });
  var [fourthCheck, setFourthCheck] = useState(() => {
    if (MaintenanceSettings[3].setting === 'true') {
      return true;
    }

    return false;
  });
  var [fifthCheck, setFifthCheck] = useState(() => {
    if (MaintenanceSettings[4].setting === 'true') {
      return true;
    }

    return false;
  });
  var [sixthCheck, setSixthCheck] = useState(() => {
    if (MaintenanceSettings[5].setting === 'true') {
      return true;
    }

    return false;
  });
  var [seventhCheck, setSeventhCheck] = useState(() => {
    if (MaintenanceSettings[6].setting === 'true') {
      return true;
    }

    return false;
  });

  // Variable that will decide whether the openedMaintenanceSettingsModal visibility is true or false
  const [openedMaintenanceSettingsModal, setOpened] = useState(false);
  const [openConfirmationModal, setOpened2] = useState(false);

  // maintenance settings
  var [sFirstCheck, setsFirstCheck] = useState(() => {
    if (MaintenanceSettings[0].setting === 'true') {
      return true;
    }
    return false;
  });

  var [sSecondCheck, setsSecondCheck] = useState(() => {
    if (MaintenanceSettings[1].setting === 'true') {
      return true;
    }

    return false;
  });
  var [sThirdCheck, setsThirdCheck] = useState(() => {
    if (MaintenanceSettings[2].setting === 'true') {
      return true;
    }

    return false;
  });
  var [sFourthCheck, setsFourthCheck] = useState(() => {
    if (MaintenanceSettings[3].setting === 'true') {
      return true;
    }

    return false;
  });
  var [sFifthCheck, setsFifthCheck] = useState(() => {
    if (MaintenanceSettings[4].setting === 'true') {
      return true;
    }

    return false;
  });
  var [sSixthCheck, setsSixthCheck] = useState(() => {
    if (MaintenanceSettings[5].setting === 'true') {
      return true;
    }

    return false;
  });
  var [sSeventhCheck, setsSeventhCheck] = useState(() => {
    if (MaintenanceSettings[6].setting === 'true') {
      return true;
    }

    return false;
  });
  var settings = [[], [], [], [], [], [], []];
  // error message
  var [errorMessage, setErrorMessage] = useState('');

  // save the settings and set the maintenance settings modal visibility to false
  function saveSettings() {
    setFirstCheck(sFirstCheck);
    setSecondCheck(sSecondCheck);
    setThirdCheck(sThirdCheck);
    setFourthCheck(sFourthCheck);
    setFifthCheck(sFifthCheck);
    setSixthCheck(sSixthCheck);
    setSeventhCheck(sSeventhCheck);
    setOpened(false);
  }

  var [totalTimeNeeded,setTotalTimeNeed] = useState(0);
  var [timeTimeNeedInStr, setTimeTimeNeedInStr] = useState('');
  function calculateAmountOfTime() {
    if (firstCheck == true) {
      totalTimeNeeded = totalTimeNeeded + MaintenanceSettings[0].time;
    }
    if (secondCheck == true) {
      totalTimeNeeded = totalTimeNeeded + MaintenanceSettings[1].time;
    }
    if (thirdCheck == true) {
      totalTimeNeeded = totalTimeNeeded + MaintenanceSettings[2].time;
    }
    if (fourthCheck == true) {
      totalTimeNeeded = totalTimeNeeded + MaintenanceSettings[3].time;
    }
    if (fifthCheck == true) {
      totalTimeNeeded = totalTimeNeeded + MaintenanceSettings[4].time;
    }
    if (sixthCheck == true) {
      totalTimeNeeded = totalTimeNeeded + MaintenanceSettings[5].time;
    }
    if (seventhCheck == true) {
      totalTimeNeeded = totalTimeNeeded + MaintenanceSettings[6].time;
    }
    if (totalTimeNeeded == 0) {
      setErrorMessage('you need to tick at least one checkbox');
    } else {
      setOpened2(true);
      setErrorMessage('');
      timeTimeNeedInStr = secondsToDhms(totalTimeNeeded);
      setTotalTimeNeed(totalTimeNeeded);
      setTimeTimeNeedInStr(secondsToDhms(totalTimeNeeded));
    }
  }
  var date = getCurrentDate('/');
  function getCurrentDate(separator = '') {
    let newDate = new Date();
    let date = newDate.getDate();
    let month = newDate.getMonth() + 1;
    let year = newDate.getFullYear();

    return `${date}${separator}${
      month < 10 ? `0${month}` : `${month}`
    }${separator}${year}`;
  }
  function secondsToDhms(seconds: number) {
    seconds = Number(seconds);
    var d = Math.floor(seconds / (3600 * 24));
    var h = Math.floor((seconds % (3600 * 24)) / 3600);
    var m = Math.floor((seconds % 3600) / 60);
    var s = Math.floor(seconds % 60);

    var dDisplay = d > 0 ? d + (d == 1 ? ' day, ' : ' days, ') : '';
    var hDisplay = h > 0 ? h + (h == 1 ? ' hour, ' : ' hrs, ') : '';
    var mDisplay = m > 0 ? m + (m == 1 ? ' minute, ' : ' mins, ') : '';
    var sDisplay = s > 0 ? s + (s == 1 ? ' second' : ' secs') : '';
    return dDisplay + hDisplay + mDisplay + sDisplay;
  }
  const navigate = useNavigate();

  return (
    <>
      <AppBase userType="admin">
        <div
          style={{
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'flex-start',
          }}
        >
          <Modal
            centered
            opened={openConfirmationModal}
            size={600}
            title={'Maintenance Confirmation'}
            styles={{
              title: {
                fontSize: '22px',
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
              },
            }}
            onClose={() => setOpened2(false)}
          >
            <div
              style={{
                display: 'flex',
                flexDirection: 'row',
                fontSize: '22px',
                justifyContent: 'space-evenly',
                backgroundColor: 'white',
              }}
            >
              Date : {date}
            </div>
            <div
              style={{
                display: 'flex',
                flexDirection: 'row',
                fontSize: '22px',
                justifyContent: 'space-evenly',
                backgroundColor: 'white',
              }}
            >
              Expected Time needed : {timeTimeNeedInStr}
            </div>
            <div
              style={{
                display: 'flex',
                flexDirection: 'row',
                justifyContent: 'space-evenly',
                backgroundColor: 'white',
              }}
            >
              <Button
                style={{
                  alignSelf: 'flex-end',
                  marginTop: '20px',
                }}
                variant="default"
                color="dark"
                size="md"
                onClick={() => setOpened2(false)}
              >
                Cancel
              </Button>
              <Button
                style={{
                  alignSelf: 'flex-end',
                  marginTop: '20px',
                }}
                variant="default"
                color="dark"
                size="md"
                onClick={()=>navigate('/performmaintenance',{state:{time:totalTimeNeeded}})}
                  
                
              >
                Confirm
              </Button>
            </div>
          </Modal>
          <Modal
            centered
            opened={openedMaintenanceSettingsModal}
            size={600}
            title={'Maintenance Settings'}
            styles={{
              title: {
                fontSize: '22px',
                fontWeight: 550,
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
              },
            }}
            onClose={() => setOpened(false)}
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
                <Table striped verticalSpacing="xs" id="maintenanceSettings">
                  <tbody>
                    <tr>
                      <td>
                        Crawl the list of files to remove permissions from
                        expired users
                      </td>
                      <td>
                        <Checkbox
                          size="md"
                          checked={sFirstCheck}
                          onChange={(event) =>
                            setsFirstCheck(event.currentTarget.checked)
                          }
                        />{' '}
                      </td>
                    </tr>

                    <tr>
                      <td>Purging orphaned files and shards</td>
                      <td>
                        <Checkbox
                          size="md"
                          checked={sSecondCheck}
                          onChange={(event) =>
                            setsSecondCheck(event.currentTarget.checked)
                          }
                        />
                      </td>
                    </tr>

                    <tr>
                      <td>Purge a user and their files</td>
                      <td>
                        {' '}
                        <Checkbox
                          size="md"
                          checked={sThirdCheck}
                          onChange={(event) =>
                            setsThirdCheck(event.currentTarget.checked)
                          }
                        />{' '}
                      </td>
                    </tr>

                    <tr>
                      <td>
                        {' '}
                        Crawl all of the files to make sure it has full replicas
                      </td>
                      <td>
                        {' '}
                        <Checkbox
                          size="md"
                          checked={sFourthCheck}
                          onChange={(event) =>
                            setsFourthCheck(event.currentTarget.checked)
                          }
                        />
                      </td>
                    </tr>

                    <tr>
                      <td>
                        Quick File Check (Only checks current versions of files
                        to see if it’s fine and is not corrupted){' '}
                      </td>
                      <td>
                        <Checkbox
                          size="md"
                          checked={sFifthCheck}
                          onChange={(event) => [
                            setsFifthCheck(event.currentTarget.checked),
                            {
                              sSixthCheck: true ? setsSixthCheck(false) : '',
                            },
                          ]}
                        />
                      </td>
                    </tr>

                    <tr>
                      <td>
                        Full File Check (Checks all fragments to ensure that
                        it’s not corrupted)
                      </td>
                      <td>
                        <Checkbox
                          size="md"
                          checked={sSixthCheck}
                          onChange={(event) => [
                            setsSixthCheck(event.currentTarget.checked),
                            {
                              sFifthCheck: true ? setsFifthCheck(false) : '',
                            },
                          ]}
                        />
                      </td>
                    </tr>

                    <tr>
                      <td>DB integrity Check</td>
                      <td>
                        <Checkbox
                          size="md"
                          checked={sSeventhCheck}
                          onChange={(event) =>
                            setsSeventhCheck(event.currentTarget.checked)
                          }
                        />
                      </td>
                    </tr>
                  </tbody>

                  <div
                    style={{
                      display: 'flex',
                      flexDirection: 'column',
                    }}
                  >
                    <Button
                      style={{
                        alignSelf: 'flex-end',
                        marginTop: '20px',
                        marginRight: '10px',
                      }}
                      variant="default"
                      color="dark"
                      size="md"
                      onClick={() => saveSettings()}
                    >
                      Save Settings
                    </Button>
                  </div>
                </Table>
              </div>
            </div>
          </Modal>
          <div className="maintenanceSettings">
            <Table striped verticalSpacing="xs" id="maintenanceSettings">
              <caption>
                {' '}
                <div style={{ marginTop: '10px' }}>
                  {' '}
                  Run Maintenance
                  <span
                    style={{
                      marginTop: '2px',
                      marginRight: '10px',
                      float: 'right',
                    }}
                  >
                    <ActionIcon onClick={() => setOpened(true)}>
                      <Settings></Settings>
                    </ActionIcon>
                  </span>
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
                      checked={firstCheck}
                      onChange={(event) =>
                        setFirstCheck(event.currentTarget.checked)
                      }
                    />{' '}
                  </td>
                </tr>

                <tr>
                  <td>Purging orphaned files and shards</td>
                  <td>
                    <Checkbox
                      size="md"
                      checked={secondCheck}
                      onChange={(event) =>
                        setSecondCheck(event.currentTarget.checked)
                      }
                    />
                  </td>
                </tr>

                <tr>
                  <td>Purge a user and their files</td>
                  <td>
                    {' '}
                    <Checkbox
                      size="md"
                      checked={thirdCheck}
                      onChange={(event) =>
                        setThirdCheck(event.currentTarget.checked)
                      }
                    />{' '}
                  </td>
                </tr>

                <tr>
                  <td>
                    {' '}
                    Crawl all of the files to make sure it has full replicas
                  </td>
                  <td>
                    {' '}
                    <Checkbox
                      size="md"
                      checked={fourthCheck}
                      onChange={(event) =>
                        setFourthCheck(event.currentTarget.checked)
                      }
                    />
                  </td>
                </tr>

                <tr>
                  <td>
                    {' '}
                    Quick File Check (Only checks current versions of files to
                    see if it’s fine and is not corrupted){' '}
                  </td>
                  <td>
                    {' '}
                    <Checkbox
                      size="md"
                      checked={fifthCheck}
                      onChange={(event) => [
                        setFifthCheck(event.currentTarget.checked),
                        {
                          sixthCheck: true ? setSixthCheck(false) : '',
                        },
                      ]}
                    />{' '}
                  </td>
                </tr>

                <tr>
                  <td>
                    {' '}
                    Full File Check (Checks all fragments to ensure that it’s
                    not corrupted){' '}
                  </td>
                  <td>
                    {' '}
                    <Checkbox
                      size="md"
                      checked={sixthCheck}
                      onChange={(event) => [
                        setSixthCheck(event.currentTarget.checked),
                        {
                          fifthCheck: true ? setFifthCheck(false) : '',
                        },
                      ]}
                    />
                  </td>
                </tr>

                <tr>
                  <td> DB integrity Check</td>
                  <td>
                    {' '}
                    <Checkbox
                      size="md"
                      checked={seventhCheck}
                      onChange={(event) =>
                        setSeventhCheck(event.currentTarget.checked)
                      }
                    />
                  </td>
                </tr>
                <tr>
                  <td>
                    <Text style={{ color: 'red' }}>{errorMessage}</Text>
                  </td>
                </tr>
              </tbody>
              <div
                style={{
                  display: 'flex',
                  flexDirection: 'column',
                }}
              >
                <Button
                  style={{
                    alignSelf: 'flex-end',
                    marginTop: '20px',
                  }}
                  variant="default"
                  color="dark"
                  size="md"
                  onClick={() => [calculateAmountOfTime()]}
                >
                  Perform Maintenance
                </Button>
              </div>
            </Table>
          </div>
        </div>
      </AppBase>
    </>
  );
}
