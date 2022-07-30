import {
  Table,
  Button,
  Text,
  ScrollArea,
  Modal,
  Textarea,
} from '@mantine/core';
import { Link } from 'react-router-dom';
import { useState } from 'react';

export interface ConsoleDetails {
  // all the data needed from database
  groupList: Array<any>;
  addObjectLabel: string;
  deleteObjectLabel: string;
  tableHeader: Array<string>;
  tableBody: Array<string>;
  caption: string;
  //whether the group can be press
  pointerEvents: boolean;
}

export function AdminConsole(props: ConsoleDetails) {
  // Variable that show all the currentSSOGroups inside the props.groupList
  const [CurrentSSOGroups, setValue] = useState(props.groupList);

  // Variable that will be added to the CurrentSSOGroups
  var [Group, addGroup] = useState('');

  // Variable that will decide whether the submit button is disabled
  var [submitBtn, setSubmitBtn] = useState(true);

  // Variable that will decide whether the addUserModel visibility is true or false
  const [openedAddUserModel, setOpened] = useState(false);

  // Variable that is bind to each specific labels
  const title = 'Add ' + props.addObjectLabel;
  const textField = 'Name of the ' + props.addObjectLabel;

  // Variable that will decide whether the errorMessage will be displayed
  var [errorMessage, setErrorMessage] = useState('');

  // Delete away specific group from CurrentSSOGroups
  const deleteGroup = (index: any) => {
    setValue(CurrentSSOGroups.filter((v, i) => i !== index));
  };

  // Add the group to the CurrentSSOGroups and set the addUserModel visibility to false
  function add() {
    setValue(CurrentSSOGroups.concat(Group));
    setOpened(false);
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
      submitBtn = false;
      setSubmitBtn(false);
    }
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
        <div style={{ marginLeft: '10px' }}>{props.tableHeader}</div>
      </th>
    </tr>
  );

  // display all the rows that is from props
  const rows = CurrentSSOGroups.map((items, index) => (
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
            pointerEvents: props.pointerEvents ? 'auto' : 'none',
          }}
          component={Link}
          to="/insidessogroup"
          variant="link"
        >
          {items}
        </Text>
        <Button
          variant="default"
          color="dark"
          size="md"
          style={{ marginRight: '15px' }}
          onClick={() => [deleteGroup(index)]}
        >
          Delete
        </Button>
      </td>
    </tr>
  ));

  return (
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
          opened={openedAddUserModel}
          onClose={() => setOpened(false)}
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
                name="password"
                id="password"
                required
                error={errorMessage}
                onChange={(event) => {
                  addGroup(event.target.value),
                    (Group = event.target.value),
                    validate();
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
              {props.caption}{' '}
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
            Add {props.addObjectLabel}
          </Button>
        </div>
      </div>
    </div>
  );
}
