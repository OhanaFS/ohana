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
import { group } from 'console';

export interface ConsoleDetails {
  groupList: Array<any>;
  addObjectLabel: string;
  deleteObjectLabel: string;
  tableHeader: Array<string>;
  tableBody: Array<string>;
  caption: string;
  pointerEvents: boolean;
}

export function AdminConsole(props: ConsoleDetails) {
  const [CurrentSSOGroups, setValue] = useState(props.groupList);
  var [Group, addGroup] = useState('');
  function add() {
    setValue(CurrentSSOGroups.concat(Group));
    setOpened(false);
  }
  var [submitBtn,setSubmitBtn]=useState(true);
  const [openedAddUserModel, setOpened] = useState(false);
  const title = 'Add ' + props.addObjectLabel;
  const textField = 'Name of the ' + props.addObjectLabel;
  var [errorMessage, setErrorMessage] = useState('');
  const deleteGroup = (index: any) => {
    setValue(CurrentSSOGroups.filter((v, i) => i !== index));
  };

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
      Group.includes('\\')||
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
      submitBtn=true;
      setSubmitBtn(true);
    } 
    else if(  Group.includes(' ')){
      errorMessage = 'No space is allowed';
      setErrorMessage('No space is allowed');
      submitBtn=true;
      setSubmitBtn(true);
    }
    else if(  Group==""){
      errorMessage = 'Details needed';
      setErrorMessage('Details needed');
      submitBtn=true;
      setSubmitBtn(true);
    }
    else {
      errorMessage = '';
      setErrorMessage('');
      submitBtn=false;
      setSubmitBtn(false);
    }
  }
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
