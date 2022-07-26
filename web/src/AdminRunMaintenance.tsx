import {
  useMantineTheme,
  Checkbox,
  Button,
  Table,
  ActionIcon,
  Modal,
} from '@mantine/core';
import { Link } from 'react-router-dom';
import { useState } from 'react';
import { Settings } from 'tabler-icons-react';
import AppBase from './components/AppBase';

export function AdminRunMaintenance() {
  const theme = useMantineTheme();
  let MaintenanceSettings = [
    { name: 'CrawlPermissions', setting: 'true' },
    { name: 'PurgeOrphanedFile', setting: 'true' },
    { name: 'PurgeUser', setting: 'false' },
    { name: 'CrawlReplicas', setting: 'true' },
    { name: 'QuickCheck', setting: 'true' },
    { name: 'FullCheck', setting: 'false' },
    { name: 'DBCheck', setting: 'true' },
  ];

  const [openedMaintenanceSettingsModal, setOpened] = useState(false);

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
                  component={Link}
                  to="/performmaintenance"
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
