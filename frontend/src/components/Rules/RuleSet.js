import * as React from 'react';
import { useLocation } from "react-router-dom";
import { Box, Fab, Typography, Alert, Stack } from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import ProgressDisplay from '../Common/ProgressDisplay';
import RuleCard from './RuleCard';
import RuleDialog from './RuleDialog';
import DeleteDialog from '../Common/DeleteDialog';

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const RuleSets = ({ user }) => {
    const [loading, setLoading] = React.useState(false)
    const [error, setError] = React.useState(null)
    const [ruleSet, setRuleSet] = React.useState(null)
    const [constants, setConstants] = React.useState(null)
    const [isEditDialogOpen, setIsEditDialogOpen] = React.useState(false);
    const [isDeleteDialogOpen, setIsDeleteDialogOpen] = React.useState(false);
    const [deleteTarget, setDeleteTarget] = React.useState(null);
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
    }, [])

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
    }, [path])

    const updateRuleSet = React.useCallback((ruleSet) => {
        setLoading(true)
        fetch(REACT_APP_API_PREFIX + '/api/v1/rulesets/' + ruleSet.id, {
            method: "PUT",
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
                    setRuleSet(content)
                    setLoading(false)
                })
            }
        }).catch(error => {
            setError(error.message)
            setLoading(false)
        });
    }, [setRuleSet])

    const deleteRule = React.useCallback((index) => {
        let newRuleSet = ruleSet
        let newRules = [...newRuleSet.rules.slice(0, index), ...newRuleSet.rules.slice(index + 1)]
        newRuleSet.rules = newRules
        updateRuleSet(newRuleSet)
    }, [updateRuleSet, ruleSet]);

    const handleEditDialogOpen = React.useCallback(() => {
        setIsEditDialogOpen(true);
    }, [setIsEditDialogOpen]);

    const handleEditDialogClose = React.useCallback(() => {
        setIsEditDialogOpen(false);
    }, [setIsEditDialogOpen]);

    const onDeleteClick = React.useCallback((index) => {
        setDeleteTarget(index)
        setIsDeleteDialogOpen(true);
    }, [setIsDeleteDialogOpen]);

    const handleDeleteDialogClose = React.useCallback(() => {
        setDeleteTarget(null);
        setIsDeleteDialogOpen(false);
    }, [setIsDeleteDialogOpen]);

    const handleDeleteDialogAccept = React.useCallback(() => {
        deleteRule(deleteTarget)
        setIsDeleteDialogOpen(false);
    }, [setIsDeleteDialogOpen, deleteRule, deleteTarget]);

    const handleEditDialogSave = React.useCallback((rule) => {
        setIsEditDialogOpen(false);
        let newRuleSet = ruleSet
        let newRules = newRuleSet.rules.slice()
        newRules.push(rule)
        newRuleSet.rules = newRules
        updateRuleSet(ruleSet)
    }, [setIsEditDialogOpen, updateRuleSet, ruleSet]);

    const onFabClick = React.useCallback(() => {
        handleEditDialogOpen()
    }, [handleEditDialogOpen]);

    const setRule = React.useCallback((rule, idx) => {
        let newRuleSet = ruleSet
        let newRules = newRuleSet.rules.slice()
        newRules[idx] = rule
        newRuleSet.rules = newRules
        updateRuleSet(newRuleSet)
    }, [updateRuleSet, ruleSet])

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
            {ruleSet !== null && constants !== null && !loading &&
                <>
                    <RuleDialog open={isEditDialogOpen} onCancel={handleEditDialogClose} onSave={handleEditDialogSave} constants={constants} />
                    <DeleteDialog open={isDeleteDialogOpen} onCancel={handleDeleteDialogClose} onDelete={handleDeleteDialogAccept} name="this rule" />
                    <Typography variant='h4' sx={{ marginBottom: 4, marginTop: 2 }} >{ruleSet.name}</Typography>
                    <Stack direction="column" alignItems="center" justifyContent="center">
                        {ruleSet.rules.map((rule, idx) => {
                            return (
                                <Box key={rule.id} spacing={{ xs: 1, sm: 2 }} maxWidth='sm'>
                                    <RuleCard rule={rule} constants={constants} onDelete={() => onDeleteClick(idx)} setRule={(updated) => setRule(updated, idx)} />
                                </Box>)
                        })}
                    </Stack>
                    <Box sx={{
                        '& > :not(style)': { m: 1 },
                        position: "fixed",
                        bottom: 16,
                        right: 16
                    }}>
                        {ruleSet.rules.length < maxRuleCount &&
                            <Fab onClick={() => onFabClick(ruleSet)} color="primary" aria-label="add">
                                <AddIcon />
                            </Fab>
                        }
                    </Box>
                </>}
        </>
    )
}

export default RuleSets;