import React from 'react';
import "./NotFoundPage.css";
import { Box, Button } from "@mui/material";

const NotFoundPage = () => {

return (
   <div id="error-page">
      <div class="content">
         <h2 class="header" data-text="404">
            404
         </h2>
         <h4 data-text="Opps! Page not found">
            Opps! Page not found
         </h4>
         <p>
            Sorry, the page you're looking for doesn't exist. If you think something is broken, report a problem.
         </p>
         <Box sx={{ flexGrow: 0, m: 2 }}>
            <Button sx={{m: 1}} variant="outlined" color="primary">Return home</Button>
            <Button sx={{m: 1}} variant="outlined" color="secondary">Report problem</Button>
         </Box>
      </div>
   </div>);
}

export default NotFoundPage;
