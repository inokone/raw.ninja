import * as React from 'react';
import { useLocation } from "react-router-dom";
import { Box, Fab, Typography, Alert, Grid } from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import ProgressDisplay from '../Common/ProgressDisplay';
import RuleCard from './RuleCard';

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const RuleSets = ({ user }) => {
    const [loading, setLoading] = React.useState(false)
    const [error, setError] = React.useState(null)
    const [ruleSet, setRuleSet] = React.useState(null)
    const [constants, setConstants] = React.useState(null)
    const location = useLocation()
    const path = location.pathname
    const maxRuleCount = 3

    const populateConstants = React.useCallback(() => {
        setLoading(true)
        fetch(REACT_APP_API_PREFIX + '/api/v1/rules/constants', {
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
                    setConstants(content)
                    setLoading(false)
                })
            }
        }).catch(error => {
            setError(error.message)
            setLoading(false)
        });
    },[])

    const populate = React.useCallback(() => {
        setLoading(true)
        fetch(REACT_APP_API_PREFIX + '/api/v1' + path, {
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
                    setRuleSet(content)
                    setLoading(false)
                })
            }
        }).catch(error => {
            setError(error.message)
            setLoading(false)
        });
    },[path])


    const onFabClick = (id) => {
        let newRuleSet = ruleSet.slice()
        newRuleSet.push({
            name: "New Rule",
            description: "Empty rule"
        })
        setRuleSet(newRuleSet)
    }

    const setRule = (rule, idx) => {
        let newRuleSet = ruleSet.slice()
        newRuleSet[idx] = rule
        setRuleSet(newRuleSet)
    }

    React.useEffect(() => {
        if (!loading && !error) {
            if (!ruleSet) {
                populate()
            }
            if (!constants) {
                populateConstants()
            }
        }
    }, [ruleSet, constants, loading, error, populate, populateConstants])

    return (
        <>
            {error && <Alert sx={{ mb: 1 }} onClose={() => setError(null)} severity="error">{error}</Alert>}
            {loading && <ProgressDisplay />}
            {ruleSet !== null && constants !== null &&
                <>
                    <Typography variant='h4' sx={{ marginBottom: 4, marginTop: 2 }} >Rule Set {ruleSet.name}</Typography>
                    <Grid container>
                        {ruleSet.rules.map((rule, idx) => {
                            return (<Grid item key={rule.id} xs={6} md={4} lg={2} xl={2}><RuleCard rule={rule} setRule={() => setRule(rule, idx)}/></Grid>)
                        })}
                    </Grid>
                    <Box sx={{
                        '& > :not(style)': { m: 1 },
                        position: "fixed",
                        bottom: 16,
                        right: 16
                    }}>
                        {ruleSet.rules.length < maxRuleCount &&
                            <Fab onClick={onFabClick} color="primary" aria-label="add">
                                <AddIcon />
                            </Fab>
                        }
                    </Box>
                </>}
        </>
    )
}

export default RuleSets;