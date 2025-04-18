import React from 'react';
import { Alert, Grid, Box, Typography } from "@mui/material";
import ProgressDisplay from '../Common/ProgressDisplay';

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const AppStats = () => {
    const [error, setError] = React.useState(null)
    const [loading, setLoading] = React.useState(false)
    const [stats, setStats] = React.useState(null)

    const formatBytes = (bytes, decimals = 2) => {
        if (!+bytes) return '0 Bytes'

        const k = 1024
        const dm = decimals < 0 ? 0 : decimals
        const sizes = ['Bytes', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB']

        const i = Math.floor(Math.log(bytes) / Math.log(k))

        return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`
    }

    const loadStats = () => {
        setLoading(true)
        fetch(REACT_APP_API_PREFIX + '/api/v1/statistics/app', {
            method: "GET",
            mode: "cors",
            credentials: "include"
        })
            .then(response => {
                if (!response.ok) {
                    response.json().then(content => {
                        setError(content.message)
                        setLoading(false)
                    });
                } else {
                    response.json().then(content => {
                        setStats(content)
                        setLoading(false)
                    })
                }
            })
            .catch(error => {
                setError(error.message)
                setLoading(false)
            });
    }

    React.useEffect(() => {
        if (!loading && !stats && !error) {
            loadStats()
        }
    }, [stats, loading, error])

    return (
        <React.Fragment>
            {error && <Alert sx={{ mb: 4 }} severity="error">{error}</Alert>}
            {loading && <ProgressDisplay />}
            {stats &&
                <Box sx={{ display: 'flex', justifyContent: 'center', borderRadius: '4px' }}>
                    <Box sx={{ bgcolor: 'rgba(0, 0, 0, 0.34)', color: 'white', mt: 10, borderRadius: '4px', width: '500px' }}>
                        <Grid container>
                            <Grid item xs={12}><Typography variant='h6' sx={{ borderRadius: '4px', bgcolor: 'rgba(0, 0, 0, 0.54)' }}>General statistics</Typography></Grid>
                        </Grid>
                        <Grid container>
                            <Grid item xs={5}><Typography>Total user count:</Typography></Grid>
                            <Grid item xs={7}><Typography>{stats.total_users}</Typography></Grid>
                        </Grid>
                        {stats.user_distribution.map((entry) => {
                            return (<Grid container key={entry.role}>
                                <Grid item xs={5}><Typography>{entry.role} users:</Typography></Grid>
                                <Grid item xs={7}><Typography>{entry.users}</Typography></Grid>
                            </Grid>)
                        })}
                        <Grid container>
                            <Grid item xs={5}><Typography>Photo count:</Typography></Grid>
                            <Grid item xs={7}><Typography>{stats.photos}</Typography></Grid>
                        </Grid>
                        <Grid container>
                            <Grid item xs={5}><Typography>Favorites:</Typography></Grid>
                            <Grid item xs={7}><Typography>{stats.favorites}</Typography></Grid>
                        </Grid>
                        <Grid container>
                            <Grid item xs={5}><Typography>Albums:</Typography></Grid>
                            <Grid item xs={7}><Typography>{stats.albums}</Typography></Grid>
                        </Grid>
                        <Grid container>
                            <Grid item xs={5}><Typography>Upload batches:</Typography></Grid>
                            <Grid item xs={7}><Typography>{stats.uploads}</Typography></Grid>
                        </Grid>
                        <Grid container>
                            <Grid item xs={5}><Typography>Used space of quota:</Typography></Grid>
                            <Grid item xs={7}><Typography>{formatBytes(stats.used_space) + " / " + formatBytes(stats.quota)}</Typography></Grid>
                        </Grid>
                    </Box>
                </Box>}
        </React.Fragment>
    );
}

export default AppStats