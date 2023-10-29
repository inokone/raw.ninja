import * as React from 'react';
import AppStats from './AppStats';
import UserTable from './UserTable';
import RoleTable from './RoleTable';
import { Typography } from "@mui/material";


const Admin = () => {

    return (
        <React.Fragment>
            <AppStats />
            <Typography variant="h4">Users</Typography>
            <UserTable />
            <Typography variant="h4">Roles</Typography>
            <RoleTable />
        </React.Fragment>
    )
}

export default Admin;