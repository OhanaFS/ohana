
import AppBase from './components/AppBase';

import { CountdownCircleTimer, useCountdown } from 'react-countdown-circle-timer'
import { Button } from '@mantine/core';
import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
export function AdminPerformMaintenance() {

  const timerProps = {
    isPlaying: true,
    size: 300,
    strokeWidth: 10
  };
  
  const renderTime = (dimension:string, time:number) => {
    return (
      <div >
        <div >{time}</div>
        <div>{dimension}</div>
      </div>
    );
  };

  const [isActive, setIsActive] = useState(true);
  const [buttonText, setButtonText] = useState('Pause');
  const navigate = useNavigate();
  function pause(){

    setIsActive(!isActive);
    if(isActive==false){
      setButtonText('Pause');
    }
    else{
      setButtonText('Unpause');
    }
    
  }
  function stop(){

    setIsActive(!isActive);
    navigate('/maintenanceresults');
  }



  

  const minuteSeconds = 60;
  const hourSeconds = 3600;
  const daySeconds = 86400;
  
  const getTimeSeconds = (time:number) => (minuteSeconds - time) | 0;
  const getTimeMinutes = (time:number) => ((time % hourSeconds) / minuteSeconds) | 0;
  const getTimeHours = (time:number) => ((time % daySeconds) / hourSeconds) | 0;
  const getTimeDays = (time:number) => (time / daySeconds) | 0;
 


  const stratTime = Date.now() / 1000; // use UNIX timestamp in seconds

  //this is where the total time 
  const endTime = stratTime + 9020; // use UNIX timestamp in seconds
  const remainingTime = endTime - stratTime;
  const days = Math.ceil(remainingTime / daySeconds);
  
  return (
    <>  

      <AppBase
        userType="admin"
        name="Alex Simmons"
        username="@alex"
        image="https://images.unsplash.com/photo-1496302662116-35cc4f36df92?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=2070&q=80"
      >
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
          border : '1px solid #ccc',
          borderRadius: '10%',
          textAlign : 'left',
          width: '400px',
          backgroundColor:'white'
          
        }}>
          <div style={{
              display: 'flex',
              justifyContent: 'center',
              flexDirection: 'row',
            
            }}><h3>Maintenance in Progress</h3></div>
          <div
            style={{
              display: 'flex',
              justifyContent: 'center',
              flexDirection: 'row',
            
            }}>
      <CountdownCircleTimer
        {...timerProps}
      
        colors="#218380"
        duration={8}
        isPlaying={isActive}
        initialRemainingTime={8}
        onComplete={() => stop()}
      >
        {({ elapsedTime, color }) => (
          <>
          <div style={{ color, display:'flex',flexDirection: 'column',textAlign:'center'
         }}>
           <div style={{color, display:'flex',flexDirection: 'column'}}>
            Time remaining:
          </div>

          <div style={{ color, display:'flex',flexDirection: 'row',marginTop:'20px'
         }}><div>
            {renderTime("hours", getTimeHours(daySeconds - elapsedTime))}
            </div>
            <div style={{marginLeft:'20px'}}>
            {renderTime("minutes", getTimeMinutes(hourSeconds - elapsedTime))}
            </div>
            <div style={{marginLeft:'20px'}} >
            {renderTime("seconds", getTimeSeconds(elapsedTime))}
            
            </div>
          </div>
         <div style={{marginTop:'20px'}}>
         % Completed: {elapsedTime/8 *100 | 0 }
         </div>
          </div>
          </>
        )}
      </CountdownCircleTimer>
      </div>
      <div style={{
              display: 'flex',
              flexDirection: 'row',
              justifyContent: 'space-between',
              marginTop:'20px',
              marginBottom:'20px'
            }}>
              <div>
         <Button  
            variant="default"
            color="dark"
            size="md"
            style={{marginLeft:'50px'}}
            onClick={() => pause()}> {buttonText}
        </Button>
        </div>
        <div>
        <Button  
            variant="default"
            color="dark"
            size="md"
            style={{marginRight:'70px'}}
            onClick={()=>stop()}> Stop
        </Button>
        </div>
        </div>
        </div>
        </div>
      </AppBase>
    </>
  );
}
