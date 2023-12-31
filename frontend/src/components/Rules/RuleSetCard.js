import * as React from 'react';
import { Card, Typography, CardContent } from "@mui/material";


const RuleSetCard = ({ ruleSet, onClick }) => {

    return (
        <Card style={{ flex: 1 }} sx={{ position: 'relative', cursor: "pointer", margin: 1, bgcolor: "lightgrey" }} onClick={onClick}>
            <CardContent>
                <Typography sx={{ fontSize: 14 }} color="text.secondary">{ruleSet.name}</Typography>
                <Typography variant="body2">{ruleSet.description}</Typography>
            </CardContent>
        </Card>
    );
}

export default RuleSetCard;