import * as React from 'react';
import FairWarningDialog from '../Common/FairWarningDialog';
import Uploads from './Uploads';
import Favorites from './Favorites';
import Statistics from './Statistics';
import { Typography } from '@mui/material';

const Dashboard = ({user}) => {

  const [isWarningDialogOpen, setIsWarningDialogOpen] = React.useState(true);

  const handleWarningDialogClose = React.useCallback(() => {
    setIsWarningDialogOpen(false);
  }, [setIsWarningDialogOpen]);

  return (
    <React.Fragment>
      {user && 
      <>
        <FairWarningDialog
          open={isWarningDialogOpen}
          onClose={handleWarningDialogClose}
        />
        <Typography variant='h2'>{"Welcome back" + (user.first_name && ", " + user.first_name) + "!"}</Typography>
        <Statistics user={user}/>
        <Uploads user={user}/>
        <Favorites user={user}/>
      </>}
    </React.Fragment>
  )
}

export default Dashboard;