import * as React from 'react';
import PhotoGrid from '../Photos/PhotoGrid';
import { Typography } from '@mui/material';


const { REACT_APP_API_PREFIX } = process.env;

const Dashboard = () => {

  const populate = () => {
    return fetch(REACT_APP_API_PREFIX + '/api/v1/search/favorites', {
      method: "GET",
      mode: "cors",
      credentials: "include"
    })
  }

  return (
    <React.Fragment>
      <Typography variant='h5' sx={{ textAlign: 'left', pt: 2, pl: 2 }}>Favorite photos</Typography>
      <PhotoGrid populator={populate} data={[]}></PhotoGrid>
    </React.Fragment>
  )
}

export default Dashboard;