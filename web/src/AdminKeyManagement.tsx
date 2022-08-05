import {
  Button,
  Modal,
  Textarea,
  ScrollArea,
  Text,
  Table,
  Checkbox,
  Radio,
} from '@mantine/core';
import { useState } from 'react';
import { Link } from 'react-router-dom';
import AppBase from './components/AppBase';

export interface Key {
  key: string;
  permissions: Array<any>;
  type: string;
  location: string;
}
export function AdminKeyManagement() {
  const data = [
    {
      id: '128c1d5d-2359-4ba1-8739-2cd30d694d67',
      permissions: ['Read', 'Write', 'Execute'],
      type: 'path',
      location: 'asd',
    },

    {
      id: '128c1eee-2359-4ba1-8739-2cd30d69sds67',
      permissions: ['Read', 'Write'],
      type: 'filename prefix',
      location: 'asd',
    },
  ];
  // Variable that show all the ApiKey
  const [ApiKey, setValue] = useState(data);

  // Variable that will be added to the key structure
  var [key, addKey] = useState('');
  var [permissions, addPermissions] = useState<string[]>([]);
  var [type, addType] = useState('');
  var [location, addLocation] = useState('');
  var [index, addIndex] = useState(Number);

  // Variable that will decide whether the submit button is disabled
  var [submitBtn, setSubmitBtn] = useState(true);

  // Variable that will decide whether the modal visibility is true or false
  const [openedAddKeyModel, setOpened] = useState(false);
  const [openedUpdateKeyModel, setOpened2] = useState(false);

  // Variable that is bind to each specific labels
  const title = 'Add Key';
  const textField = 'Name of the Key';

  // Variable that will decide the value of the errorMessage that will be displayed
  var [errorMessage, setErrorMessage] = useState('');
  var [errorMessage2, setErrorMessage2] = useState('');
  var [errorMessage3, setErrorMessage3] = useState('');
  var [errorMessage4, setErrorMessage4] = useState('');

  // Delete away specific key from apikey
  const deleteKey = (index: any) => {
    setValue(ApiKey.filter((v, i) => i !== index));
  };
  const updateKey = (index: any) => {
    addKey(ApiKey[index].id);
    addPermissions(ApiKey[index].permissions);
    addType(ApiKey[index].type);
    addLocation(ApiKey[index].location);
    addIndex(index);
    setOpened2(true);
  };

  function resetKey() {
    addKey('');
    addPermissions([]);
    addType('');
    addLocation('');
    addIndex(Number);
    setErrorMessage('');
    setErrorMessage2('');
    setErrorMessage3('');
    setErrorMessage4('');
  }
  // Add the key to the Apikey and set the addKeyModel visibility to false
  function add() {
    const vkey = [
      {
        id: key,
        permissions: permissions,
        type: type,
        location: location,
      },
    ];
    setValue([...ApiKey, ...vkey]);
    setOpened(false);
  }

  function update() {
    const vkey = [
      {
        id: key,
        permissions: permissions,
        type: type,
        location: location,
      },
    ];
    ApiKey[index].id = vkey[0].id;
    ApiKey[index].permissions = vkey[0].permissions;
    ApiKey[index].type = vkey[0].type;
    ApiKey[index].location = vkey[0].location;
    setOpened2(false);
  }
  /* Validate the textfield to check if there is any special characters
     if there is special character, the function will display error message 
     and set the submit button to false.  */
  function validate() {
    if (
      key.includes('/') ||
      key.includes('[') ||
      key.includes('!') ||
      key.includes('@') ||
      key.includes('#') ||
      key.includes('$') ||
      key.includes('%') ||
      key.includes('^') ||
      key.includes('&') ||
      key.includes('*') ||
      key.includes('(') ||
      key.includes(')') ||
      key.includes('\\') ||
      key.includes('=') ||
      key.includes('[') ||
      key.includes(']') ||
      key.includes(';') ||
      key.includes(',') ||
      key.includes('.') ||
      key.includes('<') ||
      key.includes('>') ||
      key.includes('?') ||
      key.includes('`')
    ) {
      errorMessage = 'do not include special characters';
      setErrorMessage('do not include special characters');
      submitBtn = true;
      setSubmitBtn(true);
    } else if (key.includes(' ')) {
      errorMessage = 'No space is allowed';
      setErrorMessage('No space is allowed');
      submitBtn = true;
      setSubmitBtn(true);
    } else if (key == '') {
      errorMessage = 'Details needed';
      setErrorMessage('Details needed');
      submitBtn = true;
      setSubmitBtn(true);
    } else {
      errorMessage = '';
      setErrorMessage('');
    }
    if (
      location.includes('/') ||
      location.includes('[') ||
      location.includes('!') ||
      location.includes('@') ||
      location.includes('#') ||
      location.includes('$') ||
      location.includes('%') ||
      location.includes('^') ||
      location.includes('&') ||
      location.includes('*') ||
      location.includes('(') ||
      location.includes(')') ||
      location.includes('\\') ||
      location.includes('=') ||
      location.includes('[') ||
      location.includes(']') ||
      location.includes(';') ||
      location.includes(',') ||
      location.includes('.') ||
      location.includes('<') ||
      location.includes('>') ||
      location.includes('?') ||
      location.includes('`')
    ) {
      errorMessage2 = 'do not include special characters';
      setErrorMessage2('do not include special characters');
      submitBtn = true;
      setSubmitBtn(true);
    } else if (location.includes(' ')) {
      errorMessage2 = 'No space is allowed';
      setErrorMessage2('No space is allowed');
      submitBtn = true;
      setSubmitBtn(true);
    } else if (location == '') {
      errorMessage2 = 'Details needed';
      setErrorMessage2('Details needed');
      submitBtn = true;
      setSubmitBtn(true);
    } else {
      errorMessage2 = '';
      setErrorMessage2('');
    }
    if (type == '') {
      errorMessage3 = 'Choose the type';
      setErrorMessage3('Choose the type');
    } else {
      errorMessage3 = '';
      setErrorMessage3('');
    }
    if (permissions.length === 0) {
      errorMessage4 = 'Choose the permissions';
      setErrorMessage4('Choose the permissions');
    } else {
      errorMessage4 = '';
      setErrorMessage4('');
    }

    if (
      errorMessage == '' &&
      errorMessage2 == '' &&
      errorMessage3 == '' &&
      errorMessage4 == ''
    ) {
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
        <div style={{ marginLeft: '10px' }}>Key ID</div>
      </th>
    </tr>
  );

  // display all the rows that is from props
  const rows = ApiKey.map((items, index) => (
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
          }}
        >
          {items['id']}
        </Text>
        <div>
          <Button
            variant="default"
            color="dark"
            size="md"
            style={{ marginRight: '10px' }}
            onClick={() => [updateKey(index)]}
          >
            Update
          </Button>
          <Button
            variant="default"
            color="dark"
            size="md"
            style={{}}
            onClick={() => [deleteKey(index)]}
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
            height: '50vh',
            justifyContent: 'center',
          }}
        >
          <div className="console">
            <Modal
              centered
              title="Add console"
              opened={openedAddKeyModel}
              onClose={() => [setOpened(false), resetKey()]}
              size="20%"
            >
              {
                <div
                  style={{
                    display: 'flex',
                    height: '50vh',
                    flexDirection: 'column',
                  }}
                >
                  <Textarea
                    placeholder={textField}
                    label={title}
                    defaultValue={key}
                    size="md"
                    name="password"
                    id="password"
                    required
                    error={errorMessage}
                    onChange={(event) => {
                      addKey(event.target.value),
                        (key = event.target.value),
                        validate();
                    }}
                  />
                  <Radio.Group
                    orientation="vertical"
                    label={
                      <span style={{ fontSize: '16px' }}>
                        Choose type of key
                      </span>
                    }
                    error={errorMessage3}
                    spacing="sm"
                    required
                    value={type}
                    onChange={(value) => (
                      addType(value), (type = value), validate()
                    )}
                  >
                    <Radio value="filenameprefix" label="Filename Prefix" />
                    <Radio value="path" label="Path" />
                  </Radio.Group>
                  <Checkbox.Group
                    label={
                      <span style={{ fontSize: '16px' }}>
                        Choose the Permissions
                      </span>
                    }
                    spacing="xl"
                    required
                    value={permissions}
                    error={errorMessage4}
                    onChange={(value) => {
                      addPermissions(value), (permissions = value), validate();
                    }}
                  >
                    <Checkbox value="Read" label="Read" />
                    <Checkbox value="Write" label="Write" />
                    <Checkbox value="Execute" label="Execute" />
                    <Checkbox value="Share" label="Share" />
                    <Checkbox value="Audit" label="Audit" />
                  </Checkbox.Group>

                  <Textarea
                    label="Location:"
                    radius="md"
                    size="lg"
                    error={errorMessage2}
                    defaultValue={location}
                    required
                    onChange={(event) => [
                      addLocation(event.target.value),
                      (location = event.target.value),
                      validate(),
                    ]}
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
              title="Update console"
              opened={openedUpdateKeyModel}
              onClose={() => [setOpened2(false), resetKey()]}
              size="20%"
            >
              {
                <div
                  style={{
                    display: 'flex',
                    height: '45vh',
                    flexDirection: 'column',
                  }}
                >
                  <Textarea
                    placeholder={textField}
                    label="Key"
                    defaultValue={key}
                    size="md"
                    name="password"
                    id="password"
                    required
                    error={errorMessage}
                    onChange={(event) => {
                      addKey(event.target.value),
                        (key = event.target.value),
                        validate();
                    }}
                  />
                  <Radio.Group
                    orientation="vertical"
                    label={
                      <span style={{ fontSize: '16px' }}>Type of key</span>
                    }
                    error={errorMessage3}
                    spacing="sm"
                    required
                    value={type}
                    onChange={(value) => (
                      addType(value), (type = value), validate()
                    )}
                  >
                    <Radio value="filename prefix" label="filename prefix" />
                    <Radio value="path" label="path" />
                  </Radio.Group>
                  <Checkbox.Group
                    label={
                      <span style={{ fontSize: '16px' }}>Permissions</span>
                    }
                    spacing="xl"
                    required
                    value={permissions}
                    error={errorMessage4}
                    onChange={(value) => {
                      addPermissions(value), (permissions = value), validate();
                    }}
                  >
                    <Checkbox value="Read" label="Read" />
                    <Checkbox value="Write" label="Write" />
                    <Checkbox value="Execute" label="Execute" />
                    <Checkbox value="Share" label="Share" />
                    <Checkbox value="Audit" label="Audit" />
                  </Checkbox.Group>

                  <Textarea
                    label="Location:"
                    radius="md"
                    size="lg"
                    error={errorMessage2}
                    defaultValue={location}
                    required
                    onChange={(event) => [
                      addLocation(event.target.value),
                      (location = event.target.value),
                      validate(),
                    ]}
                  />
                  <Button
                    variant="default"
                    color="dark"
                    size="md"
                    onClick={() => update()}
                    style={{
                      marginLeft: '15px',
                      alignSelf: 'flex-end',
                      marginTop: '20px',
                    }}
                  >
                    Update
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
                  {'API Key Management Console'}{' '}
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
                Add Key
              </Button>
            </div>
          </div>
        </div>
      </AppBase>
    </>
  );
}
