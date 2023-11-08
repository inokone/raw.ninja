import React from "react";
import { Grid, Typography } from "@mui/material";
import CodeIcon from "@mui/icons-material/Code";
import BuildIcon from "@mui/icons-material/Build";
import ComputerIcon from "@mui/icons-material/Computer";
import SecurityIcon from "@mui/icons-material/Security";
import AccessibilityIcon from "@mui/icons-material/Accessibility";
import CloudIcon from "@mui/icons-material/Cloud";
import calculateSpacing from "./calculateSpacing";
import useMediaQuery from "@mui/material/useMediaQuery";
import { withTheme } from "@mui/styles";
import FeatureCard from "./FeatureCard";
import useWidth from "./useWidth";

const iconSize = 30;

const features = [
    {
        color: "#00C853",
        headline: "Store image files",
        text: "You can store your RAW and processed images in the cloud anytime, anywhere. No more hassle with SD cards and pendrives.",
        icon: <CloudIcon style={{ fontSize: iconSize }} />,
        mdDelay: "0",
        smDelay: "0",
    },
    {
        color: "#6200EA",
        headline: "Wide support for RAW formats",
        text: "PhotoStore supports most camera RAW fromats including iPhone, Canon, Nikon, Sony, Fuji, Olympus and DJI cameras.",
        icon: <CodeIcon style={{ fontSize: iconSize }} />,
        mdDelay: "200",
        smDelay: "200",
    },
    {
        color: "#0091EA",
        headline: "Edit RAW online",
        text: "You can edit your photos on the web. Pre-processing your RAW photos has never been easier.",
        icon: <ComputerIcon style={{ fontSize: iconSize }} />,
        mdDelay: "400",
        smDelay: "0",
    },
    {
        color: "#d50000",
        headline: "Secure storage",
        text: "Your files are stored safe and encrypted.",
        icon: <SecurityIcon style={{ fontSize: iconSize }} />,
        mdDelay: "0",
        smDelay: "200",
    },
    {
        color: "#DD2C00",
        headline: "Automate lifecycle",
        text: "Set your own lifecycle rules for your images. You can automate achiving and disposing unnecessary images.",
        icon: <BuildIcon style={{ fontSize: iconSize }} />,
        mdDelay: "200",
        smDelay: "0",
    },
    {
        color: "#64DD17",
        headline: "Easy to start, easy to use",
        text: "You can register your free account right now and test all features.",
        icon: <AccessibilityIcon style={{ fontSize: iconSize }} />,
        mdDelay: "400",
        smDelay: "200",
    },
];

function FeatureSection(props) {
    const { theme } = props;
    const width = useWidth();
    const isWidthUpMd = useMediaQuery(theme.breakpoints.up("md"));

    return (
        <div style={{ backgroundColor: "#FFFFFF" }}>
            <div className="container-fluid lg-p-top">
                <Typography variant="h3" align="center" className="lg-mg-bottom">
                    Features
                </Typography>
                <div className="container-fluid">
                    <Grid container spacing={calculateSpacing(width, theme)}>
                        {features.map((element) => (
                            <Grid
                                item
                                xs={6}
                                md={4}
                                data-aos="zoom-in-up"
                                data-aos-delay={isWidthUpMd ? element.mdDelay : element.smDelay}
                                key={element.headline}
                            >
                                <FeatureCard
                                    Icon={element.icon}
                                    color={element.color}
                                    headline={element.headline}
                                    text={element.text}
                                />
                            </Grid>
                        ))}
                    </Grid>
                </div>
            </div>
        </div>
    );
}

FeatureSection.propTypes = {};

export default withTheme(FeatureSection);