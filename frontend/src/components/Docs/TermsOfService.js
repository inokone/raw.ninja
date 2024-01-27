import React from 'react';
import { Paper, Typography, Box } from '@mui/material';
import './TermsOfService.css'

const TermsOfService = () => {
    return (
        <Box sx={{ maxWidth: 'md', display: 'flex', mx: 'auto' }}>
            <Paper sx={{ textAlign: 'left', margin: '25px' }}>
                <Typography variant='h4' sx={{ textAlign: 'center', margin: '20px' }}>RAW.Ninja<br />Terms of Service</Typography>
                <p>These Terms of Service ("Terms") govern your use of the RAW.Ninja ("App") provided by RAW.Ninja LLC. ("Company"). By accessing or using the App, you agree to be bound by these Terms.</p>
                <ol>
                    <li><strong>License and Use:</strong> Subject to your compliance with these Terms, Company grants you a limited, non-exclusive, non-transferable, revocable license to use the App solely for your personal or business purposes.</li>
                    <li><strong>User Conduct:</strong> You agree not to engage in any of the following activities:
                        <ul>
                            <li>Violating any applicable laws or regulations.</li>
                            <li>Uploading or distributing any harmful code, malware, or unauthorized content.</li>
                            <li>Attempting to gain unauthorized access to the App or its related systems.</li>
                        </ul>
                    </li>
                    <li><strong>Data Responsibility:</strong> You understand and agree that Company is not responsible for any loss, corruption, or unauthorized access to your data stored on the App. It is your responsibility to regularly back up your data and implement necessary security measures.</li>
                    <li><strong>Intellectual Property:</strong> The App and its content, features, and functionality are owned by Company and are protected by intellectual property laws. You may not reproduce, distribute, modify, or create derivative works based on the App without Company's prior written consent.</li>
                    <li><strong>Disclaimer of Warranties:</strong> The App is provided "as is" and "as available" without any warranties of any kind, either expressed or implied. Company does not warrant that the App will be error-free, secure, or uninterrupted.</li>
                    <li><strong>Limitation of Liability:</strong> In no event shall Company be liable for any direct, indirect, incidental, special, consequential, or punitive damages arising out of or in connection with your use of the App. This includes but is not limited to loss of data, profits, or business interruption.</li>
                    <li><strong>Indemnification:</strong> You agree to indemnify and hold harmless Company from any claims, damages, losses, liabilities, costs, and expenses arising out of or in connection with your use of the App.</li>
                    <li><strong>Termination:</strong> Company may terminate your access to the App at any time without notice. Upon termination, you must cease using the App.</li>
                    <li><strong>Changes to Terms:</strong> Company reserves the right to modify or revise these Terms at any time. Your continued use of the App after any changes will signify your acceptance of such changes.</li>
                    <li><strong>Governing Law:</strong>These Terms are governed by and construed in accordance with the laws of the European Union. Any legal action or proceeding arising out of or relating to these Terms shall be brought exclusively in the courts of the European Union, and each party consents to the personal jurisdiction of such courts.</li>
                    <li><strong>Contact:</strong> For questions or concerns regarding these Terms, please contact support@raw.ninja.</li>
                </ol>

                <footer>&copy; 2023 RAW.Ninja LLC. All rights reserved.</footer>
            </Paper>
        </Box>);
}

export default TermsOfService