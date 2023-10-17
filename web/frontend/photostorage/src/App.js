import React, { useState, useEffect } from 'react';

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
import { CircularProgress } from "@mui/material";


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
          setIsLoading(false);
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
            <CircularProgress mt={5}/>
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
      </BrowserRouter>)}
    </div>
  );
}

export default App;
