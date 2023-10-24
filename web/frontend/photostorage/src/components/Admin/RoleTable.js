import React from 'react';
import { styled } from '@mui/material/styles';

import { Alert, Paper, TableContainer, TableBody, Table, TableHead, TableRow, Box } from "@mui/material";
import TableCell, { tableCellClasses } from '@mui/material/TableCell';

import ProgressDisplay from '../Common/ProgressDisplay';

const { REACT_APP_API_PREFIX } = process.env;

const StyledTableCell = styled(TableCell)(({ theme }) => ({
    [`&.${tableCellClasses.head}`]: {
        backgroundColor: theme.palette.common.black,
        color: theme.palette.common.white,
    },
    [`&.${tableCellClasses.body}`]: {
        fontSize: 14,
    },
}));

const StyledTableRow = styled(TableRow)(({ theme }) => ({
    '&:nth-of-type(odd)': {
        backgroundColor: theme.palette.action.hover,
    },
    // hide last border
    '&:last-child td, &:last-child th': {
        border: 0,
    },
}));

const RoleTable = () => {
    const [error, setError] = React.useState(null)
    const [loading, setLoading] = React.useState(true)
    const [roles, setRoles] = React.useState(null)

    const formatBytes = (bytes, decimals = 2) => {
        if (!+bytes) return '0 Bytes'

        const k = 1024
        const dm = decimals < 0 ? 0 : decimals
        const sizes = ['Bytes', 'KBytes', 'MBytes', 'GBytes', 'TBytes', 'PBytes']

        const i = Math.floor(Math.log(bytes) / Math.log(k))

        return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`
    }

    React.useEffect(() => {
        const loadRoles = () => {
            fetch(REACT_APP_API_PREFIX + '/api/v1/roles/', {
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
                            setRoles(content)
                        })
                    }
                })
                .catch(error => {
                    setError(error)
                    setLoading(false)
                });
        }

        if (!roles) {
            loadRoles()
        }
    },)

    return (
        <>
            {error !== null ? <Alert sx={{ mb: 4 }} severity="error">{error}</Alert> : null}
            {loading ? <ProgressDisplay /> :
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
                                        <StyledTableRow>
                                            <StyledTableCell>{role.name}</StyledTableCell>
                                            <StyledTableCell>{role.quota <= 0 ? 'Unlimited' : formatBytes(role.quota)}</StyledTableCell>
                                        </StyledTableRow>
                                    )
                                })}
                            </TableBody>
                        </Table>
                    </TableContainer>
                </Box>
            }
        </>
    );
}

export default RoleTable