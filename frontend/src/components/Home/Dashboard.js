import * as React from 'react';
import PropTypes from "prop-types";
import FairWarningDialog from '../Common/FairWarningDialog';
import Uploads from './Uploads';
import Favorites from './Favorites';
import Statistics from './Statistics';
import { Typography, Box } from '@mui/material';
import Welcome from './Welcome';

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
            <Welcome />
          }
          {!isNewUser &&
          <Box>
            <Typography variant='h2' sx={{ marginBottom: 8, marginTop: 2 }}>{"Welcome back" + (user.first_name && ", " + user.first_name) + "!"}</Typography>
            <Typography align='left' variant='h6' sx={{ margin: 2 }}>It's great to see you again. Here are your latest storage statistics:</Typography>
            <Box maxWidth='lg' sx={{margin: 'auto'}}>
              <Statistics user={user} onDataLoaded={handleStatsLoaded} />
            </Box>
          </Box>}
          {uploadCount > 0 && <Typography align='left' variant='h6' sx={{ margin: 2, marginTop: 10 }}>your recent upload batches:</Typography>}
          <Uploads user={user} onDataLoaded={handleUploadsLoaded} />
          {favoriteCount > 0 && <Typography align='left' variant='h6' sx={{ margin: 2, marginTop: 10 }}>and your latest favorite photos:</Typography>}
          <Favorites user={user} onDataLoaded={handleFavoritesLoaded} />
        </>}
    </React.Fragment>
  )
}

Dashboard.propTypes = {
  user: PropTypes.object
};

export default Dashboard;