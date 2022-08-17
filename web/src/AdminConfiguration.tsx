import { Button, Textarea, Checkbox, Text, Divider } from '@mantine/core';
import { useState } from 'react';
import AppBase from './components/AppBase';

export function AdminConfiguration() {
  //function will be rotate key
  function rotateKey() {}

  var [location, addLocation] = useState('');
  var [errorMessage, setErrorMessage] = useState('');
  var [rotateButton, setButton] = useState(true);
  const allowedChar = /^[A-Za-z0-9\s]*$/;
  const space = /^\s*$/;
  function validate() {
    //if the location is blank
    if (space.test(location) == true) {
      errorMessage = 'Do not leave blank';
      setErrorMessage('Do not leave blank');
      rotateButton = true;
      setButton(true);
    } else if (
      //if the location contains other than letter and digit
      allowedChar.test(location) == false
    ) {
      errorMessage = 'do not include special characters';
      setErrorMessage('do not include special characters');
      rotateButton = true;
      setButton(true);
    } else {
      // all test are pass
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
            alignItems: 'center',
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
            Settings
          </caption>
          <div className="flex flex-col">
            <Divider
              my="xs"
              label="Rotate Key"
              variant="dotted"
              labelPosition="center"
            />
            <Textarea
              label="Specify the file/directory location and the system will
            auto rotate the key"
              radius="md"
              size="lg"
              error={errorMessage}
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

          <div className="flex flex-col w-full mt-5">
            <Divider
              my="xs"
              label="Set Redundancy Level"
              variant="dotted"
              labelPosition="center"
            />
            <Text>test</Text>
          </div>
        </div>
      </div>
    </AppBase>
  );
}
