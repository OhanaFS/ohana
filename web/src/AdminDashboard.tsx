import { useState } from 'react';
import { Button, Card, Table, Text, ScrollArea, Modal } from '@mantine/core';
import {
  Area,
  AreaChart,
  Cell,
  Legend,
  Pie,
  PieChart,
  Tooltip,
  XAxis,
  YAxis,
  ResponsiveContainer,
} from 'recharts';
import AppBase from './components/AppBase';
import '../src/assets/styles.css';

export function AdminDashboard() {
  // Pie chart
  const barColors = ['#1f77b4', '#ff0000'];
  const RADIAN = Math.PI / 180;
  const renderCustomizedLabel = ({
    cx,
    cy,
    midAngle,
    innerRadius,
    outerRadius,
    percent,
    index,
  }: any) => {
    const radius = innerRadius + (outerRadius - innerRadius) * 0.5;
    const x = cx + radius * Math.cos(-midAngle * RADIAN);
    const y = cy + radius * Math.sin(-midAngle * RADIAN);
    return (
      <text
        x={x}
        y={y}
        fill="white"
        textAnchor={x > cx ? 'start' : 'end'}
        dominantBaseline="central"
      >
        {`${(percent * 100).toFixed(0)}%`}
      </text>
    );
  };

  // all the logs
  const [logsModal, setOpened] = useState(false);

  var logsDetails = [
    {
      'Date and time': '09/16/2019, 14:07',
      Node: 'Peter',
      Change: 'Added a node ip address 45.2.1.6',
    },
    {
      'Date and time': '09/16/2019, 14:07',
      Node: 'Peter',
      Change: 'Added a node ip address 95.2.2.6',
    },
    {
      'Date and time': '09/16/2019, 14:09',
      Node: 'Peter',
      Change: 'Added a node ip address 125.2.1.6',
    },
    {
      'Date and time': '09/16/2019, 14:10',
      Node: 'Peter',
      Change: 'Added a node ip address 125.2.1.6',
    },
    {
      'Date and time': '09/16/2019, 14:10',
      Node: 'Peter',
      Change: 'Added a node ip address 125.2.1.7',
    },
    {
      'Date and time': '09/16/2019, 14:10',
      Node: 'Peter',
      Change: 'Added a node ip address 125.2.1.8',
    },
    {
      'Date and time': '09/16/2019, 14:10',
      Node: 'Peter',
      Change: 'Added a node ip address 125.2.1.9',
    },
    {
      'Date and time': '09/16/2019, 14:10',
      Node: 'Peter',
      Change: 'Added a node ip address 125.2.1.10',
    },
    {
      'Date and time': '09/16/2019, 14:10',
      Node: 'Peter',
      Change: 'Added a node ip address 125.2.1.11',
    },
    {
      'Date and time': '09/16/2019, 14:10',
      Node: 'Peter',
      Change: 'Added a node ip address 125.2.1.12',
    },
    {
      'Date and time': '09/16/2019, 14:10',
      Node: 'Peter',
      Change: 'Added a node ip address 125.2.1.13',
    },
    {
      'Date and time': '09/16/2019, 14:10',
      Node: 'Peter',
      Change: 'Added a node ip address 125.2.1.14',
    },
  ];

  // data for logs
  const [logs, setlogs] = useState(logsDetails);

  // data for clusterhealth
  const ClusterHealthChartData = [
    {
      name: 'No of Healthy Nodes',
      value: 600,
    },
    {
      name: 'No of Unhealthy Nodes',
      value: 400,
    },
  ];

  // data for diskusage
  const DiskUsageChartData = [
    {
      name: 'Empty',
      value: 600,
    },
    {
      name: 'Filled',
      value: 400,
    },
  ];
  const NodesStatus = [
    {
      name: 'Online',
      value: 600,
    },
    {
      name: 'Offline',
      value: 400,
    },
  ];
  // data for new user
  const NewUserChartData = [
    {
      Date: 'jan 20',
      'Total Data Used': 4000,
    },
    {
      Date: 'feb 20',
      'Total Data Used': 3000,
    },
    {
      Date: 'mar 20',
      'Total Data Used': 2000,
    },
    {
      Date: 'apr 20',
      'Total Data Used': 2780,
    },
    {
      Date: 'may 20',
      'Total Data Used': 1890,
    },
    {
      Date: 'jun 20',
      'Total Data Used': 2390,
    },
    {
      Date: 'july 20',
      'Total Data Used': 3490,
    },
  ];

  // data for new user
  const SizeOfFiles = [
    {
      Date: 'jan 20',
      'Total bytes': 4000,
    },
    {
      Date: 'feb 20',
      'Total bytes': 3000,
    },
    {
      Date: 'mar 20',
      'Total bytes': 2000,
    },
    {
      Date: 'apr 20',
      'Total bytes': 2780,
    },
    {
      Date: 'may 20',
      'Total bytes': 1890,
    },
    {
      Date: 'jun 20',
      'Total bytes': 2390,
    },
    {
      Date: 'july 20',
      'Total bytes': 3490,
    },
  ];

  
  // data for new files
  const NewFileChartData = [
    {
      Date: 'jan 20',
      'Total New File Stored': 1000,
    },
    {
      Date: 'feb 20',
      'Total New File Stored': 2500,
    },
    {
      Date: 'mar 20',
      'Total New File Stored': 2000,
    },
    {
      Date: 'apr 20',
      'Total New File Stored': 2780,
    },
    {
      Date: 'may 20',
      'Total New File Stored': 1890,
    },
    {
      Date: 'jun 20',
      'Total New File Stored': 2390,
    },
    {
      Date: 'july 20',
      'Total New File Stored': 3490,
    },
  ];

  // table header
  const ths = (
    <tr>
      <th
        style={{
          width: '15%',
          textAlign: 'left',
          fontWeight: '700',
          fontSize: '16px',
          color: 'black',
        }}
      >
        Date and time
      </th>
      <th
        style={{
          width: '10%',
          textAlign: 'left',
          fontWeight: '700',
          fontSize: '16px',
          color: 'black',
        }}
      >
        {' '}
        Node
      </th>
      <th
        style={{
          width: '30%',
          textAlign: 'left',
          fontWeight: '700',
          fontSize: '16px',
          color: 'black',
        }}
      >
        Message
      </th>
    </tr>
  );

  // display the recent 4 row from log
  const recentRows = logs.map((items, index) =>
    index < 4 ? (
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
          {items['Date and time']}
        </td>
        <td
          width="10%"
          style={{
            textAlign: 'left',
            fontWeight: '400',
            fontSize: '16px',
            color: 'black',
          }}
        >
          {items['Node']}
        </td>
        <td
          width="30%"
          style={{
            textAlign: 'left',
            fontWeight: '400',
            fontSize: '16px',
            color: 'black',
          }}
        >
          {items['Change']}
        </td>
      </tr>
    ) : (
      ''
    )
  );

  // display the all the row from log
  const rows = logs.map((items, index) => (
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
        {items['Date and time']}
      </td>
      <td
        width="10%"
        style={{
          textAlign: 'left',
          fontWeight: '400',
          fontSize: '16px',
          color: 'black',
        }}
      >
        {items['Node']}
      </td>
      <td
        width="30%"
        style={{
          textAlign: 'left',
          fontWeight: '400',
          fontSize: '16px',
          color: 'black',
        }}
      >
        {items['Change']}
      </td>
    </tr>
  ));




  // function to download all the logs
  function downloadLogs() {
    const fileData = JSON.stringify(logs);
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
  return (
    <AppBase userType="admin">
      <Modal
        centered
        size={600}
        opened={logsModal}
        title={
          <span style={{ fontSize: '22px', fontWeight: 550 }}> All Logs</span>
        }
        onClose={() => setOpened(false)}
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
              <Table captionSide="top" verticalSpacing="sm">
                <thead style={{}}>{ths}</thead>
                <tbody>{rows}</tbody>
              </Table>
            </ScrollArea>
            <Button
              variant="default"
              color="dark"
              size="md"
              style={{ alignSelf: 'flex-end' }}
              onClick={() => downloadLogs()}
            >
              Download Logs
            </Button>
          </div>
        </div>
      </Modal>
      <div
        style={{
          display: 'flex',
          flexDirection: 'column',
          justifyContent: 'center',
       
        }}
      >
        <div style={{
          display: 'flex',
          flexDirection: 'row',
          flexWrap: 'wrap',
          justifyContent: 'space-evenly',
        }}>
        <Card
          className="dashboardCard"
     
          shadow="sm"
          p="xl"
        >
          <Text weight={700}>Total Data Used:</Text>
          <ResponsiveContainer width="100%" height={220}>
            <AreaChart
              data={NewUserChartData}
              margin={{ top: 20, right: 10, left: -10, bottom: 0 }}
            >
              <defs>
                <linearGradient id="color" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#8884d8" stopOpacity={0.8} />
                  <stop offset="95%" stopColor="#8884d8" stopOpacity={0.0} />
                </linearGradient>
              </defs>
              //delete this if dont want mouseover
              <Tooltip></Tooltip>
              <XAxis dataKey="Date" />
              <YAxis dataKey="Total Data Used" />
              //use this if want axis CartesianGrid strokeDasharray="1 1"
              <Area
                type="monotone"
                dataKey="Total Data Used"
                stroke="#8884d8"
                fillOpacity={1}
                fill="url(#color)"
              />
            </AreaChart>
          </ResponsiveContainer>
        </Card>

        <Card
          className="dashboardCard"
      
          shadow="sm"
          p="xl"
        >
          <Text weight={700}>Total File Stored:</Text>
          <ResponsiveContainer width="100%" height={220}>
            <AreaChart
              data={NewFileChartData}
              margin={{ top: 20, right: 10, left: -10, bottom: 0 }}
            >
              <defs>
                <linearGradient id="color" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#8884d8" stopOpacity={0.8} />
                  <stop offset="95%" stopColor="#8884d8" stopOpacity={0.0} />
                </linearGradient>
              </defs>
              //delete this if dont want mouseover
              <Tooltip></Tooltip>
              <XAxis dataKey="Date" />
              <YAxis dataKey="Total New File Stored" />
              //use this if want axis CartesianGrid strokeDasharray="1 1"
              <Area
                type="monotone"
                dataKey="Total New File Stored"
                stroke="#8884d8"
                fillOpacity={1}
                fill="url(#color)"
              />
            </AreaChart>
          </ResponsiveContainer>
        </Card>
        <Card
          className="dashboardCard"
          p="xl"
        >
          <Text weight={700}>Total files size stored (not incl. replicas):</Text>
          <ResponsiveContainer width="100%" height={220}>
            <AreaChart
              data={SizeOfFiles}
              margin={{ top: 20, right: 10, left: -10, bottom: 0 }}
            >
              <defs>
                <linearGradient id="color" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#8884d8" stopOpacity={0.8} />
                  <stop offset="95%" stopColor="#8884d8" stopOpacity={0.0} />
                </linearGradient>
              </defs>
              //delete this if dont want mouseover
              <Tooltip></Tooltip>
              <XAxis dataKey="Date" />
              <YAxis dataKey="Total bytes" />
              //use this if want axis CartesianGrid strokeDasharray="1 1"
              <Area
                type="monotone"
                dataKey="Total bytes"
                stroke="#8884d8"
                fillOpacity={1}
                fill="url(#color)"
              />
            </AreaChart>
          </ResponsiveContainer>
        </Card>
        </div>
         <div style={{
          display: 'flex',
          flexDirection: 'row',
          flexWrap: 'wrap',
          justifyContent: 'space-evenly',
        
        }}>
        <Card     className="dashboardCard"  shadow="sm" p="xl">
          <Text weight={700}> Total Disk usage: </Text>
          <div style={{ marginTop: '-10px' }}>
            <ResponsiveContainer width="100%" height={250}>
              <PieChart>
                <Pie
                  data={DiskUsageChartData}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={renderCustomizedLabel}
                  outerRadius={100}
                  fill="#8884d8"
                  dataKey="value"
                >
                  {DiskUsageChartData.map((entry, index) => (
                    <Cell
                      key={`cell-${index}`}
                      fill={barColors[index % barColors.length]}
                    />
                  ))}
                </Pie>
                <Legend layout="horizontal" />
                <Tooltip></Tooltip>
              </PieChart>
            </ResponsiveContainer>
          </div>
        </Card>

        <Card     className="dashboardCard"  shadow="sm" p="xl">
          <Text
            style={{ marginTop: '-10px', marginBottom: '10px' }}
            weight={700}
          >
            {' '}
            Cluster Health : {' '}
          </Text>
          <div style={{ marginTop: '-10px' }}>
            <ResponsiveContainer width="100%" height={250}>
              <PieChart>
                <Pie
                  data={ClusterHealthChartData}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={renderCustomizedLabel}
                  outerRadius={100}
                  fill="#8884d8"
                  dataKey="value"
                >
                  {DiskUsageChartData.map((entry, index) => (
                    <Cell
                      key={`cell-${index}`}
                      fill={barColors[index % barColors.length]}
                    />
                  ))}
                </Pie>
                <Legend layout="horizontal" />
                <Tooltip></Tooltip>
              </PieChart>
            </ResponsiveContainer>
          </div>
        </Card>

        <Card     className="dashboardCard"  shadow="sm" p="xl">
          <Text
            style={{ marginTop: '-10px', marginBottom: '10px' }}
            weight={700}
          >
            {' '}
         Nodes Status  :{' '}
          </Text>
          <div style={{ marginTop: '-10px' }}>
            <ResponsiveContainer width="100%" height={250}>
              <PieChart>
                <Pie
                  data={NodesStatus}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={renderCustomizedLabel}
                  outerRadius={100}
                  fill="#8884d8"
                  dataKey="value"
                >
                  {DiskUsageChartData.map((entry, index) => (
                    <Cell
                      key={`cell-${index}`}
                      fill={barColors[index % barColors.length]}
                    />
                  ))}
                </Pie>
                <Legend layout="horizontal" />
                <Tooltip></Tooltip>
              </PieChart>
            </ResponsiveContainer>
          </div>
        </Card>            

     </div>
     <div style={{
          display: 'flex',
          flexDirection: 'row',
          flexWrap: 'wrap',
          justifyContent: 'center',
       
        }}>
        <Card className="dashboardLogsCard" shadow="sm" p="xl">
          <ScrollArea style={{ height: '90%', width: '100%' }}>
            <Table
              captionSide="top"
              striped
              highlightOnHover
              verticalSpacing="xs"
            >
              <caption
                style={{
                  textAlign: 'left',
                  fontWeight: '600',
                  fontSize: '24px',
                  color: 'black',
                  marginLeft: '2%',
                }}
              >
                {' '}
                Logs
              </caption>
              <thead>{ths}</thead>
              <tbody>{recentRows}</tbody>
            </Table>
          </ScrollArea>
          <Button
            variant="default"
            color="dark"
            size="md"
            style={{ textAlign: 'right' }}
            onClick={() => setOpened(true)}
          >
            View All Logs
          </Button>
        </Card>
        </div>
      </div>
    </AppBase>
  );
}
