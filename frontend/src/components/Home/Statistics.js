import * as React from 'react';
import { Grid, Alert } from '@mui/material';
import ProgressDisplay from '../Common/ProgressDisplay';
import SpaceChart from './SpaceChart';
import AggregatedChart from './AggregatedChart';
import UploadChart from './UploadChart';
import { makeStyles } from '@mui/styles';

const useStyles = makeStyles((theme) => ({
    chart: {
        bgcolor: theme.palette.primary.light,
        borderRadius: '4px'
    }
}));

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const Statistics = ({ onDataLoaded }) => {
    const classes = useStyles();
    const [loading, setLoading] = React.useState(false)
    const [error, setError] = React.useState(null)
    const [stats, setStats] = React.useState(null)

    const populate = React.useCallback(() => {
        fetch(REACT_APP_API_PREFIX + '/api/v1/statistics/user', {
            method: "GET",
            mode: "cors",
            credentials: "include"
        }).then(response => {
            if (!response.ok) {
                response.json().then(content => {
                    setError(content.message)
                });
            } else {
                response.json().then(content => {
                    setLoading(false)
                    setStats(content)
                    if (onDataLoaded) {
                        onDataLoaded(content)
                    }
                })
            }
        }).catch(error => {
            setError(error.message)
            setLoading(false)
        });
    }, [onDataLoaded])

    React.useEffect(() => {
        if (!loading && !stats && !error) {
            populate()
        }
    }, [stats, loading, error, populate])

    return (
        <>
            {error && <Alert sx={{ mb: 1 }} onClose={() => setError(null)} severity="error">{error}</Alert>}
            {loading && <ProgressDisplay />}
            {stats !== null &&
                <Grid container spacing={2} alignContent={'baseline'}>
                    <Grid item xs={12} md={4} display="flex" justifyContent="center" alignItems="center" className={classes.chart}>
                        <SpaceChart usedSpace={stats.used_space} quota={stats.quota} />
                    </Grid>
                    <Grid item xs={12} md={4} display="flex" justifyContent="center" alignItems="center" className={classes.chart}>
                        <AggregatedChart photos={stats.photos} favorites={stats.favorites} albums={stats.albums} />
                    </Grid>
                    <Grid item xs={12} md={4} display="flex" justifyContent="center" alignItems="center" className={classes.chart}>
                        <UploadChart uploads={stats.uploads} />
                    </Grid>
                </Grid>}
        </>
    )
}

export default Statistics;