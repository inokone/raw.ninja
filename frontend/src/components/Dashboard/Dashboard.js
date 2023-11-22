import * as React from 'react';
import PhotoGrid from '../Photos/PhotoGrid';
import { Typography } from '@mui/material';
import FairWarningDialog from '../Common/FairWarningDialog';


const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const Dashboard = ({user}) => {

  const [isWarningDialogOpen, setIsWarningDialogOpen] = React.useState(true);

  const handleWarningDialogClose = React.useCallback(() => {
    setIsWarningDialogOpen(false);
  }, [setIsWarningDialogOpen]);

  const populate = () => {
    if (!user) {
      return new Promise((resolve, reject) => {
        reject("User not set!")
      })
    }
    return fetch(REACT_APP_API_PREFIX + '/api/v1/search/favorites', {
      method: "GET",
      mode: "cors",
      credentials: "include"
    })
  }

  return (
    <React.Fragment>
      <FairWarningDialog
        open={isWarningDialogOpen}
        onClose={handleWarningDialogClose}
      />
      <Typography variant='h5' sx={{ textAlign: 'left', pt: 2, pl: 2 }}>Favorite photos</Typography>
      <PhotoGrid populator={populate} data={[]}></PhotoGrid>
    </React.Fragment>
  )
}

export default Dashboard;