import { ScrollArea, Button,Table } from '@mantine/core';
import { Link } from 'react-router-dom';
import AppBase from './components/AppBase';

export function AdminMaintenanceLogs() {
  const maintenanceLogss = [
    {
      'Maintenance date': '26/05/2023',
      'Start Time': '5:55',
      'End Time': '22:02',
      'Total Files': '34706',
      'Maintenance Type': 'Adaptive',
    },
    {
      'Maintenance date': '25/03/2023',
      'Start Time': '23:54',
      'End Time': '2:35',
      'Total Files': '98616',
      'Maintenance Type': 'Perfective',
    },
    {
      'Maintenance date': '11/06/2022',
      'Start Time': '0:39',
      'End Time': '16:31',
      'Total Files': '16484',
      'Maintenance Type': 'Perfective',
    },
    {
      'Maintenance date': '16/12/2021',
      'Start Time': '3:11',
      'End Time': '8:02',
      'Total Files': '95333',
      'Maintenance Type': 'Perfective',
    },
    {
      'Maintenance date': '01/07/2021',
      'Start Time': '16:13',
      'End Time': '7:54',
      'Total Files': '41738',
      'Maintenance Type': 'Preventive',
    },
    {
      'Maintenance date': '23/02/2023',
      'Start Time': '13:27',
      'End Time': '19:43',
      'Total Files': '58872',
      'Maintenance Type': 'Corrective',
    },
    {
      'Maintenance date': '04/11/2021',
      'Start Time': '10:40',
      'End Time': '21:39',
      'Total Files': '74169',
      'Maintenance Type': 'Corrective',
    },
    {
      'Maintenance date': '16/07/2021',
      'Start Time': '22:08',
      'End Time': '9:16',
      'Total Files': '41815',
      'Maintenance Type': 'Adaptive',
    },
    {
      'Maintenance date': '22/06/2023',
      'Start Time': '11:01',
      'End Time': '18:28',
      'Total Files': '15746',
      'Maintenance Type': 'Preventive',
    },
    {
      'Maintenance date': '01/10/2021',
      'Start Time': '5:37',
      'End Time': '8:51',
      'Total Files': '89951',
      'Maintenance Type': 'Adaptive',
    },
    {
      'Maintenance date': '09/10/2022',
      'Start Time': '21:51',
      'End Time': '0:07',
      'Total Files': '2723',
      'Maintenance Type': 'Preventive',
    },
    {
      'Maintenance date': '18/05/2022',
      'Start Time': '17:49',
      'End Time': '19:24',
      'Total Files': '44373',
      'Maintenance Type': 'Adaptive',
    },
    {
      'Maintenance date': '28/09/2021',
      'Start Time': '13:18',
      'End Time': '1:56',
      'Total Files': '35616',
      'Maintenance Type': 'Corrective',
    },
    {
      'Maintenance date': '29/11/2021',
      'Start Time': '2:07',
      'End Time': '17:29',
      'Total Files': '70070',
      'Maintenance Type': 'Preventive',
    },
    {
      'Maintenance date': '12/03/2023',
      'Start Time': '7:38',
      'End Time': '2:12',
      'Total Files': '81428',
      'Maintenance Type': 'Preventive',
    },
    {
      'Maintenance date': '16/11/2021',
      'Start Time': '14:06',
      'End Time': '10:37',
      'Total Files': '48063',
      'Maintenance Type': 'Perfective',
    },
    {
      'Maintenance date': '21/04/2022',
      'Start Time': '18:35',
      'End Time': '15:22',
      'Total Files': '26741',
      'Maintenance Type': 'Adaptive',
    },
    {
      'Maintenance date': '24/09/2021',
      'Start Time': '0:49',
      'End Time': '3:50',
      'Total Files': '64875',
      'Maintenance Type': 'Corrective',
    },
    {
      'Maintenance date': '15/05/2022',
      'Start Time': '13:26',
      'End Time': '9:14',
      'Total Files': '53077',
      'Maintenance Type': 'Corrective',
    },
    {
      'Maintenance date': '31/03/2022',
      'Start Time': '14:56',
      'End Time': '17:52',
      'Total Files': '57823',
      'Maintenance Type': 'Adaptive',
    },
    {
      'Maintenance date': '07/12/2022',
      'Start Time': '21:47',
      'End Time': '22:19',
      'Total Files': '8513',
      'Maintenance Type': 'Corrective',
    },
    {
      'Maintenance date': '04/05/2023',
      'Start Time': '5:30',
      'End Time': '3:02',
      'Total Files': '94245',
      'Maintenance Type': 'Preventive',
    },
    {
      'Maintenance date': '25/10/2022',
      'Start Time': '1:13',
      'End Time': '21:08',
      'Total Files': '57802',
      'Maintenance Type': 'Adaptive',
    },
    {
      'Maintenance date': '23/06/2022',
      'Start Time': '11:22',
      'End Time': '0:51',
      'Total Files': '16987',
      'Maintenance Type': 'Corrective',
    },
    {
      'Maintenance date': '03/06/2022',
      'Start Time': '4:10',
      'End Time': '12:26',
      'Total Files': '95198',
      'Maintenance Type': 'Corrective',
    },
    {
      'Maintenance date': '25/08/2021',
      'Start Time': '10:26',
      'End Time': '12:16',
      'Total Files': '65752',
      'Maintenance Type': 'Preventive',
    },
    {
      'Maintenance date': '24/11/2021',
      'Start Time': '18:02',
      'End Time': '13:49',
      'Total Files': '75780',
      'Maintenance Type': 'Corrective',
    },
    {
      'Maintenance date': '16/11/2021',
      'Start Time': '22:17',
      'End Time': '2:01',
      'Total Files': '32230',
      'Maintenance Type': 'Corrective',
    },
    {
      'Maintenance date': '26/08/2022',
      'Start Time': '20:19',
      'End Time': '0:17',
      'Total Files': '97923',
      'Maintenance Type': 'Perfective',
    },
    {
      'Maintenance date': '09/02/2023',
      'Start Time': '21:06',
      'End Time': '17:45',
      'Total Files': '37292',
      'Maintenance Type': 'Adaptive',
    },
    {
      'Maintenance date': '14/12/2022',
      'Start Time': '16:55',
      'End Time': '11:10',
      'Total Files': '50793',
      'Maintenance Type': 'Corrective',
    },
    {
      'Maintenance date': '14/06/2023',
      'Start Time': '14:22',
      'End Time': '17:49',
      'Total Files': '35406',
      'Maintenance Type': 'Perfective',
    },
    {
      'Maintenance date': '10/01/2023',
      'Start Time': '15:43',
      'End Time': '16:09',
      'Total Files': '54998',
      'Maintenance Type': 'Preventive',
    },
    {
      'Maintenance date': '05/01/2023',
      'Start Time': '23:44',
      'End Time': '0:42',
      'Total Files': '70297',
      'Maintenance Type': 'Corrective',
    },
    {
      'Maintenance date': '27/07/2022',
      'Start Time': '0:09',
      'End Time': '14:48',
      'Total Files': '47194',
      'Maintenance Type': 'Perfective',
    },
    {
      'Maintenance date': '29/10/2022',
      'Start Time': '0:13',
      'End Time': '16:25',
      'Total Files': '21992',
      'Maintenance Type': 'Adaptive',
    },
    {
      'Maintenance date': '18/09/2022',
      'Start Time': '11:10',
      'End Time': '18:43',
      'Total Files': '9062',
      'Maintenance Type': 'Preventive',
    },
    {
      'Maintenance date': '28/11/2021',
      'Start Time': '16:33',
      'End Time': '7:56',
      'Total Files': '35509',
      'Maintenance Type': 'Corrective',
    },
    {
      'Maintenance date': '26/06/2022',
      'Start Time': '22:37',
      'End Time': '0:01',
      'Total Files': '46064',
      'Maintenance Type': 'Adaptive',
    },
    {
      'Maintenance date': '25/12/2021',
      'Start Time': '4:41',
      'End Time': '23:05',
      'Total Files': '60845',
      'Maintenance Type': 'Perfective',
    },
    {
      'Maintenance date': '26/09/2021',
      'Start Time': '9:20',
      'End Time': '8:19',
      'Total Files': '78916',
      'Maintenance Type': 'Corrective',
    },
    {
      'Maintenance date': '14/10/2021',
      'Start Time': '23:41',
      'End Time': '1:19',
      'Total Files': '36573',
      'Maintenance Type': 'Corrective',
    },
    {
      'Maintenance date': '07/01/2023',
      'Start Time': '8:05',
      'End Time': '16:22',
      'Total Files': '20964',
      'Maintenance Type': 'Preventive',
    },
    {
      'Maintenance date': '05/06/2023',
      'Start Time': '0:14',
      'End Time': '17:27',
      'Total Files': '33820',
      'Maintenance Type': 'Preventive',
    },
    {
      'Maintenance date': '17/09/2021',
      'Start Time': '23:46',
      'End Time': '19:49',
      'Total Files': '70783',
      'Maintenance Type': 'Corrective',
    },
    {
      'Maintenance date': '14/06/2022',
      'Start Time': '1:52',
      'End Time': '11:34',
      'Total Files': '61721',
      'Maintenance Type': 'Corrective',
    },
    {
      'Maintenance date': '26/11/2022',
      'Start Time': '17:37',
      'End Time': '7:25',
      'Total Files': '59886',
      'Maintenance Type': 'Adaptive',
    },
    {
      'Maintenance date': '25/03/2023',
      'Start Time': '4:30',
      'End Time': '11:12',
      'Total Files': '94366',
      'Maintenance Type': 'Preventive',
    },
    {
      'Maintenance date': '04/02/2022',
      'Start Time': '11:47',
      'End Time': '23:51',
      'Total Files': '50188',
      'Maintenance Type': 'Adaptive',
    },
    {
      'Maintenance date': '29/03/2023',
      'Start Time': '19:13',
      'End Time': '17:12',
      'Total Files': '52022',
      'Maintenance Type': 'Adaptive',
    },
    {
      'Maintenance date': '15/02/2022',
      'Start Time': '7:42',
      'End Time': '1:28',
      'Total Files': '15100',
      'Maintenance Type': 'Adaptive',
    },
    {
      'Maintenance date': '13/04/2022',
      'Start Time': '15:01',
      'End Time': '5:39',
      'Total Files': '70405',
      'Maintenance Type': 'Adaptive',
    },
    {
      'Maintenance date': '11/06/2022',
      'Start Time': '11:41',
      'End Time': '2:00',
      'Total Files': '8216',
      'Maintenance Type': 'Corrective',
    },
    {
      'Maintenance date': '03/04/2023',
      'Start Time': '4:57',
      'End Time': '15:03',
      'Total Files': '18713',
      'Maintenance Type': 'Perfective',
    },
    {
      'Maintenance date': '14/11/2021',
      'Start Time': '21:37',
      'End Time': '3:58',
      'Total Files': '61030',
      'Maintenance Type': 'Perfective',
    },
    {
      'Maintenance date': '01/07/2022',
      'Start Time': '10:31',
      'End Time': '21:48',
      'Total Files': '16119',
      'Maintenance Type': 'Preventive',
    },
    {
      'Maintenance date': '25/07/2022',
      'Start Time': '12:03',
      'End Time': '15:37',
      'Total Files': '65747',
      'Maintenance Type': 'Perfective',
    },
    {
      'Maintenance date': '30/10/2022',
      'Start Time': '7:52',
      'End Time': '19:15',
      'Total Files': '37902',
      'Maintenance Type': 'Perfective',
    },
    {
      'Maintenance date': '10/11/2021',
      'Start Time': '22:34',
      'End Time': '2:05',
      'Total Files': '77358',
      'Maintenance Type': 'Adaptive',
    },
    {
      'Maintenance date': '31/12/2021',
      'Start Time': '9:25',
      'End Time': '19:27',
      'Total Files': '29050',
      'Maintenance Type': 'Preventive',
    },
    {
      'Maintenance date': '18/07/2021',
      'Start Time': '1:52',
      'End Time': '14:26',
      'Total Files': '59928',
      'Maintenance Type': 'Preventive',
    },
    {
      'Maintenance date': '12/05/2022',
      'Start Time': '10:16',
      'End Time': '23:43',
      'Total Files': '87393',
      'Maintenance Type': 'Adaptive',
    },
    {
      'Maintenance date': '10/01/2023',
      'Start Time': '7:51',
      'End Time': '18:55',
      'Total Files': '94458',
      'Maintenance Type': 'Preventive',
    },
    {
      'Maintenance date': '23/04/2022',
      'Start Time': '22:21',
      'End Time': '23:56',
      'Total Files': '39736',
      'Maintenance Type': 'Preventive',
    },
    {
      'Maintenance date': '21/01/2023',
      'Start Time': '20:38',
      'End Time': '3:16',
      'Total Files': '4369',
      'Maintenance Type': 'Perfective',
    },
    {
      'Maintenance date': '11/10/2021',
      'Start Time': '14:17',
      'End Time': '19:56',
      'Total Files': '98134',
      'Maintenance Type': 'Adaptive',
    },
    {
      'Maintenance date': '10/07/2021',
      'Start Time': '15:21',
      'End Time': '16:20',
      'Total Files': '34215',
      'Maintenance Type': 'Perfective',
    },
    {
      'Maintenance date': '18/12/2022',
      'Start Time': '9:20',
      'End Time': '22:54',
      'Total Files': '90812',
      'Maintenance Type': 'Perfective',
    },
    {
      'Maintenance date': '01/05/2023',
      'Start Time': '20:14',
      'End Time': '9:43',
      'Total Files': '10904',
      'Maintenance Type': 'Adaptive',
    },
    {
      'Maintenance date': '31/01/2022',
      'Start Time': '2:59',
      'End Time': '21:43',
      'Total Files': '79911',
      'Maintenance Type': 'Adaptive',
    },
    {
      'Maintenance date': '20/11/2022',
      'Start Time': '4:02',
      'End Time': '8:01',
      'Total Files': '18697',
      'Maintenance Type': 'Perfective',
    },
    {
      'Maintenance date': '11/02/2023',
      'Start Time': '12:13',
      'End Time': '11:59',
      'Total Files': '51100',
      'Maintenance Type': 'Preventive',
    },
    {
      'Maintenance date': '06/06/2022',
      'Start Time': '23:42',
      'End Time': '3:01',
      'Total Files': '10617',
      'Maintenance Type': 'Adaptive',
    },
    {
      'Maintenance date': '30/01/2022',
      'Start Time': '8:25',
      'End Time': '9:28',
      'Total Files': '43850',
      'Maintenance Type': 'Corrective',
    },
    {
      'Maintenance date': '01/11/2021',
      'Start Time': '5:47',
      'End Time': '5:38',
      'Total Files': '69661',
      'Maintenance Type': 'Corrective',
    },
    {
      'Maintenance date': '30/05/2022',
      'Start Time': '6:55',
      'End Time': '7:52',
      'Total Files': '3031',
      'Maintenance Type': 'Perfective',
    },
    {
      'Maintenance date': '01/08/2022',
      'Start Time': '4:51',
      'End Time': '23:27',
      'Total Files': '49711',
      'Maintenance Type': 'Corrective',
    },
    {
      'Maintenance date': '21/02/2022',
      'Start Time': '10:21',
      'End Time': '15:31',
      'Total Files': '93587',
      'Maintenance Type': 'Adaptive',
    },
    {
      'Maintenance date': '15/12/2021',
      'Start Time': '13:14',
      'End Time': '22:11',
      'Total Files': '31428',
      'Maintenance Type': 'Perfective',
    },
    {
      'Maintenance date': '26/12/2022',
      'Start Time': '2:16',
      'End Time': '10:26',
      'Total Files': '64560',
      'Maintenance Type': 'Perfective',
    },
    {
      'Maintenance date': '27/09/2021',
      'Start Time': '14:11',
      'End Time': '17:58',
      'Total Files': '16352',
      'Maintenance Type': 'Corrective',
    },
    {
      'Maintenance date': '28/01/2023',
      'Start Time': '0:51',
      'End Time': '16:14',
      'Total Files': '20736',
      'Maintenance Type': 'Corrective',
    },
    {
      'Maintenance date': '17/01/2022',
      'Start Time': '6:01',
      'End Time': '22:15',
      'Total Files': '98602',
      'Maintenance Type': 'Adaptive',
    },
    {
      'Maintenance date': '11/05/2023',
      'Start Time': '17:00',
      'End Time': '12:28',
      'Total Files': '48534',
      'Maintenance Type': 'Preventive',
    },
    {
      'Maintenance date': '09/05/2023',
      'Start Time': '14:54',
      'End Time': '19:19',
      'Total Files': '96591',
      'Maintenance Type': 'Perfective',
    },
    {
      'Maintenance date': '19/11/2021',
      'Start Time': '9:01',
      'End Time': '13:54',
      'Total Files': '38409',
      'Maintenance Type': 'Adaptive',
    },
    {
      'Maintenance date': '19/03/2022',
      'Start Time': '3:41',
      'End Time': '3:52',
      'Total Files': '58648',
      'Maintenance Type': 'Corrective',
    },
    {
      'Maintenance date': '14/07/2021',
      'Start Time': '7:37',
      'End Time': '15:43',
      'Total Files': '61943',
      'Maintenance Type': 'Corrective',
    },
    {
      'Maintenance date': '21/05/2022',
      'Start Time': '1:33',
      'End Time': '5:39',
      'Total Files': '14096',
      'Maintenance Type': 'Corrective',
    },
    {
      'Maintenance date': '08/12/2021',
      'Start Time': '1:51',
      'End Time': '17:13',
      'Total Files': '16657',
      'Maintenance Type': 'Corrective',
    },
    {
      'Maintenance date': '25/04/2023',
      'Start Time': '18:32',
      'End Time': '10:44',
      'Total Files': '59591',
      'Maintenance Type': 'Corrective',
    },
    {
      'Maintenance date': '21/06/2023',
      'Start Time': '23:33',
      'End Time': '8:33',
      'Total Files': '26750',
      'Maintenance Type': 'Preventive',
    },
    {
      'Maintenance date': '26/05/2023',
      'Start Time': '16:29',
      'End Time': '23:52',
      'Total Files': '58137',
      'Maintenance Type': 'Adaptive',
    },
    {
      'Maintenance date': '08/09/2021',
      'Start Time': '15:21',
      'End Time': '7:10',
      'Total Files': '64637',
      'Maintenance Type': 'Perfective',
    },
    {
      'Maintenance date': '06/06/2023',
      'Start Time': '10:22',
      'End Time': '5:09',
      'Total Files': '22539',
      'Maintenance Type': 'Corrective',
    },
    {
      'Maintenance date': '19/01/2023',
      'Start Time': '23:32',
      'End Time': '0:17',
      'Total Files': '14255',
      'Maintenance Type': 'Perfective',
    },
    {
      'Maintenance date': '04/02/2022',
      'Start Time': '3:52',
      'End Time': '0:49',
      'Total Files': '79820',
      'Maintenance Type': 'Corrective',
    },
    {
      'Maintenance date': '07/03/2023',
      'Start Time': '13:49',
      'End Time': '4:00',
      'Total Files': '17822',
      'Maintenance Type': 'Perfective',
    },
    {
      'Maintenance date': '14/03/2023',
      'Start Time': '10:30',
      'End Time': '1:47',
      'Total Files': '17411',
      'Maintenance Type': 'Preventive',
    },
    {
      'Maintenance date': '21/08/2021',
      'Start Time': '6:47',
      'End Time': '21:01',
      'Total Files': '30851',
      'Maintenance Type': 'Corrective',
    },
  ];
  // variable that show all the logs inside the props.groupList
  const logsHeader =  [
    "Maintenance date", "Total Files","Start Time","End Time","Maintenance Type"
  ]

  // display table header that is from props
  const ths = logsHeader.map((items) => (
    <th style={{ fontWeight: '600' }}>{items}</th>
  ));

  // display all the rows that is from props
  const rows = maintenanceLogss.map((items) => (
    <tr>
      <td>{items['Maintenance date']}</td>
      <td>{items['Total Files']}</td>
      <td>{items['Start Time']}</td>
      <td>{items['End Time']}</td>
      <td>{items['Maintenance Type']}</td>
    </tr>
  ));
  return (
    <>
      <AppBase userType="admin">
      <div
      style={{
        display: 'flex',
        justifyContent: 'center',
        height: '85vh',
      }}
    >
      <div className="maintenanceLogs">
        <ScrollArea
          style={{
            height: '85%',
            width: '90%',
            marginTop: '10px',
            marginLeft: '20px',
          }}
        >
          <Table
            id="maintenanceLogsTable"
            captionSide="top"
            striped
            highlightOnHover
            verticalSpacing="sm"
          >
            <caption
              style={{
                textAlign: 'center',
                fontWeight: '600',
                fontSize: '24px',
                color: 'black',
              }}
            >
            Maintenance Records
            </caption>
            <thead>{ths}</thead>
            <tbody>{rows}</tbody>
          </Table>
        </ScrollArea>

        <div
          style={{
            display: 'flex',
            flexDirection: 'column',
          }}
        >
          <Button
            variant="default"
            color="dark"
            size="md"
            style={{
              alignSelf: 'flex-end',
              marginRight: '15px',
              marginTop: '10px',
            }}
            component={Link}
            to="/runmaintenance"
          >
            Perform Maintenance
          </Button>
        </div>
      </div>
    </div>
      </AppBase>
    </>
  );
}

export default AdminMaintenanceLogs;
function useState(groupList: any): [any, any] {
  throw new Error('Function not implemented.');
}

