import React from 'react';
import AppBase from "./AppShell";
import { Table } from '@mantine/core';

const UserFiles = () => {
    const elements = [
        { fileName: 'SampleFile.txt', lastModified: '10 mins ago', size: '104 Kb'},
        { fileName: 'Wot.pdf', lastModified: '1 day ago', size: '104 Mb'},
        { fileName: 'Omg.mp4', lastModified: '10 days ago', size: '1.5 Gb'},
      ];
      const rows = elements.map((element) => (
        <tr key={element.fileName}>
          <td>{element.fileName}</td>
          <td>{element.lastModified}</td>
          <td>{element.size}</td>
        </tr>
      ));
    return (
        <AppBase name='Cute Guy' username='@person' image='https://images.unsplash.com/photo-1496302662116-35cc4f36df92?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=2070&q=80'>
            <Table>
                <thead>
                    <tr>
                    <th>File</th>
                    <th>Last Modified</th>
                    <th>Size</th>
                    </tr>
                </thead>
                <tbody>{rows}</tbody>
            </Table>
        </AppBase>
      );
  };
    
  export default UserFiles;
  
  
