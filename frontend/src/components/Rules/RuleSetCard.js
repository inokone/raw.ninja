import * as React from 'react';
import { Card, Typography, CardContent, CardActions, Button } from "@mui/material";
import { useTheme } from '@mui/styles';

const RuleSetCard = ({ ruleSet, onClick, onDelete }) => {
    const theme = useTheme();

    return (
        <Card style={{ flex: 1 }} sx={{
            position: 'relative', 
            cursor: "pointer", 
            backgroundColor: 'white',
            color: 'text.secondary',
            height: "100%",
            display: "flex",
            flexDirection: "column",
        }}>
            <Typography sx={{ cursor: 'pointer', backgroundColor: theme.palette.secondary.light, color: theme.palette.common.white, padding: 0.5 }}
                variant="h6"
            >{ruleSet.name}</Typography>
            <CardContent onClick={onClick} sx={{ cursor: 'pointer', display: 'flex', flexDirection: 'column', height: "100%", }}>
                <Typography sx={{ mb: 2 }} variant="body2" component="div">{ruleSet.description}</Typography>
            </CardContent>
            <CardActions disableSpacing sx={{ mt: "auto"}}>
                <Button color="secondary" onClick={onDelete}>
                    Delete
                </Button>
                <Typography variant="body1" sx={{ marginLeft: "auto" }}>{ruleSet.rules.length} {ruleSet.rules.length === 1? "rule": "rules"}</Typography>
            </CardActions>
        </Card>
    );
}

export default RuleSetCard;