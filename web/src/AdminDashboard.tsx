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
import {
  useQueryGetnumOfFiles,
  useQueryGetnumOfHistoricalFiles,
  useQueryGetstorageUsedWithParity,
  useQueryGethistoricalStorageUsedWithParity,
  useQueryGetserverLogs,
  useQueryGethistoricalStorageUsed,
  useQueryGetstorageUsed,
} from './api/cluster';

function formatBytes(bytes: number, decimals = 2) {
  if (bytes === 0) return '0 Bytes';

  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}

export function AdminDashboard() {
  // Pie chart
  const barColors = ['#1f77b4', '#ff0000'];
  const RADIAN = Math.PI / 180;

  const qGetnumfiles = useQueryGetnumOfFiles();
  const qGetStorageUsed = useQueryGetstorageUsedWithParity();
  const qGetStorageUsedWOParity = useQueryGetstorageUsed();
  const qHistoricalDataUsage = useQueryGethistoricalStorageUsedWithParity(
    1,
    '',
    ''
  );
  const qHistoricalStorageUsedWOParity = useQueryGethistoricalStorageUsed(
    1,
    '',
    ''
  );

  const qNumberOfHistoricalFiles = useQueryGetnumOfHistoricalFiles(1, '', '');

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

  const qServerLogs = useQueryGetserverLogs(0, '', '', '');
  console.log(qServerLogs.data);

  // data for logs
  const [logs, setlogs] = useState(qServerLogs.data);

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
  const recentRows = qServerLogs?.data?.map((items, index) =>
    index < 4 ? (
      <tr key={index}>
        <td
          width="15%"
          style={{
            textAlign: 'left',
            fontWeight: '400',
            fontSize: '16px',
            color: 'black',
          }}
        >
          {items.TimeStamp}
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
          {items.ServerName}
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
          {items.Message}
        </td>
      </tr>
    ) : (
      ''
    )
  );

  // display the all the row from log
  const rows = qServerLogs?.data?.map((items, index) => (
    <tr key={index}>
      <td
        width="15%"
        style={{
          textAlign: 'left',
          fontWeight: '400',
          fontSize: '16px',
          color: 'black',
        }}
      >
        {items.TimeStamp}
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
        {items.ServerName}
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
        {items.Message}
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
        <div
          style={{
            display: 'flex',
            flexDirection: 'row',
            flexWrap: 'wrap',
            justifyContent: 'space-evenly',
          }}
        >
          <div
            style={{
              display: 'flex',
              flexDirection: 'column',
              flexWrap: 'wrap',
              justifyContent: 'space-evenly',
            }}
          >
            <Card className="dashboardCard" shadow="sm" p="xl">
              <Text weight={700}>
                Total Data Used:{' '}
                {qGetStorageUsed.data
                  ? formatBytes(qGetStorageUsed.data)
                  : null}
              </Text>
              <ResponsiveContainer width="100%" height={220}>
                <AreaChart
                  data={qHistoricalDataUsage.data}
                  margin={{ top: 20, right: 10, left: -10, bottom: 0 }}
                >
                  <defs>
                    <linearGradient id="color" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor="#8884d8" stopOpacity={0.8} />
                      <stop
                        offset="95%"
                        stopColor="#8884d8"
                        stopOpacity={0.0}
                      />
                    </linearGradient>
                  </defs>
                  //delete this if dont want mouseover
                  <Tooltip></Tooltip>
                  <XAxis dataKey="date" />
                  <YAxis dataKey="value" />
                  //use this if want axis CartesianGrid strokeDasharray="1 1"
                  <Area
                    type="monotone"
                    dataKey="value"
                    stroke="#8884d8"
                    fillOpacity={1}
                    fill="url(#color)"
                  />
                </AreaChart>
              </ResponsiveContainer>
            </Card>
            <Card className="dashboardCard" shadow="sm" p="xl">
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
          </div>
          <div
            style={{
              display: 'flex',
              flexDirection: 'column',
              flexWrap: 'wrap',
              justifyContent: 'space-evenly',
            }}
          >
            <Card className="dashboardCard" shadow="sm" p="xl">
              <Text weight={700}>Total File Stored: {qGetnumfiles.data}</Text>
              <ResponsiveContainer width="100%" height={220}>
                <AreaChart
                  data={qNumberOfHistoricalFiles.data}
                  margin={{ top: 20, right: 10, left: -10, bottom: 0 }}
                >
                  <defs>
                    <linearGradient id="color" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor="#8884d8" stopOpacity={0.8} />
                      <stop
                        offset="95%"
                        stopColor="#8884d8"
                        stopOpacity={0.0}
                      />
                    </linearGradient>
                  </defs>
                  //delete this if dont want mouseover
                  <Tooltip></Tooltip>
                  <XAxis dataKey="date" />
                  <YAxis dataKey="value" />
                  //use this if want axis CartesianGrid strokeDasharray="1 1"
                  <Area
                    type="monotone"
                    dataKey="value"
                    stroke="#8884d8"
                    fillOpacity={1}
                    fill="url(#color)"
                  />
                </AreaChart>
              </ResponsiveContainer>
            </Card>
            <Card className="dashboardCard" shadow="sm" p="xl">
              <Text
                style={{ marginTop: '-10px', marginBottom: '10px' }}
                weight={700}
              >
                {' '}
                Cluster Health :{' '}
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
          </div>

          <div
            style={{
              display: 'flex',
              flexDirection: 'column',
              flexWrap: 'wrap',
              justifyContent: 'space-evenly',
            }}
          >
            <Card className="dashboardCard" p="xl">
              <Text weight={700}>
                Total data used (not incl. replicas):{' '}
                {qGetStorageUsedWOParity.data
                  ? formatBytes(qGetStorageUsedWOParity.data)
                  : null}
              </Text>
              <ResponsiveContainer width="100%" height={220}>
                <AreaChart
                  data={qHistoricalStorageUsedWOParity.data}
                  margin={{ top: 20, right: 10, left: -10, bottom: 0 }}
                >
                  <defs>
                    <linearGradient id="color" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor="#8884d8" stopOpacity={0.8} />
                      <stop
                        offset="95%"
                        stopColor="#8884d8"
                        stopOpacity={0.0}
                      />
                    </linearGradient>
                  </defs>
                  //delete this if dont want mouseover
                  <Tooltip></Tooltip>
                  <XAxis dataKey="date" />
                  <YAxis dataKey="value" />
                  //use this if want axis CartesianGrid strokeDasharray="1 1"
                  <Area
                    type="monotone"
                    dataKey="value"
                    stroke="#8884d8"
                    fillOpacity={1}
                    fill="url(#color)"
                  />
                </AreaChart>
              </ResponsiveContainer>
            </Card>

            <Card className="dashboardCard" shadow="sm" p="xl">
              <Text
                style={{ marginTop: '-10px', marginBottom: '10px' }}
                weight={700}
              >
                {' '}
                Nodes Status :{' '}
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

          <div
            style={{
              display: 'flex',
              flexDirection: 'row',
              flexWrap: 'wrap',
              justifyContent: 'center',
            }}
          >
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
                  {(qServerLogs.data ?? []).length === 0 ? (
                    <Text className=" ml-2 mt-2 mb-5">Nothing here!</Text>
                  ) : (
                    <tbody>{recentRows}</tbody>
                  )}
                </Table>
              </ScrollArea>
              {(qServerLogs.data ?? []).length > 4 ? (
                <Button
                  variant="default"
                  color="dark"
                  size="md"
                  style={{ textAlign: 'right' }}
                  onClick={() => setOpened(true)}
                >
                  View All Logs
                </Button>
              ) : null}
            </Card>
          </div>
        </div>
      </div>
    </AppBase>
  );
}
