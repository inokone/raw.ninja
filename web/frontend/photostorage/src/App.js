import React, { useState } from 'react';

import './App.css';
import ResponsiveAppBar from './components/Common/AppBar';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import Dashboard from './components/Dashboard/Dashboard';
import Preferences from './components/Preferences/Preferences';
import Upload from './components/Upload/Upload';
import PhotoList from './components/Photos/PhotoList';
import UserProfile from './components/Auth/UserProfile';
import PhotoDisplay from './components/Photos/PhotoDisplay';
import RegisterForm from './components/Auth/Register';
import Login from './components/Auth/Login';
import Logout from './components/Auth/Logout';
import ProtectedRoute from './components/Common/ProtectedRoute';
import NotFoundPage from './components/Common/NotFoundPage';
import SearchResult from './components/Search/SearchResult';

const App = () => {
  const [user, setUser] = useState(null);
  const [query, setQuery] = useState(null);

  return (
    <div className="App">
      <BrowserRouter>
        <ResponsiveAppBar user={user} setQuery={setQuery} />
        <header className="App-header">
          <div className="wrapper">
            <Routes>
              <Route path="/login" element={<Login setUser={setUser} />} />
              <Route path="/logout" element={<Logout setUser={setUser} />} />
              <Route path="/register" element={<RegisterForm />} />
              <Route element={<ProtectedRoute user={user} redirect="/login" />}>
                <Route path="/" element={<Dashboard user={user} />} />
                <Route path="/upload" element={<Upload user={user} />} />
                <Route path="/photos" element={<PhotoList user={user} />} />
                <Route path="/photos/:photosId" element={<PhotoDisplay user={user} />} />
                <Route path="/users/:userId" element={<UserProfile user={user} />} />
                <Route path="/profile" element={<Preferences user={user} />} />
                <Route path="/search" element={<SearchResult query={query} />} />
              </Route>
              <Route path="*" element={<NotFoundPage />} />
            </Routes>
          </div>
        </header>
      </BrowserRouter>
    </div>
  );
}

export default App;
