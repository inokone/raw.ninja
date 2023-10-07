import React, {useState} from 'react';

import './App.css';
import ResponsiveAppBar from './AppBar';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import Dashboard from './components/Dashboard/Dashboard';
import Preferences from './components/Preferences/Preferences';
import Upload from './components/Upload/Upload';
import Gallery from './components/Gallery/Gallery';
import UserProfile from './components/Auth/UserProfile';
import PhotoDisplay from './components/Gallery/PhotoDisplay';
import RegisterForm from './components/Auth/Register';
import Login from './components/Auth/Login';

function App() {
  const [token, setToken] = useState();
  return (
    <div className="App">
      <ResponsiveAppBar />
      <header className="App-header">
      <div className="wrapper">
        <BrowserRouter>
          <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/login" element={<Login setToken={setToken} />} />
          <Route path="/register" element={<RegisterForm />} />          
          <Route path="/upload" element={<Upload />} />
          <Route path="/photos" element={<Gallery />} />
          <Route path="/photos/:photosId" element={<PhotoDisplay />} />
          <Route path="/users/:userId" element={<UserProfile />} />
          <Route path="/preferences" element={<Preferences />} />
          </Routes>
        </BrowserRouter>
     </div>
      </header>
    </div>
  );
}

export default App;
