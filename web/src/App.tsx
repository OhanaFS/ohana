//import AppBase from "./AppBase";
import { BrowserRouter as Router, Routes, Route} from 'react-router-dom';
import AppBase from "./AppShell"; //switched to responsive base
import UserFiles from './userFiles';

export default function Demo() {
  return (
    <Router>
      <Routes>
        <Route path='/' element={<AppBase name='Cute Guy' username='@person' image='https://images.unsplash.com/photo-1496302662116-35cc4f36df92?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=2070&q=80'><></></AppBase>} />
        <Route path='/files' element={<UserFiles />} />
      </Routes>
    </Router>
  );
}