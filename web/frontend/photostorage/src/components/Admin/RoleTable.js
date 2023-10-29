import React from 'react';
import { styled } from '@mui/material/styles';

import { Alert, Paper, TableContainer, TableBody, Table, TableHead, TableRow, Box } from "@mui/material";
import TableCell, { tableCellClasses } from '@mui/material/TableCell';

import ProgressDisplay from '../Common/ProgressDisplay';
import EditableTableCell from './EditableTableCell';

const { REACT_APP_API_PREFIX } = process.env;

const StyledTableRow = styled(TableRow)(({ theme }) => ({
    '&:nth-of-type(odd)': {
        backgroundColor: theme.palette.action.hover,
    },
    // hide last border
    '&:last-child td, &:last-child th': {
        border: 0,
    },
}));

const StyledTableCell = styled(TableCell)(({ theme }) => ({
    [`&.${tableCellClasses.head}`]: {
        backgroundColor: theme.palette.common.black,
        color: theme.palette.common.white,
    },
    [`&.${tableCellClasses.body}`]: {
        fontSize: 14,
    },
}));

const RoleTable = () => {
    const [error, setError] = React.useState(null)
    const [loading, setLoading] = React.useState(false)
    const [roles, setRoles] = React.useState(null)

    const formatBytes = (bytes, decimals = 2) => {
        if (!+bytes) return '0 Bytes'

        const k = 1024
        const dm = decimals < 0 ? 0 : decimals
        const sizes = ['Bytes', 'KBytes', 'MBytes', 'GBytes', 'TBytes', 'PBytes']

        const i = Math.floor(Math.log(bytes) / Math.log(k))

        return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`
    }

    const handleCellEdit = (role, field, value) => {
        let newRoles = roles.slice()
        for (let i = 0; i < newRoles.length; i++) {
            if (newRoles[i].id === role.id) {
                newRoles[i][field] = value
            }
        }
        setRoles(newRoles)
        return fetch(REACT_APP_API_PREFIX + '/api/v1/roles/' + role.id, {
            method: "PATCH",
            mode: "cors",
            credentials: "include",
            body: JSON.stringify({
                'id': role.id,
                field: value
            })
        })
    }

    const loadRoles = () => {
        setLoading(true)
        fetch(REACT_APP_API_PREFIX + '/api/v1/roles/', {
            method: "GET",
            mode: "cors",
            credentials: "include"
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error(response.status + ": " + response.statusText);
                } else {
                    response.json().then(content => {
                        setRoles(content)
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
        if (!loading && !roles && !error) {
            loadRoles()
        }
    }, [roles, error, loading])

    return (
        <React.Fragment>
            {error && <Alert sx={{ mb: 4 }} severity="error">{error}</Alert>}
            {loading && <ProgressDisplay />}
            {roles &&
                <Box sx={{ display: 'flex', justifyContent: 'center', borderRadius: '4px', pb: 4 }}>
                    <TableContainer component={Paper} style={{ width: 1200 }}>
                        <Table style={{ width: 1200 }}>
                            <TableHead>
                                <TableRow>
                                    <StyledTableCell>Name</StyledTableCell>
                                    <StyledTableCell>Quota</StyledTableCell>
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {roles.map((role) => {
                                    return (
                                        <StyledTableRow key={role.id}>
                                            <EditableTableCell
                                                value={role.name}
                                                formatter={(value) => { return value }}
                                                onCellEdit={(value) => handleCellEdit(role, 'name', value)}>
                                            </EditableTableCell>
                                            <EditableTableCell
                                                value={role.quota}
                                                formatter={(value) => { return value <= 0 ? 'Unlimited' : formatBytes(value) }}
                                                onCellEdit={(value) => handleCellEdit(role, 'quota', value)}>
                                            </EditableTableCell>
                                        </StyledTableRow>
                                    )
                                })}
                            </TableBody>
                        </Table>
                    </TableContainer>
                </Box>}
        </React.Fragment>
    );
}

export default RoleTable