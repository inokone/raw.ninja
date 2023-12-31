import * as React from 'react';
import { useNavigate } from "react-router-dom";
import { Box, Fab, Typography, Alert, Grid } from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import ProgressDisplay from '../Common/ProgressDisplay';
import RuleSetCard from './RuleSetCard';

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const RuleSets = ({ user }) => {
    const navigate = useNavigate()
    const [loading, setLoading] = React.useState(false)
    const [error, setError] = React.useState(null)
    const [ruleSets, setRuleSets] = React.useState(null)

    const populate = () => {
        setLoading(true)
        fetch(REACT_APP_API_PREFIX + '/api/v1/rulesets/', {
            method: "GET",
            mode: "cors",
            credentials: "include"
        }).then(response => {
            if (!response.ok) {
                response.json().then(content => {
                    setError(content.message)
                    setLoading(false)
                });
            } else {
                response.json().then(content => {
                    if (content === null) {
                        content = []
                    }
                    setRuleSets(content)
                    setLoading(false)
                })
            }
        }).catch(error => {
            setError(error.message)
            setLoading(false)
        });
    }

    const onRuleSetClick = (id) => {
        navigate("/rulesets/" + id)
    }

    const onFabClick = (id) => {
        navigate("/rulesets/create")
    }

    React.useEffect(() => {
        if (!loading && !ruleSets && !error) {
            populate()
        }
    }, [ruleSets, loading, error])

    return (
        <>
            {error && <Alert sx={{ mb: 1 }} onClose={() => setError(null)} severity="error">{error}</Alert>}
            {loading && <ProgressDisplay />}
            {ruleSets !== null &&
                <>
                    <Typography variant='h4' sx={{ marginBottom: 4, marginTop: 2 }} >Rule Sets</Typography>
                    <Grid container>
                        {ruleSets.map((ruleSet) => {
                            return (<Grid item key={ruleSet.id} xs={6} md={4} lg={2} xl={2}><RuleSetCard ruleSet={ruleSet} onClick={onRuleSetClick} /></Grid>)
                        })}
                    </Grid>
                    <Box sx={{
                        '& > :not(style)': { m: 1 },
                        position: "fixed",
                        bottom: 16,
                        right: 16
                    }}>
                        <Fab onClick={onFabClick} color="primary" aria-label="add">
                            <AddIcon />
                        </Fab>
                    </Box>
                </>}
        </>
    )
}

export default RuleSets;