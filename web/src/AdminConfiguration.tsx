import { Button, Textarea, Checkbox, Text } from '@mantine/core';
import { useState } from 'react';
import AppBase from './components/AppBase';

export function AdminConfiguration() {
  //function will be rotate key
  function rotateKey() {}

  var [location, addLocation] = useState('');
  var [errorMessage, setErrorMessage] = useState('');
  var [rotateButton,setButton] = useState(true);
  function validate() {
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
      errorMessage = 'do not include special characters';
      setErrorMessage('do not include special characters');
      rotateButton = true;
      setButton(true);
    } else if (location.includes(' ')) {
      errorMessage = 'No space is allowed';
      setErrorMessage('No space is allowed');
      rotateButton = true;
      setButton(true);
    } else if (location == '') {
      errorMessage = 'Details needed';
      setErrorMessage('Details needed');
      rotateButton = true;
      setButton(true);
    } else {
      errorMessage = '';
      setErrorMessage('');
      rotateButton = false;
      setButton(false);
    }
  }
  return (
    <AppBase userType="admin">
      <div
        style={{
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'flex-start',
          height: '100%',
        }}
      >
        <div
          style={{
            display: 'flex',
            border: '1px solid #ccc',
            flexDirection: 'column',
            justifyContent: 'center',
            alignItems: 'flex-start',
            width: '90%',
            backgroundColor: 'white',
            borderRadius: '10px',
            padding: '20px',
            maxWidth: '500px',
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
          >
            Rotate Key
          </caption>
          <Textarea
            label="Specify the file/directory location and the system will
            auto rotate the key"
            radius="md"
            size="lg"
            error= {errorMessage}
            onChange={(event) => {
              addLocation(event.target.value),
                (location = event.target.value),
                validate();
            }}
          />
          <div
            style={{
              display: 'flex',
              flexDirection: 'row',
              margin: '20px 0',
            }}
          >
            <Text>Master Key :</Text>
            <Checkbox style={{ marginLeft: '10px' }}> </Checkbox>
          </div>
          <Button
            variant="default"
            color="dark"
            size="md"
            disabled={rotateButton}
            style={{ alignSelf: 'flex-end' }}
            onClick={() => rotateKey()}
          >
            Rotate Key
          </Button>
        </div>
      </div>
    </AppBase>
  );
}
