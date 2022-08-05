import AppBase from './components/AppBase';

import {
  CountdownCircleTimer,
  useCountdown,
} from 'react-countdown-circle-timer';
import { Button, Modal, ScrollArea, Table } from '@mantine/core';
import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useSetState } from '@mantine/hooks';
import { ChevronsDownLeft } from 'tabler-icons-react';

export interface Maintenance {
 
  maintenanceTiming: number;
  maintenanceSettings:Array<any>
}


export function AdminPerformMaintenance(props:Maintenance) {
 

  //labels
  const [timeRemaining, setTime] = useState(' Time remaining: ');
  const [percentageCompleted, setPercent] = useState('% Completed: ');
  const [elapsedTime, setElapsedTime] = useState(0);

  const timerProps = {
    isPlaying: true,
    size: 300,
    strokeWidth: 10,
  };

  //get total maintenance of time needed based on how much steps
  const time = props.maintenanceTiming;

  const minuteSeconds = time;
  const [pauseBtn, setPauseBtn] = useState(false);
  const [stopBtn, setStopBtn] = useState(false);
  const getTimeSeconds = (time: number) => (minuteSeconds - time) | 0;

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
    var end =
      seconds == 0
        ? [
            'Maintenance Completed. Gathering Results',
            setTime(''),
            setLogs(''),
            setPauseBtn(true),
            setStopBtn(true),
          ]
        : '';
    return dDisplay + hDisplay + mDisplay + sDisplay + end;
  }

  const [logs, setLogs] = useState('Maintenance Logs : ');
  const [listLogs,setListlogs]= useState([{log:''}]);

  // convert the number of seconds into day hour month and seconds
  const renderTime = (dimension: string, time: number) => {
    return (
      <div
        style={{
          display: 'flex',
          justifyContent: 'center',
          flexDirection: 'row',
        }}
      >
        {secondsToDhms(time)}
      </div>
    );
  };

  const [isActive, setIsActive] = useState(true);
  const [buttonText, setButtonText] = useState('Pause');
  const [maintenanceStatus, setMaintenanceStatus] = useState('Not Completed');
  const navigate = useNavigate();
  function pause() {
    setIsActive(!isActive);
    if (isActive == false) {
      setButtonText('Pause');
    } else {
      setButtonText('Unpause');
    }
    const vLogs = [
      {
       log:buttonText
      },
    ];
    setListlogs([...listLogs, ...vLogs])
  }
  function stop() {
    setIsActive(false);
    setMaintenanceModal(true);
    setMaintenanceStatus('Not Completed');
    const vLogs = [
      {
       log:'stop'
      },
    ];
    setListlogs([...listLogs, ...vLogs])
  }
  //modal
  const [maintenanceModal, setMaintenanceModal] = useState(false);



  const recentRows = listLogs.map((items, index) =>
    index <= elapsedTime ? (
      <tr>
        <td
          width="15%"
          style={{
            textAlign: 'left',
            fontWeight: '400',
            fontSize: '16px',
            color: 'black',
          }}
        >
          {items.log}
        </td>
      </tr>
    ) : (
      ''
    )
  );

  function getTotalTime(seconds: number) {
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

  //get the date and time of maintenance
  function getCurrentDate(separator = '') {
    let newDate = new Date();
    let date = newDate.getDate();
    let month = newDate.getMonth() + 1;
    let year = newDate.getFullYear();

    return `${year}${separator}${
      month < 10 ? `0${month}` : `${month}`
    }${separator}${date}`;
  }
  var datalogs;
  function downloadLogs() {
    datalogs = listLogs.map( function(v) {
      return Object.values( v );
   });
   
    const fileData = JSON.stringify(
      'Maintenance Status: ' +
        maintenanceStatus +
        ', Time taken: ' +
        getTotalTime(elapsedTime) +
        ', ' +
        'Date: ' +
        getCurrentDate('/') +
        ', Maintenance logs: ' +
        datalogs
    );

    const blob = new Blob([fileData], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.download = 'logs.txt';
    link.href = url;
    link.click();

    /* after download, delete away all the logs?
    setlogs(current =>
      current.filter(logs => {
        return null;
      }),
    );
     */
  }

  function setComplete() {
    setMaintenanceStatus('Completed');
    setMaintenanceModal(true);
  }

  return (
    <>
      <AppBase userType="admin">
        <Modal
          centered
          size={600}
          opened={maintenanceModal}
          title={
            <span style={{ fontSize: '22px', fontWeight: 550 }}>
              {' '}
              Maintenance Status Report
            </span>
          }
          //
          onClose={() => [setMaintenanceModal(false), navigate('/maintenance')]}
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
              <ScrollArea
                style={{
                  height: '500px',
                  width: '100%',
                  marginTop: '1%',
                }}
              >
                <Table captionSide="top" verticalSpacing="xs">
                  <thead>
                    <div style={{ fontSize: '18px' }}>
                      Date:{' '}
                      <span style={{ fontWeight: 500 }}>
                        {' '}
                        {getCurrentDate('/')}
                      </span>
                    </div>
                    <div style={{ fontSize: '18px' }}>
                      Status:{' '}
                      <span style={{ fontWeight: 500 }}>
                        {' '}
                        {maintenanceStatus}
                      </span>
                    </div>
                    <div style={{ fontSize: '18px' }}>
                      Time taken:{' '}
                      <span style={{ fontWeight: 500 }}>
                        {getTotalTime(elapsedTime)}{' '}
                      </span>
                    </div>
                    <div
                      style={{
                        fontSize: '20px',
                        fontWeight: 500,
                        textAlign: 'center',
                      }}
                    >
                      Maintenance logs:{' '}
                    </div>
                  </thead>
                  <tbody>{recentRows}</tbody>
                </Table>
              </ScrollArea>
              <Button
                variant="default"
                color="dark"
                size="md"
                style={{ alignSelf: 'flex-end' }}
                onClick={() => downloadLogs()}
              >
                Download Report
              </Button>
            </div>
          </div>
        </Modal>
        <div
          style={{
            display: 'flex',
            justifyContent: 'center',
            flexDirection: 'column',
            alignItems: 'center',
            height: '100%',
          }}
        >
          <div
            style={{
              display: 'flex',
              flexDirection: 'column',
              border: '1px solid #ccc',
              borderRadius: '10%',
              textAlign: 'left',
              width: '400px',
              backgroundColor: 'white',
            }}
          >
            <div
              style={{
                display: 'flex',
                justifyContent: 'center',
                flexDirection: 'row',
                marginTop: '20px',
              }}
            >
              <span
                style={{
                  fontWeight: 500,
                  fontSize: '22px',
                }}
              >
                Maintenance in Progress
              </span>
            </div>
            <div
              style={{
                display: 'flex',
                justifyContent: 'center',
                flexDirection: 'row',
                marginTop: '20px',
              }}
            >
              <CountdownCircleTimer
                {...timerProps}
                colors="#218380"
                duration={time}
                isPlaying={isActive}
                initialRemainingTime={time}
                onComplete={() => setComplete()}
              >
                {({ elapsedTime, color }) => (
                  <>
                    {setElapsedTime(Math.floor(elapsedTime))}

                    {setLogs(
               
                    'Maintenance Logs : ' + listLogs[listLogs.length-1].log
                    )}
               
                    <div
                      style={{
                        color,
                        display: 'flex',
                        flexDirection: 'column',
                        textAlign: 'center',
                      }}
                    >
                      <div
                        style={{
                          color,
                          display: 'flex',
                          flexDirection: 'column',
                        }}
                      >
                        {timeRemaining}
                      </div>

                      <div
                        style={{
                          color,
                          display: 'flex',
                          flexDirection: 'row',
                          marginTop: '20px',
                        }}
                      >
                        <div style={{ marginLeft: '20px' }}>
                          {renderTime('seconds', getTimeSeconds(elapsedTime))}
                        </div>
                      </div>

                      <div style={{ marginTop: '20px' }}>
                        {percentageCompleted} {((elapsedTime / time) * 100) | 0}
                      </div>
                    </div>
                  </>
                )}
              </CountdownCircleTimer>
            </div>
            <div
              style={{
                display: 'flex',
                justifyContent: 'center',
                flexDirection: 'row',
                marginTop: '20px',
                marginBottom: '20px',
              }}
            >
              {logs}
            </div>
            <div
              style={{
                display: 'flex',
                justifyContent: 'space-between',
                flexDirection: 'row',
                marginTop: '20px',
                marginBottom: '20px',
              }}
            >
              <div>
                <Button
                  variant="default"
                  color="dark"
                  size="md"
                  disabled={pauseBtn}
                  style={{ marginLeft: '50px' }}
                  onClick={() => pause()}
                >
                  {' '}
                  {buttonText}
                </Button>
              </div>
              <div>
                <Button
                  variant="default"
                  color="dark"
                  size="md"
                  disabled={stopBtn}
                  style={{ marginRight: '70px' }}
                  onClick={() => stop()}
                >
                  {' '}
                  Stop
                </Button>
              </div>
            </div>
          </div>
        </div>
      </AppBase>
    </>
  );
}
