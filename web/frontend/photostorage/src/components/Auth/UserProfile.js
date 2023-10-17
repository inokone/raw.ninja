import React from 'react';
import { CircularProgress, Alert, Grid, Box, Typography } from "@mui/material";

const { REACT_APP_API_PREFIX } = process.env;

const Profile = () => {
  const [error, setError] = React.useState(null)
  const [loading, setLoading] = React.useState(true)
  const [stats, setStats] = React.useState(null)

  const formatBytes = (bytes, decimals = 2) => {
    if (!+bytes) return '0 Bytes'

    const k = 1024
    const dm = decimals < 0 ? 0 : decimals
    const sizes = ['Bytes', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB']

    const i = Math.floor(Math.log(bytes) / Math.log(k))

    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`
  }

  React.useEffect(() => {
    const loadStats = () => {
      fetch(REACT_APP_API_PREFIX + '/api/v1/statistics/user', {
        method: "GET",
        mode: "cors",
        credentials: "include"
      })
        .then(response => {
          if (!response.ok) {
            if (response.status !== 200) {
              setError(response.status + ": " + response.statusText);
            } else {
              response.json().then(content => setError(content.message))
            }
            setLoading(false)
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

    loadStats()
  },)

  return (
    <>
      {error !== null ? <Alert sx={{ mb: 4 }} severity="error">{error}</Alert> : null}
      {loading ? <CircularProgress mt={10} /> :
        <Box sx={{ display: 'flex', justifyContent: 'center', borderRadius: '4px', pb: 4 }}>
          <Box sx={{ bgcolor: 'rgba(0, 0, 0, 0.34)', color: 'white', mt: 10, borderRadius: '4px', width: '500px' }}>
            <Grid container>
              <Grid xs={12}><Typography variant='h6' sx={{ borderRadius: '4px', bgcolor: 'rgba(0, 0, 0, 0.54)' }}>Profile</Typography></Grid>
            </Grid>
            <Grid container>
              <Grid xs={5}><Typography>E-mail:</Typography></Grid>
              <Grid xs={7}><Typography>{stats.email}</Typography></Grid>
            </Grid>
            <Grid container>
              <Grid xs={5}><Typography>Phone:</Typography></Grid>
              <Grid xs={7}><Typography>{stats.phone}</Typography></Grid>
            </Grid>
            <Grid container>
              <Grid xs={5}><Typography>Registered on:</Typography></Grid>
              <Grid xs={7}><Typography>{new Date(stats.registration_date * 1000).toLocaleDateString()}</Typography></Grid>
            </Grid>
            <Grid container>
              <Grid xs={5}><Typography>Photo count:</Typography></Grid>
              <Grid xs={7}><Typography>{stats.photos}</Typography></Grid>
            </Grid>
            <Grid container>
              <Grid xs={5}><Typography>Favorites:</Typography></Grid>
              <Grid xs={7}><Typography>{stats.favorites}</Typography></Grid>
            </Grid>
            <Grid container>
              <Grid xs={5}><Typography>User space:</Typography></Grid>
              <Grid xs={7}><Typography>{formatBytes(stats.used_space)}</Typography></Grid>
            </Grid>
          </Box>
        </Box>}
    </>
  );
}
export default Profile; 