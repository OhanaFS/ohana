
import AppBase from './components/AppBase';

import { CountdownCircleTimer, useCountdown } from 'react-countdown-circle-timer'
import { Button } from '@mantine/core';
import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
export function AdminPerformMaintenance() {

  const steps = [
  'Cleaning Server 1',
  'Finish Cleaning Server 1' ,
  'Added Server 2 ', 
  'Finish Cleaning Server 2',
  'Added Server 3 ', 
  'Finish Cleaning Server 3',
  'Added Server 4 ', 
  'Finish Cleaning Server 4',
  'Added Server 5 ', 
  'Finish Cleaning Server 5',
  'Added Server 6 ', 
  'Finish Cleaning Server 6',
  'Added Server 7 ', 
  'Finish Cleaning Server 7',
  ];

  const [elapsedTime,setElapsedTime] = useState(0);
 
  const timerProps = {
    isPlaying: true,
    size: 300,
    strokeWidth: 10
  };

  //get the maintenance of time needed based on how much steps
  const time = steps.length;
  const minuteSeconds = time;
  const getTimeSeconds = (time:number) => (minuteSeconds - time) | 0;

  function secondsToDhms(seconds: number) {
    seconds = Number(seconds);
    var d = Math.floor(seconds / (3600 * 24));
    var h = Math.floor(seconds % (3600 * 24) / 3600);
    var m = Math.floor(seconds % 3600 / 60);
    var s = Math.floor(seconds % 60);
    var dDisplay = d > 0 ? d + (d == 1 ? " day, " : " days, ") : "";
    var hDisplay = h > 0 ? h + (h == 1 ? " hour, " : " hrs, ") : "";
    var mDisplay = m > 0 ? m + (m == 1 ? " minute, " : " mins, ") : "";
    var sDisplay = s > 0 ? s + (s == 1 ? " second" : " secs") : "";
    return dDisplay + hDisplay + mDisplay + sDisplay;
  }
  
// convert the number of seconds into day hour month and seconds
  const renderTime = (dimension:string, time:number) => {
    return (
      <div >
        <div>{secondsToDhms(time)}</div>
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
              marginTop:'20px'
              
            }}> 
            <span style={{
                fontWeight:500,
                fontSize:'22px'
            }}>
              Maintenance in Progress
            </span>
            
            </div>
          <div
            style={{
              display: 'flex',
              justifyContent: 'center',
              flexDirection: 'row',
              marginTop:'20px'
            }}>
      <CountdownCircleTimer
        {...timerProps}
      
        colors="#218380"
        duration={time}
        isPlaying={isActive}
        initialRemainingTime={time}
        onComplete={() => stop()}
      >
        {({ elapsedTime, color }) => (
          <>
            {setElapsedTime(Math.floor(elapsedTime))}
     
        
          <div style={{ color, display:'flex',flexDirection: 'column',textAlign:'center'
         }}>
           <div style={{color, display:'flex',flexDirection: 'column'}}>
            Time remaining:
          </div>

          <div style={{ color, display:'flex',flexDirection: 'row',marginTop:'20px'
         }}>
            <div style={{marginLeft:'20px'}} >
            {renderTime("seconds", getTimeSeconds(elapsedTime))}
            
            </div>
          
          </div>
          
         <div style={{marginTop:'20px'}}>
          
         % Completed: {elapsedTime/(time) *100 | 0 }
         </div>
   
          </div>
          </>
        )}
      </CountdownCircleTimer>
      </div>
      <div style={{
              display: 'flex',
              justifyContent: 'center',
              flexDirection: 'row',
              marginTop:'20px',
              marginBottom:'20px'
            }}>

            Maintenance Logs : {steps[Math.floor(elapsedTime)]}
          
      </div>
      <div style={{
              display: 'flex',
              justifyContent: 'space-between',
              flexDirection: 'row',
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
