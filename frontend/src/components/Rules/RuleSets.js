import * as React from 'react';
import { useNavigate } from "react-router-dom";
import { Box, Fab, Typography, Alert, Grid } from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import ProgressDisplay from '../Common/ProgressDisplay';
import RuleSetCard from './RuleSetCard';
import RuleSetDialog from './RuleSetDialog';
import DeleteDialog from '../Common/DeleteDialog';
import RuleDocs from '../Docs/RuleDocs';

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const RuleSets = () => {
    const navigate = useNavigate()
    const [loading, setLoading] = React.useState(false)
    const [error, setError] = React.useState(null)
    const [ruleSets, setRuleSets] = React.useState(null)
    const [isEditDialogOpen, setIsEditDialogOpen] = React.useState(false);
    const [isDeleteDialogOpen, setIsDeleteDialogOpen] = React.useState(false);
    const [deleteTarget, setDeleteTarget] = React.useState(null);

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

    const createRuleSet = React.useCallback((ruleSet) => {
        setLoading(true)
        fetch(REACT_APP_API_PREFIX + '/api/v1/rulesets/', {
            method: "POST",
            mode: "cors",
            credentials: "include",
            body: JSON.stringify(ruleSet)
        }).then(response => {
            if (!response.ok) {
                response.json().then(content => {
                    setError(content.message)
                    setLoading(false)
                });
            } else {
                response.json().then(content => {
                    console.log(ruleSets)
                    let newRuleSets = ruleSets.slice()
                    newRuleSets.push(content)
                    setRuleSets(newRuleSets)
                    setLoading(false)
                })
            }
        }).catch(error => {
            setError(error.message)
            setLoading(false)
        });
    }, [ruleSets])

    const deleteRuleSet = React.useCallback((id) => {
        setLoading(true)
        fetch(REACT_APP_API_PREFIX + '/api/v1/rulesets/' + id, {
            method: "DELETE",
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
                    setLoading(false)
                    navigate(0)
                })
            }
        }).catch(error => {
            setError(error.message)
            setLoading(false)
        });
    }, [navigate])

    const handleEditDialogOpen = React.useCallback(() => {
        setIsEditDialogOpen(true);
    }, [setIsEditDialogOpen]);

    const handleEditDialogClose = React.useCallback(() => {
        setIsEditDialogOpen(false);
    }, [setIsEditDialogOpen]);

    const handleEditDialogSave = React.useCallback((ruleSet) => {
        setIsEditDialogOpen(false);
        createRuleSet(ruleSet)
    }, [setIsEditDialogOpen, createRuleSet]);

    const onDeleteClick = React.useCallback((ruleSet) => {
        setDeleteTarget(ruleSet)
        setIsDeleteDialogOpen(true);
    }, [setIsDeleteDialogOpen]);

    const handleDeleteDialogClose = React.useCallback(() => {
        setDeleteTarget(null);
        setIsDeleteDialogOpen(false);
    }, [setIsDeleteDialogOpen]);

    const handleDeleteDialogAccept = React.useCallback(() => {
        deleteRuleSet(deleteTarget)
        setIsDeleteDialogOpen(false);
    }, [setIsDeleteDialogOpen, deleteRuleSet, deleteTarget]);

    const onRuleSetClick = (id) => {
        navigate("/rulesets/" + id)
    }

    const onFabClick = (id) => {
        handleEditDialogOpen()
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
            <RuleSetDialog open={isEditDialogOpen} onCancel={handleEditDialogClose} onSave={handleEditDialogSave} />
            <DeleteDialog open={isDeleteDialogOpen} onCancel={handleDeleteDialogClose} onDelete={handleDeleteDialogAccept} name="this rule set" />
            {ruleSets && ruleSets.length === 0 && 
                <RuleDocs />
            }
            {!loading && ruleSets !== null && ruleSets.length > 0 &&
                <>
                    <Typography variant='h4' sx={{ marginBottom: 4, marginTop: 2 }} >Lifecycle Rule Sets</Typography>
                    <Grid container spacing={2} sx={{ padding: 1 }}>
                        {ruleSets.map((ruleSet) => {
                            return (<Grid item key={ruleSet.id} xs={12} sm={6} md={4} lg={2} xl={1}><RuleSetCard ruleSet={ruleSet} onClick={() => onRuleSetClick(ruleSet.id)} onDelete={() => onDeleteClick(ruleSet.id)} /></Grid>)
                        })}
                    </Grid>
                </>}
            {!loading &&
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
            }
        </>
    )
}

export default RuleSets;