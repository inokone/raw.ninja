import React from 'react';
import { Alert, Grid, Box, Typography, Button } from "@mui/material";
import ProgressDisplay from '../Common/ProgressDisplay';
import ChangePassword from './ChangePassword';

const { REACT_APP_API_PREFIX } = process.env;

const Profile = ({user}) => {
  const [error, setError] = React.useState(null)
  const [loading, setLoading] = React.useState(false)
  const [stats, setStats] = React.useState(null)
  const [resendSuccess, setResendSuccess] = React.useState(false)

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
          response.json().then(content => {
            setError(content.message)
          });
        } else {
          response.json().then(content => {
            setLoading(false)
            setStats(content)
          })
        }
      })
      .catch(error => {
        setError(error.message)
        setLoading(false)
      });
  }

  const handleResendClick = () => {
    setError(null)
    fetch(REACT_APP_API_PREFIX + '/api/v1/account/resend', {
      method: "POST",
      mode: "cors",
      credentials: "include",
      body: JSON.stringify({
        "email": user.email
      })
    })
      .then(response => {
        if (!response.ok) {
          response.json().then(content => {
            setError(content.message)
          });
        } else {
          setResendSuccess(true)
        }
      })
      .catch(error => {
        setError(error.message)
      });
  }

  React.useEffect(() => {
    if (!stats && !error && !loading) {
      loadStats()
    }
  }, [stats, loading, error])

  return (
    <React.Fragment>
      {error && <Alert sx={{ mb: 4 }} severity="error">{error}</Alert>}
      {loading && <ProgressDisplay /> }
      {user && user.status === "registered" &&
        <Box sx={{ display: 'flex', justifyContent: 'center', borderRadius: '4px' }}>
          <Box sx={{ mt: 5, borderRadius: '4px', width: '500px' }}>
            <Alert severity='warning'>Your e-mail address hasn't been confirmed yet. Some features of the application will not be available to you. 
              <Button onClick={handleResendClick}>Resend confirmation</Button>
            </Alert>
            {resendSuccess && <Alert onClose={() => { setResendSuccess(null) }}>Confirmation e-mail sent.</Alert>}
          </Box>
        </Box>
      }
      
      {stats &&
        <Box sx={{ display: 'flex', justifyContent: 'center', borderRadius: '4px', pb: 4 }}>
          <Box sx={{ bgcolor: 'rgba(0, 0, 0, 0.34)', color: 'white', mt: 5, borderRadius: '4px', width: '500px' }}>
            <Grid container>
              <Grid item xs={12}><Typography variant='h6' sx={{ borderRadius: '4px', bgcolor: 'rgba(0, 0, 0, 0.54)' }}>Profile</Typography></Grid>
            </Grid>
            <Grid container>
              <Grid item xs={5}><Typography>Firstname:</Typography></Grid>
              <Grid item xs={7}><Typography>{stats.first_name}</Typography></Grid>
            </Grid>
            <Grid container>
              <Grid item xs={5}><Typography>Lastname:</Typography></Grid>
              <Grid item xs={7}><Typography>{stats.last_name}</Typography></Grid>
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
      
      <Typography variant='h6'>Change Password</Typography>
      <Box sx={{ display: 'flex', justifyContent: 'center', borderRadius: '4px', pb: 4 }}>
        <ChangePassword user={user}/>
      </Box>
    </React.Fragment>
  );
}
export default Profile; 