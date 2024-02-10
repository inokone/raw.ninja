import React, { useState, useEffect } from 'react';

import './App.css';
import ResponsiveAppBar from './components/Common/AppBar';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import Dashboard from './components/Home/Dashboard';
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
import CookieConsent from './components/Common/CookieConsent';
import CookieRulesDialog from './components/Common/CookieRulesDialog';
import UploadDisplay from './components/Upload/UploadDisplay';
import CreateAlbum from './components/Album/CreateAlbum';
import AlbumDisplay from './components/Album/AlbumDisplay';
import AlbumRating from './components/Album/AlbumRating';
import AlbumList from './components/Album/AlbumList';
import RuleSets from './components/Rules/RuleSets';
import RuleSet from './components/Rules/RuleSet';
import AddPhotos from './components/Album/AddPhotos';
import Docs from './components/Docs/Docs';
import PrivacyPolicy from './components/Docs/PrivacyPolicy';
import TermsOfService from './components/Docs/TermsOfService';
import RatingGallery from './components/Photos/RatingGallery';
import UploadRating from './components/Upload/UploadRating';

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const App = () => {
  const [user, setUser] = useState(null);
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
          <ResponsiveAppBar user={user} />
          <header className="App-header">
            <div className="wrapper">
              <Routes>
                <Route path="/login" element={<Login setUser={setUser} />} />
                <Route path="/logout" element={<Logout setUser={setUser} />} />
                <Route path="/password/reset" element={<ResetPassword />} />
                <Route path="/password/recover" element={<RecoverPassword />} />
                <Route path="/signup" element={<SignupForm />} />
                <Route path="/terms" element={<TermsOfService />} />
                <Route path="/privacy" element={<PrivacyPolicy />} />
                <Route path="/docs" element={<Docs />} />
                <Route path="/" element={<Landing />} />
                <Route path="/confirm" element={<EmailConfirmation />} />
                <Route element={<ProtectedRoute user={user} redirect="/" />}>
                  <Route path="/home" element={<Dashboard user={user} />} />
                  <Route path="/upload" element={<Upload user={user} />} />
                  <Route path="/photos" element={<PhotoList user={user} />} />
                  <Route path="/ratings" element={<RatingGallery user={user} />} />
                  <Route path="/albums" element={<AlbumList user={user} />} />
                  <Route path="/albums/create" element={<CreateAlbum user={user} />} />
                  <Route path="/albums/:albumId/add" element={<AddPhotos user={user} />} />
                  <Route path="/albums/:albumId" element={<AlbumDisplay user={user} />} />
                  <Route path="/albums/:albumId/ratings" element={<AlbumRating user={user} />} />
                  <Route path="/editor/:photoId" element={<Photopea />} />
                  <Route path="/photos/:photosId" element={<PhotoDisplay user={user} />} />
                  <Route path="/uploads/:uploadId" element={<UploadDisplay user={user} />} />
                  <Route path="/uploads/:uploadId/ratings" element={<UploadRating user={user} />} />
                  <Route path="/rulesets" element={<RuleSets user={user} />} />
                  <Route path="/rulesets/:ruleSetId" element={<RuleSet user={user} />} />
                  <Route path="/profile" element={<UserProfile user={user} />} />
                  <Route path="/search" element={<SearchResult />} />
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
