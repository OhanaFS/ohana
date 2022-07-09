import { useState } from 'react';
import { Button, Card, Table, Text, ScrollArea } from '@mantine/core';

import { Area, AreaChart, Cell, Legend, Pie, PieChart, Tooltip, XAxis, YAxis, ResponsiveContainer } from 'recharts';

import AppBase from './components/AppBase';
import '../src/assets/styles.css';

export function AdminDashboard() {
  const [status, setStatus] = useState(["Healthy"]);
  const barColors = ["#1f77b4", "#ff0000"]
  const RADIAN = Math.PI / 180;

  const renderCustomizedLabel = ({ cx, cy, midAngle, innerRadius, outerRadius, percent, index }: any) => {
    const radius = innerRadius + (outerRadius - innerRadius) * 0.5;
    const x = cx + radius * Math.cos(-midAngle * RADIAN);
    const y = cy + radius * Math.sin(-midAngle * RADIAN);
    return (
      <text x={x} y={y} fill="white" textAnchor={x > cx ? 'start' : 'end'} dominantBaseline="central">
        {`${(percent * 100).toFixed(0)}%`}
      </text>
    );
  };

  // Show only 4 recent logs
  const logs = [
    {
      "Date and time": "09/16/2019, 14:07",
      "Account": "End-user",
      "User": "Peter",
      "Change": "Added a node ip address 45.2.1.6"
    },
    {
      "Date and time": "09/16/2019, 14:07",
      "Account": "End-user",
      "User": "Peter",
      "Change": "Added a node ip address 95.2.2.6"
    },
    {
      "Date and time": "09/16/2019, 14:09",
      "Account": "End-user",
      "User": "Peter",
      "Change": "Added a node ip address 125.2.1.6"
    },
    {
      "Date and time": "09/16/2019, 14:10",
      "Account": "End-user",
      "User": "Peter",
      "Change": "Added a node ip address 125.2.1.6"
    },


  ];
  const ClusterHealthChartData = [

    {
      "name": "No of Healthy Nodes",
      "value": 600
    },
    {
      "name": "No of Unhealthy Nodes",
      "value": 400
    },
  ]
  const DiskUsageChartData = [


    {
      "name": "Empty",
      "value": 600
    },
    {
      "name": "Filled",
      "value": 400
    },
  ]

  const NewUserChartData = [
    {
      "Date": "jan 20",
      "Total Data Used": 4000,

    },
    {
      "Date": "feb 20",
      "Total Data Used": 3000,

    },
    {
      "Date": "mar 20",
      "Total Data Used": 2000,

    },
    {
      "Date": "apr 20",
      "Total Data Used": 2780,

    },
    {
      "Date": "may 20",
      "Total Data Used": 1890,

    },
    {
      "Date": "jun 20",
      "Total Data Used": 2390,

    },
    {
      "Date": "july 20",
      "Total Data Used": 3490,

    }
  ]

  const NewFileChartData = [
    {
      "Date": "jan 20",
      "Total New File Stored": 1000,

    },
    {
      "Date": "feb 20",
      "Total New File Stored": 2500,

    },
    {
      "Date": "mar 20",
      "Total New File Stored": 2000,

    },
    {
      "Date": "apr 20",
      "Total New File Stored": 2780,

    },
    {
      "Date": "may 20",
      "Total New File Stored": 1890,

    },
    {
      "Date": "jun 20",
      "Total New File Stored": 2390,

    },
    {
      "Date": "july 20",
      "Total New File Stored": 3490,

    }
  ]

  const ths = (
    <tr >
      <th style={{ width: "15%", textAlign: "left", fontWeight: "700", fontSize: "16px", color: "black" }}>Date and time</th>
      <th style={{ width: "10%", textAlign: "left", fontWeight: "700", fontSize: "16px", color: "black" }}>User</th>
      <th style={{ width: "10%", textAlign: "left", fontWeight: "700", fontSize: "16px", color: "black" }}>	Account</th>
      <th style={{ width: "30%", textAlign: "left", fontWeight: "700", fontSize: "16px", color: "black" }}>Change</th>

    </tr>
  );
  const rows = logs.map((items) => (
    <tr >
      <td width="15%" style={{ textAlign: "left", fontWeight: "400", fontSize: "16px", color: "black" }}>{items["Date and time"]}</td>
      <td width="10%" style={{ textAlign: "left", fontWeight: "400", fontSize: "16px", color: "black" }}>{items["User"]}</td>
      <td width="10%" style={{ textAlign: "left", fontWeight: "400", fontSize: "16px", color: "black" }}>{items["Account"]}</td>
      <td width="30%" style={{ textAlign: "left", fontWeight: "400", fontSize: "16px", color: "black" }}>{items["Change"]}</td>

    </tr>
  ));
  const dashboardCard = {
    width: '365px',
    border: '0px',
    margin: '10px',
    height: '300px'
  }
  return (
    <AppBase userType="admin" name='Alex Simmons' username='@alex' image='https://images.unsplash.com/photo-1496302662116-35cc4f36df92?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=2070&q=80'>

      <div style={{
        display: "flex",
        flexDirection: "row",
        flexWrap: "wrap",
        justifyContent: "center",
        alignItems: 'flex-start'
      }}>

        <Card className='dashboardCard' style={dashboardCard}
          shadow="sm" p="xl">
          <Text weight={700}>Total Data Used:</Text>
          <ResponsiveContainer width='100%' height={220}>
            <AreaChart data={NewUserChartData}
              margin={{ top: 20, right: 10, left: -10, bottom: 0 }}>
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
              <Area type="monotone" dataKey="Total Data Used" stroke="#8884d8" fillOpacity={1} fill="url(#color)" />
            </AreaChart>
          </ResponsiveContainer>
        </Card>

        <Card className='dashboardCard' style={dashboardCard}
          shadow="sm" p="xl">
          <Text weight={700}>Total File Stored:</Text>
          <ResponsiveContainer width='100%' height={220}>
            <AreaChart data={NewFileChartData}
              margin={{ top: 20, right: 10, left: -10, bottom: 0 }}>
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
              <Area type="monotone" dataKey="Total New File Stored" stroke="#8884d8" fillOpacity={1} fill="url(#color)" />
            </AreaChart>
          </ResponsiveContainer>
        </Card>

        <Card style={dashboardCard} shadow="sm" p="xl">
          <Text weight={700}>  Total Disk usage:   </Text>
          <div style={{ marginTop: '-10px' }}>
            <ResponsiveContainer width='100%' height={250}>
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
                    <Cell key={`cell-${index}`} fill={barColors[index % barColors.length]} />
                  ))}
                </Pie>
                <Legend layout="horizontal" />
                <Tooltip></Tooltip>
              </PieChart>
            </ResponsiveContainer>
          </div>
        </Card>

        <Card style={dashboardCard} shadow="sm" p="xl">
          <Text style={{ marginTop: '-10px', marginBottom: '10px' }} weight={700}> Cluster Health : {status} </Text>
          <div style={{ marginTop: '-10px' }}>
            <ResponsiveContainer width='100%' height={250} >
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
                    <Cell key={`cell-${index}`} fill={barColors[index % barColors.length]} />
                  ))}
                </Pie>
                <Legend layout="horizontal" />
                <Tooltip></Tooltip>
              </PieChart>
            </ResponsiveContainer>
          </div>
        </Card>

        <Card className='dashboardLogsCard' shadow="sm" p="xl">
          <ScrollArea style={{ height: "95%", width: "100%" }}>
            <Table captionSide="top" striped highlightOnHover verticalSpacing="xs" >
              <caption style={{ textAlign: "left", fontWeight: "600", fontSize: "24px", color: "black", marginLeft: "2%" }}> Logs</caption>
              <thead>{ths}</thead>
              <tbody>{rows}</tbody>
            </Table>
          </ScrollArea>
          <Button variant="default" color="dark" size="md" style={{ textAlign: "right", marginTop: '1%' }}>
            View All Logs
          </Button>
        </Card>
      </div>
    </AppBase>
  );
}