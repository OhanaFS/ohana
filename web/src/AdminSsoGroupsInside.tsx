import {
  Modal,
  Textarea,
  Button,
  ScrollArea,
  Table,
  Text,
  TextInput,
} from '@mantine/core';
import { useState } from 'react';
import { Link, useLocation } from 'react-router-dom';
import AppBase from './components/AppBase';

export function AdminSsoGroupsInside() {
  const state = useLocation();
  const userList: any = state.state;
  // Variable that show all the currentSSOGroups inside the props.groupList
  const [CurrentSSOGroups, setValue] = useState(userList.users);

  // Variable that will be added to the CurrentSSOGroups
  var [Group, addGroup] = useState('');
  var [limit, addLimit] = useState(Number);
  var [index, addIndex] = useState(Number);

  // Variable that will decide whether the submit button is disabled
  var [submitBtn, setSubmitBtn] = useState(true);

  // Variable that will decide whether the addUserModel visibility is true or false
  const [openedAddUserModel, setOpened] = useState(false);
  const [openedUpdateUserModel, setOpened2] = useState(false);

  // Variable that is bind to each specific labels
  const title = 'Add User ';
  const textField = 'Name of the User';

  // Variable that will decide whether the errorMessage will be displayed
  var [errorMessage, setErrorMessage] = useState('');
  var [errorMessage2, setErrorMessage2] = useState('');

  // Delete away specific group from CurrentSSOGroups
  const deleteUser = (index: any) => {
    setValue(CurrentSSOGroups.filter((v: any, i: any) => i !== index));
  };

  // Add the group to the CurrentSSOGroups and set the addUserModel visibility to false
  function add() {
    const vgroup = [
      {
        name: Group,
        limit: limit,
      },
    ];
    setValue([...CurrentSSOGroups, ...vgroup]);
    setOpened(false);
  }
  function updateUser(index: any) {
    addGroup(CurrentSSOGroups[index].name);
    addLimit(CurrentSSOGroups[index].limit);
    addIndex(index);
    setOpened2(true);
  }
  /* Validate the textfield to check if there is any special characters
    if there is special character, the function will display error message 
    and set the submit button to false.  */
  function validate() {
    if (
      Group.includes('/') ||
      Group.includes('[') ||
      Group.includes('!') ||
      Group.includes('@') ||
      Group.includes('#') ||
      Group.includes('$') ||
      Group.includes('%') ||
      Group.includes('^') ||
      Group.includes('&') ||
      Group.includes('*') ||
      Group.includes('(') ||
      Group.includes(')') ||
      Group.includes('\\') ||
      Group.includes('=') ||
      Group.includes('[') ||
      Group.includes(']') ||
      Group.includes(';') ||
      Group.includes(',') ||
      Group.includes('.') ||
      Group.includes('<') ||
      Group.includes('>') ||
      Group.includes('?') ||
      Group.includes('`')
    ) {
      errorMessage = 'do not include special characters';
      setErrorMessage('do not include special characters');
      submitBtn = true;
      setSubmitBtn(true);
    } else if (Group.includes(' ')) {
      errorMessage = 'No space is allowed';
      setErrorMessage('No space is allowed');
      submitBtn = true;
      setSubmitBtn(true);
    } else if (Group == '') {
      errorMessage = 'Details needed';
      setErrorMessage('Details needed');
      submitBtn = true;
      setSubmitBtn(true);
    } else {
      errorMessage = '';
      setErrorMessage('');
    }
    if (isNaN(limit) == false && limit !== 0) {
      errorMessage2 = '';
      setErrorMessage2('');
    } else if (limit === 0) {
      errorMessage2 = 'limit cannot be blank and  must be more than 0';
      setErrorMessage2('limit cannot be blank and must be more than 0');
    } else {
      errorMessage2 = 'Enter number only';
      setErrorMessage2('Enter number only');
      submitBtn = true;

      setSubmitBtn(true);
    }
    if (errorMessage === '' && errorMessage2 === '') {
      submitBtn = false;
      setSubmitBtn(false);
    }
  }
  function reset() {
    Group = '';
    limit = 0;
    index = 0;
  }
  function updateLimit() {
    const vkey = [
      {
        limit: limit,
      },
    ];
    CurrentSSOGroups[index].limit = vkey[0].limit;
    setOpened2(false);
  }
  // display table header that is from props
  const ths = (
    <tr>
      <th
        style={{
          textAlign: 'left',
          fontWeight: '700',
          fontSize: '16px',
          color: 'black',
        }}
      >
        <div style={{ marginLeft: '10px' }}>
          List of Users inside this group
        </div>
      </th>
    </tr>
  );

  // display all the rows that is from props
  const rows = CurrentSSOGroups.map((items: any, index: any) => (
    <tr key={index}>
      <td
        style={{
          display: 'flex',
          justifyContent: 'space-between',
        }}
      >
        <Text
          color="dark"
          style={{
            marginLeft: '10px',
            pointerEvents: false ? 'auto' : 'none',
          }}
          component={Link}
          to="/insidessogroup"
          variant="link"
        >
          {items.name}
        </Text>
        <div>
          <Button
            variant="default"
            color="dark"
            size="md"
            style={{ marginRight: '15px' }}
            onClick={() => [updateUser(index)]}
          >
            Update
          </Button>
          <Button
            variant="default"
            color="dark"
            size="md"
            style={{ marginRight: '15px' }}
            onClick={() => [deleteUser(index)]}
          >
            Delete
          </Button>
        </div>
      </td>
    </tr>
  ));
  return (
    <>
      <AppBase userType="admin">
        <div
          style={{
            display: 'flex',
            height: '80vh',
            justifyContent: 'center',
          }}
        >
          <div className="console">
            <Modal
              centered
              title="Add User console"
              opened={openedAddUserModel}
              onClose={() => [setOpened(false), reset()]}
            >
              {
                <div
                  style={{
                    display: 'flex',
                    height: '25vh',
                    flexDirection: 'column',
                  }}
                >
                  <Textarea
                    placeholder={textField}
                    label={title}
                    size="md"
                    required
                    error={errorMessage}
                    onChange={(event) => {
                      addGroup(event.target.value),
                        (Group = event.target.value),
                        validate();
                    }}
                  />
                  <TextInput
                    placeholder="quotas limits for each person"
                    label="Quota"
                    size="md"
                    error={errorMessage2}
                    required
                    onChange={(event) => {
                      [
                        addLimit(Number(event.target.value)),
                        (limit = Number(event.target.value)),
                        validate(),
                      ];
                    }}
                  />
                  <Button
                    variant="default"
                    color="dark"
                    size="md"
                    disabled={submitBtn}
                    onClick={() => add()}
                    style={{
                      marginLeft: '15px',
                      alignSelf: 'flex-end',
                      marginTop: '20px',
                    }}
                  >
                    Submit
                  </Button>
                </div>
              }
            </Modal>
            <Modal
              centered
              title="Update User console"
              opened={openedUpdateUserModel}
              onClose={() => [setOpened2(false), reset()]}
            >
              {
                <div
                  style={{
                    display: 'flex',
                    height: '10vh',
                    flexDirection: 'column',
                  }}
                >
                  <TextInput
                    placeholder="quotas limits for each person"
                    label="Quota"
                    size="md"
                    error={errorMessage2}
                    value={limit}
                    required
                    onChange={(event) => {
                      [
                        addLimit(Number(event.target.value)),
                        (limit = Number(event.target.value)),
                        validate(),
                      ];
                    }}
                  />
                  <Button
                    variant="default"
                    color="dark"
                    size="md"
                    disabled={submitBtn}
                    onClick={() => updateLimit()}
                    style={{
                      marginLeft: '15px',
                      alignSelf: 'flex-end',
                      marginTop: '20px',
                    }}
                  >
                    Submit
                  </Button>
                </div>
              }
            </Modal>
            <ScrollArea
              style={{
                height: '90%',
                width: '100%',
                marginTop: '1%',
              }}
            >
              <Table captionSide="top" verticalSpacing="sm">
                <caption
                  style={{
                    textAlign: 'center',
                    fontWeight: '600',
                    fontSize: '24px',
                    color: 'black',
                    marginTop: '5px',
                  }}
                >
                  User Management Console
                </caption>
                <thead style={{}}>{ths}</thead>
                <tbody>{rows}</tbody>
              </Table>
            </ScrollArea>

            <div
              style={{
                display: 'flex',
                flexDirection: 'row',
                justifyContent: 'space-between',
              }}
            >
              <Button
                variant="default"
                color="dark"
                size="md"
                onClick={() => setOpened(true)}
                style={{ marginLeft: '15px' }}
              >
                Add User
              </Button>
            </div>
          </div>
        </div>
        );
      </AppBase>
    </>
  );
}
