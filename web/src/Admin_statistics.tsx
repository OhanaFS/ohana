import Admin_navigation from './Admin_navigation';
import AppBase from './components/AppBase';

import { Button, Card, Table, Text, useMantineTheme, Grid, ScrollArea } from '@mantine/core';

import { useState } from 'react';
import { useScrollIntoView } from '@mantine/hooks';

import { Area, AreaChart, Cell, Legend, Pie, PieChart, Tooltip, XAxis, YAxis, ResponsiveContainer } from 'recharts';







export function Admin_statistics() {


  const theme = useMantineTheme();

  const [status, setStatus] = useState(["Healthy"]);




  const { scrollableRef } = useScrollIntoView();

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




  const logs = [
    {
      "Date and time": "09/16/2019, 14:02",
      "Account": "End-user",
      "User": "Tom",
      "Change": "Added a node ip address 196.125.14.1"

    },
    {
      "Date and time": "09/16/2019, 14:03",
      "Account": "End-user",
      "User": "Mary",
      "Change": "Added a node ip address 68.1.4.5"
    },
    {
      "Date and time": "09/16/2019, 14:04",
      "Account": "End-user",
      "User": "Peter",
      "Change": "Added a node ip address 51.2.1.6"
    },
    {
      "Date and time": "09/16/2019, 14:05",
      "Account": "End-user",
      "User": "Peter",
      "Change": "Added a node ip address 52.2.1.6"
    },
    {
      "Date and time": "09/16/2019, 14:06",
      "Account": "End-user",
      "User": "Peter",
      "Change": "Added a node ip address 12.2.1.6"
    },
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
    {
      "Date and time": "09/16/2019, 14:11",
      "Account": "End-user",
      "User": "Peter",
      "Change": "Added a node ip address 120.2.1.6"
    },
    {
      "Date and time": "09/16/2019, 14:04",
      "Account": "End-user",
      "User": "Peter",
      "Change": "Added a node ip address 135.2.1.6"
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
  return (





    <AppBase userType="admin" name='Alex Simmons' username='@alex' image='https://images.unsplash.com/photo-1496302662116-35cc4f36df92?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=2070&q=80'>

      <Text style={{ marginLeft: "10%", width: "" }}> </Text>


      <Table style={{ marginLeft: "10%", width: "180vh" }}>
        <tr> <h2 style={{ fontWeight: "600", fontSize: "22px", color: "black", marginLeft: "5%" }}>DASHBOARD</h2></tr>
        <tr style={{}}>
          <td style={{ width: "20%" }}>
            <Card style={{ marginLeft: "3%", height: '350px', border: '1px solid ', marginTop: "0%", width: "95%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.colors.gray[0] }}
              shadow="sm"
              p="xl"      >

              <Card.Section style={{ textAlign: 'left', marginLeft: "1%", marginTop: "1%" }}>
                <Text weight={700}> Total Data Used:  </Text>
              </Card.Section>


              <AreaChart width={500} height={250} data={NewUserChartData}
                margin={{ top: 20, right: 50, left: 0, bottom: 0 }}>
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

            </Card>

          </td>
          <td style={{ width: "20%" }}>
            <Card style={{ marginLeft: "3%", height: '350px', border: '1px solid ', marginTop: "1%", width: "95%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.colors.gray[0] }}
              shadow="sm"
              p="xl"

            >
              <Card.Section style={{ textAlign: 'left', marginLeft: "1%", marginTop: "0%" }}>
                <Text weight={700}>Total File Stored: </Text>

              </Card.Section>

              <AreaChart width={500} height={250} data={NewFileChartData}
                margin={{ top: 20, right: 50, left: 0, bottom: 0 }}>
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

            </Card>

          </td>
          <td rowSpan={2} >
            <Card style={{ marginLeft: "3%", height: '700px', border: '1px solid ', marginTop: "0%", width: "95%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.colors.gray[0] }}
              shadow="sm"
              p="xl"

            >
              <ScrollArea style={{ height: "95%", width: "103%" }}>
                <Table captionSide="top" striped highlightOnHover verticalSpacing="xs" >

                  <caption style={{ textAlign: "left", fontWeight: "600", fontSize: "24px", color: "black", marginLeft: "2%" }}> Logs</caption>





                  <thead>{ths}</thead>
                  <tbody>{rows}</tbody>


                </Table>
              </ScrollArea>
              <Button variant="default" color="dark" size="md" style={{ textAlign: "right", marginLeft: '75%', marginTop: '1%' }}>
                Export logs
              </Button>
            </Card>

          </td>
        </tr>
        <tr>
          <td>
            <Card style={{ marginLeft: "3%", height: '350px', border: '1px solid ', marginTop: "3%", width: "95%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.colors.gray[0] }}
              shadow="sm"
              p="xl"

            >
              <Card.Section style={{ textAlign: 'left', marginLeft: "1%", marginTop: "1%" }}>
                <Text weight={700}>  Total Disk usage:   </Text>

              </Card.Section>

              <ResponsiveContainer width={400} height={280} >
                <PieChart width={400} height={280}>

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
                  <Legend layout="vertical" />
                  <Tooltip></Tooltip>
                </PieChart>

              </ResponsiveContainer>

            </Card>

          </td>
          <td>
            <Card style={{ marginLeft: "3%", height: '350px', border: '1px solid ', marginTop: "3%", width: "95%", background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.colors.gray[0] }}
              shadow="sm"
              p="xl"

            >
              <Card.Section style={{ textAlign: 'left', marginLeft: "1%", marginTop: "1%" }}>
                <Text weight={700}> Cluster Health :  </Text>

              </Card.Section>
              <Grid >
                <Grid.Col span={8} style={{ marginLeft: "" }}>

                  <ResponsiveContainer width={300} height={280} >
                    <PieChart width={300} height={280}>

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
                      <Legend layout="vertical" />
                      <Tooltip></Tooltip>
                    </PieChart>

                  </ResponsiveContainer>

                </Grid.Col>
                <Grid.Col span={4} style={{ marginTop: "15%" }}>

                  <Text  >Status: {status} </Text>





                </Grid.Col>


              </Grid>



            </Card>

          </td>


        </tr>

      </Table>

    </AppBase>





  );
}





