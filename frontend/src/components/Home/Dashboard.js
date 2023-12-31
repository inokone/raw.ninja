import * as React from 'react';
import FairWarningDialog from '../Common/FairWarningDialog';
import Uploads from './Uploads';
import Favorites from './Favorites';
import Statistics from './Statistics';
import { Typography } from '@mui/material';

const Dashboard = ({ user }) => {

  const [isWarningDialogOpen, setIsWarningDialogOpen] = React.useState(true);
  const [isNewUser, setIsNewUser] = React.useState(false);
  const [uploadCount, setUploadCount] = React.useState(0);
  const [favoriteCount, setFavoriteCount] = React.useState(0);


  const handleWarningDialogClose = React.useCallback(() => {
    setIsWarningDialogOpen(false);
  }, [setIsWarningDialogOpen]);

  const handleStatsLoaded = (stats) => {
    setIsNewUser(stats.photos === 0)
  }

  const handleUploadsLoaded = (uploads) => {
    setUploadCount(uploads.length)
  }

  const handleFavoritesLoaded = (favorites) => {
    setFavoriteCount(favorites.length)
  }

  return (
    <React.Fragment>
      {user &&
        <>
          <FairWarningDialog
            open={isWarningDialogOpen}
            onClose={handleWarningDialogClose}
          />
          {isNewUser &&
            <>
              <Typography variant='h2' sx={{ marginBottom: 8, marginTop: 2 }}>{"Welcome to RAW.Ninja!"}</Typography>
              <Typography align='left' sx={{ margin: 2 }}>We're excited to have you on board as a new member of our photo-centric community.</Typography>
              <Typography align='left' sx={{ margin: 2 }}>Your journey with us begins now, as you embark on a seamless experience for storing and managing your raw photos. Dive into our user-friendly platform to securely upload and organize your precious data in their purest form. Explore the possibilities of preserving your visual stories, and if you ever have questions or need assistance, our team is here to help.</Typography>
              <Typography align='left' sx={{ margin: 2, marginBottom: 10 }}>Thank you for entrusting us with your photos, and we look forward to being a part of your photography journey on RAW.Ninja!</Typography>
            </>}
          {!isNewUser && <Typography variant='h2' sx={{ marginBottom: 8, marginTop: 2 }}>{"Welcome back" + (user.first_name && ", " + user.first_name) + "!"}</Typography>}
          <Typography align='left' variant='h6' sx={{ margin: 2 }}>It's great to see you again. Here are your latest storage statistics:</Typography>
          <Statistics user={user} onDataLoaded={handleStatsLoaded} />
          {uploadCount > 0 && <Typography align='left' variant='h6' sx={{ margin: 2, marginTop: 10 }}>your recent upload batches:</Typography>}
          <Uploads user={user} onDataLoaded={handleUploadsLoaded} />
          {favoriteCount > 0 && <Typography align='left' variant='h6' sx={{ margin: 2, marginTop: 10 }}>and your latest favorite photos:</Typography>}
          <Favorites user={user} onDataLoaded={handleFavoritesLoaded} />
        </>}
    </React.Fragment>
  )
}

export default Dashboard;