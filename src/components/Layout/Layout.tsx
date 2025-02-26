import React, { useState } from 'react';
import { Box, Toolbar, useMediaQuery, useTheme } from '@mui/material';
import Header from './Header';
import Sidebar from './Sidebar';

interface LayoutProps {
  children: React.ReactNode;
}

const SIDEBAR_WIDTH = 240;

const Layout: React.FC<LayoutProps> = ({ children }) => {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
  const [isSidebarOpen, setIsSidebarOpen] = useState(!isMobile);
  
  const handleSidebarToggle = () => {
    setIsSidebarOpen(!isSidebarOpen);
  };
  
  return (
    <Box sx={{ display: 'flex', height: '100vh' }}>
      <Header onSidebarToggle={handleSidebarToggle} />
      
      {!isMobile && (
        <Sidebar
          width={SIDEBAR_WIDTH}
          open={isSidebarOpen}
          onClose={isMobile ? handleSidebarToggle : undefined}
        />
      )}
      
      <Box
        component="main"
        sx={{
          flexGrow: 1,
          p: 3,
          width: { sm: `calc(100% - ${SIDEBAR_WIDTH}px)` },
          height: '100vh',
          overflow: 'auto',
          backgroundColor: (theme) => theme.palette.background.default,
        }}
      >
        <Toolbar />
        {children}
      </Box>
    </Box>
  );
};

export default Layout; 