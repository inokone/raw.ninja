import React, { useState, useEffect } from 'react';

import './App.css';
import ResponsiveAppBar from './components/Common/AppBar';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import Dashboard from './components/Dashboard/Dashboard';
import Preferences from './components/Preferences/Preferences';
import Admin from './components/Admin/Admin';
import Upload from './components/Upload/Upload';
import PhotoList from './components/Photos/PhotoList';
import UserProfile from './components/Account/UserProfile';
import ResetPassword from './components/Account/ResetPassword';
import EmailConfirmation from './components/Account/EmailConfirmation';
import PhotoDisplay from './components/Photos/PhotoDisplay';
import SignupForm from './components/Account/Signup';
import Login from './components/Auth/Login';
import Logout from './components/Auth/Logout';
import ProtectedRoute from './components/Common/ProtectedRoute';
import NotFoundPage from './components/Common/NotFoundPage';
import SearchResult from './components/Search/SearchResult';
import ProgressDisplay from './components/Common/ProgressDisplay';
import RecoverPassword from './components/Account/RecoverPassword';
import Photopea from './components/Editor/Photopea';
import Landing from './components/Landing/Landing';
import TermsOfUse from './components/Common/TermsOfUse';
import CookieConsent from './components/Common/CookieConsent';
import CookieRulesDialog from './components/Common/CookieRulesDialog';

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const App = () => {
  const [user, setUser] = useState(null);
  const [query, setQuery] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isCookieRulesDialogOpen, setIsCookieRulesDialogOpen] = useState(false);

  useEffect(() => {
    fetch(REACT_APP_API_PREFIX + '/api/v1/account/profile', {
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

  const handleCookieRulesDialogOpen = React.useCallback(() => {
    setIsCookieRulesDialogOpen(true);
  }, [setIsCookieRulesDialogOpen]);

  const handleCookieRulesDialogClose = React.useCallback(() => {
    setIsCookieRulesDialogOpen(false);
  }, [setIsCookieRulesDialogOpen]);

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
        {!isCookieRulesDialogOpen && (
          <CookieConsent
            handleCookieRulesDialogOpen={handleCookieRulesDialogOpen}
          />
        )}
        <CookieRulesDialog
          open={isCookieRulesDialogOpen}
          onClose={handleCookieRulesDialogClose}
        />
        <ResponsiveAppBar user={user} setQuery={setQuery} />
        <header className="App-header">
          <div className="wrapper">
            <Routes>
              <Route path="/login" element={<Login setUser={setUser} />} />
              <Route path="/logout" element={<Logout setUser={setUser} />} />
              <Route path="/password/reset" element={<ResetPassword />} />
              <Route path="/password/recover" element={<RecoverPassword />} />
              <Route path="/signup" element={<SignupForm />} />
              <Route path="/terms" element={<TermsOfUse />} />
              <Route path="/" element={<Landing />} />
              <Route path="/confirm" element={<EmailConfirmation />} />
                  <Route element={<ProtectedRoute user={user} redirect="/" />}>
                <Route path="/home" element={<Dashboard user={user} />} />
                <Route path="/upload" element={<Upload user={user} />} />
                <Route path="/photos" element={<PhotoList user={user} />} />
                <Route path="/editor/:photoId" element={<Photopea/>} />
                <Route path="/photos/:photosId" element={<PhotoDisplay user={user} />} />
                <Route path="/users/:userId" element={<Preferences user={user} />} />
                <Route path="/profile" element={<UserProfile user={user} />} />
                <Route path="/search" element={<SearchResult query={query} />} />
              </Route>
              <Route element={<ProtectedRoute user={user} target="admin" redirect="/" />}>
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
