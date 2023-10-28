import React from 'react';
import { Alert, Grid, Box, Typography } from "@mui/material";
import ProgressDisplay from '../Common/ProgressDisplay';

const { REACT_APP_API_PREFIX } = process.env;

const Profile = () => {
  const [error, setError] = React.useState(null)
  const [loading, setLoading] = React.useState(false)
  const [stats, setStats] = React.useState(null)

  const formatBytes = (bytes, decimals = 2) => {
    let negative = (bytes < 0)
    if (negative) {
      bytes = -bytes 
    }
    if (!+bytes) return '0 Bytes'

    const k = 1024
    const dm = decimals < 0 ? 0 : decimals
    const sizes = ['Bytes', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB']

    const i = Math.floor(Math.log(bytes) / Math.log(k))

    let displayText = `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`
    return negative ? "-" + displayText : displayText
  }

  const loadStats = () => {
    fetch(REACT_APP_API_PREFIX + '/api/v1/statistics/user', {
      method: "GET",
      mode: "cors",
      credentials: "include"
    })
      .then(response => {
        if (!response.ok) {
          throw new Error(response.status + ": " + response.statusText);
        } else {
          response.json().then(content => {
            setLoading(false)
            setStats(content)
          })
        }
      })
      .catch(error => {
        setError(error)
        setLoading(false)
      });
  }

  React.useEffect(() => {
    if (!stats && !error && !loading) {
      loadStats()
    }
  }, [stats, loading, error])

  return (
    <>
      {error && <Alert sx={{ mb: 4 }} severity="error">{error}</Alert>}
      {loading && <ProgressDisplay /> }
      {stats &&
        <Box sx={{ display: 'flex', justifyContent: 'center', borderRadius: '4px', pb: 4 }}>
          <Box sx={{ bgcolor: 'rgba(0, 0, 0, 0.34)', color: 'white', mt: 10, borderRadius: '4px', width: '500px' }}>
            <Grid container>
              <Grid item xs={12}><Typography variant='h6' sx={{ borderRadius: '4px', bgcolor: 'rgba(0, 0, 0, 0.54)' }}>Profile</Typography></Grid>
            </Grid>
            <Grid container>
              <Grid item xs={5}><Typography>E-mail:</Typography></Grid>
              <Grid item xs={7}><Typography>{stats.email}</Typography></Grid>
            </Grid>
            <Grid container>
              <Grid item xs={5}><Typography>Subscription model:</Typography></Grid>
              <Grid item xs={7}><Typography>{stats.role}</Typography></Grid>
            </Grid>
            <Grid container>
              <Grid item xs={5}><Typography>Phone:</Typography></Grid>
              <Grid item xs={7}><Typography>{stats.phone}</Typography></Grid>
            </Grid>
            <Grid container>
              <Grid item xs={5}><Typography>Registered on:</Typography></Grid>
              <Grid item xs={7}><Typography>{new Date(stats.registration_date * 1000).toLocaleDateString()}</Typography></Grid>
            </Grid>
            <Grid container>
              <Grid item xs={5}><Typography>Photo count:</Typography></Grid>
              <Grid item xs={7}><Typography>{stats.photos}</Typography></Grid>
            </Grid>
            <Grid container>
              <Grid item xs={5}><Typography>Favorites:</Typography></Grid>
              <Grid item xs={7}><Typography>{stats.favorites}</Typography></Grid>
            </Grid>
            <Grid container>
              <Grid item xs={5}><Typography>User space:</Typography></Grid>
              <Grid item xs={7}><Typography>{formatBytes(stats.used_space)}</Typography></Grid>
            </Grid>
            <Grid container>
              <Grid item xs={5}><Typography>Available space:</Typography></Grid>
              <Grid item xs={7}><Typography>{stats.quota <= 0 ? 'Unlimited' : formatBytes(stats.available_space)}</Typography></Grid>
            </Grid>
          </Box>
        </Box>}
    </>
  );
}
export default Profile; 