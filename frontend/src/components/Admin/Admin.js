import * as React from 'react';
import AppStats from './AppStats';
import UserTable from './UserTable';
import RoleTable from './RoleTable';
import { Typography, Box } from "@mui/material";


const Admin = ({user}) => {

    return (
        <React.Fragment>
            {user &&
                <React.Fragment>
                    <AppStats />
                    <Box>
                        <Typography variant="h4" sx={{ marginTop: 6, marginBottom: 2 }}>Users</Typography>
                        <Box maxWidth="lg" sx={{ margin: 'auto' }}>
                            <UserTable/>
                        </Box>
                        <Typography variant="h4" sx={{marginTop: 6, marginBottom: 2}}>Roles</Typography>
                        <Box maxWidth="sm" sx={{ margin: 'auto' }}>
                            <RoleTable />
                        </Box>
                    </Box>
                </React.Fragment>
            }
        </React.Fragment>
    )
}

export default Admin;