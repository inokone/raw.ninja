import React, { useState, useEffect } from 'react';

import './App.css';
import ResponsiveAppBar from './components/Common/AppBar';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import Dashboard from './components/Dashboard/Dashboard';
import Preferences from './components/Preferences/Preferences';
import Admin from './components/Admin/Admin';
import Upload from './components/Upload/Upload';
import PhotoList from './components/Photos/PhotoList';
import UserProfile from './components/Auth/UserProfile';
import ResetPassword from './components/Auth/ResetPassword';
import EmailConfirmation from './components/Auth/EmailConfirmation';
import PhotoDisplay from './components/Photos/PhotoDisplay';
import RegisterForm from './components/Auth/Register';
import Login from './components/Auth/Login';
import Logout from './components/Auth/Logout';
import ProtectedRoute from './components/Common/ProtectedRoute';
import NotFoundPage from './components/Common/NotFoundPage';
import SearchResult from './components/Search/SearchResult';
import ProgressDisplay from './components/Common/ProgressDisplay';


const { REACT_APP_API_PREFIX } = process.env;

const App = () => {
  const [user, setUser] = useState(null);
  const [query, setQuery] = useState(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    fetch(REACT_APP_API_PREFIX + '/api/v1/auth/profile', {
      method: "GET",
      mode: "cors",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      }
    })
      .then(response => {
        if (response.ok) {
          response.json().then(content => {
            setUser(content)
            setIsLoading(false);
          })
        } else {
          throw new Error(response.status + ": " + response.statusText);
        }
      }).catch(() => {
        setIsLoading(false);
      });
  }, [])

  return (
    <div className="App">
      {isLoading ? (
        <header className="App-header">
          <div className="wrapper">
            <ProgressDisplay />
          </div>
        </header>
      ) : (
      <BrowserRouter>
        <ResponsiveAppBar user={user} setQuery={setQuery} />
        <header className="App-header">
          <div className="wrapper">
            <Routes>
              <Route path="/login" element={<Login setUser={setUser} />} />
              <Route path="/logout" element={<Logout setUser={setUser} />} />
              <Route path="/password" element={<ResetPassword />} />
              <Route path="/register" element={<RegisterForm />} />
              <Route path="/confirm" element={<EmailConfirmation />} />
              <Route element={<ProtectedRoute user={user} redirect="/login" />}>
                <Route path="/" element={<Dashboard user={user} />} />
                <Route path="/upload" element={<Upload user={user} />} />
                <Route path="/photos" element={<PhotoList user={user} />} />
                <Route path="/photos/:photosId" element={<PhotoDisplay user={user} />} />
                <Route path="/users/:userId" element={<Preferences user={user} />} />
                <Route path="/profile" element={<UserProfile user={user} />} />
                <Route path="/search" element={<SearchResult query={query} />} />
              </Route>
              <Route element={<ProtectedRoute user={user} target="admin" redirect="/login" />}>
                <Route path="/admin" element={<Admin user={user} />} />
              </Route>
              <Route path="*" element={<NotFoundPage />} />
            </Routes>
          </div>
        </header>
      </BrowserRouter>)}
    </div>
  );
}

export default App;
