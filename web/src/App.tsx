//import AppBase from "./AppBase";
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import AppBase from "./components/AppShell"; //switched to responsive base
import { VFSBrowser } from './components/userFiles';

export default function Demo() {
  return (
    <Router>
      <Routes>
        <Route path='/files' element={<AppBase name='Alex' username='@alex' image='https://images.unsplash.com/photo-1496302662116-35cc4f36df92?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=2070&q=80'><></></AppBase>} />
        <Route path='/' element={<VFSBrowser />} />
      </Routes>
    </Router>
  );
}